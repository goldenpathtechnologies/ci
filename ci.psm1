function Invoke-Ci {
    $exitArgs = @("-v", "--version", "-h", "--help")
    $ciExe = "$home\Documents\WindowsPowerShell\Modules\ci\ci.exe"

    if ($args | Where-Object { $exitArgs -contains $_ }) {
        & $ciExe $args
    } else {
        $output = & $ciExe $args

        if ($?) {
            if ((Get-Item $output) -is [System.IO.DirectoryInfo]) {
                Set-Location -Path $output.ToString()
            } else {
                Write-Host "$output is not a directory."

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
