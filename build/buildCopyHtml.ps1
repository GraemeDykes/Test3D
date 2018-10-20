

#https://stackoverflow.com/questions/4037939/powershell-says-execution-of-scripts-is-disabled-on-this-system/13696614

# --------------------------------------------------------

Write-Host "Copy wasm_exec_3.html to deploy."
copy ../webserver/wasm_exec_3.html ../deploy/

Write-Host "Copy main1.js to deploy."
copy ../webserver/js/main1.js ../deploy/js/

Write-Host "Copy main2.js to deploy."
copy ../webserver/js/main2.js ../deploy/js/

Write-Host "Copy main3.js to deploy."
copy ../webserver/js/main3.js ../deploy/js/

Write-Host "Copy main4.js to deploy."
copy ../webserver/js/main4.js ../deploy/js/

# --------------------------------------------------------

Write-Host "Finished"
