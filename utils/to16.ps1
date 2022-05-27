param([string]$Name)

Get-Content $Name | Set-Content -Encoding UTF-16LE ..\..\reactivedrop_public_src\reactivedrop\resource\$Name
