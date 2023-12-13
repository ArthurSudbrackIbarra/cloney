#!/bin/bash

# set -e: exit as soon as any command fails.
set -e

# Colors.
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

# Check if this script is being run as root.
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}Error: Please run this script as root.${NC}"
    exit 1
fi

# Get OS and Architecture in lowercase.
OPERATING_SYSTEM=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCHITECTURE=$(uname -m | tr '[:upper:]' '[:lower:]')
echo "Operating System: $OPERATING_SYSTEM"
echo "Architecture: $ARCHITECTURE"

# Define file name according to OS and Architecture.
if [ "$OPERATING_SYSTEM" = "linux" ] && { [ "$ARCHITECTURE" = "x86_64" ] || [ "$ARCHITECTURE" = "amd64" ]; }; then
    FILE_NAME="cloney-linux-amd64"
    elif [ "$OPERATING_SYSTEM" = "linux" ] && { [ "$ARCHITECTURE" = "aarch64" ] || [ "$ARCHITECTURE" = "arm64" ]; }; then
    FILE_NAME="cloney-linux-arm64"
    elif [ "$OPERATING_SYSTEM" = "darwin" ] && { [ "$ARCHITECTURE" = "x86_64" ] || [ "$ARCHITECTURE" = "amd64" ]; }; then
    FILE_NAME="cloney-darwin-amd64"
    elif [ "$OPERATING_SYSTEM" = "darwin" ] && { [ "$ARCHITECTURE" = "aarch64" ] || [ "$ARCHITECTURE" = "arm64" ]; }; then
    FILE_NAME="cloney-darwin-arm64"
else
    echo -e "${RED}Error: Unsupported operating system and/or architecture.${NC}"
    exit 1
fi

# Define other variables.
#! CLONEY_VERSION is set automatically during the pipeline that tags the release (.github/workflows/auto_tag.yaml).
#! Keep this value as it is.
CLONEY_VERSION="X.X.X"
BINARY_LOCATION="/usr/local/bin/cloney"

# Download Cloney Zip.
curl -A "Cloney Download Script" -OL \
"https://github.com/ArthurSudbrackIbarra/cloney/releases/download/$CLONEY_VERSION/$FILE_NAME.zip" ||
{
    echo -e "${RED}Error: Failed to download Cloney. Please check your internet connection.${NC}"
    exit 1
}

# Unzip Cloney.
unzip -o "$FILE_NAME.zip" ||
{
    echo -e "${RED}Error: Failed to unzip Cloney. Please install the zip package.${NC}"
    exit 1
}

# Move Cloney to /usr/local/bin.
mv -f "$FILE_NAME/cloney" $BINARY_LOCATION ||
{
    echo -e "${RED}Error: Failed to move Cloney to $BINARY_LOCATION.${NC}"
    exit 1
}

# Make Cloney executable.
chmod +x $BINARY_LOCATION ||
{
    echo -e "${RED}Error: Failed to make Cloney binary executable.${NC}"
    exit 1
}

# Remove trash.
rm -rf "$FILE_NAME.zip" $FILE_NAME

echo
echo -e "${GREEN}Cloney $CLONEY_VERSION was successfully installed!${NC}"
