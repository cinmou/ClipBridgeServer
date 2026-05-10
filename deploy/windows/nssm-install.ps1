$AppDir = "C:\ClipBridge"
$Binary = Join-Path $AppDir "clipbridge-server-windows-amd64.exe"
$Config = Join-Path $AppDir "config.yaml"
$ServiceName = "ClipBridgeServer"

Write-Host "Installing $ServiceName with NSSM..."
Write-Host "Adjust paths first if your installation directory is different."

nssm install $ServiceName $Binary -config $Config
nssm set $ServiceName AppDirectory $AppDir
nssm set $ServiceName Start SERVICE_AUTO_START
nssm set $ServiceName DisplayName "ClipBridgeServer"
nssm set $ServiceName Description "ClipBridge self-hosted clipboard server"
nssm start $ServiceName
