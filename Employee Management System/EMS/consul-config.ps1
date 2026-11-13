
# Employee Management System - Consul Configuration Check

$ErrorActionPreference = "Stop"

$ConsulAddress = "http://127.0.0.1:8500"

Write-Host "`n========================================"
Write-Host " EMS - Consul Configuration Check"
Write-Host "========================================"

# Step 1: Check Consul.
Write-Host "`n[1/3] Checking Consul..."

try {
    $Leader = Invoke-RestMethod `
        -Uri "$ConsulAddress/v1/status/leader" `
        -TimeoutSec 5

    if ([string]::IsNullOrWhiteSpace([string]$Leader)) {
        throw "Consul has no elected leader."
    }

    Write-Host "Consul is running." -ForegroundColor Green
    Write-Host "Leader: $Leader"
}
catch {
    Write-Host "ERROR: Consul is unavailable." -ForegroundColor Red
    Write-Host $_.Exception.Message
    exit 1
}

# Step 2: Check registered services.
Write-Host "`n[2/3] Checking registered EMS services..."

try {
    $RegisteredServices = Invoke-RestMethod `
        -Uri "$ConsulAddress/v1/agent/services" `
        -Method Get

    $ExpectedServices = @(
        "department-service",
        "employee-service",
        "auth-service"
    )

    foreach ($Name in $ExpectedServices) {

        $Matches = @(
            $RegisteredServices.PSObject.Properties |
                ForEach-Object { $_.Value } |
                Where-Object { $_.Service -eq $Name }
        )

        if ($Matches.Count -gt 0) {
            foreach ($Service in $Matches) {
                Write-Host "$Name : Registered on port $($Service.Port)" -ForegroundColor Green
            }
        }
        else {
            Write-Host "$Name : Not registered" -ForegroundColor Yellow
        }
    }
}
catch {
    Write-Host "ERROR: Unable to read registered services." -ForegroundColor Red
    Write-Host $_.Exception.Message
    exit 1
}

# Step 3: Check health.
Write-Host "`n[3/3] Checking service health..."

$Ports = @{
    "department-service" = 8081
    "employee-service"   = 8082
    "auth-service"       = 8083
    "api-gateway"        = 8080
}

foreach ($Service in $Ports.Keys) {
    try {
        $Response = Invoke-RestMethod `
            -Uri "http://127.0.0.1:$($Ports[$Service])/health" `
            -TimeoutSec 5

        Write-Host "$Service : $($Response.status)" -ForegroundColor Green
    }
    catch {
        Write-Host "$Service : DOWN" -ForegroundColor Red
    }
}

Write-Host "`nConsul check completed." -ForegroundColor Cyan
