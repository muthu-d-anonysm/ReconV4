package main

import (
	"bufio"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fatih/color"
)

const Version = "4.2.0"

type Config struct {
	Domain         string
	DiscordWebhook string
	Concurrency    int
	OutputDir      string
	StartTime      time.Time
	ScanMonth      string
	UserAgent      string
}

type ScanResults struct {
	Domain             string         `json:"domain"`
	ScanDate           string         `json:"scan_date"`
	Duration           string         `json:"duration"`
	TotalSubdomains    int            `json:"total_subdomains"`
	ResolvedSubdomains int            `json:"resolved_subdomains"`
	LiveHosts          int            `json:"live_hosts"`
	TotalURLs          int            `json:"total_urls"`
	JSFiles            int            `json:"js_files"`
	SourceMaps         int            `json:"source_maps"`
	Secrets            int            `json:"secrets"`
	Vulnerabilities    int            `json:"vulnerabilities"`
	CriticalVulns      int            `json:"critical_vulns"`
	CloudAssets        int            `json:"cloud_assets"`
	SensitiveFiles     int            `json:"sensitive_files"`
	GFMatches          int            `json:"gf_matches"`
	TechStack          map[string]int `json:"tech_stack"`
}

type SecretsFinding struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	File        string `json:"file"`
	FileURL     string `json:"file_url"`
	Line        int    `json:"line"`
	Secret      string `json:"secret"`
}

type SourceMapInfo struct {
	FilePath    string `json:"file_path"`
	OriginalURL string `json:"original_url"`
	SourceMapURL string `json:"source_map_url"`
}

var (
	ctx    context.Context
	cancel context.CancelFunc
)

func main() {
	domain := flag.String("d", "", "Target domain (required)")
	concurrency := flag.Int("c", 0, "Concurrency level (100-800)")
	flag.Parse()

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		color.Red("\n[!] Scan interrupted by user")
		cancel()
		os.Exit(0)
	}()

	printBanner()

	config := getConfig(*domain, *concurrency)

	color.Cyan("\n[✓] Configuration saved")
	color.White("[*] Auto-detected optimal settings:")
	color.White("    - Concurrency: %d", config.Concurrency)
	color.White("    - System: %d cores, Platform: %s", runtime.NumCPU(), runtime.GOOS)
	color.White("\n[*] Starting comprehensive scan for: %s", config.Domain)
	color.White("[*] Estimated time: 45-60 minutes")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	results := &ScanResults{
		Domain:    config.Domain,
		ScanDate:  time.Now().Format("2006-01-02 15:04:05"),
		TechStack: make(map[string]int),
	}

	if err := runPhase1(ctx, config, results); err != nil {
		handleError("Phase 1", err, config)
	}

	if err := runPhase2(ctx, config, results); err != nil {
		handleError("Phase 2", err, config)
	}

	if err := runPhase3(ctx, config, results); err != nil {
		handleError("Phase 3", err, config)
	}

	if err := runPhase4(ctx, config, results); err != nil {
		handleError("Phase 4", err, config)
	}

	if err := runPhase45(ctx, config, results); err != nil {
		handleError("Phase 4.5", err, config)
	}

	if err := runPhase5(ctx, config, results); err != nil {
		handleError("Phase 5", err, config)
	}

	if err := runPhase6(ctx, config, results); err != nil {
		handleError("Phase 6", err, config)
	}

	if err := runPhase65(ctx, config, results); err != nil {
		handleError("Phase 6.5", err, config)
	}

	if err := runPhase7(ctx, config, results); err != nil {
		handleError("Phase 7", err, config)
	}

	if err := runPhase8(ctx, config, results); err != nil {
		handleError("Phase 8", err, config)
	}

	results.Duration = time.Since(config.StartTime).Round(time.Second).String()
	saveResults(config, results)

	printSummary(config, results)

	if config.DiscordWebhook != "" {
		sendDiscordNotification(config, results)
	}

	color.Green("\n[✓] Press Enter to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func printBanner() {
	banner := `
██████╗ ███████╗ ██████╗ ██████╗ ███╗   ██╗██╗   ██╗██╗  ██╗
██╔══██╗██╔════╝██╔════╝██╔═══██╗████╗  ██║██║   ██║██║  ██║
██████╔╝█████╗  ██║     ██║   ██║██╔██╗ ██║██║   ██║███████║
██╔══██╗██╔══╝  ██║     ██║   ██║██║╚██╗██║╚██╗ ██╔╝╚════██║
██║  ██║███████╗╚██████╗╚██████╔╝██║ ╚████║ ╚████╔╝      ██║
╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝  ╚═══╝       ╚═╝
                                                    v%s
        Comprehensive Reconnaissance Automation Tool
                    By: Muthu D
	`
	color.Cyan(banner, Version)
}

func getConfig(domain string, concurrency int) *Config {
	reader := bufio.NewReader(os.Stdin)

	if domain == "" {
		color.Yellow("\n[?] Enter target domain: ")
		domain, _ = reader.ReadString('\n')
		domain = strings.TrimSpace(domain)
	}

	color.Yellow("[?] Discord Webhook URL (press Enter to skip): ")
	webhook, _ := reader.ReadString('\n')
	webhook = strings.TrimSpace(webhook)

	if concurrency == 0 {
		cpuCount := runtime.NumCPU()
		concurrency = cpuCount * 4
		if concurrency > 50 {
			concurrency = 50
		}
		if concurrency < 20 {
			concurrency = 20
		}
	} else {
		if concurrency < 100 {
			concurrency = 100
		}
		if concurrency > 800 {
			concurrency = 800
		}
	}

	domainClean := strings.ReplaceAll(domain, ".", "_")
	scanMonth := time.Now().Format("2006-01")
	outputDir := filepath.Join(".results", domainClean, scanMonth)
	os.MkdirAll(outputDir, 0755)

	return &Config{
		Domain:         domain,
		DiscordWebhook: webhook,
		Concurrency:    concurrency,
		OutputDir:      outputDir,
		StartTime:      time.Now(),
		ScanMonth:      scanMonth,
		UserAgent:      "ReconV4/4.2.0 Professional Bug Bounty Research",
	}
}

func runPhase1(ctx context.Context, config *Config, results *ScanResults) error {
	color.Blue("\n[Phase 1/9] Subdomain Enumeration ⏳")

	subdomainDir := filepath.Join(config.OutputDir, "subdomains")
	os.MkdirAll(subdomainDir, 0755)

	allSubdomains := make(map[string]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	tools := []struct {
		name string
		cmd  []string
	}{
		{"subfinder", []string{"subfinder", "-d", config.Domain, "-silent"}},
		{"assetfinder", []string{"assetfinder", "--subs-only", config.Domain}},
		{"findomain", []string{"findomain", "-t", config.Domain, "-q"}},
	}

	for _, tool := range tools {
		wg.Add(1)
		go func(t struct {
			name string
			cmd  []string
		}) {
			defer wg.Done()

			output, err := runCommand(ctx, t.cmd[0], t.cmd[1:]...)
			if err != nil {
				color.Red("  %s failed: %v", t.name, err)
				return
			}

			subs := strings.Split(strings.TrimSpace(output), "\n")
			
			rawFile := filepath.Join(subdomainDir, fmt.Sprintf("raw_%s.txt", t.name))
			saveToFile(rawFile, subs)

			mu.Lock()
			for _, sub := range subs {
				sub = strings.TrimSpace(sub)
				if sub != "" && strings.Contains(sub, ".") {
					allSubdomains[sub] = true
				}
			}
			mu.Unlock()

			color.Green("  Running %s... ✓ (%d found)", t.name, len(subs))
		}(tool)
	}

	wg.Wait()

	subList := make([]string, 0, len(allSubdomains))
	for sub := range allSubdomains {
		subList = append(subList, sub)
	}
	sort.Strings(subList)

	allSubFile := filepath.Join(subdomainDir, "all_subdomains.txt")
	saveToFile(allSubFile, subList)

	results.TotalSubdomains = len(subList)
	color.Cyan("  Total unique: %d\n", len(subList))

	return nil
}

func runPhase2(ctx context.Context, config *Config, results *ScanResults) error {
	color.Blue("\n[Phase 2/9] DNS Resolution ⏳")

	dnsDir := filepath.Join(config.OutputDir, "dns")
	os.MkdirAll(dnsDir, 0755)

	subdomainFile := filepath.Join(config.OutputDir, "subdomains", "all_subdomains.txt")
	subdomains, err := readLines(subdomainFile)
	if err != nil {
		return err
	}

	color.White("  Validating %d subdomains...", len(subdomains))

	tempFile := filepath.Join(dnsDir, "temp_input.txt")
	saveToFile(tempFile, subdomains)
	defer os.Remove(tempFile)

	resolvedFile := filepath.Join(dnsDir, "resolved_subdomains.txt")
	dnsDetailsFile := filepath.Join(dnsDir, "dns_details.json")

	cmd := exec.CommandContext(ctx, "dnsx",
		"-l", tempFile,
		"-o", resolvedFile,
		"-json",
		"-resp",
		"-t", fmt.Sprintf("%d", config.Concurrency),
	)

	output, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(string(output), "") {
		return fmt.Errorf("dnsx failed: %v", err)
	}

	ioutil.WriteFile(dnsDetailsFile, output, 0644)

	resolved, _ := readLines(resolvedFile)
	results.ResolvedSubdomains = len(resolved)

	percentage := (float64(results.ResolvedSubdomains) / float64(results.TotalSubdomains)) * 100
	color.Green("  Resolved: %d (%.1f%%)\n", results.ResolvedSubdomains, percentage)

	return nil
}

func runPhase3(ctx context.Context, config *Config, results *ScanResults) error {
	color.Blue("\n[Phase 3/9] Live Host Detection ⏳")

	httpxDir := filepath.Join(config.OutputDir, "httpx")
	os.MkdirAll(httpxDir, 0755)

	resolvedFile := filepath.Join(config.OutputDir, "dns", "resolved_subdomains.txt")

	allResolved, err := readLines(resolvedFile)
	if err != nil {
		return err
	}

	filtered := []string{}
	for _, line := range allResolved {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(line), &data); err == nil {
			if host, ok := data["host"].(string); ok {
				if !strings.HasPrefix(host, "*.") && !strings.HasPrefix(host, "*") {
					filtered = append(filtered, host)
				}
			}
		}
	}

	filteredFile := filepath.Join(httpxDir, "filtered_hosts.txt")
	saveToFile(filteredFile, filtered)

	color.White("  Probing %d hosts (filtered from %d)...", len(filtered), len(allResolved))

	liveFile := filepath.Join(httpxDir, "live_hosts.txt")
	jsonFile := filepath.Join(httpxDir, "httpx_results.json")

	cmd := exec.CommandContext(ctx, "httpx",
		"-l", filteredFile,
		"-o", liveFile,
		"-json",
		"-tech-detect",
		"-status-code",
		"-title",
		"-cdn",
		"-threads", fmt.Sprintf("%d", config.Concurrency),
		"-rate-limit", "150",
		"-timeout", "20",
		"-retries", "2",
		"-silent",
		"-header", "User-Agent: Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
	)

	output, err := cmd.CombinedOutput()
	if err != nil && len(output) == 0 {
		return fmt.Errorf("httpx failed: %v", err)
	}

	ioutil.WriteFile(jsonFile, output, 0644)

	// Extract full URLs with schemes and save to live_hosts_urls.txt
	liveHostURLs := []string{}
	techStack := make(map[string]int)
	
	for _, line := range strings.Split(string(output), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(line), &data); err == nil {
			// Extract URL with scheme
			if url, ok := data["url"].(string); ok {
				liveHostURLs = append(liveHostURLs, url)
			}
			
			// Extract tech stack
			if techs, ok := data["tech"].([]interface{}); ok {
				for _, tech := range techs {
					if techStr, ok := tech.(string); ok {
						techStack[techStr]++
					}
				}
			}
		}
	}
	
	// Save live URLs with schemes for URL crawlers
	liveURLsFile := filepath.Join(httpxDir, "live_hosts_urls.txt")
	saveToFile(liveURLsFile, liveHostURLs)
	
	results.TechStack = techStack

	techFile := filepath.Join(httpxDir, "tech_stack.txt")
	techList := []string{}
	for tech, count := range techStack {
		techList = append(techList, fmt.Sprintf("%s (%d)", tech, count))
	}
	saveToFile(techFile, techList)

	liveHosts, _ := readLines(liveFile)
	results.LiveHosts = len(liveHosts)
	percentage := (float64(results.LiveHosts) / float64(len(filtered))) * 100

	color.Green("  Live: %d (%.1f%%)", results.LiveHosts, percentage)

	if len(techStack) > 0 {
		color.Cyan("  Tech stack detected: %s\n", formatTechStack(techStack, 3))
	}

	return nil
}

func runPhase4(ctx context.Context, config *Config, results *ScanResults) error {
	color.Blue("\n[Phase 4/9] URL Discovery ⏳")

	urlsDir := filepath.Join(config.OutputDir, "urls")
	os.MkdirAll(urlsDir, 0755)

	// Use live_hosts_urls.txt which contains full URLs with schemes
	liveURLsFile := filepath.Join(config.OutputDir, "httpx", "live_hosts_urls.txt")
	liveURLs, err := readLines(liveURLsFile)
	if err != nil {
		return err
	}

	allURLs := make(map[string]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	tools := []struct {
		name string
		run  func(context.Context, *Config, string) ([]string, error)
	}{
		{"waybackurls", runWaybackurls},
		{"gau", runGau},
		{"katana", runKatana},
	}

	for _, tool := range tools {
		wg.Add(1)
		go func(t struct {
			name string
			run  func(context.Context, *Config, string) ([]string, error)
		}) {
			defer wg.Done()

			urls := []string{}
			for _, urlStr := range liveURLs {
				hostURLs, err := t.run(ctx, config, urlStr)
				if err != nil {
					continue
				}
				urls = append(urls, hostURLs...)
			}

			mu.Lock()
			for _, url := range urls {
				allURLs[url] = true
			}
			mu.Unlock()

			outputFile := filepath.Join(urlsDir, fmt.Sprintf("%s_urls.txt", t.name))
			saveToFile(outputFile, urls)
			color.Green("  Running %s... ✓ (%d URLs)", t.name, len(urls))
		}(tool)
	}

	wg.Wait()

	urlList := make([]string, 0, len(allURLs))
	for url := range allURLs {
		urlList = append(urlList, url)
	}
	sort.Strings(urlList)

	allURLFile := filepath.Join(urlsDir, "all_urls.txt")
	saveToFile(allURLFile, urlList)

	results.TotalURLs = len(urlList)
	color.Cyan("  Total: %d unique URLs\n", len(urlList))

	return nil
}

func runGau(ctx context.Context, config *Config, urlStr string) ([]string, error) {
	cmd := exec.CommandContext(ctx, "gau",
		"--threads", "100",
		"--subs",
		urlStr,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	
	urls := strings.Split(strings.TrimSpace(string(output)), "\n")
	var filtered []string
	for _, url := range urls {
		url = strings.TrimSpace(url)
		if url != "" && !strings.HasPrefix(url, "time=") && strings.HasPrefix(url, "http") {
			filtered = append(filtered, url)
		}
	}
	return filtered, nil
}

func runWaybackurls(ctx context.Context, config *Config, urlStr string) ([]string, error) {
	cmd := exec.CommandContext(ctx, "waybackurls", urlStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	
	urls := strings.Split(strings.TrimSpace(string(output)), "\n")
	var filtered []string
	for _, url := range urls {
		url = strings.TrimSpace(url)
		if url != "" && strings.HasPrefix(url, "http") {
			filtered = append(filtered, url)
		}
	}
	return filtered, nil
}

func runKatana(ctx context.Context, config *Config, urlStr string) ([]string, error) {
	cmd := exec.CommandContext(ctx, "katana",
		"-u", urlStr,
		"-jc",
		"-d", "3",
		"-silent",
		"-c", fmt.Sprintf("%d", config.Concurrency/10),
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	
	urls := strings.Split(strings.TrimSpace(string(output)), "\n")
	var filtered []string
	for _, url := range urls {
		url = strings.TrimSpace(url)
		if url != "" && strings.HasPrefix(url, "http") {
			filtered = append(filtered, url)
		}
	}
	return filtered, nil
}

func runPhase45(ctx context.Context, config *Config, results *ScanResults) error {
	color.Blue("\n[Phase 4.5/9] JavaScript Analysis ⏳")

	jsDir := filepath.Join(config.OutputDir, "js_files")
	os.MkdirAll(jsDir, 0755)
	os.MkdirAll(filepath.Join(jsDir, "sourcemaps"), 0755)

	secretsDir := filepath.Join(config.OutputDir, "secrets")
	os.MkdirAll(secretsDir, 0755)

	allURLFile := filepath.Join(config.OutputDir, "urls", "all_urls.txt")
	allURLs, err := readLines(allURLFile)
	if err != nil {
		return err
	}

	jsURLs := []string{}
	jsPattern := regexp.MustCompile(`\.(js|jsx|ts)(\?|#|$)|webpack|bundle`)
	for _, url := range allURLs {
		if jsPattern.MatchString(url) {
			jsURLs = append(jsURLs, url)
		}
	}

	jsURLFile := filepath.Join(jsDir, "js_urls.txt")
	saveToFile(jsURLFile, jsURLs)

	color.White("  Found %d JS files", len(jsURLs))

	downloadedFiles := downloadJSFiles(ctx, jsURLs, jsDir, config.Concurrency)
	results.JSFiles = len(downloadedFiles)
	color.Green("  Downloaded %d files", len(downloadedFiles))

	sourceMaps := extractSourceMaps(ctx, downloadedFiles, jsDir)
	results.SourceMaps = len(sourceMaps)
	
	// Save source map info with full URLs
	if len(sourceMaps) > 0 {
		sourceMapInfoFile := filepath.Join(jsDir, "sourcemaps", "sourcemap_urls.json")
		jsonData, _ := json.MarshalIndent(sourceMaps, "", "  ")
		ioutil.WriteFile(sourceMapInfoFile, jsonData, 0644)
		color.Green("  Found %d source maps 🗺️", len(sourceMaps))
	}

	// Extract source map files for endpoint extraction
	sourceMapFiles := []string{}
	for _, sm := range sourceMaps {
		sourceMapFiles = append(sourceMapFiles, sm.FilePath)
	}

	endpoints := extractEndpoints(downloadedFiles, sourceMapFiles)
	endpointFile := filepath.Join(jsDir, "js_endpoints.txt")
	saveToFile(endpointFile, endpoints)
	color.Cyan("  Extracted %d endpoints", len(endpoints))

	secrets := scanForSecrets(ctx, jsDir, secretsDir, downloadedFiles)
	results.Secrets = len(secrets)
	
	// Save secrets with full file URLs
	if len(secrets) > 0 {
		secretsWithURLFile := filepath.Join(secretsDir, "secrets_with_urls.json")
		jsonData, _ := json.MarshalIndent(secrets, "", "  ")
		ioutil.WriteFile(secretsWithURLFile, jsonData, 0644)
		color.Yellow("  Secret scanning... %d secrets found ⚠️\n", len(secrets))
	} else {
		color.Green("  Secret scanning... 0 secrets found ✓\n")
	}

	return nil
}

func downloadJSFiles(ctx context.Context, urls []string, outputDir string, concurrency int) []string {
	downloaded := []string{}
	var mu sync.Mutex
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for _, url := range urls {
		if url == "" {
			continue
		}

		wg.Add(1)
		go func(u string) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			select {
			case <-ctx.Done():
				return
			default:
			}

			hash := md5.Sum([]byte(u))
			filename := hex.EncodeToString(hash[:]) + ".js"
			filepath := filepath.Join(outputDir, filename)

			client := &http.Client{Timeout: 30 * time.Second}
			resp, err := client.Get(u)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				return
			}

			out, err := os.Create(filepath)
			if err != nil {
				return
			}
			defer out.Close()

			_, err = io.Copy(out, resp.Body)
			if err != nil {
				return
			}

			metaFile := filepath + ".meta"
			ioutil.WriteFile(metaFile, []byte(u), 0644)

			mu.Lock()
			downloaded = append(downloaded, filepath)
			mu.Unlock()
		}(url)
	}

	wg.Wait()
	return downloaded
}

func extractSourceMaps(ctx context.Context, jsFiles []string, jsDir string) []SourceMapInfo {
	sourceMaps := []SourceMapInfo{}
	sourceMapPattern := regexp.MustCompile(`sourceMappingURL=([^\s]+)`)

	for _, jsFile := range jsFiles {
		content, err := ioutil.ReadFile(jsFile)
		if err != nil {
			continue
		}

		matches := sourceMapPattern.FindAllStringSubmatch(string(content), -1)
		for _, match := range matches {
			if len(match) > 1 {
				mapURL := match[1]

				metaFile := jsFile + ".meta"
				originalURL, _ := ioutil.ReadFile(metaFile)
				baseURL := string(originalURL)

				if !strings.HasPrefix(mapURL, "http") {
					if strings.HasPrefix(mapURL, "/") {
						parts := strings.Split(baseURL, "/")
						if len(parts) > 2 {
							mapURL = parts[0] + "//" + parts[2] + mapURL
						}
					} else {
						lastSlash := strings.LastIndex(baseURL, "/")
						if lastSlash != -1 {
							mapURL = baseURL[:lastSlash+1] + mapURL
						}
					}
				}

				client := &http.Client{Timeout: 30 * time.Second}
				resp, err := client.Get(mapURL)
				if err != nil {
					continue
				}
				defer resp.Body.Close()

				if resp.StatusCode == 200 {
					hash := md5.Sum([]byte(mapURL))
					mapFile := filepath.Join(jsDir, "sourcemaps", hex.EncodeToString(hash[:])+".map")

					out, err := os.Create(mapFile)
					if err != nil {
						continue
					}
					io.Copy(out, resp.Body)
					out.Close()

					ioutil.WriteFile(mapFile+".meta", []byte(mapURL), 0644)
					
					sourceMaps = append(sourceMaps, SourceMapInfo{
						FilePath:     mapFile,
						OriginalURL:  baseURL,
						SourceMapURL: mapURL,
					})
				}
			}
		}
	}

	return sourceMaps
}

func extractEndpoints(jsFiles []string, sourceMaps []string) []string {
	endpoints := make(map[string]bool)

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`["'](/[a-zA-Z0-9_\-/.]+)["']`),
		regexp.MustCompile(`["'](https?://[^\s"']+)["']`),
		regexp.MustCompile(`["'](api/[a-zA-Z0-9_\-/.]+)["']`),
		regexp.MustCompile(`["'](/v[0-9]/[a-zA-Z0-9_\-/.]+)["']`),
	}

	allFiles := append(jsFiles, sourceMaps...)

	for _, file := range allFiles {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}

		for _, pattern := range patterns {
			matches := pattern.FindAllStringSubmatch(string(content), -1)
			for _, match := range matches {
				if len(match) > 1 {
					endpoint := match[1]
					if len(endpoint) > 5 && len(endpoint) < 200 {
						endpoints[endpoint] = true
					}
				}
			}
		}
	}

	result := []string{}
	for endpoint := range endpoints {
		result = append(result, endpoint)
	}
	sort.Strings(result)

	return result
}

func scanForSecrets(ctx context.Context, jsDir, secretsDir string, jsFiles []string) []SecretsFinding {
	gitleaksFile := filepath.Join(secretsDir, "gitleaks_findings.json")

	cmd := exec.CommandContext(ctx, "gitleaks", "detect",
		"--source", jsDir,
		"--report-path", gitleaksFile,
		"--report-format", "json",
		"--no-git",
	)

	cmd.CombinedOutput()

	findings := []SecretsFinding{}
	data, err := ioutil.ReadFile(gitleaksFile)
	if err != nil {
		return findings
	}

	var gitleaksOutput []map[string]interface{}
	if err := json.Unmarshal(data, &gitleaksOutput); err != nil {
		return findings
	}

	// Create file to URL mapping
	fileURLMap := make(map[string]string)
	for _, jsFile := range jsFiles {
		metaFile := jsFile + ".meta"
		if urlData, err := ioutil.ReadFile(metaFile); err == nil {
			fileURLMap[jsFile] = string(urlData)
		}
	}

	categories := map[string][]SecretsFinding{
		"apikeys":       {},
		"tokens":        {},
		"awskeys":       {},
		"passwords":     {},
		"privatekeys":   {},
		"dbcredentials": {},
	}

	for _, finding := range gitleaksOutput {
		filePath := getString(finding, "File")
		fileURL := fileURLMap[filePath]
		if fileURL == "" {
			fileURL = filePath
		}

		secret := SecretsFinding{
			Type:        getString(finding, "RuleID"),
			Description: getString(finding, "Description"),
			File:        filePath,
			FileURL:     fileURL,
			Line:        getInt(finding, "StartLine"),
			Secret:      getString(finding, "Secret"),
		}
		findings = append(findings, secret)

		ruleID := strings.ToLower(secret.Type)
		if strings.Contains(ruleID, "api") || strings.Contains(ruleID, "key") {
			categories["apikeys"] = append(categories["apikeys"], secret)
		} else if strings.Contains(ruleID, "token") {
			categories["tokens"] = append(categories["tokens"], secret)
		} else if strings.Contains(ruleID, "aws") {
			categories["awskeys"] = append(categories["awskeys"], secret)
		} else if strings.Contains(ruleID, "password") {
			categories["passwords"] = append(categories["passwords"], secret)
		} else if strings.Contains(ruleID, "private") {
			categories["privatekeys"] = append(categories["privatekeys"], secret)
		}
	}

	for category, secrets := range categories {
		if len(secrets) > 0 {
			categoryFile := filepath.Join(secretsDir, category+".json")
			jsonData, _ := json.MarshalIndent(secrets, "", "  ")
			ioutil.WriteFile(categoryFile, jsonData, 0644)
		}
	}

	return findings
}

func runPhase5(ctx context.Context, config *Config, results *ScanResults) error {
	color.Blue("\n[Phase 5/9] Vulnerability Scanning ⏳")

	nucleiDir := filepath.Join(config.OutputDir, "nuclei")
	os.MkdirAll(nucleiDir, 0755)

	liveURLsFile := filepath.Join(config.OutputDir, "httpx", "live_hosts_urls.txt")

	color.White("  Updating nuclei templates...")
	runCommand(ctx, "nuclei", "-update-templates")

	color.White("  Scanning hosts...")

	findingsFile := filepath.Join(nucleiDir, "findings.json")

	cmd := exec.CommandContext(ctx, "nuclei",
		"-l", liveURLsFile,
		"-j",
		"-o", findingsFile,
		"-tags", "cve,misconfig,exposure,takeover,cms,vuln,default-login",
		"-severity", "critical,high,medium",
		"-c", fmt.Sprintf("%d", 100),
		"-rate-limit", "150",
		"-timeout", "10",
		"-retries", "1",
		"-silent",
	)

	output, err := cmd.CombinedOutput()
	if err != nil && len(output) == 0 {
		return fmt.Errorf("nuclei failed: %v", err)
	}

	findings := []map[string]interface{}{}
	data, _ := ioutil.ReadFile(findingsFile)

	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}

		var finding map[string]interface{}
		if err := json.Unmarshal([]byte(line), &finding); err == nil {
			findings = append(findings, finding)
		}
	}

	categories := map[string][]map[string]interface{}{
		"cves":         {},
		"misconfigs":   {},
		"takeovers":    {},
		"cmsvulns":     {},
		"exposures":    {},
		"defaultcreds": {},
	}

	criticalCount := 0
	for _, finding := range findings {
		tags := getString(finding, "tags")
		severity := strings.ToLower(getString(finding, "severity"))

		if severity == "critical" {
			criticalCount++
		}

		if strings.Contains(tags, "cve") {
			categories["cves"] = append(categories["cves"], finding)
		} else if strings.Contains(tags, "misconfig") {
			categories["misconfigs"] = append(categories["misconfigs"], finding)
		} else if strings.Contains(tags, "takeover") {
			categories["takeovers"] = append(categories["takeovers"], finding)
		} else if strings.Contains(tags, "cms") {
			categories["cmsvulns"] = append(categories["cmsvulns"], finding)
		} else if strings.Contains(tags, "exposure") {
			categories["exposures"] = append(categories["exposures"], finding)
		} else if strings.Contains(tags, "default-login") {
			categories["defaultcreds"] = append(categories["defaultcreds"], finding)
		}
	}

	for category, items := range categories {
		if len(items) > 0 {
			categoryFile := filepath.Join(nucleiDir, category+".json")
			jsonData, _ := json.MarshalIndent(items, "", "  ")
			ioutil.WriteFile(categoryFile, jsonData, 0644)
		}
	}

	results.Vulnerabilities = len(findings)
	results.CriticalVulns = criticalCount

	if len(findings) > 0 {
		color.Yellow("  Found %d vulnerabilities (%d critical) ⚠️\n", len(findings), criticalCount)
	} else {
		color.Green("  Found 0 vulnerabilities ✓\n")
	}

	return nil
}

func runPhase6(ctx context.Context, config *Config, results *ScanResults) error {
	color.Blue("\n[Phase 6/9] Cloud Asset Discovery ⏳")

	cloudDir := filepath.Join(config.OutputDir, "cloud")
	os.MkdirAll(cloudDir, 0755)

	domain := strings.ReplaceAll(config.Domain, ".", "")
	domainDash := strings.ReplaceAll(config.Domain, ".", "-")

	keywords := []string{
		config.Domain,
		domain,
		domainDash,
	}

	keywordsFile := filepath.Join(cloudDir, "keywords.txt")
	saveToFile(keywordsFile, keywords)

	bucketsFile := filepath.Join(cloudDir, "buckets.txt")

	color.White("  Searching S3/Azure/GCP...")

	cmd := exec.CommandContext(ctx, "cloud_enum",
		"-k", keywordsFile,
		"-l", bucketsFile,
		"-t", "100",
	)

	cmd.CombinedOutput()

	buckets, _ := readLines(bucketsFile)
	results.CloudAssets = len(buckets)

	if len(buckets) > 0 {
		color.Green("  Found %d cloud assets 🌩️\n", len(buckets))
	} else {
		color.Cyan("  Found 0 cloud assets\n")
	}

	return nil
}

func runPhase65(ctx context.Context, config *Config, results *ScanResults) error {
	color.Blue("\n[Phase 6.5/9] Sensitive File Discovery ⏳")

	sensitiveDir := filepath.Join(config.OutputDir, "sensitive")
	os.MkdirAll(sensitiveDir, 0755)

	allURLFile := filepath.Join(config.OutputDir, "urls", "all_urls.txt")
	allURLs, err := readLines(allURLFile)
	if err != nil {
		return err
	}

	color.White("  Scanning %d URLs for sensitive extensions...", len(allURLs))

	categories := map[string][]string{
		"configfiles":   {".xml", ".config", ".conf", ".ini", ".yaml", ".yml", ".toml", ".properties", ".env"},
		"documents":     {".doc", ".docx", ".xls", ".xlsx", ".pdf", ".txt", ".csv", ".ppt", ".pptx"},
		"backups":       {".bak", ".backup", ".old", ".orig", ".save", ".swp", ".tmp"},
		"databases":     {".sql", ".db", ".sqlite", ".mdb", ".accdb"},
		"sourcecontrol": {".git", ".svn", ".gitignore", ".htaccess", ".htpasswd"},
		"archives":      {".zip", ".tar", ".gz", ".rar", ".7z", ".bz2"},
		"logs":          {".log", ".logs", ".trace"},
		"keyscerts":     {".key", ".pem", ".crt", ".cer", ".p12", ".pfx"},
	}

	categoryEmojis := map[string]string{
		"configfiles":   "⚙️",
		"documents":     "📄",
		"backups":       "💾",
		"databases":     "🗄️",
		"sourcecontrol": "🔀",
		"archives":      "📦",
		"logs":          "📝",
		"keyscerts":     "🔑",
	}

	foundFiles := make(map[string][]string)
	allSensitive := []string{}

	for _, url := range allURLs {
		urlLower := strings.ToLower(url)
		for category, extensions := range categories {
			for _, ext := range extensions {
				if strings.HasSuffix(urlLower, ext) || strings.Contains(urlLower, ext+"?") {
					foundFiles[category] = append(foundFiles[category], url)
					allSensitive = append(allSensitive, url)
					break
				}
			}
		}
	}

	totalCount := 0
	for category, files := range foundFiles {
		if len(files) > 0 {
			categoryFile := filepath.Join(sensitiveDir, category+".txt")
			saveToFile(categoryFile, files)

			emoji := categoryEmojis[category]
			color.Green("  %s %s: %d found", emoji, strings.ReplaceAll(category, "", " "), len(files))
			totalCount += len(files)
		}
	}

	allSensitive = deduplicateStrings(allSensitive)
	allSensitiveFile := filepath.Join(sensitiveDir, "all_sensitive.txt")
	saveToFile(allSensitiveFile, allSensitive)

	summaryData := map[string]interface{}{
		"total_sensitive_files": len(allSensitive),
		"categories":            foundFiles,
	}
	summaryFile := filepath.Join(sensitiveDir, "summary.json")
	jsonData, _ := json.MarshalIndent(summaryData, "", "  ")
	ioutil.WriteFile(summaryFile, jsonData, 0644)

	results.SensitiveFiles = len(allSensitive)

	if totalCount > 0 {
		color.Cyan("  Total: %d sensitive files found\n", totalCount)
	} else {
		color.Green("  Total: 0 sensitive files found\n")
	}

	return nil
}

func runPhase7(ctx context.Context, config *Config, results *ScanResults) error {
	color.Blue("\n[Phase 7/9] Pattern Filtering ⏳")

	gfDir := filepath.Join(config.OutputDir, "gf_patterns")
	os.MkdirAll(gfDir, 0755)

	allURLFile := filepath.Join(config.OutputDir, "urls", "all_urls.txt")

	patterns := []string{"xss", "sqli", "ssrf", "ssti", "lfi", "redirect", "rce", "idor", "debug", "api"}

	color.White("  Applying GF patterns...")

	totalMatches := 0
	for _, pattern := range patterns {
		outputFile := filepath.Join(gfDir, pattern+"_candidates.txt")

		cmd := exec.CommandContext(ctx, "sh", "-c",
			fmt.Sprintf("cat %s | gf %s > %s 2>/dev/null", allURLFile, pattern, outputFile))
		cmd.Run()

		matches, _ := readLines(outputFile)
		if len(matches) > 0 {
			totalMatches += len(matches)
		}
	}

	results.GFMatches = totalMatches

	if totalMatches > 0 {
		color.Green("  %d potential vulnerability candidates 🎯\n", totalMatches)
	} else {
		color.Cyan("  0 pattern matches\n")
	}

	return nil
}

func runPhase8(ctx context.Context, config *Config, results *ScanResults) error {
	color.Blue("\n[Phase 8/9] Monthly Comparison ⏳")

	domainClean := strings.ReplaceAll(config.Domain, ".", "_")
	resultsBase := filepath.Join(".results", domainClean)

	currentMonth := time.Now()
	prevMonth := currentMonth.AddDate(0, -1, 0)
	prevMonthStr := prevMonth.Format("2006-01")

	prevScanDir := filepath.Join(resultsBase, prevMonthStr)
	prevResultsFile := filepath.Join(prevScanDir, "scan_results.json")

	if _, err := os.Stat(prevResultsFile); os.IsNotExist(err) {
		color.Cyan("  No previous scan detected (first scan)\n")
		return nil
	}

	color.White("  Previous scan detected: %s", prevMonthStr)
	color.White("  Generating comparison...")

	prevData, err := ioutil.ReadFile(prevResultsFile)
	if err != nil {
		return err
	}

	var prevResults ScanResults
	json.Unmarshal(prevData, &prevResults)

	comparison := map[string]interface{}{
		"current_scan":  config.ScanMonth,
		"previous_scan": prevMonthStr,
		"differences": map[string]int{
			"subdomains":      results.TotalSubdomains - prevResults.TotalSubdomains,
			"live_hosts":      results.LiveHosts - prevResults.LiveHosts,
			"urls":            results.TotalURLs - prevResults.TotalURLs,
			"vulnerabilities": results.Vulnerabilities - prevResults.Vulnerabilities,
			"secrets":         results.Secrets - prevResults.Secrets,
		},
	}

	comparisonDir := filepath.Join(resultsBase, "comparison")
	os.MkdirAll(comparisonDir, 0755)

	comparisonFile := filepath.Join(comparisonDir,
		fmt.Sprintf("%s_vs_%s.json", config.ScanMonth, prevMonthStr))
	jsonData, _ := json.MarshalIndent(comparison, "", "  ")
	ioutil.WriteFile(comparisonFile, jsonData, 0644)

	diffs := comparison["differences"].(map[string]int)

	color.Green("  Comparison generated")
	if diffs["subdomains"] > 0 {
		color.Cyan("    +%d new subdomains", diffs["subdomains"])
	}
	if diffs["urls"] > 0 {
		color.Cyan("    +%d new URLs", diffs["urls"])
	}
	if diffs["vulnerabilities"] > 0 {
		color.Yellow("    +%d new vulnerabilities ⚠️", diffs["vulnerabilities"])
	}

	fmt.Println()
	return nil
}

func runCommand(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, scanner.Err()
}

func saveToFile(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	return writer.Flush()
}

func deduplicateStrings(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range input {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	sort.Strings(result)
	return result
}

func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getInt(m map[string]interface{}, key string) int {
	if val, ok := m[key]; ok {
		if num, ok := val.(float64); ok {
			return int(num)
		}
	}
	return 0
}

func formatTechStack(techStack map[string]int, limit int) string {
	type kv struct {
		Key   string
		Value int
	}

	var sorted []kv
	for k, v := range techStack {
		sorted = append(sorted, kv{k, v})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})

	result := []string{}
	for i, item := range sorted {
		if i >= limit {
			break
		}
		result = append(result, fmt.Sprintf("%s (%d)", item.Key, item.Value))
	}

	return strings.Join(result, ", ")
}

func handleError(phase string, err error, config *Config) {
	color.Red("  [!] %s error: %v", phase, err)
	color.Yellow("  [*] Continuing with next phase...")
}

func saveResults(config *Config, results *ScanResults) {
	resultsFile := filepath.Join(config.OutputDir, "scan_results.json")
	jsonData, _ := json.MarshalIndent(results, "", "  ")
	ioutil.WriteFile(resultsFile, jsonData, 0644)
}

func printSummary(config *Config, results *ScanResults) {
	fmt.Println()
	color.Green("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	color.Green("SCAN COMPLETE! 🎉")
	color.Green("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	color.Cyan("\nDuration: %s", results.Duration)
	color.White("\n📊 Results Summary:")
	color.White("  Target: %s", results.Domain)
	color.White("  Subdomains: %d (%d live)", results.TotalSubdomains, results.LiveHosts)
	color.White("  URLs: %d", results.TotalURLs)

	if results.Vulnerabilities > 0 {
		if results.CriticalVulns > 0 {
			color.Yellow("  Vulnerabilities: %d (%d critical) ⚠️", results.Vulnerabilities, results.CriticalVulns)
		} else {
			color.White("  Vulnerabilities: %d", results.Vulnerabilities)
		}
	}

	if results.Secrets > 0 {
		color.Yellow("  Secrets: %d ⚠️", results.Secrets)
	}

	if results.SensitiveFiles > 0 {
		color.Yellow("  Sensitive Files: %d ⚠️", results.SensitiveFiles)
	}

	if results.CloudAssets > 0 {
		color.White("  Cloud Assets: %d", results.CloudAssets)
	}

	color.Cyan("\n📁 %s", config.OutputDir)
}

func sendDiscordNotification(config *Config, results *ScanResults) {
	if config.DiscordWebhook == "" {
		return
	}

	emoji := "✅"
	if results.CriticalVulns > 0 || results.Secrets > 0 {
		emoji = "⚠️"
	}

	message := fmt.Sprintf(`%s **RECONV4 SCAN COMPLETE** %s

**Target:** %s
**Duration:** %s

**📊 Results:**
• **Subdomains:** %d total, %d live (%.1f%% live)
• **URLs Discovered:** %d
• **JS Files:** %d
• **Source Maps:** %d
• **Endpoints Extracted:** Available in results
• **Vulnerabilities:** %d%s
• **Secrets Found:** %d%s
• **Sensitive Files:** %d
• **Cloud Assets:** %d
• **GF Pattern Matches:** %d

**🔍 Tech Stack Detected:**
%s

**📂 Results Location:**
%s`,
		emoji,
		emoji,
		results.Domain,
		results.Duration,
		results.TotalSubdomains,
		results.LiveHosts,
		float64(results.LiveHosts)/float64(results.TotalSubdomains)*100,
		results.TotalURLs,
		results.JSFiles,
		results.SourceMaps,
		results.Vulnerabilities,
		func() string {
			if results.CriticalVulns > 0 {
				return fmt.Sprintf(" (%d critical)", results.CriticalVulns)
			}
			return ""
		}(),
		results.Secrets,
		func() string {
			if results.Secrets > 0 {
				return " ⚠️"
			}
			return ""
		}(),
		results.SensitiveFiles,
		results.CloudAssets,
		results.GFMatches,
		formatTechStack(results.TechStack, 5),
		config.OutputDir,
	)

	payload := map[string]interface{}{
		"content": message,
	}

	jsonData, _ := json.Marshal(payload)

	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("POST", config.DiscordWebhook, strings.NewReader(string(jsonData)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == 204 || resp.StatusCode == 200 {
			color.Green("\n[✓] Discord notification sent")
		}
	}
}
