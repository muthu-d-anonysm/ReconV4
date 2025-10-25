# Changelog

All notable changes to Reconv4 will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [4.0.0] - 2025-10-25

### 🎉 Initial Release

Complete rewrite from Reconv3 with major architectural improvements.

### Added
- **9 Automated Phases** for comprehensive reconnaissance
- **Phase 1**: Subdomain Enumeration (subfinder, assetfinder, findomain)
- **Phase 2**: DNS Resolution with dnsx
- **Phase 3**: Live Host Detection with httpx and tech stack detection
- **Phase 4**: URL Discovery (gau, waybackurls, katana)
- **Phase 4.5**: JavaScript Analysis with source map extraction
- **Phase 5**: Vulnerability Scanning with nuclei and auto-categorization
- **Phase 6**: Cloud Asset Discovery (S3, Azure, GCP)
- **Phase 6.5**: Sensitive File Discovery with 8 categories
- **Phase 7**: GF Pattern Filtering for vulnerability candidates
- **Phase 8**: Monthly Comparison for delta detection
- Auto-detection of optimal system settings
- Discord webhook integration for real-time notifications
- Comprehensive error handling and crash-proof design
- Intelligent retry logic for network operations
- Progress tracking with emoji indicators
- Colored terminal output for better UX
- Structured JSON output for all findings
- Monthly timeline tracking
- Tech stack detection and categorization

### Technical Improvements
- Concurrent execution with goroutines for maximum speed
- Smart concurrency tuning per operation type
- Context-based cancellation for graceful shutdown
- MD5-based filename generation for downloaded JS files
- Source map URL resolution (relative and absolute)
- Categorized secret detection (API keys, tokens, AWS keys, etc.)
- Custom pattern matching for S3 buckets, Firebase, internal IPs
- Vulnerability auto-categorization (CVEs, misconfigs, takeovers, etc.)
- Deduplication of URLs, subdomains, and findings
- Memory-efficient file processing
- Rate limiting to prevent blocking

### Performance
- Process 1000 domains in ~3.5-4 hours
- Optimized for 12-core systems with 16GB RAM
- 40-50 concurrent operations (auto-tuned)
- Batch processing instead of per-subdomain operations
- Parallel tool execution where possible

### Security Features
- Gitleaks integration for secret scanning
- Custom regex patterns for sensitive data
- Source map analysis for hidden endpoints
- API endpoint extraction from JavaScript
- Certificate and key file detection
- Backup and config file identification
- Database file discovery
- Cloud bucket enumeration

### Documentation
- Comprehensive README.md with examples
- Detailed USAGE.md with workflows
- Installation script for all dependencies
- Troubleshooting guide
- Bug bounty integration examples
- Performance tuning guidelines

### File Structure
```
reconv4/
├── main.go              # Core application
├── go.mod               # Go dependencies
├── install.sh           # Installation script
├── README.md            # Main documentation
├── USAGE.md             # Usage examples
├── CHANGELOG.md         # Version history
├── LICENSE              # MIT License
└── .gitignore          # Git ignore rules
```

### Output Structure
```
.results/
└── {domain}/
    ├── {YYYY-MM}/
    │   ├── subdomains/
    │   ├── dns/
    │   ├── httpx/
    │   ├── urls/
    │   ├── js_files/
    │   ├── secrets/
    │   ├── nuclei/
    │   ├── cloud/
    │   ├── sensitive/
    │   ├── gf_patterns/
    │   └── scan_results.json
    ├── comparison/
    └── timeline.json
```

### Dependencies
- Go 1.21+
- subfinder v2
- assetfinder
- findomain
- dnsx
- httpx
- gau v2
- waybackurls
- katana
- nuclei v3
- gf
- gitleaks v8
- cloud_enum
- GF-Patterns

### System Requirements
- OS: Linux (Ubuntu/Debian/Kali) or macOS
- CPU: 4+ cores (12+ recommended)
- RAM: 8GB minimum (16GB recommended)
- Disk: 10GB free space

---

## Roadmap

### [4.1.0] - Planned
- [ ] Web UI dashboard for results visualization
- [ ] API endpoint for programmatic access
- [ ] PDF/HTML report generation
- [ ] Multi-target parallel scanning
- [ ] Slack integration
- [ ] Telegram bot integration

### [4.2.0] - Planned
- [ ] Docker containerization
- [ ] Kubernetes deployment manifests
- [ ] Database backend (PostgreSQL/MongoDB)
- [ ] Historical trend analysis
- [ ] Machine learning for false positive reduction
- [ ] Custom template support for nuclei

### [5.0.0] - Future
- [ ] Distributed scanning architecture
- [ ] Real-time collaborative workspace
- [ ] Integration with SIEM systems
- [ ] Advanced correlation engine
- [ ] Automated exploitation (ethical, with consent)
- [ ] CI/CD pipeline integration

---

## Version History

### v4.0.0 (2025-10-25)
- Initial production release
- Complete rewrite from Reconv3
- 9-phase automated workflow
- Discord integration
- Monthly comparison tracking

### v3.x (Previous)
- Legacy version (deprecated)
- Basic reconnaissance features
- Limited tool integration
- No monthly tracking

---

## Breaking Changes from v3.x

### Architecture
- **Complete rewrite**: v4.0.0 is not compatible with v3.x
- **New output structure**: Results are organized by month
- **Different tool set**: Added katana, dnsx, cloud_enum
- **New phases**: Added JavaScript analysis, sensitive file discovery

### Migration Guide
1. Export v3.x results if needed
2. Uninstall v3.x: `sudo rm /usr/local/bin/reconv3`
3. Install v4.0.0: `./install.sh`
4. Run first scan: `reconv4`

### Configuration
- v3.x used config files, v4.0.0 uses interactive prompts
- Discord webhooks replace email notifications
- No manual configuration needed (auto-detection)

---

## Known Issues

### v4.0.0
- cloud_enum may not be available on all systems (optional dependency)
- Some tools may require API keys for full functionality (optional)
- Large scans (5000+ domains) may require manual tuning

### Workarounds
- If cloud_enum fails, install manually from GitHub
- Add API keys to tool configs for better results
- For large scans, reduce concurrency in main.go

---

## Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### How to Report Bugs
1. Check [existing issues](https://github.com/yourusername/reconv4/issues)
2. Create new issue with:
   - OS and version
   - Go version
   - Full error message
   - Steps to reproduce

### Feature Requests
1. Open an issue with `[Feature Request]` prefix
2. Describe the use case
3. Explain expected behavior

---

## Credits

### Reconv4 Author
- Security Researcher & Bug Bounty Hunter

### Tool Authors
Special thanks to the creators of:
- ProjectDiscovery team (subfinder, dnsx, httpx, katana, nuclei)
- Tom Hudson (assetfinder, waybackurls, gf)
- Findomain (Eduard Tolosa)
- Gitleaks (Zachary Rice)
- cloud_enum (initstring)
- All GF pattern contributors

### Inspiration
- OWASP Testing Guide
- Bug Bounty Methodology v4
- HackerOne community best practices

---

## License

MIT License - see [LICENSE](LICENSE) file for details.

---

**Last Updated**: October 25, 2025
**Current Version**: 4.0.0
