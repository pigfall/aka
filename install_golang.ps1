<#
.SYNOPSIS
<<<<<<< HEAD
    Downloads and installs GoLang, then adds it to the current user's PATH environment variable.
    Allows specifying the download region (Global or CN).
    Displays download progress during the process.

.DESCRIPTION
    This script automates the process of downloading a specified version of GoLang,
    extracting it to a target directory, and configuring the current user's PATH
    environment variable to include the GoLang binary directory.
    It checks for existing GoLang installations and provides options to proceed.
    The download URL can be explicitly set for "Global" (golang.org) or "CN" (golang.google.cn) regions.
    The Go installation will be placed directly in the specified InstallPath (e.g., $HOME/tools/go),
    avoiding an extra 'go' subdirectory.
=======
    Downloads and installs GoLang, then adds it to the system's PATH environment variable.

.DESCRIPTION
    This script automates the process of downloading a specified version of GoLang,
    extracting it to a target directory, and configuring the system's PATH
    environment variable to include the GoLang binary directory.
    It checks for existing GoLang installations and provides options to proceed.
>>>>>>> 73f95b8 (up)

.PARAMETER GoVersion
    Specifies the version of GoLang to download (e.g., "1.22.4").
    If not provided, the script will use a default version.

.PARAMETER InstallPath
    Specifies the directory where GoLang will be installed.
<<<<<<< HEAD
    If not provided, it defaults to "$HOME/tools/go".

.PARAMETER Region
    Specifies the download region for GoLang.
    "Global" uses golang.org, "CN" uses golang.google.cn.
    Defaults to "Global".
=======
    If not provided, it defaults to "C:\Go".
>>>>>>> 73f95b8 (up)

.NOTES
    Author: Gemini
    Date: July 4, 2025
<<<<<<< HEAD
    This script modifies the user's environment variables, which typically does not
    require Administrator privileges. However, creating the installation directory
    (e.g., "$HOME/tools/go") might require appropriate user permissions if the
    parent directories are restricted.
=======
    Requires Administrator privileges to modify system environment variables.
>>>>>>> 73f95b8 (up)
#>

param(
    [string]$GoVersion = "1.22.4", # Default GoLang version
<<<<<<< HEAD
    [string]$InstallPath = "$HOME/tools/go", # Default installation path, set to user's home directory
    [ValidateSet("Global", "CN")][string]$Region = "Global" # Parameter for specifying region
)

Write-Host "Starting GoLang installation script..." -ForegroundColor Green
Write-Host "Target GoLang Version: ${GoVersion}" -ForegroundColor Cyan
Write-Host "Installation Path: ${InstallPath}" -ForegroundColor Cyan
Write-Host "Download Region: ${Region}" -ForegroundColor Cyan

# --- Determine download URL based on specified region ---
$baseDownloadUrl = "https://golang.org/dl/"
if (${Region} -eq "CN") {
    Write-Host "Using download URL for CN: golang.google.cn" -ForegroundColor Green
    $baseDownloadUrl = "https://golang.google.cn/dl/"
} else {
    Write-Host "Using global download URL: golang.org" -ForegroundColor Green
}

$goDownloadUrl = "${baseDownloadUrl}go${GoVersion}.windows-amd64.zip"
$zipFileName = "go${GoVersion}.windows-amd64.zip"
$downloadPath = Join-Path $env:TEMP ${zipFileName}

# --- Resolve the full installation path ---
# This ensures $HOME is correctly expanded before path operations
$resolvedInstallPath = Convert-Path ${InstallPath}

# --- Check if GoLang is already installed at the target path ---
if (Test-Path ${resolvedInstallPath}) {
    Write-Host "A directory already exists at '${resolvedInstallPath}'." -ForegroundColor Yellow
=======
    [string]$InstallPath = "~/tools/go" # Default installation path
)

# --- Function to check for Administrator privileges ---
function Test-IsAdministrator {
    $currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
    $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

# --- Check for Administrator privileges ---
if (-not (Test-IsAdministrator)) {
    Write-Host "This script requires Administrator privileges to modify system environment variables." -ForegroundColor Red
    Write-Host "Please run PowerShell as Administrator and try again." -ForegroundColor Red
    exit 1
}

Write-Host "Starting GoLang installation script..." -ForegroundColor Green
Write-Host "Target GoLang Version: $GoVersion" -ForegroundColor Cyan
Write-Host "Installation Path: $InstallPath" -ForegroundColor Cyan

# --- Define GoLang download URL ---
$goDownloadUrl = "https://golang.org/dl/go$GoVersion.windows-amd64.zip"
$zipFileName = "go$GoVersion.windows-amd64.zip"
$downloadPath = Join-Path $env:TEMP $zipFileName

# --- Check if GoLang is already installed at the target path ---
if (Test-Path $InstallPath) {
    Write-Host "A directory already exists at '$InstallPath'." -ForegroundColor Yellow
>>>>>>> 73f95b8 (up)
    $response = Read-Host "Do you want to remove the existing directory and proceed with installation? (Y/N)"
    if ($response -ne 'Y' -and $response -ne 'y') {
        Write-Host "Installation aborted by user." -ForegroundColor Red
        exit 0
    } else {
<<<<<<< HEAD
        Write-Host "Removing existing directory '${resolvedInstallPath}'..." -ForegroundColor Yellow
        Remove-Item ${resolvedInstallPath} -Recurse -Force -ErrorAction SilentlyContinue
        if (Test-Path ${resolvedInstallPath}) {
=======
        Write-Host "Removing existing directory '$InstallPath'..." -ForegroundColor Yellow
        Remove-Item $InstallPath -Recurse -Force -ErrorAction SilentlyContinue
        if (Test-Path $InstallPath) {
>>>>>>> 73f95b8 (up)
            Write-Host "Failed to remove existing directory. Please ensure no files are in use." -ForegroundColor Red
            exit 1
        }
    }
}

# --- Create the installation directory if it doesn't exist ---
<<<<<<< HEAD
if (-not (Test-Path ${resolvedInstallPath})) {
    Write-Host "Creating installation directory: ${resolvedInstallPath}" -ForegroundColor Green
    try {
        New-Item -ItemType Directory -Path ${resolvedInstallPath} -Force | Out-Null
    }
    catch {
        Write-Host "Error creating installation directory '${resolvedInstallPath}': ${_.Exception.Message}" -ForegroundColor Red
        Write-Host "Please ensure you have appropriate permissions to create directories in your home path." -ForegroundColor Red
        exit 1
    }
}

# --- Download GoLang with progress ---
Write-Host "Downloading GoLang from '${goDownloadUrl}' to '${downloadPath}'..." -ForegroundColor Green
try {
    $webClient = New-Object System.Net.WebClient

    # Register the DownloadProgressChanged event
    $webClient.add_DownloadProgressChanged({
        param($sender, $e)
        # Clear the current line and write progress
        Write-Progress -Activity "Downloading GoLang" -Status "Progress: ${e.ProgressPercentage}%" -PercentComplete ${e.ProgressPercentage}
    })

    # Download the file
    $webClient.DownloadFile(${goDownloadUrl}, ${downloadPath})
    Write-Host "`nGoLang downloaded successfully." -ForegroundColor Green # New line after progress bar
}
catch {
    Write-Host "`nError downloading GoLang from '${goDownloadUrl}': ${_.Exception.Message}" -ForegroundColor Red
    exit 1
}

# --- Extract GoLang archive and move contents ---
Write-Host "Extracting GoLang to temporary location and moving contents to '${resolvedInstallPath}'..." -ForegroundColor Green
$tempExtractPath = Join-Path $env:TEMP "go_temp_extract"
try {
    # Create temporary extraction directory
    New-Item -ItemType Directory -Path ${tempExtractPath} -Force | Out-Null

    # Extract the archive to the temporary directory
    Expand-Archive -Path ${downloadPath} -DestinationPath ${tempExtractPath} -Force

    # The zip typically extracts to a 'go' subdirectory inside the temporary path
    $extractedGoDir = Join-Path ${tempExtractPath} "go"

    # Move contents of the 'go' subdirectory to the final install path
    Move-Item -Path "${extractedGoDir}\*" -Destination ${resolvedInstallPath} -Force

    Write-Host "GoLang extracted and moved successfully." -ForegroundColor Green
}
catch {
    Write-Host "Error extracting or moving GoLang: ${_.Exception.Message}" -ForegroundColor Red
    exit 1
}
finally {
    # Clean up temporary extraction directory
    if (Test-Path ${tempExtractPath}) {
        Write-Host "Cleaning up temporary extraction directory..." -ForegroundColor Green
        Remove-Item ${tempExtractPath} -Recurse -Force -ErrorAction SilentlyContinue
    }
}

# --- Clean up downloaded zip file ---
if (Test-Path ${downloadPath}) {
    Write-Host "Cleaning up downloaded zip file..." -ForegroundColor Green
    Remove-Item ${downloadPath} -Force -ErrorAction SilentlyContinue
}

# --- Set GoLang environment variables for the current user ---
# Now the bin directory is directly under the resolvedInstallPath
$goBinPath = Join-Path ${resolvedInstallPath} "bin"

# Get current user PATH
$currentUserPath = [System.Environment]::GetEnvironmentVariable("Path", "User")

# Remove existing Go path if it's already there to avoid duplicates and ensure it's at the front
if (${currentUserPath} -like "*${goBinPath}*") {
    Write-Host "Removing existing GoLang path from user PATH to re-add at front..." -ForegroundColor Yellow
    # Replace Join-String with -join operator for broader PowerShell compatibility
    $currentUserPath = ($currentUserPath -split ';') | Where-Object { $_ -ne ${goBinPath} } -join ';'
}

Write-Host "Adding '${goBinPath}' to the front of the current user's PATH environment variable..." -ForegroundColor Green
try {
    # Prepend the new path. Ensure paths are separated by semicolon.
    # Handle cases where currentUserPath might be empty after removal or initially
    if ([string]::IsNullOrEmpty(${currentUserPath})) {
        [System.Environment]::SetEnvironmentVariable("Path", ${goBinPath}, "User")
    } else {
        [System.Environment]::SetEnvironmentVariable("Path", "${goBinPath};${currentUserPath}", "User")
    }
    Write-Host "User PATH updated successfully. You may need to restart your console/system for changes to take effect." -ForegroundColor Green
}
catch {
    Write-Host "Error setting user PATH: ${_.Exception.Message}" -ForegroundColor Red
    exit 1
=======
if (-not (Test-Path $InstallPath)) {
    Write-Host "Creating installation directory: $InstallPath" -ForegroundColor Green
    New-Item -ItemType Directory -Path $InstallPath | Out-Null
}

# --- Download GoLang ---
Write-Host "Downloading GoLang from '$goDownloadUrl' to '$downloadPath'..." -ForegroundColor Green
try {
    Invoke-WebRequest -Uri $goDownloadUrl -OutFile $downloadPath -UseBasicParsing
    Write-Host "GoLang downloaded successfully." -ForegroundColor Green
}
catch {
    Write-Host "Error downloading GoLang: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# --- Extract GoLang archive ---
Write-Host "Extracting GoLang to '$InstallPath'..." -ForegroundColor Green
try {
    Expand-Archive -Path $downloadPath -DestinationPath $InstallPath -Force
    Write-Host "GoLang extracted successfully." -ForegroundColor Green
}
catch {
    Write-Host "Error extracting GoLang: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# --- Clean up downloaded zip file ---
if (Test-Path $downloadPath) {
    Write-Host "Cleaning up downloaded zip file..." -ForegroundColor Green
    Remove-Item $downloadPath -Force -ErrorAction SilentlyContinue
}

# --- Set GoLang environment variables ---
$goBinPath = Join-Path $InstallPath "go\bin" # The zip extracts to a 'go' subdirectory

# Get current system PATH
$currentPath = [System.Environment]::GetEnvironmentVariable("Path", "Machine")

# Check if GoLang path is already in system PATH
if ($currentPath -notlike "*$goBinPath*") {
    Write-Host "Adding '$goBinPath' to system PATH environment variable..." -ForegroundColor Green
    try {
        # Append the new path. Ensure paths are separated by semicolon.
        [System.Environment]::SetEnvironmentVariable("Path", "$currentPath;$goBinPath", "Machine")
        Write-Host "System PATH updated successfully. You may need to restart your console/system for changes to take effect." -ForegroundColor Green
    }
    catch {
        Write-Host "Error setting system PATH: $($_.Exception.Message)" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "'$goBinPath' is already in the system PATH." -ForegroundColor Yellow
}

# --- Verify installation ---
Write-Host "Verifying GoLang installation..." -ForegroundColor Green
Write-Host "Please open a NEW PowerShell or Command Prompt window and run 'go version' to confirm." -ForegroundColor Yellow

# Optionally, try to run go version in the current session (might not work immediately without new session)
try {
    Write-Host "Attempting to run 'go version' in current session (may require new session for full effect):" -ForegroundColor DarkCyan
    & "$goBinPath\go.exe" version
}
catch {
    Write-Host "Could not run 'go version' in current session. Please open a new console." -ForegroundColor Yellow
>>>>>>> 73f95b8 (up)
}

Write-Host "GoLang installation script finished." -ForegroundColor Green

