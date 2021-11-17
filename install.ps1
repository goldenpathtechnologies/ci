# TODO: Consider publishing to the PSGallery in the future, see below:
#  https://docs.microsoft.com/en-us/powershell/module/powershellget/publish-module?view=powershell-7.1
# TODO: Create a manifest file for this module. See below comments for more information.

$relModulePath = "$home\Documents\WindowsPowerShell\Modules\ci"
$ciExe = "$relModulePath\ci.exe"

if (Test-Path $relModulePath) {
    $currentVersion = $(& $ciExe -v | Select-String -Pattern "Version: (\d+\.\d+\.\d+).*").Matches.groups[1].value
    $newVersion = $(.\bin\ci.exe -v | Select-String -Pattern "Version: (\d+\.\d+\.\d+).*").Matches.groups[1].value

    function Get-IntVersion {
        $data = $(Write-Output $@).Split(".")
        return "{0:d}{1:d3}{2:d3}" -f $data[0],$data[1],$data[2]
    }

    if ($(Get-IntVersion $currentVersion) -gt $(Get-IntVersion $newVersion)) {
        Write-Host "A newer version of ci (v$currentVersion) is already installed."
        Write-Host "Please first uninstall v$currentVersion if you would like to install v$newVersion."
        exit 0
    } elseif ($(Get-IntVersion $currentVersion) -eq $(Get-IntVersion $newVersion)) {
        Write-Host "The installed version of ci ($currentVersion) is up to date."
        exit 0
    } else {
        .\uninstall.ps1
    }
}

# Step 1: Create the module directory
# - mkdir $home\Documents\WindowsPowerShell\Modules\ci
$modulePath = $(New-Item -Path $home\Documents\WindowsPowerShell\Modules\ci -ItemType "directory").FullName

# Step 2: Download files to the module directory
# - ci.exe
# - ci.psm1 (contains the function that calls the exe and the export call for that fn)
# - ci.psd1 (module manifest file, see https://adamtheautomator.com/powershell-modules/)
# - README.md
# - LICENSE.txt
# - Above files are the basics, there may be more that need to be downloaded
#   but there may be ones that will be created such as log files and other
#   files that persist information crucial to program operation.
Copy-Item -Path .\bin\ci.exe -Destination $modulePath
Copy-Item -Path .\ci.psm1 -Destination $modulePath
Copy-Item -Path .\LICENSE -Destination $modulePath
Copy-Item -Path .\CHANGELOG.md -Destination $modulePath
