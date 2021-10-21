$exitArgs = @("-v", "--version", "-h", "--help")

if ($args | Where-Object { $exitArgs -contains $_ }) {
    # TODO: Do not use relative paths for the script.
    # TODO: Use an environment variable for the executable name or ensure the name is consistent across dev and prod.
    .\bin\go_build_ci__Windows_.exe $args
} else {
    $output = .\bin\go_build_ci__Windows_.exe $args

    if ($?) {
        if ((Get-Item $output) -is [System.IO.DirectoryInfo]) {
            Set-Location -Path $output.ToString()
        } else {
            Write-Host $output + " is not a directory."
            exit 1
        }
    } else {
        $output
        exit 1
    }
}
