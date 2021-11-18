# TODO: Consider publishing to the PSGallery in the future, see below:
#  https://docs.microsoft.com/en-us/powershell/module/powershellget/publish-module?view=powershell-7.1
# TODO: Create a manifest file for this module. See below comments for more information.

$relModulePath = "$home\Documents\WindowsPowerShell\Modules\ci"
$ciExe = "$relModulePath\ci.exe"

if (Test-Path $relModulePath) {
    function Get-OrdinalVersion {
        $pattern = "Version: (?:(\d+)\.(\d+)\.(\d+))(?:(?:-[a-zA-Z]+\.)(\d+)){0,1}"
        $data = $(Write-Output $args[0] | Select-String -Pattern $pattern).Matches.groups
        return "{0:d}{1:d3}{2:d3}{3:d3}" -f [int]$data[1].value,[int]$data[2].value,[int]$data[3].value,[int]$data[4].value
    }

    function Get-StringVersion {
        $pattern = "Version: (\d+\.\d+\.\d+(-[a-zA-Z]+\.\d+){0,1})"
        return $(Write-Output $args[0] | Select-String -Pattern $pattern).Matches.groups[1].value
    }

    $currentVersionData = $(& $ciExe -v)
    $newVersionData = $(.\bin\ci.exe -v)
    $currentVersionText = $(Get-StringVersion $currentVersionData)
    $newVersionText = $(Get-StringVersion $newVersionData)
    $currentVersion = $(Get-OrdinalVersion $currentVersionData)
    $newVersion = $(Get-OrdinalVersion $newVersionData)

    if ($currentVersion -gt $newVersion) {
        Write-Host "ci is already installed at v$currentVersionText."
        Write-Host "Please uninstall the current version before installing v$newVersionText."
        exit 0
    } elseif ($currentVersion -eq $newVersion) {
        Write-Host "The installed version of ci ($currentVersionText) is up to date."
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
