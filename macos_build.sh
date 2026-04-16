#!/bin/bash

echo "Preparing .app directory... 🔄"
mkdir -p ./build
mkdir -p ./build/icons.iconset
mkdir -p ./build/SAM.app/Contents/MacOS
mkdir -p ./build/SAM.app/Contents/Resources
echo "Finished preparing .app directory ✅"

echo "Building binary... 🔄"
go build -tags main -o ./build/SAM.app/Contents/MacOS/sam . && \
   echo "Finished building binary ✅" || (echo echo "Binary build failed ❌"; exit 1)

echo "Generating iconset... 🔄"
for SIZE in 16 32 64 128 256 512; do
  sips -z $SIZE $SIZE Logo.png --out ./build/icons.iconset/logo_${SIZE}x${SIZE}.png
done

for SIZE in 32 64 128 256 512 1024; do
  sips -z $SIZE $SIZE Logo@2x.png --out ./build/icons.iconset/icon_$(($SIZE / 2))x$(($SIZE / 2))@2x.png
done

iconutil -c icns -o ./build/SAM.app/Contents/Resources/icon.icns ./build/icons.iconset
echo "Finished generating Iconset ✅"

echo "Building .app... 🔄"
cp ./Info.plist ./build/SAM.app/Contents
echo "Finished building .app ✅"

