
# Employee Management System - Start All Services
#Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass

$ErrorActionPreference = "Stop"

$ProjectRoot = $PSScriptRoot
$ConsulAddress = "127.0.0.1:8500"
$ConsulExe = "D:\Tenarai\consul_2.1.0_windows_amd64\consul.exe"

$BinDirectory = Join-Path $ProjectRoot "bin"
$PidDirectory = Join-Path $ProjectRoot ".pids"

$Services = @(
    @{ Name = "department-service"; Path = "cmd/department-service"; Port = 8081 },
    @{ Name = "employee-service";   Path = "cmd/employee-service";   Port = 8082 },
    @{ Name = "auth-service";       Path = "cmd/auth-service";       Port = 8083 },
    @{ Name = "api-gateway";        Path = "cmd/api-gateway";        Port = 8080 }
)

function Test-TcpPort {
    param([string]$HostName, [int]$Port)

    $Client = New-Object System.Net.Sockets.TcpClient

    try {
        $Result = $Client.BeginConnect($HostName, $Port, $null, $null)

        if (-not $Result.AsyncWaitHandle.WaitOne(1000, $false)) {
            return $false
        }

        $Client.EndConnect($Result)
        return $true
    }
    catch {
        return $false
    }
    finally {
        $Client.Close()
    }
}

function Wait-ForHealth {
    param([int]$Port, [int]$TimeoutSeconds = 60)

    $Deadline = (Get-Date).AddSeconds($TimeoutSeconds)

    while ((Get-Date) -lt $Deadline) {
        try {
            $Response = Invoke-RestMethod `
                -Uri "http://127.0.0.1:$Port/health" `
                -TimeoutSec 3

            if ($Response.status -eq "UP") {
                return $true
            }
        }
        catch {
            Start-Sleep -Seconds 2
        }
    }

    return $false
}

Write-Host "`n=== EMS STARTUP ===" -ForegroundColor Cyan

# Check root .env.
$EnvFile = Join-Path $ProjectRoot ".env"

if (-not (Test-Path $EnvFile)) {
    throw ".env file not found: $EnvFile"
}

# Check service entry points.
foreach ($Service in $Services) {
    $MainFile = Join-Path $ProjectRoot "$($Service.Path)/main.go"

    if (-not (Test-Path $MainFile)) {
        throw "Service main.go not found: $MainFile"
    }
}

# Start Consul if it is not running.
if (-not (Test-TcpPort -HostName "127.0.0.1" -Port 8500)) {

    if (-not (Test-Path $ConsulExe)) {
        throw "Consul executable not found: $ConsulExe"
    }

    Write-Host "Starting Consul..." -ForegroundColor Yellow

    Start-Process `
        -FilePath $ConsulExe `
        -ArgumentList @("agent", "-dev", "-client=127.0.0.1") `
        -WorkingDirectory (Split-Path $ConsulExe)

    $Deadline = (Get-Date).AddSeconds(30)

    while ((Get-Date) -lt $Deadline) {
        if (Test-TcpPort -HostName "127.0.0.1" -Port 8500) {
            break
        }

        Start-Sleep -Seconds 2
    }
}

try {
    $Leader = Invoke-RestMethod `
        -Uri "http://$ConsulAddress/v1/status/leader" `
        -TimeoutSec 5

    if ([string]::IsNullOrWhiteSpace([string]$Leader)) {
        throw "Consul has no elected leader."
    }
}
catch {
    throw "Consul is not ready: $($_.Exception.Message)"
}

Write-Host "Consul is ready." -ForegroundColor Green

# Create directories.
New-Item -ItemType Directory -Force -Path $BinDirectory | Out-Null
New-Item -ItemType Directory -Force -Path $PidDirectory | Out-Null

# Build services.
Write-Host "`nBuilding EMS services..." -ForegroundColor Yellow

Push-Location $ProjectRoot

try {
    foreach ($Service in $Services) {

        $Executable = Join-Path $BinDirectory "$($Service.Name).exe"

        go build -o $Executable "./$($Service.Path)"

        if ($LASTEXITCODE -ne 0) {
            throw "Build failed: $($Service.Name)"
        }

        Write-Host "Built $($Service.Name)." -ForegroundColor Green
    }
}
finally {
    Pop-Location
}

# Start services.
foreach ($Service in $Services) {

    $Port = [int]$Service.Port
    $PidFile = Join-Path $PidDirectory "$($Service.Name).pid"
    $Executable = Join-Path $BinDirectory "$($Service.Name).exe"

    Write-Host "`nStarting $($Service.Name)..." -ForegroundColor Cyan

    if (Test-TcpPort -HostName "127.0.0.1" -Port $Port) {

        if (Wait-ForHealth -Port $Port -TimeoutSeconds 5) {
            Write-Host "$($Service.Name) is already running." -ForegroundColor Yellow
            continue
        }

        throw "Port $Port is occupied by an unhealthy or unknown process."
    }

    $Process = Start-Process `
        -FilePath $Executable `
        -WorkingDirectory $ProjectRoot `
        -PassThru

    Set-Content -Path $PidFile -Value $Process.Id

    Write-Host "Started with PID $($Process.Id)." -ForegroundColor Yellow

    if (Wait-ForHealth -Port $Port -TimeoutSeconds 60) {
        Write-Host "$($Service.Name) is UP." -ForegroundColor Green
    }
    else {
        Write-Warning "$($Service.Name) did not become healthy."
        throw "Startup failed for $($Service.Name). Check the process."
    }
}

Write-Host "`n=== EMS STARTUP COMPLETED ===" -ForegroundColor Cyan
Write-Host "Consul UI:    http://localhost:8500/ui/"
Write-Host "API Gateway:  http://localhost:8080"
Write-Host "Health check: .\check-health.ps1"
