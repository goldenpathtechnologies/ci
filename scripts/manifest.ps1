New-ModuleManifest -Path .\ci.psd1 `
-Author "Daryl Wright" `
-CompanyName "Golden Path Technologies Inc." `
-RootModule "ci.psm1" `
-FunctionsToExport @("Invoke-Ci") `
-AliasesToExport @("ci") `
-Copyright "(c) Golden Path Technologies Inc." `
-GUID "7f8bcbba-4245-4659-be3b-008de6b6a974" `
-ModuleVersion "0.0.0" `
-Description "A utility that makes traversing directories in the terminal quick and easy" `
-FileList @("bin/ci.exe", "ci.psm1", "LICENSE", "CHANGELOG.md") `
-LicenseUri "https://github.com/goldenpathtechnologies/ci/blob/main/LICENSE" `
-ProjectUri "https://github.com/goldenpathtechnologies/ci" 

$content = Get-Content -Path .\ci.psd1
Out-File -InputObject $content -FilePath .\ci.psd1 -Encoding utf8