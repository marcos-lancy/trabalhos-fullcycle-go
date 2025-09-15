# Script PowerShell para testar o rate limiter
# Execute este script após iniciar o servidor

Write-Host "Testing Rate Limiter..." -ForegroundColor Green
Write-Host "=======================" -ForegroundColor Green

# Teste 1: Health check
Write-Host "1. Testing health check..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get
    $response | ConvertTo-Json
} catch {
    Write-Host "Error: $_" -ForegroundColor Red
}
Write-Host ""

# Teste 2: Rate limiting por IP
Write-Host "2. Testing IP rate limiting (limit: 5 req/s)..." -ForegroundColor Yellow
for ($i = 1; $i -le 7; $i++) {
    Write-Host "Request $i:" -ForegroundColor Cyan
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:8080/api/data" -Method Get
        Write-Host "HTTP Status: $($response.StatusCode)" -ForegroundColor Green
        $response.Content | ConvertFrom-Json | ConvertTo-Json
    } catch {
        if ($_.Exception.Response.StatusCode -eq 429) {
            Write-Host "HTTP Status: 429 (Rate Limited)" -ForegroundColor Red
            $_.Exception.Response.GetResponseStream() | ForEach-Object { 
                $reader = New-Object System.IO.StreamReader($_)
                $reader.ReadToEnd() | ConvertFrom-Json | ConvertTo-Json
            }
        } else {
            Write-Host "Error: $_" -ForegroundColor Red
        }
    }
    Write-Host ""
    Start-Sleep -Milliseconds 100
}

Write-Host "Waiting 2 seconds before next test..." -ForegroundColor Yellow
Start-Sleep -Seconds 2

# Teste 3: Rate limiting por token
Write-Host "3. Testing token rate limiting (limit: 10 req/s)..." -ForegroundColor Yellow
$headers = @{"API_KEY" = "test-token"}
for ($i = 1; $i -le 12; $i++) {
    Write-Host "Request $i:" -ForegroundColor Cyan
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:8080/api/data" -Method Get -Headers $headers
        Write-Host "HTTP Status: $($response.StatusCode)" -ForegroundColor Green
        $response.Content | ConvertFrom-Json | ConvertTo-Json
    } catch {
        if ($_.Exception.Response.StatusCode -eq 429) {
            Write-Host "HTTP Status: 429 (Rate Limited)" -ForegroundColor Red
            $_.Exception.Response.GetResponseStream() | ForEach-Object { 
                $reader = New-Object System.IO.StreamReader($_)
                $reader.ReadToEnd() | ConvertFrom-Json | ConvertTo-Json
            }
        } else {
            Write-Host "Error: $_" -ForegroundColor Red
        }
    }
    Write-Host ""
    Start-Sleep -Milliseconds 100
}

Write-Host "Waiting 2 seconds before next test..." -ForegroundColor Yellow
Start-Sleep -Seconds 2

# Teste 4: Verificar status do rate limit
Write-Host "4. Checking rate limit status..." -ForegroundColor Yellow
Write-Host "IP status:" -ForegroundColor Cyan
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/api/rate-limit/status" -Method Get
    $response | ConvertTo-Json
} catch {
    Write-Host "Error: $_" -ForegroundColor Red
}
Write-Host ""
Write-Host "Token status:" -ForegroundColor Cyan
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/api/rate-limit/status" -Method Get -Headers $headers
    $response | ConvertTo-Json
} catch {
    Write-Host "Error: $_" -ForegroundColor Red
}
Write-Host ""

# Teste 5: Reset rate limit
Write-Host "5. Resetting rate limit..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/api/rate-limit/reset" -Method Post
    $response | ConvertTo-Json
} catch {
    Write-Host "Error: $_" -ForegroundColor Red
}
Write-Host ""

# Teste 6: Verificar se reset funcionou
Write-Host "6. Testing after reset..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/api/data" -Method Get
    $response | ConvertTo-Json
} catch {
    Write-Host "Error: $_" -ForegroundColor Red
}
Write-Host ""

Write-Host "Rate limiter test completed!" -ForegroundColor Green
