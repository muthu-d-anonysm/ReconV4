#!/bin/bash

# Reconv4 Quick Start Script
# For first-time users

set -e

clear

cat << "EOF"
╔════════════════════════════════════════════════════════════╗
║                                                            ║
║           Welcome to Reconv4 Quick Start!                  ║
║                                                            ║
║  This script will help you get started with Reconv4       ║
║                                                            ║
╚════════════════════════════════════════════════════════════╝
EOF

echo ""
echo "This quick start will:"
echo "  1. Check if Reconv4 is installed"
echo "  2. Help you run your first scan"
echo "  3. Show you how to analyze results"
echo ""

read -p "Press Enter to continue..."

# Check if reconv4 is installed
echo ""
echo "[*] Checking Reconv4 installation..."

if ! command -v reconv4 &> /dev/null; then
    echo "❌ Reconv4 is not installed!"
    echo ""
    echo "Please run the installation script first:"
    echo "  chmod +x install.sh"
    echo "  ./install.sh"
    exit 1
fi

echo "✅ Reconv4 is installed!"
echo ""

# Check required tools
echo "[*] Checking required tools..."

TOOLS=("subfinder" "httpx" "nuclei" "gitleaks")
MISSING=()

for tool in "${TOOLS[@]}"; do
    if command -v $tool &> /dev/null; then
        echo "  ✅ $tool"
    else
        echo "  ❌ $tool"
        MISSING+=($tool)
    fi
done

if [ ${#MISSING[@]} -gt 0 ]; then
    echo ""
    echo "⚠️  Some tools are missing. Please run install.sh"
    exit 1
fi

echo ""
echo "✅ All tools are installed!"
echo ""

# Example domain selection
cat << EOF
═══════════════════════════════════════════════════════════

Let's run your first scan!

For this example, you can use:
  - A domain you have permission to scan
  - A bug bounty target
  - A test domain (example.com - for testing only)

⚠️  WARNING: Only scan domains you have permission to test!

═══════════════════════════════════════════════════════════
EOF

echo ""
read -p "Enter a domain to scan (or 'demo' for example.com): " DOMAIN

if [ "$DOMAIN" == "demo" ]; then
    DOMAIN="example.com"
    echo ""
    echo "⚠️  Using example.com for demonstration"
    echo "    This is only for testing the workflow"
fi

echo ""
read -p "Do you want Discord notifications? (y/n): " DISCORD_CHOICE

WEBHOOK=""
if [[ "$DISCORD_CHOICE" =~ ^[Yy]$ ]]; then
    echo ""
    echo "To get a Discord webhook:"
    echo "  1. Go to your Discord server"
    echo "  2. Server Settings → Integrations → Webhooks"
    echo "  3. Create webhook and copy URL"
    echo ""
    read -p "Enter Discord Webhook URL (or press Enter to skip): " WEBHOOK
fi

# Confirm scan
echo ""
echo "═══════════════════════════════════════════════════════════"
echo "Ready to start scan!"
echo ""
echo "  Target: $DOMAIN"
if [ -n "$WEBHOOK" ]; then
    echo "  Discord: Enabled"
else
    echo "  Discord: Disabled"
fi
echo ""
echo "  Estimated time: 45-60 minutes"
echo "═══════════════════════════════════════════════════════════"
echo ""
read -p "Start scan now? (y/n): " CONFIRM

if [[ ! "$CONFIRM" =~ ^[Yy]$ ]]; then
    echo "Scan cancelled."
    exit 0
fi

# Run scan
echo ""
echo "[*] Starting Reconv4..."
echo ""

if [ -n "$WEBHOOK" ]; then
    echo -e "${DOMAIN}\n${WEBHOOK}" | reconv4
else
    echo -e "${DOMAIN}\n" | reconv4
fi

# After scan completes
echo ""
echo "═══════════════════════════════════════════════════════════"
echo "✅ Scan completed!"
echo "═══════════════════════════════════════════════════════════"
echo ""

# Show results location
DOMAIN_CLEAN=$(echo $DOMAIN | tr '.' '_')
MONTH=$(date +%Y-%m)
RESULTS_DIR=".results/${DOMAIN_CLEAN}/${MONTH}"

if [ -d "$RESULTS_DIR" ]; then
    echo "📁 Results saved to: $RESULTS_DIR"
    echo ""

    # Quick summary
    echo "📊 Quick Summary:"
    echo ""

    if [ -f "${RESULTS_DIR}/scan_results.json" ]; then
        echo "  Results:"
        cat "${RESULTS_DIR}/scan_results.json" | grep -E "(total_subdomains|live_hosts|total_urls|vulnerabilities|secrets)" | head -5
    fi

    echo ""
    echo "═══════════════════════════════════════════════════════════"
    echo "Next Steps:"
    echo "═══════════════════════════════════════════════════════════"
    echo ""
    echo "1. Review critical findings:"
    echo "   cat ${RESULTS_DIR}/nuclei/findings.json"
    echo ""
    echo "2. Check for secrets:"
    echo "   cat ${RESULTS_DIR}/secrets/gitleaks_findings.json"
    echo ""
    echo "3. View sensitive files:"
    echo "   cat ${RESULTS_DIR}/sensitive/all_sensitive.txt"
    echo ""
    echo "4. Explore URLs:"
    echo "   cat ${RESULTS_DIR}/urls/all_urls.txt"
    echo ""
    echo "5. Read full guide:"
    echo "   cat USAGE.md"
    echo ""
    echo "═══════════════════════════════════════════════════════════"
    echo ""
    echo "Happy hunting! 🎯"
    echo ""
fi
