#!/usr/bin/env bash

echo  "Dumping ${1} to ascii hex..."
hexdump -ve '1/1 "%.2X"' $1 > tmp.bin

echo "Blanket patching..."
sed -i 's/^((\*|\#)|)(git\.zertex\.space|github\.com)/oCM8RMBN945a4gAyxaN1/g' ./tmp.bin

echo "Repacking"
xxd -r -p tmp.bin patched-$1