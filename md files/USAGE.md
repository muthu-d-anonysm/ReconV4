# Reconv4 - Detailed Usage Guide

## Table of Contents
1. [Basic Usage](#basic-usage)
2. [Understanding Output](#understanding-output)
3. [Advanced Workflows](#advanced-workflows)
4. [Integration Examples](#integration-examples)
5. [Tips & Tricks](#tips--tricks)

---

## Basic Usage

### First Run

```bash
# Start the tool
reconv4

# You'll see:
[?] Enter target domain: example.com
[?] Discord Webhook URL (press Enter to skip): 

# That's it! The tool handles everything else automatically.
```

### With Discord Notifications

```bash
reconv4

[?] Enter target domain: bugcrowd.com
[?] Discord Webhook URL: https://discord.com/api/webhooks/YOUR_WEBHOOK_URL
```

Discord will receive updates after each phase and a final summary.

---

## Understanding Output

### Directory Structure Explained

```
.results/example_com/2025-10/
```

Each scan creates a timestamped folder (YYYY-MM format) with the following structure:

#### 1. Subdomains (`subdomains/`)

```bash
# All unique subdomains found
cat subdomains/all_subdomains.txt

# Raw output from each tool
cat subdomains/raw_subfinder.txt
cat subdomains/raw_assetfinder.txt
cat subdomains/raw_findomain.txt
```

**Example output:**
```
api.example.com
dev.example.com
staging.example.com
```

#### 2. DNS (`dns/`)

```bash
# Only DNS-resolved subdomains
cat dns/resolved_subdomains.txt

# Detailed DNS information (A, AAAA, CNAME records)
cat dns/dns_details.json
```

**Example JSON:**
```json
{
  "host": "api.example.com",
  "a": ["192.168.1.100"],
  "aaaa": ["2001:db8::1"],
  "cname": []
}
```

#### 3. HTTP(X) (`httpx/`)

```bash
# Live HTTP/HTTPS hosts
cat httpx/live_hosts.txt

# Full details with status codes, tech stack
cat httpx/httpx_results.json

# Technology summary
cat httpx/tech_stack.txt
```

**Tech stack example:**
```
WordPress: 12
React: 8
Nginx: 45
Cloudflare: 23
```

#### 4. URLs (`urls/`)

```bash
# All discovered URLs
cat urls/all_urls.txt

# By tool
cat urls/gau_urls.txt
cat urls/waybackurls.txt
cat urls/katana_urls.txt
```

**Example URLs:**
```
https://example.com/api/v1/users
https://example.com/admin/login
https://example.com/static/js/bundle.js
```

#### 5. JavaScript Files (`js_files/`)

```bash
# List of JS URLs
cat js_files/js_urls.txt

# Downloaded JS files (MD5 hash names)
ls js_files/*.js

# Original URL metadata
cat js_files/HASH.js.meta

# Extracted API endpoints
cat js_files/js_endpoints.txt

# Source maps
ls js_files/sourcemaps/*.map
```

**Endpoints example:**
```
/api/v1/authenticate
/api/v2/users/{id}
/graphql
```

#### 6. Secrets (`secrets/`)

```bash
# All findings from gitleaks
cat secrets/gitleaks_findings.json

# Categorized secrets
cat secrets/api_keys.json
cat secrets/tokens.json
cat secrets/aws_keys.json
cat secrets/passwords.json

# Custom patterns (S3 buckets, Firebase, IPs)
cat secrets/custom_findings.json
```

**Critical! Review immediately:**
```bash
# Quick check for high-value secrets
jq '.[] | select(.Description | contains("AWS"))' secrets/gitleaks_findings.json
```

#### 7. Vulnerabilities (`nuclei/`)

```bash
# All findings
cat nuclei/findings.json

# Categorized by type
cat nuclei/cves.json           # CVE-based vulnerabilities
cat nuclei/misconfigs.json     # Misconfigurations
cat nuclei/takeovers.json      # Subdomain takeover opportunities
cat nuclei/cms_vulns.json      # CMS-specific vulnerabilities
cat nuclei/exposures.json      # Exposed sensitive info
cat nuclei/default_creds.json  # Default credentials
```

**Finding critical CVEs:**
```bash
jq '.[] | select(.severity=="critical")' nuclei/findings.json
```

#### 8. Cloud Assets (`cloud/`)

```bash
# Discovered S3/Azure/GCP buckets
cat cloud/buckets.txt

# Keywords used for enumeration
cat cloud/keywords.txt
```

**Example buckets:**
```
example-prod-backups.s3.amazonaws.com (PUBLIC ACCESS!)
example-dev.blob.core.windows.net
```

#### 9. Sensitive Files (`sensitive/`)

```bash
# All sensitive files
cat sensitive/all_sensitive.txt

# By category
cat sensitive/config_files.txt      # .xml, .yaml, .env
cat sensitive/backups.txt           # .bak, .backup, .old
cat sensitive/databases.txt         # .sql, .db
cat sensitive/keys_certs.txt        # .key, .pem, .crt
cat sensitive/source_control.txt    # .git, .htaccess
cat sensitive/archives.txt          # .zip, .tar
cat sensitive/logs.txt              # .log files

# Summary JSON
cat sensitive/summary.json
```

#### 10. GF Patterns (`gf_patterns/`)

```bash
# URLs matching vulnerability patterns
cat gf_patterns/xss_candidates.txt
cat gf_patterns/sqli_candidates.txt
cat gf_patterns/ssrf_candidates.txt
cat gf_patterns/ssti_candidates.txt
cat gf_patterns/lfi_candidates.txt
cat gf_patterns/redirect_candidates.txt
cat gf_patterns/rce_candidates.txt
cat gf_patterns/idor_candidates.txt
```

**These are CANDIDATES - test manually!**

#### 11. Scan Results (`scan_results.json`)

```bash
# Complete summary
cat scan_results.json
```

**Example:**
```json
{
  "domain": "example.com",
  "scan_date": "2025-10-25 19:00:00",
  "duration": "38m42s",
  "total_subdomains": 2345,
  "resolved_subdomains": 1023,
  "live_hosts": 456,
  "total_urls": 5234,
  "js_files": 189,
  "source_maps": 3,
  "secrets": 12,
  "vulnerabilities": 23,
  "critical_vulns": 2,
  "sensitive_files": 115,
  "cloud_assets": 3
}
```

---

## Advanced Workflows

### 1. Bug Bounty Hunting Workflow

#### Step 1: Initial Reconnaissance
```bash
reconv4
# Enter: target.com
```

#### Step 2: Prioritize Critical Findings
```bash
DOMAIN="target_com"
MONTH="2025-10"
BASE=".results/${DOMAIN}/${MONTH}"

# Check critical vulnerabilities
echo "=== CRITICAL VULNERABILITIES ==="
jq '.[] | select(.severity=="critical")' ${BASE}/nuclei/findings.json

# Check secrets
echo "\n=== SECRETS FOUND ==="
cat ${BASE}/secrets/api_keys.json
cat ${BASE}/secrets/aws_keys.json

# Check sensitive files
echo "\n=== SENSITIVE FILES ==="
cat ${BASE}/sensitive/keys_certs.txt
cat ${BASE}/sensitive/backups.txt
```

#### Step 3: Manual Validation
```bash
# Test XSS candidates
while read url; do
    echo "Testing: $url"
    # Use your favorite XSS payloads
done < ${BASE}/gf_patterns/xss_candidates.txt

# Test SQLi candidates
while read url; do
    echo "Testing: $url"
    sqlmap -u "$url" --batch --risk 3
done < ${BASE}/gf_patterns/sqli_candidates.txt
```

#### Step 4: Report Findings
Use the categorized nuclei findings for structured reports.

### 2. Monthly Monitoring Workflow

```bash
# Run monthly
reconv4
# Enter same domain as last month

# Review comparison
DOMAIN="target_com"
cat .results/${DOMAIN}/comparison/2025-10_vs_2025-09.json

# Extract NEW findings only
jq '.differences' .results/${DOMAIN}/comparison/2025-10_vs_2025-09.json
```

**Focus on:**
- New subdomains (expanded attack surface)
- New URLs (new functionality)
- New vulnerabilities (fresh bugs!)

### 3. Multi-Domain Scanning

Create a wrapper script:

```bash
#!/bin/bash
# multi-scan.sh

DOMAINS=(
    "example.com"
    "test.com"
    "demo.com"
)

for domain in "${DOMAINS[@]}"; do
    echo "Scanning $domain..."

    # Auto-input for reconv4
    echo -e "${domain}\n" | reconv4

    echo "Completed: $domain"
    sleep 60  # Rate limit between domains
done
```

### 4. Continuous Monitoring (Cron)

```bash
# Add to crontab
crontab -e

# Run monthly on 1st at 2 AM
0 2 1 * * /path/to/reconv4-wrapper.sh >> /var/log/reconv4.log 2>&1
```

**Wrapper script:**
```bash
#!/bin/bash
# reconv4-wrapper.sh

export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin

DOMAIN="target.com"
WEBHOOK="https://discord.com/api/webhooks/YOUR_WEBHOOK"

echo -e "${DOMAIN}\n${WEBHOOK}" | /usr/local/bin/reconv4
```

---

## Integration Examples

### 1. Burp Suite Integration

```bash
# Export live hosts to Burp
DOMAIN="target_com"
MONTH="2025-10"

cat .results/${DOMAIN}/${MONTH}/httpx/live_hosts.txt > burp-targets.txt

# Import burp-targets.txt into Burp Suite Target Scope
```

### 2. Feed URLs to Custom Tools

```bash
# Send all URLs to your custom fuzzer
cat .results/target_com/2025-10/urls/all_urls.txt | your-fuzzer

# Send JS endpoints to custom analyzer
cat .results/target_com/2025-10/js_files/js_endpoints.txt | endpoint-analyzer
```

### 3. Slack Notification (Alternative to Discord)

Modify `main.go` to add Slack webhook support:

```go
func sendSlackNotification(config Config, results *ScanResults) {
    // Similar to sendDiscordNotification but for Slack format
    // Payload format: {"text": "message"}
}
```

### 4. Elasticsearch/Kibana Integration

Export to JSON for Elasticsearch:

```bash
# Index scan results
DOMAIN="target_com"
MONTH="2025-10"

curl -X POST "localhost:9200/reconv4-scans/_doc" \
  -H 'Content-Type: application/json' \
  -d @.results/${DOMAIN}/${MONTH}/scan_results.json
```

---

## Tips & Tricks

### 1. Quickly Find High-Value Targets

```bash
DOMAIN="target_com"
MONTH="2025-10"
BASE=".results/${DOMAIN}/${MONTH}"

# Subdomains with "admin", "dev", "staging"
grep -E "(admin|dev|staging|test)" ${BASE}/subdomains/all_subdomains.txt

# URLs with sensitive keywords
grep -iE "(admin|api|dashboard|internal)" ${BASE}/urls/all_urls.txt

# Config files
cat ${BASE}/sensitive/config_files.txt | grep -E "(.env|config|database)"
```

### 2. Extract Specific Technologies

```bash
# Find all WordPress sites
jq '.tech[] | select(. == "WordPress")' ${BASE}/httpx/httpx_results.json

# Find all React apps
grep -i "react" ${BASE}/httpx/tech_stack.txt
```

### 3. Identify Subdomain Takeover Candidates

```bash
# Check nuclei takeover findings
cat ${BASE}/nuclei/takeovers.json

# Manual verification
while read subdomain; do
    dig $subdomain CNAME
done < ${BASE}/subdomains/all_subdomains.txt | grep -i "cloudfront\|herokuapp\|s3"
```

### 4. Find API Endpoints

```bash
# From URLs
grep -E "/api/|/v[0-9]+/" ${BASE}/urls/all_urls.txt

# From JS files
cat ${BASE}/js_files/js_endpoints.txt | grep -i "api"

# From GF patterns
cat ${BASE}/gf_patterns/api_candidates.txt
```

### 5. Quick Security Score

```bash
# Generate a security score
CRITICAL=$(jq 'length' ${BASE}/nuclei/findings.json 2>/dev/null || echo 0)
SECRETS=$(jq 'length' ${BASE}/secrets/gitleaks_findings.json 2>/dev/null || echo 0)
SENSITIVE=$(wc -l < ${BASE}/sensitive/all_sensitive.txt 2>/dev/null || echo 0)

SCORE=$((100 - CRITICAL*10 - SECRETS*5 - SENSITIVE/10))

echo "Security Score: $SCORE/100"
echo "  - Critical Vulns: $CRITICAL (-10 each)"
echo "  - Secrets Found: $SECRETS (-5 each)"
echo "  - Sensitive Files: $SENSITIVE (-0.1 each)"
```

### 6. Compare Two Scans Manually

```bash
# Compare subdomains between months
diff .results/target_com/2025-09/subdomains/all_subdomains.txt \
     .results/target_com/2025-10/subdomains/all_subdomains.txt

# Find NEW subdomains
comm -13 <(sort .results/target_com/2025-09/subdomains/all_subdomains.txt) \
         <(sort .results/target_com/2025-10/subdomains/all_subdomains.txt)

# Find REMOVED subdomains
comm -23 <(sort .results/target_com/2025-09/subdomains/all_subdomains.txt) \
         <(sort .results/target_com/2025-10/subdomains/all_subdomains.txt)
```

### 7. Export to CSV for Excel/Sheets

```bash
# Export vulnerabilities to CSV
echo "Host,Severity,Template,Title" > vulns.csv
jq -r '.[] | [.host, .severity, .template_id, .info.name] | @csv' \
  ${BASE}/nuclei/findings.json >> vulns.csv

# Export subdomains with status
echo "Subdomain,Status,Title" > subdomains.csv
jq -r '[.url, .status_code, .title] | @csv' \
  ${BASE}/httpx/httpx_results.json >> subdomains.csv
```

### 8. Automate Report Generation

```bash
#!/bin/bash
# generate-report.sh

DOMAIN=$1
MONTH=$2
BASE=".results/${DOMAIN}/${MONTH}"

cat > report.md << EOF
# Security Scan Report: $DOMAIN
Date: $(date)

## Summary
- Total Subdomains: $(wc -l < ${BASE}/subdomains/all_subdomains.txt)
- Live Hosts: $(wc -l < ${BASE}/httpx/live_hosts.txt)
- URLs Discovered: $(wc -l < ${BASE}/urls/all_urls.txt)
- Vulnerabilities: $(jq 'length' ${BASE}/nuclei/findings.json)
- Secrets Found: $(jq 'length' ${BASE}/secrets/gitleaks_findings.json)

## Critical Findings
$(jq -r '.[] | select(.severity=="critical") | "- \(.template_id): \(.info.name)"' ${BASE}/nuclei/findings.json)

## Recommendations
1. Patch critical vulnerabilities immediately
2. Rotate exposed secrets
3. Secure sensitive file access
EOF

echo "Report generated: report.md"
```

### 9. Monitor Specific Patterns

```bash
# Watch for specific vulnerabilities
watch -n 300 '
  tail -1 .results/*/*/nuclei/findings.json | \
  jq ".[] | select(.template_id | contains("sql-injection"))"
'
```

### 10. Backup Results

```bash
# Compress and backup monthly
tar -czf reconv4-backup-$(date +%Y-%m-%d).tar.gz .results/
aws s3 cp reconv4-backup-$(date +%Y-%m-%d).tar.gz s3://my-backup-bucket/
```

---

## Performance Tuning

### For Faster Scans

Edit `main.go` and increase concurrency:

```go
// Phase 1: Subdomain Enumeration
concurrency := 100  // Increase from 50

// Phase 3: HTTP Probing
"-threads", "100",  // Increase from 50

// Phase 5: Nuclei
"-c", "50",  // Increase from 25
```

### For More Thorough Scans

```go
// Phase 4: Katana crawling depth
"-d", "5",  // Increase from 3 for deeper crawling

// Phase 5: Include all severities
"-severity", "critical,high,medium,low,info",
```

---

## Troubleshooting Common Issues

### Issue: No subdomains found

**Solution:**
```bash
# Test tools individually
subfinder -d example.com
assetfinder --subs-only example.com
findomain -t example.com
```

### Issue: Nuclei finds nothing

**Solution:**
```bash
# Update templates
nuclei -update-templates

# Test manually
echo "https://example.com" | nuclei -silent
```

### Issue: Out of memory

**Solution:**
```bash
# Reduce concurrency in main.go
concurrency := 20  // Lower value

# Or increase system swap
sudo fallocate -l 8G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

---

## Best Practices

1. **Always have permission** before scanning
2. **Start with small scopes** to test setup
3. **Review findings manually** - tools have false positives
4. **Keep tools updated** regularly
5. **Back up results** before cleanup
6. **Use Discord notifications** for long scans
7. **Run monthly** to catch new vulnerabilities
8. **Document your findings** properly
9. **Respect rate limits** to avoid blocking
10. **Combine with manual testing** for best results

---

Happy hunting! 🎯
