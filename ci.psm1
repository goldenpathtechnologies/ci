function Invoke-Ci {
    $exitArgs = @("-v", "--version", "-h", "--help")
    $ciExe = "$home\Documents\WindowsPowerShell\Modules\ci\ci.exe"

    if ($args | Where-Object { $exitArgs -contains $_ }) {
        & $ciExe $args
    } else {
        $output = & $ciExe $args

        if ($? -and $null -ne $output) {
            if ((Get-Item $output) -is [System.IO.DirectoryInfo]) {
                Set-Location -Path $output.ToString()
            } else {
                Write-Host "$output is not a directory."

                throw
            }
        } elseif ($null -eq $output) {
            Write-Host "Program forcefully exited"
            
            return
        } else {
            $output

            throw
        }
    }
}

New-Alias -Name ci -Value Invoke-Ci

Export-ModuleMember -Alias * -Function *
