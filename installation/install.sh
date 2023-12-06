#!/bin/bash

# set -e: exit as soon as any command fails.
set -e

# Colors.
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

# Check if this script is being run as root.
if [ "$EUID" -ne 0 ]
then echo -e "${RED}Error: Please run this script as root.${NC}"
    exit 1
fi

# Get OS and Architecture in lowercase.
OPERATING_SYSTEM=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCHITECTURE=$(uname -m | tr '[:upper:]' '[:lower:]')
echo "Operating System: $OPERATING_SYSTEM"
echo "Architecture: $ARCHITECTURE"

# Define file name according to OS and Architecture.
if [ "$OPERATING_SYSTEM" = "linux" ] && { [ "$ARCHITECTURE" = "x86_64" ] || [ "$ARCHITECTURE" = "amd64" ]; }
then
    FILE_NAME="cloney-linux-amd64"
elif [ "$OPERATING_SYSTEM" = "linux" ] && { [ "$ARCHITECTURE" = "aarch64" ] || [ "$ARCHITECTURE" = "arm64" ]; }
then
    FILE_NAME="cloney-linux-arm64"
elif [ "$OPERATING_SYSTEM" = "darwin" ] && { [ "$ARCHITECTURE" = "x86_64" ] || [ "$ARCHITECTURE" = "amd64" ]; }
then
    FILE_NAME="cloney-darwin-amd64"
elif [ "$OPERATING_SYSTEM" = "darwin" ] && { [ "$ARCHITECTURE" = "aarch64" ] || [ "$ARCHITECTURE" = "arm64" ]; }
then
    FILE_NAME="cloney-darwin-arm64"
else
    echo -e "${RED}Error: Unsupported Operating System and/or Architecture.${NC}"
    exit 1
fi

# Define other variables.
CLONEY_VERSION="0.2.0" # Change to 1.0.0 when releasing.
BINARY_LOCATION="/usr/local/bin/cloney"

# Download Cloney Zip.
curl -A "Cloney Download Script" -OL \
"https://github.com/ArthurSudbrackIbarra/cloney/releases/download/$CLONEY_VERSION/$FILE_NAME.zip"
echo "Downloaded Cloney $CLONEY_VERSION."

# Unzip Cloney.
unzip -o $FILE_NAME.zip
echo "Unzipped Cloney $CLONEY_VERSION."

# Move Cloney to /usr/local/bin.
mv -f $FILE_NAME/cloney $BINARY_LOCATION
echo "Moved Cloney $CLONEY_VERSION to $BINARY_LOCATION."

# Make Cloney executable.
chmod +x $BINARY_LOCATION

# Remove Trash.
rm -rf $FILE_NAME.zip $FILE_NAME
echo "Removed trash from installation."

echo
echo -e "${GREEN}Cloney $CLONEY_VERSION was successfully installed!${NC}"
