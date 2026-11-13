
# Employee Management System - Health Check

$Services = @(
    @{ Name = "Department Service"; Port = 8081 },
    @{ Name = "Employee Service";   Port = 8082 },
    @{ Name = "Auth Service";       Port = 8083 },
    @{ Name = "API Gateway";        Port = 8080 }
)

Write-Host "`n========== EMS HEALTH CHECK ==========" -ForegroundColor Cyan

$AllHealthy = $true

foreach ($Service in $Services) {

    $Healthy = $false

    try {
        $Response = Invoke-RestMethod `
            -Uri "http://localhost:$($Service.Port)/health" `
            -Method Get `
            -TimeoutSec 5

        if ($Response.status -eq "UP") {
            $Healthy = $true
        }
    }
    catch {
        $Healthy = $false
    }

    if ($Healthy) {
        Write-Host "$($Service.Name) : UP" -ForegroundColor Green
    }
    else {
        Write-Host "$($Service.Name) : DOWN" -ForegroundColor Red
        $AllHealthy = $false
    }
}

Write-Host "`n======================================="

if ($AllHealthy) {
    Write-Host "All EMS services are healthy." -ForegroundColor Green
}
else {
    Write-Host "One or more services are unavailable." -ForegroundColor Red
}
