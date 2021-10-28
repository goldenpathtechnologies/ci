# TODO: Install CI as a PowerShell module. Use the following as a guide:
#  https://docs.microsoft.com/en-us/powershell/scripting/developer/module/how-to-write-a-powershell-script-module?view=powershell-7.1
#  This file will be used to install the module.
#  I may also consider publishing to the PSGallery in the future, see below:
#  https://docs.microsoft.com/en-us/powershell/module/powershellget/publish-module?view=powershell-7.1

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
Copy-Item -Path .\bin\go_build_ci.exe -Destination $modulePath
Copy-Item -Path .\ci.psm1 -Destination $modulePath
