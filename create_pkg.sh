#!/bin/bash

# Set variables
APP_NAME="File Modification Tracker"
VERSION="1.0.0"
BUNDLE_ID="com.example.file-mod-tracker"
OUTPUT_DIR="./dist"

# Create necessary directories
mkdir -p "${OUTPUT_DIR}/flat/root/Applications/${APP_NAME}.app/Contents/MacOS"
mkdir -p "${OUTPUT_DIR}/flat/root/Applications/${APP_NAME}.app/Contents/Resources"

# Copy binary and resources
cp ./bin/file-mod-tracker "${OUTPUT_DIR}/flat/root/Applications/${APP_NAME}.app/Contents/MacOS/"
cp ./config.yaml "${OUTPUT_DIR}/flat/root/Applications/${APP_NAME}.app/Contents/Resources/"

# Create Info.plist
cat > "${OUTPUT_DIR}/flat/root/Applications/${APP_NAME}.app/Contents/Info.plist" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>file-mod-tracker</string>
    <key>CFBundleIdentifier</key>
    <string>${BUNDLE_ID}</string>
    <key>CFBundleName</key>
    <string>${APP_NAME}</string>
    <key>CFBundleVersion</key>
    <string>${VERSION}</string>
</dict>
</plist>
EOF

# Build the package
pkgbuild --root "${OUTPUT_DIR}/flat/root" \
         --identifier "${BUNDLE_ID}" \
         --version "${VERSION}" \
         --install-location "/" \
         "${OUTPUT_DIR}/${APP_NAME}.pkg"

echo "Package created at ${OUTPUT_DIR}/${APP_NAME}.pkg"
