#!/bin/bash

echo "Preparing .app directory... 🔄"
rm -rf ./build
mkdir -p ./build
mkdir -p ./build/icons.iconset
mkdir -p ./build/SAM.app/Contents/MacOS
mkdir -p ./build/SAM.app/Contents/Resources
echo "Finished preparing .app directory ✅"

echo "Building binary... 🔄"
go build -tags main -o ./build/SAM.app/Contents/MacOS/sam . && \
   echo "Finished building binary ✅" || (echo echo "Binary build failed ❌"; exit 1)

echo "Generating iconset... 🔄"
for SIZE in 16 32  128 256 512; do
  sips -z $SIZE $SIZE ./assets/Logo.png --out ./build/icons.iconset/icon_${SIZE}x${SIZE}.png
done

for SIZE in 16 32 128 256 512; do
  DOUBLE=$((SIZE * 2))
  sips -z $DOUBLE $DOUBLE ./assets/Logo@2x.png --out ./build/icons.iconset/icon_${SIZE}x${SIZE}@2x.png
done

iconutil -c icns -o ./build/SAM.app/Contents/Resources/icon.icns ./build/icons.iconset && \
  echo "Finished generating Iconset ✅" || { echo "Iconset generation failed ❌"; exit 1; }

echo "Building .app... 🔄"
cp ./Info.plist ./build/SAM.app/Contents
touch ./build/SAM.app
echo "Finished building .app ✅"

