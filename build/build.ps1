

#Reference https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.management/stop-process?view=powershell-6

Write-Host "Change to webserver folder."
cd ../webserver

Write-Host "Build webserver."
go build

Write-Host "Stop webserver."
Stop-Process -Name "webserver"

Write-Host "Copy webserver to deploy."
copy webserver.exe ../deploy/

cd ../deploy

Write-Host "Start webserver"
Start-Process "./webserver.exe"

cd ../build

# --------------------------------------------------------

Write-Host "Build wasm"

cd ../wasm

$prevGOOS = $Env:GOOS
$prevGOARCH = $Env:GOARCH

$Env:GOOS = "js"
$Env:GOARCH="wasm"

go build -o main.wasm hello.go

$Env:GOOS = $prevGOOS
$Env:GOARCH = $prevGOARCH

Write-Host "Copy wasm to deploy"
copy main.wasm ../deploy/ws/

cd ../build

# --------------------------------------------------------

Write-Host "Finished"
