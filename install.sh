#!/bin/bash

# Reconv4 Installation Script
# Automated installation of all required tools

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}"
cat << "EOF"
██████╗ ███████╗ ██████╗ ██████╗ ███╗   ██╗██╗   ██╗██╗  ██╗
██╔══██╗██╔════╝██╔════╝██╔═══██╗████╗  ██║██║   ██║██║  ██║
██████╔╝█████╗  ██║     ██║   ██║██╔██╗ ██║██║   ██║███████║
██╔══██╗██╔══╝  ██║     ██║   ██║██║╚██╗██║╚██╗ ██╔╝╚════██║
██║  ██║███████╗╚██████╗╚██████╔╝██║ ╚████║ ╚████╔╝      ██║
╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝  ╚═══╝       ╚═╝
                    Installation Script v4.0
EOF
echo -e "${NC}"

# Check if running as root
if [ "$EUID" -eq 0 ]; then 
   echo -e "${RED}[!] Please do not run this script as root${NC}"
   echo -e "${YELLOW}[*] Run as regular user. It will ask for sudo when needed.${NC}"
   exit 1
fi

# Detect OS
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    OS="linux"
    echo -e "${GREEN}[✓] Detected OS: Linux${NC}"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    OS="macos"
    echo -e "${GREEN}[✓] Detected OS: macOS${NC}"
else
    echo -e "${RED}[!] Unsupported OS: $OSTYPE${NC}"
    exit 1
fi

# Check for Go installation
echo -e "\n${BLUE}[*] Checking Go installation...${NC}"
if ! command -v go &> /dev/null; then
    echo -e "${YELLOW}[!] Go not found. Installing Go 1.21...${NC}"

    if [ "$OS" == "linux" ]; then
        cd /tmp
        wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
        sudo rm -rf /usr/local/go
        sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz

        # Add to PATH
        if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
            echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
            echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
        fi

        export PATH=$PATH:/usr/local/go/bin
        export PATH=$PATH:$HOME/go/bin

        rm go1.21.6.linux-amd64.tar.gz
    elif [ "$OS" == "macos" ]; then
        if ! command -v brew &> /dev/null; then
            echo -e "${YELLOW}[!] Homebrew not found. Installing Homebrew...${NC}"
            /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        fi
        brew install go
    fi

    echo -e "${GREEN}[✓] Go installed successfully${NC}"
else
    GO_VERSION=$(go version | awk '{print $3}')
    echo -e "${GREEN}[✓] Go already installed: $GO_VERSION${NC}"
fi

# Install system dependencies
echo -e "\n${BLUE}[*] Installing system dependencies...${NC}"
if [ "$OS" == "linux" ]; then
    sudo apt-get update -qq
    sudo apt-get install -y -qq git wget curl python3 python3-pip build-essential
elif [ "$OS" == "macos" ]; then
    brew install git wget curl python3
fi
echo -e "${GREEN}[✓] System dependencies installed${NC}"

# Create tools directory
mkdir -p ~/tools
cd ~/tools

# Function to install Go tools
install_go_tool() {
    local tool_name=$1
    local tool_repo=$2

    echo -e "${YELLOW}[*] Installing $tool_name...${NC}"
    if command -v $tool_name &> /dev/null; then
        echo -e "${GREEN}[✓] $tool_name already installed${NC}"
    else
        go install -v $tool_repo@latest
        echo -e "${GREEN}[✓] $tool_name installed${NC}"
    fi
}

# Install reconnaissance tools
echo -e "\n${BLUE}[*] Installing reconnaissance tools...${NC}"

install_go_tool "subfinder" "github.com/projectdiscovery/subfinder/v2/cmd/subfinder"
install_go_tool "assetfinder" "github.com/tomnomnom/assetfinder"
install_go_tool "findomain" "github.com/Findomain/Findomain"
install_go_tool "dnsx" "github.com/projectdiscovery/dnsx/cmd/dnsx"
install_go_tool "httpx" "github.com/projectdiscovery/httpx/cmd/httpx"
install_go_tool "gau" "github.com/lc/gau/v2/cmd/gau"
install_go_tool "waybackurls" "github.com/tomnomnom/waybackurls"
install_go_tool "katana" "github.com/projectdiscovery/katana/cmd/katana"
install_go_tool "nuclei" "github.com/projectdiscovery/nuclei/v3/cmd/nuclei"
install_go_tool "gf" "github.com/tomnomnom/gf"

# Install gitleaks
echo -e "${YELLOW}[*] Installing gitleaks...${NC}"
if command -v gitleaks &> /dev/null; then
    echo -e "${GREEN}[✓] gitleaks already installed${NC}"
else
    if [ "$OS" == "linux" ]; then
        cd /tmp
        wget https://github.com/gitleaks/gitleaks/releases/download/v8.18.1/gitleaks_8.18.1_linux_x64.tar.gz
        tar -xzf gitleaks_8.18.1_linux_x64.tar.gz
        sudo mv gitleaks /usr/local/bin/
        rm gitleaks_8.18.1_linux_x64.tar.gz
    elif [ "$OS" == "macos" ]; then
        brew install gitleaks
    fi
    echo -e "${GREEN}[✓] gitleaks installed${NC}"
fi

# Install cloud_enum
echo -e "${YELLOW}[*] Installing cloud_enum...${NC}"
if [ -d "~/tools/cloud_enum" ]; then
    echo -e "${GREEN}[✓] cloud_enum already installed${NC}"
else
    cd ~/tools
    git clone https://github.com/initstring/cloud_enum.git
    cd cloud_enum
    pip3 install -r requirements.txt --break-system-packages

    # Create wrapper script
    sudo tee /usr/local/bin/cloud_enum > /dev/null << 'WRAPPER'
#!/bin/bash
python3 ~/tools/cloud_enum/cloud_enum.py "$@"
WRAPPER
    sudo chmod +x /usr/local/bin/cloud_enum

    echo -e "${GREEN}[✓] cloud_enum installed${NC}"
fi

# Install GF patterns
echo -e "\n${BLUE}[*] Installing GF patterns...${NC}"
if [ -d "~/.gf" ]; then
    echo -e "${GREEN}[✓] GF patterns already installed${NC}"
else
    mkdir -p ~/.gf
    cd ~/.gf
    git clone https://github.com/1ndianl33t/Gf-Patterns.git
    mv Gf-Patterns/*.json .
    rm -rf Gf-Patterns
    echo -e "${GREEN}[✓] GF patterns installed${NC}"
fi

# Update nuclei templates
echo -e "\n${BLUE}[*] Updating nuclei templates...${NC}"
nuclei -update-templates
echo -e "${GREEN}[✓] Nuclei templates updated${NC}"

# Build Reconv4
echo ""
echo -e "${BLUE}[*] Building Reconv4...${NC}"

# Save the original directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

if [ ! -f "main.go" ]; then
    echo -e "${RED}[!] main.go not found in $SCRIPT_DIR${NC}"
    echo -e "${YELLOW}[*] Current directory: $(pwd)${NC}"
    echo -e "${YELLOW}[*] Looking for main.go in: $SCRIPT_DIR${NC}"
    exit 1
fi

echo -e "${BLUE}[*] Downloading Go dependencies...${NC}"
go mod tidy
go mod download

echo -e "${BLUE}[*] Building binary...${NC}"
go build -o reconv4 main.go

if [ -f "reconv4" ]; then
    sudo mv reconv4 /usr/local/bin/
    sudo chmod +x /usr/local/bin/reconv4
    echo -e "${GREEN}[✓] Reconv4 built and installed to /usr/local/bin/reconv4${NC}"
else
    echo -e "${RED}[!] Failed to build Reconv4${NC}"
    exit 1
fi

echo -e "${BLUE}[*] Downloading Go dependencies...${NC}"
go mod tidy
go mod download

echo -e "${BLUE}[*] Building binary...${NC}"
go build -o reconv4 main.go

if [ -f "reconv4" ]; then
    sudo mv reconv4 /usr/local/bin/
    sudo chmod +x /usr/local/bin/reconv4
    echo -e "${GREEN}[✓] Reconv4 built and installed to /usr/local/bin/reconv4${NC}"
else
    echo -e "${RED}[!] Failed to build Reconv4${NC}"
    exit 1
fi

echo -e "${BLUE}[*] Downloading Go dependencies...${NC}"
go mod tidy
go mod download

echo -e "${BLUE}[*] Building binary...${NC}"
go build -o reconv4 main.go

if [ -f "reconv4" ]; then
    sudo mv reconv4 /usr/local/bin/
    sudo chmod +x /usr/local/bin/reconv4
    echo -e "${GREEN}[✓] Reconv4 built and installed to /usr/local/bin/reconv4${NC}"
else
    echo -e "${RED}[!] Failed to build Reconv4${NC}"
    exit 1
fi

# Final verification
echo -e "\n${BLUE}[*] Verifying installation...${NC}"

TOOLS=(
    "subfinder"
    "assetfinder"
    "findomain"
    "dnsx"
    "httpx"
    "gau"
    "waybackurls"
    "katana"
    "nuclei"
    "gf"
    "gitleaks"
    "cloud_enum"
    "reconv4"
)

MISSING=()

for tool in "${TOOLS[@]}"; do
    if command -v $tool &> /dev/null; then
        echo -e "${GREEN}  ✓ $tool${NC}"
    else
        echo -e "${RED}  ✗ $tool${NC}"
        MISSING+=($tool)
    fi
done

if [ ${#MISSING[@]} -eq 0 ]; then
    echo -e "\n${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}[✓] Installation completed successfully!${NC}"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "\n${YELLOW}[*] Usage: reconv4${NC}"
    echo -e "${YELLOW}[*] The tool will guide you through the scan process${NC}"
    echo -e "\n${BLUE}[*] Please reload your shell or run:${NC}"
    echo -e "${YELLOW}    source ~/.bashrc${NC}"
else
    echo -e "\n${RED}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${RED}[!] Installation incomplete. Missing tools:${NC}"
    for tool in "${MISSING[@]}"; do
        echo -e "${RED}    - $tool${NC}"
    done
    echo -e "${RED}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${YELLOW}[*] Please install missing tools manually${NC}"
fi
