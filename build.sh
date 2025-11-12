#!/bin/bash

# IMAP Migrator - Build Script
# Automatically installs Go (if needed) and builds the executable

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Minimum Go version required
MIN_GO_VERSION="1.19"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  IMAP Migrator - Build Script${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Function to compare version numbers
version_ge() {
    # Returns 0 if $1 >= $2
    printf '%s\n%s\n' "$2" "$1" | sort -V -C
}

# Function to check if Go is installed
check_go() {
    if command -v go &> /dev/null; then
        GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        echo -e "${GREEN}✓${NC} Go is installed: version ${GO_VERSION}"
        
        if version_ge "$GO_VERSION" "$MIN_GO_VERSION"; then
            echo -e "${GREEN}✓${NC} Go version is sufficient (>= ${MIN_GO_VERSION})"
            return 0
        else
            echo -e "${YELLOW}⚠${NC} Go version ${GO_VERSION} is too old (need >= ${MIN_GO_VERSION})"
            return 1
        fi
    else
        echo -e "${YELLOW}⚠${NC} Go is not installed"
        return 1
    fi
}

# Function to detect Linux distribution
detect_distro() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        echo "$ID"
    elif [ -f /etc/lsb-release ]; then
        . /etc/lsb-release
        echo "$DISTRIB_ID" | tr '[:upper:]' '[:lower:]'
    else
        echo "unknown"
    fi
}

# Function to install Go on Linux
install_go_linux() {
    echo ""
    echo -e "${BLUE}Installing Go...${NC}"
    
    DISTRO=$(detect_distro)
    echo "Detected distribution: $DISTRO"
    
    case "$DISTRO" in
        ubuntu|debian|linuxmint|pop)
            echo "Installing Go via apt..."
            sudo apt update
            sudo apt install -y golang-go
            ;;
        fedora|rhel|centos)
            echo "Installing Go via dnf/yum..."
            if command -v dnf &> /dev/null; then
                sudo dnf install -y golang
            else
                sudo yum install -y golang
            fi
            ;;
        arch|manjaro)
            echo "Installing Go via pacman..."
            sudo pacman -Sy --noconfirm go
            ;;
        opensuse*)
            echo "Installing Go via zypper..."
            sudo zypper install -y go
            ;;
        *)
            echo -e "${YELLOW}Distribution not recognized. Installing Go manually...${NC}"
            install_go_manual
            return
            ;;
    esac
    
    # Verify installation
    if command -v go &> /dev/null; then
        GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        echo -e "${GREEN}✓${NC} Go ${GO_VERSION} installed successfully!"
    else
        echo -e "${RED}✗${NC} Go installation failed via package manager. Trying manual installation..."
        install_go_manual
    fi
}

# Function to manually install Go (download from official site)
install_go_manual() {
    echo ""
    echo -e "${BLUE}Installing Go manually from official source...${NC}"
    
    # Detect architecture
    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64)
            GO_ARCH="amd64"
            ;;
        aarch64|arm64)
            GO_ARCH="arm64"
            ;;
        armv7l)
            GO_ARCH="armv6l"
            ;;
        *)
            echo -e "${RED}✗${NC} Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac
    
    # Download latest stable Go
    GO_VERSION="1.23.3"  # Update this to latest stable version
    GO_TAR="go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    GO_URL="https://go.dev/dl/${GO_TAR}"
    
    echo "Downloading Go ${GO_VERSION} for ${GO_ARCH}..."
    wget -q --show-progress "$GO_URL" -O "/tmp/${GO_TAR}"
    
    echo "Extracting Go..."
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf "/tmp/${GO_TAR}"
    
    # Add to PATH if not already there
    if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
    fi
    
    # Export for current session
    export PATH=$PATH:/usr/local/go/bin
    export PATH=$PATH:$HOME/go/bin
    
    # Cleanup
    rm "/tmp/${GO_TAR}"
    
    # Verify installation
    if command -v go &> /dev/null; then
        GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        echo -e "${GREEN}✓${NC} Go ${GO_VERSION} installed successfully!"
        echo -e "${YELLOW}Note:${NC} Please restart your terminal or run: source ~/.bashrc"
    else
        echo -e "${RED}✗${NC} Go installation failed"
        exit 1
    fi
}

# Main installation logic
echo "Checking Go installation..."
if ! check_go; then
    echo ""
    read -p "Do you want to install/upgrade Go now? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        install_go_linux
    else
        echo -e "${RED}✗${NC} Go is required to build this project. Exiting."
        exit 1
    fi
fi

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Building IMAP Migrator${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Tidy and download dependencies
echo "Resolving dependencies..."
go mod tidy
echo -e "${GREEN}✓${NC} Dependencies resolved"

echo ""
echo "Downloading dependencies..."
go mod download
echo -e "${GREEN}✓${NC} Dependencies downloaded"

# Build the executable
echo ""
echo "Compiling..."
go build -o imap-migrator *.go

if [ -f "imap-migrator" ]; then
    echo -e "${GREEN}✓${NC} Build successful!"
    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}  Build Complete!${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo -e "Executable created: ${GREEN}./imap-migrator${NC}"
    echo ""
    echo "To run the migrator:"
    echo -e "  ${BLUE}./imap-migrator${NC}"
    echo ""
    echo "Make sure you have:"
    echo -e "  1. Created ${YELLOW}accounts.csv${NC} with your account details"
    echo -e "  2. Configured ${YELLOW}config.json${NC} with your preferences"
    echo ""
    echo "For help, see:"
    echo "  - README.md (or README-ptbr.md)"
    echo "  - QUICKSTART.md (or QUICKSTART-ptbr.md)"
    echo ""
else
    echo -e "${RED}✗${NC} Build failed"
    exit 1
fi
