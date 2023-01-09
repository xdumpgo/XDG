sudo docker run --rm -e IDUG=0:0 -v /home/zertex/go:/media/sf_GOPATH0 -e GO11MODULE=off -e GOPATH=/home/user/work:/media/sf_GOPATH0 -i therecipe/qt:windows_64_static qtdeploy -debug -ldflags= -tags= build windows /media/sf_GOPATH0/src/git.quartzinc.dev/Zertex/XDGv2/cmd/qtui
sudo chown -R zertex:zertex ./
cp deploy/windows/qtui.exe ./xdumpgo.exe
export QT_PLUGIN_PATH=Qt/5.13.2/gcc_64/plugins
export CGO_CXXFLAGS=-std=c++11
qtdeploy -debug build linux
cp deploy/linux/qtui ./xdumpgo
upx -9 xdumpgo.exe
upx -9 xdumpgo
zip -r XDG.zip xdumpgo xdumpgo.exe static views start.bat start.sh
echo "Done!"
