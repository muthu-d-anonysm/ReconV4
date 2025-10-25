#!/bin/bash

# Reconv4 Installation Verification Script
# Checks all dependencies and system requirements

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}"
cat << "EOF"
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║          Reconv4 Installation Verification                ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

PASS=0
FAIL=0
WARN=0

# Function to check command
check_tool() {
    local tool=$1
    local required=$2

    if command -v $tool &> /dev/null; then
        version=$(${tool} --version 2>&1 | head -1 || echo "unknown")
        echo -e "  ${GREEN}✓${NC} $tool ${GREEN}INSTALLED${NC} ($version)"
        ((PASS++))
        return 0
    else
        if [ "$required" == "required" ]; then
            echo -e "  ${RED}✗${NC} $tool ${RED}MISSING (REQUIRED)${NC}"
            ((FAIL++))
        else
            echo -e "  ${YELLOW}!${NC} $tool ${YELLOW}MISSING (OPTIONAL)${NC}"
            ((WARN++))
        fi
        return 1
    fi
}

echo ""
echo -e "${BLUE}[1/6] Checking System Requirements...${NC}"
echo ""

# Check OS
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    echo -e "  ${GREEN}✓${NC} Operating System: Linux"
    ((PASS++))
elif [[ "$OSTYPE" == "darwin"* ]]; then
    echo -e "  ${GREEN}✓${NC} Operating System: macOS"
    ((PASS++))
else
    echo -e "  ${YELLOW}!${NC} Operating System: $OSTYPE (Untested)"
    ((WARN++))
fi

# Check CPU cores
CORES=$(nproc 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo 1)
if [ $CORES -ge 4 ]; then
    echo -e "  ${GREEN}✓${NC} CPU Cores: $CORES (Recommended: 4+)"
    ((PASS++))
else
    echo -e "  ${YELLOW}!${NC} CPU Cores: $CORES (Recommended: 4+)"
    ((WARN++))
fi

# Check RAM
if command -v free &> /dev/null; then
    RAM_GB=$(free -g | awk '/^Mem:/{print $2}')
    if [ $RAM_GB -ge 8 ]; then
        echo -e "  ${GREEN}✓${NC} RAM: ${RAM_GB}GB (Recommended: 8GB+)"
        ((PASS++))
    else
        echo -e "  ${YELLOW}!${NC} RAM: ${RAM_GB}GB (Recommended: 8GB+)"
        ((WARN++))
    fi
fi

# Check disk space
DISK_FREE=$(df -h . | awk 'NR==2 {print $4}' | sed 's/G//')
if [ ${DISK_FREE%.*} -ge 10 ]; then
    echo -e "  ${GREEN}✓${NC} Disk Space: ${DISK_FREE}G free (Recommended: 10GB+)"
    ((PASS++))
else
    echo -e "  ${YELLOW}!${NC} Disk Space: ${DISK_FREE}G free (Recommended: 10GB+)"
    ((WARN++))
fi

echo ""
echo -e "${BLUE}[2/6] Checking Go Installation...${NC}"
echo ""

check_tool "go" "required"

if command -v go &> /dev/null; then
    # Check GOPATH
    if [ -d "$HOME/go/bin" ]; then
        echo -e "  ${GREEN}✓${NC} GOPATH: $HOME/go"
        ((PASS++))
    else
        echo -e "  ${YELLOW}!${NC} GOPATH: Not found (will be created)"
        ((WARN++))
    fi

    # Check PATH
    if echo $PATH | grep -q "$HOME/go/bin"; then
        echo -e "  ${GREEN}✓${NC} Go bin in PATH"
        ((PASS++))
    else
        echo -e "  ${RED}✗${NC} Go bin NOT in PATH"
        echo -e "      ${YELLOW}Add to ~/.bashrc: export PATH=\$PATH:\$HOME/go/bin${NC}"
        ((FAIL++))
    fi
fi

echo ""
echo -e "${BLUE}[3/6] Checking Subdomain Enumeration Tools...${NC}"
echo ""

check_tool "subfinder" "required"
check_tool "assetfinder" "required"
check_tool "findomain" "required"

echo ""
echo -e "${BLUE}[4/6] Checking DNS/HTTP Tools...${NC}"
echo ""

check_tool "dnsx" "required"
check_tool "httpx" "required"

echo ""
echo -e "${BLUE}[5/6] Checking URL Discovery & Analysis Tools...${NC}"
echo ""

check_tool "gau" "required"
check_tool "waybackurls" "required"
check_tool "katana" "required"
check_tool "gf" "required"

# Check GF patterns
if [ -d "$HOME/.gf" ] && [ "$(ls -A $HOME/.gf/*.json 2>/dev/null | wc -l)" -gt 0 ]; then
    pattern_count=$(ls -1 $HOME/.gf/*.json 2>/dev/null | wc -l)
    echo -e "  ${GREEN}✓${NC} GF Patterns: $pattern_count patterns installed"
    ((PASS++))
else
    echo -e "  ${RED}✗${NC} GF Patterns: Not installed"
    ((FAIL++))
fi

echo ""
echo -e "${BLUE}[6/6] Checking Security Analysis Tools...${NC}"
echo ""

check_tool "nuclei" "required"
check_tool "gitleaks" "required"
check_tool "cloud_enum" "optional"

# Check nuclei templates
if [ -d "$HOME/.nuclei-templates" ]; then
    template_count=$(find $HOME/.nuclei-templates -name "*.yaml" 2>/dev/null | wc -l)
    if [ $template_count -gt 100 ]; then
        echo -e "  ${GREEN}✓${NC} Nuclei Templates: $template_count templates"
        ((PASS++))
    else
        echo -e "  ${YELLOW}!${NC} Nuclei Templates: $template_count (run: nuclei -update-templates)"
        ((WARN++))
    fi
else
    echo -e "  ${RED}✗${NC} Nuclei Templates: Not found"
    ((FAIL++))
fi

echo ""
echo -e "${BLUE}[7/7] Checking Reconv4 Installation...${NC}"
echo ""

if [ -f "main.go" ]; then
    echo -e "  ${GREEN}✓${NC} main.go found in current directory"
    ((PASS++))
else
    echo -e "  ${RED}✗${NC} main.go not found"
    ((FAIL++))
fi

if [ -f "go.mod" ]; then
    echo -e "  ${GREEN}✓${NC} go.mod found"
    ((PASS++))
else
    echo -e "  ${RED}✗${NC} go.mod not found"
    ((FAIL++))
fi

check_tool "reconv4" "optional"

# Summary
echo ""
echo "═══════════════════════════════════════════════════════════"
echo -e "${BLUE}VERIFICATION SUMMARY${NC}"
echo "═══════════════════════════════════════════════════════════"
echo ""
echo -e "  ${GREEN}✓ Passed:${NC}  $PASS"
echo -e "  ${YELLOW}! Warnings:${NC} $WARN"
echo -e "  ${RED}✗ Failed:${NC}  $FAIL"
echo ""

if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}✅ VERIFICATION PASSED!${NC}"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo "All required tools are installed and configured correctly."
    echo ""

    if [ ! -f "/usr/local/bin/reconv4" ]; then
        echo -e "${YELLOW}[*] To complete installation, run:${NC}"
        echo "    go build -o reconv4 main.go"
        echo "    sudo mv reconv4 /usr/local/bin/"
        echo ""
    fi

    echo -e "${GREEN}Ready to start scanning!${NC}"
    echo ""
    echo "Usage:"
    echo "  ./quickstart.sh  (interactive guide)"
    echo "  OR"
    echo "  reconv4          (direct usage)"
    echo ""

    if [ $WARN -gt 0 ]; then
        echo -e "${YELLOW}Note: There are $WARN warnings, but they won't prevent scanning.${NC}"
        echo ""
    fi

    exit 0
else
    echo -e "${RED}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${RED}❌ VERIFICATION FAILED!${NC}"
    echo -e "${RED}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "${RED}$FAIL required component(s) are missing.${NC}"
    echo ""
    echo "To fix:"
    echo "  1. Run the installation script:"
    echo "     chmod +x install.sh"
    echo "     ./install.sh"
    echo ""
    echo "  2. Reload your shell:"
    echo "     source ~/.bashrc"
    echo ""
    echo "  3. Run verification again:"
    echo "     ./verify.sh"
    echo ""
    exit 1
fi
