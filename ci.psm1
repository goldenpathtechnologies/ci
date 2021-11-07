function Invoke-Ci {
    $exitArgs = @("-v", "--version", "-h", "--help")
    $ciExe = "$home\Documents\WindowsPowerShell\Modules\ci\ci.exe"

    if ($args | Where-Object { $exitArgs -contains $_ }) {
        # TODO: Do not use relative paths for the script.
        # TODO: Use an environment variable for the executable name or ensure the name is consistent across dev and prod.
        & $ciExe $args
    } else {
        $output = & $ciExe $args

        if ($?) {
            if ((Get-Item $output) -is [System.IO.DirectoryInfo]) {
                Set-Location -Path $output.ToString()
            } else {
                Write-Host $output + " is not a directory."

                throw
            }
        } else {
            $output

            throw
        }
    }
}

New-Alias -Name ci -Value Invoke-Ci

Export-ModuleMember -Alias * -Function *
