# TODO: Consider publishing to the PSGallery in the future, see below:
#  https://docs.microsoft.com/en-us/powershell/module/powershellget/publish-module?view=powershell-7.1

$origLocation = Get-Location
$installLocation = $origLocation

if ($origLocation -like "*\scripts") {
    $installLocation = Get-Location | Split-Path -Parent
}

if (!(Test-Path "$installLocation\bin\ci.exe")) {
    Write-Host "ci executable not present, unable to install"
    exit 1
}

Set-Location $installLocation

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

$modulePath = $(New-Item -Path $home\Documents\WindowsPowerShell\Modules\ci -ItemType "directory").FullName

Copy-Item .\bin\ci.exe,.\scripts\ci.psm1,.\scripts\ci.psd1,.\LICENSE,.\CHANGELOG.md -Destination $modulePath

Set-Location $origLocation
