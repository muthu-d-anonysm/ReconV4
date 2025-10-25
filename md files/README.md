# Reconv4 - Comprehensive Reconnaissance Automation Tool

[![Version](https://img.shields.io/badge/version-4.0.0-blue.svg)](https://github.com/yourusername/reconv4)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.21+-00ADD8.svg)](https://go.dev/)

**Reconv4** is a production-ready, crash-proof reconnaissance automation tool designed for bug bounty hunters and security researchers. It performs comprehensive subdomain enumeration, vulnerability scanning, JavaScript analysis, secret detection, and monthly comparison tracking - all in one unified workflow.

## 🎯 Key Features

- **9 Automated Phases** - Complete reconnaissance workflow from subdomain enumeration to monthly comparison
- **15+ Integrated Tools** - Best-in-class tools orchestrated seamlessly
- **Crash-Proof Design** - Intelligent error handling and graceful degradation
- **High Performance** - Process 1000 domains in ~3.5-4 hours
- **Monthly Tracking** - Automatic comparison with previous scans
- **Discord Notifications** - Real-time progress updates
- **Zero Configuration** - Auto-detects optimal settings for your system

## 📊 Workflow Overview

```
Phase 1: Subdomain Enumeration (subfinder, assetfinder, findomain)
Phase 2: DNS Resolution (dnsx)
Phase 3: Live Host Detection (httpx with tech stack detection)
Phase 4: URL Discovery (gau, waybackurls, katana)
Phase 4.5: JavaScript Analysis (JS files, source maps, endpoints, secrets)
Phase 5: Vulnerability Scanning (nuclei with auto-categorization)
Phase 6: Cloud Asset Discovery (cloud_enum for S3/Azure/GCP)
Phase 6.5: Sensitive File Discovery (config files, backups, keys, etc.)
Phase 7: GF Pattern Filtering (XSS, SQLi, SSRF, etc.)
Phase 8: Monthly Comparison (delta detection from previous scans)
```

## 🚀 Quick Start

### Installation

```bash
# Clone repository
git clone https://github.com/yourusername/reconv4.git
cd reconv4

# Run installation script (installs all dependencies)
chmod +x install.sh
./install.sh

# Reload shell
source ~/.bashrc
```

### Usage

```bash
# Simple - just run the tool
reconv4

# You'll be prompted for:
# 1. Target domain (e.g., example.com)
# 2. Discord webhook URL (optional)
```

That's it! The tool auto-configures everything else.

## 📋 System Requirements

- **OS**: Linux (Ubuntu/Debian/Kali) or macOS
- **CPU**: 4+ cores (12+ cores recommended)
- **RAM**: 8GB minimum (16GB recommended)
- **Disk**: 10GB free space
- **Go**: 1.21+ (auto-installed by install.sh)

## 🔧 Installed Tools

The installation script automatically installs:

### Subdomain Enumeration
- `subfinder` - Fast passive subdomain enumeration
- `assetfinder` - Find domains and subdomains
- `findomain` - Cross-platform subdomain enumerator

### DNS & HTTP
- `dnsx` - Fast DNS toolkit
- `httpx` - HTTP toolkit with tech detection

### URL Discovery
- `gau` - Fetch known URLs from AlienVault's OTX, Wayback Machine, etc.
- `waybackurls` - Fetch URLs from Wayback Machine
- `katana` - Next-generation crawling framework

### Vulnerability Scanning
- `nuclei` - Fast vulnerability scanner with 5000+ templates

### Secret Detection
- `gitleaks` - SAST tool for detecting hardcoded secrets

### Cloud Discovery
- `cloud_enum` - Multi-cloud OSINT tool (S3, Azure, GCP)

### Pattern Matching
- `gf` - Wrapper for grep with security patterns
- `Gf-Patterns` - Community security patterns

## 📁 Output Structure

```
.results/
└── example_com/
    ├── 2025-10/                    # Current scan
    │   ├── subdomains/
    │   │   ├── all_subdomains.txt
    │   │   ├── raw_subfinder.txt
    │   │   ├── raw_assetfinder.txt
    │   │   └── raw_findomain.txt
    │   ├── dns/
    │   │   ├── resolved_subdomains.txt
    │   │   └── dns_details.json
    │   ├── httpx/
    │   │   ├── live_hosts.txt
    │   │   ├── httpx_results.json
    │   │   └── tech_stack.txt
    │   ├── urls/
    │   │   ├── all_urls.txt
    │   │   ├── gau_urls.txt
    │   │   ├── waybackurls.txt
    │   │   └── katana_urls.txt
    │   ├── js_files/
    │   │   ├── *.js (downloaded files)
    │   │   ├── js_urls.txt
    │   │   ├── js_endpoints.txt
    │   │   └── sourcemaps/*.map
    │   ├── secrets/
    │   │   ├── gitleaks_findings.json
    │   │   ├── api_keys.json
    │   │   ├── tokens.json
    │   │   ├── aws_keys.json
    │   │   └── custom_findings.json
    │   ├── nuclei/
    │   │   ├── findings.json
    │   │   ├── cves.json
    │   │   ├── misconfigs.json
    │   │   ├── takeovers.json
    │   │   └── cms_vulns.json
    │   ├── cloud/
    │   │   ├── buckets.txt
    │   │   └── keywords.txt
    │   ├── sensitive/
    │   │   ├── all_sensitive.txt
    │   │   ├── config_files.txt
    │   │   ├── backups.txt
    │   │   ├── databases.txt
    │   │   ├── keys_certs.txt
    │   │   └── summary.json
    │   ├── gf_patterns/
    │   │   ├── xss_candidates.txt
    │   │   ├── sqli_candidates.txt
    │   │   ├── ssrf_candidates.txt
    │   │   └── ... (10+ patterns)
    │   └── scan_results.json       # Complete summary
    ├── 2025-09/                    # Previous scan
    ├── comparison/
    │   └── 2025-10_vs_2025-09.json
    └── timeline.json
```

## 🎨 Example Output

```
██████╗ ███████╗ ██████╗ ██████╗ ███╗   ██╗██╗   ██╗██╗  ██╗
██╔══██╗██╔════╝██╔════╝██╔═══██╗████╗  ██║██║   ██║██║  ██║
██████╔╝█████╗  ██║     ██║   ██║██╔██╗ ██║██║   ██║███████║
██╔══██╗██╔══╝  ██║     ██║   ██║██║╚██╗██║╚██╗ ██╔╝╚════██║
██║  ██║███████╗╚██████╗╚██████╔╝██║ ╚████║ ╚████╔╝      ██║
╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝  ╚═══╝       ╚═╝

[?] Enter target domain: example.com
[?] Discord Webhook URL (press Enter to skip): 

[✓] Configuration saved
[*] Auto-detected optimal settings:
    - Concurrency: 40
    - System: 12 cores, Platform: linux
[*] Starting comprehensive scan for: example.com
[*] Estimated time: 45-60 minutes

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

[Phase 1/9] Subdomain Enumeration ⏳
  Running subfinder... ✓ (1,234 found)
  Running assetfinder... ✓ (987 found)
  Running findomain... ✓ (1,456 found)
  Total unique: 2,345

[Phase 2/9] DNS Resolution ⏳
  Validating 2,345 subdomains...
  Resolved: 1,023 (43.6%)

[Phase 3/9] Live Host Detection ⏳
  Probing hosts...
  Live: 456 (44.6%)
  Tech stack detected: WordPress (12), React (8), Nginx (45)

[Phase 4/9] URL Discovery ⏳
  Running gau... ✓ (3,456 URLs)
  Running waybackurls... ✓ (2,134 URLs)
  Running katana... ✓ (1,890 URLs)
  Total: 5,234 unique URLs

[Phase 4.5/9] JavaScript Analysis ⏳
  Found 234 JS files
  Downloaded 189 files
  Found 3 source maps ✓
  Extracted 678 endpoints
  Secret scanning... 12 secrets found ⚠️

[Phase 5/9] Vulnerability Scanning ⏳
  Updating nuclei templates...
  Scanning hosts...
  Found 23 vulnerabilities (2 critical ⚠️)

[Phase 6/9] Cloud Asset Discovery ⏳
  Searching S3/Azure/GCP...
  Found 3 cloud assets ✓

[Phase 6.5/9] Sensitive File Discovery ⏳
  Scanning 5,234 URLs for sensitive extensions...
  📋 config files: 23 found
  💾 backups: 8 found
  🔐 keys/certs: 2 found
  Total: 115 sensitive files found

[Phase 7/9] Pattern Filtering ⏳
  Applying GF patterns...
  89 potential vulnerability candidates ✓

[Phase 8/9] Monthly Comparison ⏳
  Previous scan detected (2025-09)
  Generating comparison...
  Comparison generated ✓
    +45 new subdomains
    +234 new URLs
    +3 new vulnerabilities ⚠️

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ SCAN COMPLETE!
Duration: 38m 42s

📊 Results Summary:
  Target: example.com
  Subdomains: 2,345 (456 live)
  URLs: 5,234
  Vulnerabilities: 23 (2 critical ⚠️)
  Secrets: 12 ⚠️
  Sensitive Files: 115 ⚠️
  Cloud Assets: 3

Results: .results/example_com/2025-10/

🔔 Discord notification sent ✓

Press Enter to exit...
```

## ⚙️ Advanced Configuration

### Custom Concurrency

By default, Reconv4 auto-detects optimal concurrency based on your CPU cores. To override:

Edit `main.go` and modify the `getConfig()` function:

```go
concurrency := 50 // Force 50 concurrent operations
```

### Tool-Specific Adjustments

Each phase has different concurrency settings optimized for the operation type:

- **Subdomain enum**: 50 (I/O bound)
- **DNS resolution**: 40 (Network bound)
- **HTTP probing**: 50 (Network bound)
- **URL crawling**: 30 (CPU + Network)
- **Nuclei scanning**: 25 (CPU intensive)

## 🐛 Troubleshooting

### Tool Not Found

If a tool is missing after installation:

```bash
# Reload shell
source ~/.bashrc

# Or add to PATH manually
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin

# Verify tool installation
which subfinder
which nuclei
```

### Nuclei Templates Not Updating

```bash
# Manually update templates
nuclei -update-templates

# Check templates location
ls ~/.nuclei-templates/
```

### GF Patterns Not Working

```bash
# Reinstall GF patterns
rm -rf ~/.gf
mkdir ~/.gf
cd ~/.gf
git clone https://github.com/1ndianl33t/Gf-Patterns.git
mv Gf-Patterns/*.json .
```

### Cloud Enum Issues

```bash
# Test cloud_enum
cloud_enum -h

# If not working, reinstall
cd ~/tools
git clone https://github.com/initstring/cloud_enum.git
cd cloud_enum
pip3 install -r requirements.txt
```

## 📈 Performance Optimization

### For Large-Scale Scans (1000+ domains)

1. **Increase system resources**
   - 16+ CPU cores
   - 32GB RAM
   - SSD storage

2. **Adjust rate limits**
   - Edit httpx rate-limit in `main.go` (default: 150)
   - Edit nuclei rate-limit (default: 150)

3. **Run in background**
   ```bash
   nohup reconv4 > reconv4.log 2>&1 &
   ```

### HP Victus Optimization (Your Setup)

Your system (12 cores, 16GB RAM, RTX GPU) is well-suited for Reconv4:

- **Expected performance**: 1000 domains in ~3.5-4 hours
- **Optimal concurrency**: 40-50 (auto-detected)
- **Memory usage**: ~8-10GB during peak operations

## 🔐 Security Considerations

### API Keys & Tokens

Some tools benefit from API keys for better results:

**Subfinder** (`~/.config/subfinder/provider-config.yaml`):
```yaml
shodan:
  - API_KEY_HERE
censys:
  - API_ID_HERE
  - API_SECRET_HERE
virustotal:
  - API_KEY_HERE
```

**GitHub Token** for waybackurls/gau:
```bash
export GITHUB_TOKEN=your_token_here
```

### Responsible Usage

- ⚠️ **Only scan domains you have permission to test**
- 🚫 **Do not scan government or military domains without authorization**
- ✅ **Respect rate limits and robots.txt**
- ✅ **Use responsibly in bug bounty programs**

## 📚 Integration with Bug Bounty Workflow

### Recommended Workflow

1. **Initial Scan**
   ```bash
   reconv4
   # Enter: target.com
   ```

2. **Review Critical Findings**
   ```bash
   # Check critical vulnerabilities
   cat .results/target_com/2025-10/nuclei/findings.json | jq '.[] | select(.severity=="critical")'

   # Check secrets
   cat .results/target_com/2025-10/secrets/gitleaks_findings.json

   # Check sensitive files
   cat .results/target_com/2025-10/sensitive/keys_certs.txt
   ```

3. **Manual Validation**
   - Validate vulnerabilities with Burp Suite
   - Test sensitive endpoints manually
   - Verify cloud bucket permissions

4. **Monthly Rescans**
   - Run Reconv4 monthly
   - Automatic comparison shows new attack surface
   - Focus on new subdomains/URLs/vulns

### Discord Integration

Get real-time updates on your Discord server:

1. Create a webhook in Discord Server Settings → Integrations
2. Copy webhook URL
3. Paste when prompted by Reconv4

## 🤝 Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

This tool integrates and orchestrates many excellent open-source projects:

- [ProjectDiscovery](https://github.com/projectdiscovery) - subfinder, dnsx, httpx, katana, nuclei
- [Tom Hudson](https://github.com/tomnomnom) - assetfinder, waybackurls, gf
- [Findomain](https://github.com/Findomain/Findomain)
- [Gitleaks](https://github.com/gitleaks/gitleaks)
- [cloud_enum](https://github.com/initstring/cloud_enum)
- All GF pattern contributors

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/yourusername/reconv4/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/reconv4/discussions)
- **Twitter**: [@yourusername](https://twitter.com/yourusername)

## 🗺️ Roadmap

- [ ] Web UI dashboard
- [ ] API endpoint for CI/CD integration
- [ ] Kubernetes/Docker deployment
- [ ] Multi-target parallel scanning
- [ ] Custom nuclei template support
- [ ] Automated report generation (PDF/HTML)
- [ ] Slack/Telegram integration
- [ ] Database backend for historical tracking

---

**⚡ Built for speed. Designed for scale. Optimized for bug bounties.**

*Made with ❤️ by a security researcher, for security researchers*
