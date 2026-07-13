$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $PSScriptRoot
$cmdDir = Join-Path $root "cmd\program-launch-manager"
$distDir = Join-Path $root "dist"
$outputName = (-join @([char]0x7A0B, [char]0x5E8F, [char]0x542F, [char]0x52A8, [char]0x7BA1, [char]0x7406, [char]0x5668)) + ".exe"
$outExe = Join-Path $distDir $outputName
$syso = Join-Path $cmdDir "app.syso"
$generatedRc = Join-Path $cmdDir "app.generated.rc"
$pngIcon = Join-Path $root "ICO.png"
$icoIcon = Join-Path $root "ICO.ico"
$buildTime = Get-Date
$displayVersion = $buildTime.ToString("yyyyMMdd_HHmm")
$numericVersion = "{0},{1},{2},{3}" -f $buildTime.Year, $buildTime.Month, $buildTime.Day, ($buildTime.Hour * 100 + $buildTime.Minute)

New-Item -ItemType Directory -Force -Path $distDir | Out-Null

if (Test-Path $pngIcon) {
    $needsIcon = -not (Test-Path $icoIcon)
    if (-not $needsIcon) {
        $needsIcon = (Get-Item $pngIcon).LastWriteTimeUtc -gt (Get-Item $icoIcon).LastWriteTimeUtc
    }
    if ($needsIcon) {
        $env:ICON_SRC = $pngIcon
        $env:ICON_DST = $icoIcon
        @'
import os
from PIL import Image

src = os.environ["ICON_SRC"]
dst = os.environ["ICON_DST"]
img = Image.open(src).convert("RGBA")
base_size = 1024
img.thumbnail((base_size, base_size), Image.Resampling.LANCZOS)
base = Image.new("RGBA", (base_size, base_size), (0, 0, 0, 0))
base.alpha_composite(img, ((base_size - img.width) // 2, (base_size - img.height) // 2))
sizes = [(16, 16), (24, 24), (32, 32), (48, 48), (64, 64), (128, 128), (256, 256)]
base.save(dst, format="ICO", sizes=sizes)
'@ | python -
    }
}

if (-not (Test-Path $icoIcon)) {
    throw "Icon file not found: $icoIcon"
}

Push-Location $cmdDir
try {
    $rc = Get-Content -Raw -Encoding UTF8 "app.rc"
    $rc = $rc -replace '(?m)^FILEVERSION\s+.*$', "FILEVERSION $numericVersion"
    $rc = $rc -replace '(?m)^PRODUCTVERSION\s+.*$', "PRODUCTVERSION $numericVersion"
    $rc = $rc -replace 'VALUE "FileVersion",\s*"[^"]*"', "VALUE `"FileVersion`", `"$displayVersion`""
    $rc = $rc -replace 'VALUE "ProductVersion",\s*"[^"]*"', "VALUE `"ProductVersion`", `"$displayVersion`""
    [System.IO.File]::WriteAllText($generatedRc, $rc, [System.Text.UTF8Encoding]::new($false))
    windres --codepage=65001 -O coff -F pe-x86-64 -i app.generated.rc -o app.syso
} finally {
    Remove-Item -LiteralPath $generatedRc -Force -ErrorAction SilentlyContinue
    Pop-Location
}

Push-Location $root
try {
    go test ./...
    go build -buildvcs=false -ldflags "-H windowsgui -X program-launch-manager/internal/winui.appVersion=$displayVersion" -o $outExe .\cmd\program-launch-manager
} finally {
    Pop-Location
}

if (-not (Test-Path $outExe)) {
    throw "Build output not found: $outExe"
}

Write-Host "Built: $outExe"
Write-Host "Version: $displayVersion"
if (Test-Path $syso) {
    Write-Host "Embedded resource: $syso"
}
