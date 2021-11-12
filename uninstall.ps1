$relModulePath = "$home\Documents\WindowsPowerShell\Modules\ci"

if (Test-Path $relModulePath) {
    Remove-Item -Path $relModulePath -Recurse -Force
} else {
    Write-Host "ci is not installed"
}
