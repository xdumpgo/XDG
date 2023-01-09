echo "Building binary"
GOOS=windows go build -ldflags="-s -w" -trimpath
echo "Patching binary"
./patch2.sh xdumpgo.exe
echo "Creating zip"
mv patched-xdumpgo.exe xdumpgo.exe
rm XDG.zip
upx -9 xdumpgo.exe
zip -r XDG.zip xdumpgo.exe static views
