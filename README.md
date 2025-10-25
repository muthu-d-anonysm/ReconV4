
# ReconV4 - Comprehensive Reconnaissance Automation Tool

ReconV4 is a professional reconnaissance automation framework designed for bug bounty hunters, security researchers, and penetration testers to automate subdomain enumeration, URL discovery, JavaScript secrets detection, cloud asset enumeration, vulnerability scanning, and more.

---

## Features

- Fast and efficient **subdomain enumeration** using Subfinder, Assetfinder, Findomain  
- Accurate **live host detection** with httpx integration  
- Parallelized **URL discovery** via Waybackurls, GAU, and Katana  
- Deep **JavaScript secrets and Source Map analysis** with full URL output  
- Integrated **cloud asset discovery** (S3, Azure, GCP) using cloud_enum with domain keyword bruteforcing  
- Powerful **vulnerability scanning** with Nuclei v3, multiple severity filters, and auto-categorization  
- Comprehensive **secret analysis** using Gitleaks, with full URLs and categorized outputs  
- Configurable concurrency and runtime options for balanced performance  
- Optional Discord notifications for scan completion alerts  
- Structured output directory by domain and scan month  

---

## Requirements

- Go v1.20 or higher  
- External tools installed and present in `$PATH`:  
  `subfinder`, `assetfinder`, `findomain`, `httpx`, `waybackurls`, `gau`, `katana`, `nuclei`, `gitleaks`, `cloud_enum`  
- Internet connectivity for API-based enumeration and scanning  

---

## Installation

Clone the repository and build the tool:

```

git clone https://github.com/muthu-d-anonysm/ReconV4.git
cd ReconV4
go build -o reconv4 main.go

```

---

## Usage

```

./reconv4 -d example.com -c 100

```

- `-d`: Target domain (required)  
- `-c`: Concurrency level (default 100, max 800)  

---

## Output Structure

Results are stored under `.results/<domain>/<yyyy-mm>/` with subdirectories:

- `subdomains/` - All enumerated subdomains  
- `dns/` - DNS resolution details  
- `httpx/` - Live hosts and technology stack detection  
- `urls/` - Comprehensive URL list discovered via crawling  
- `js_files/` - Downloaded JS files and source maps for secret scanning  
- `secrets/` - JSON categorized secrets with associated URLs  
- `cloud/` - Cloud assets found by cloud_enum (buckets, containers)  
- `nuclei/` - Vulnerability scan reports separated by categories  

---

## Tool Details and Notes

- Cloud asset discovery uses `cloud_enum` with domain keyword bruteforcing (`-k keywords.txt -l buckets.txt -t 10`)  
- URL crawlers run in parallel by tool but sequential per URL internally  
- Nuclei uses concurrency flag (`-c`) defaulting to 25 but configurable via global concurrency divided proportionally  
- GAU concurrency typically set to 100 threads for fast URL retrieval  
- Full URLs with HTTP/HTTPS schemes are used for probing and scanning tools for accuracy  
- Clean domain names without scheme prefixes are used for cloud enumeration keywords  
- Secret scanning integrates gitleaks output processed with original URLs preserved  

---

## Contributing

Pull requests, issues, and feature requests are welcome! Please follow coding standards and provide testing where applicable.

---

## License

MIT License


*Developed by Muthu D for professional bug bounty hunters and security researchers.*

```
