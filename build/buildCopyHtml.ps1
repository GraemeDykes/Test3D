

#https://stackoverflow.com/questions/4037939/powershell-says-execution-of-scripts-is-disabled-on-this-system/13696614

# --------------------------------------------------------

Write-Host "Copy webserver to deploy."
copy ../webserver/wasm_exec_3.html ../deploy/

# --------------------------------------------------------

Write-Host "Finished"
