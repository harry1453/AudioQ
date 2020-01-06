@echo off

echo Please make sure rsrc and upx are installed before running this script.
echo To install rsrc, run "go install github.com/akavel/rsrc"
echo To install upx, download the latest release from https://github.com/upx/upx/releases/latest and place the executable in a directory on your PATH.

rsrc -manifest audioQ.manifest -o audioQ.syso

go build -ldflags="-s -w -H windowsgui" -o audioQ.exe

upx --compress-icons=3 --ultra-brute -9 audioQ.exe
