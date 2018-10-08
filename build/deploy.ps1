

# https://stackoverflow.com/questions/4037939/powershell-says-execution-of-scripts-is-disabled-on-this-system/13696614

Write-Host "Copying files to cloud server."

# Write-Host "main.wasm"
& 'C:\Program Files\PuTTY\pscp' -pw ossTrich '..\deploy\wasm\main.wasm' neuro@178.128.1.110:/home/neuro/test3d/wasm

# Write-Host "test1.htm"
& 'C:\Program Files\PuTTY\pscp' -pw ossTrich '..\deploy\test1.htm' neuro@178.128.1.110:/home/neuro/test3d/

& 'C:\Program Files\PuTTY\pscp' -pw ossTrich '..\deploy\wasm_exec_3.html' neuro@178.128.1.110:/home/neuro/test3d/

# --------------------------------------------------------

Write-Host "Finished"
