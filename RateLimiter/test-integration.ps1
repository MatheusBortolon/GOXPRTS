# Rate Limiter Integration Tests (PowerShell)
# Script de teste de integração do Rate Limiter para Windows

Write-Host "===================================" -ForegroundColor Cyan
Write-Host "Rate Limiter - Testes de Integração" -ForegroundColor Cyan
Write-Host "===================================" -ForegroundColor Cyan
Write-Host ""

$BaseUrl = "http://localhost:8080"

# Função auxiliar para fazer requisições
function Invoke-TestRequest {
    param(
        [string]$Url,
        [hashtable]$Headers = @{}
    )
    
    try {
        $response = Invoke-WebRequest -Uri $Url -Headers $Headers -UseBasicParsing -ErrorAction Stop
        return $response.StatusCode
    }
    catch {
        if ($_.Exception.Response) {
            return [int]$_.Exception.Response.StatusCode
        }
        return 0
    }
}

# Teste 1: Health Check
Write-Host "1. Testando endpoint de health..." -ForegroundColor Yellow
$statusCode = Invoke-TestRequest -Url "$BaseUrl/health"
if ($statusCode -eq 200) {
    Write-Host "✓ Health check passou" -ForegroundColor Green
}
else {
    Write-Host "✗ Health check falhou (Status: $statusCode)" -ForegroundColor Red
    exit 1
}
Write-Host ""

# Teste 2: Limite por IP
Write-Host "2. Testando limite por IP (5 req/s)..." -ForegroundColor Yellow
$successCount = 0
$failCount = 0

for ($i = 1; $i -le 10; $i++) {
    $statusCode = Invoke-TestRequest -Url "$BaseUrl/"
    if ($statusCode -eq 200) {
        $successCount++
        Write-Host "Requisição ${i}: OK (200)" -ForegroundColor Green
    }
    elseif ($statusCode -eq 429) {
        $failCount++
        Write-Host "Requisição ${i}: Bloqueada (429)" -ForegroundColor Red
    }
    Start-Sleep -Milliseconds 100
}

Write-Host ""
Write-Host "Resultado: $successCount sucesso, $failCount bloqueadas"
if ($successCount -eq 5 -and $failCount -eq 5) {
    Write-Host "✓ Limite por IP funcionando corretamente" -ForegroundColor Green
}
else {
    Write-Host "⚠ Resultado inesperado. Esperado: 5 sucesso, 5 bloqueadas" -ForegroundColor Yellow
}
Write-Host ""

# Aguardar reset
Write-Host "Aguardando reset do contador..."
Start-Sleep -Seconds 2

# Teste 3: Limite por token
Write-Host "3. Testando limite por token abc123 (10 req/s)..." -ForegroundColor Yellow
$successCount = 0
$failCount = 0

for ($i = 1; $i -le 15; $i++) {
    $headers = @{ "API_KEY" = "abc123" }
    $statusCode = Invoke-TestRequest -Url "$BaseUrl/" -Headers $headers
    if ($statusCode -eq 200) {
        $successCount++
        Write-Host "Requisição ${i}: OK (200)" -ForegroundColor Green
    }
    elseif ($statusCode -eq 429) {
        $failCount++
        Write-Host "Requisição ${i}: Bloqueada (429)" -ForegroundColor Red
    }
    Start-Sleep -Milliseconds 100
}

Write-Host ""
Write-Host "Resultado: $successCount sucesso, $failCount bloqueadas"
if ($successCount -eq 10 -and $failCount -eq 5) {
    Write-Host "✓ Limite por token funcionando corretamente" -ForegroundColor Green
}
else {
    Write-Host "⚠ Resultado inesperado. Esperado: 10 sucesso, 5 bloqueadas" -ForegroundColor Yellow
}
Write-Host ""

# Aguardar reset
Write-Host "Aguardando reset do contador..."
Start-Sleep -Seconds 2

# Teste 4: Prioridade de token
Write-Host "4. Testando prioridade de token sobre IP..." -ForegroundColor Yellow
Write-Host "Fazendo 8 requisições com token (limite 10) - deve passar"
$successCount = 0

for ($i = 1; $i -le 8; $i++) {
    $headers = @{ "API_KEY" = "xyz789" }
    $statusCode = Invoke-TestRequest -Url "$BaseUrl/" -Headers $headers
    if ($statusCode -eq 200) {
        $successCount++
    }
    Start-Sleep -Milliseconds 100
}

if ($successCount -eq 8) {
    Write-Host "✓ Token xyz789 permite mais requisições que limite de IP" -ForegroundColor Green
}
else {
    Write-Host "✗ Falha no teste de prioridade" -ForegroundColor Red
}
Write-Host ""

# Teste 5: Mensagem de bloqueio
Write-Host "5. Testando resposta de bloqueio..." -ForegroundColor Yellow
Start-Sleep -Seconds 2

# Esgotar o limite
for ($i = 1; $i -le 6; $i++) {
    Invoke-TestRequest -Url "$BaseUrl/" | Out-Null
    Start-Sleep -Milliseconds 100
}

try {
    $response = Invoke-WebRequest -Uri "$BaseUrl/" -UseBasicParsing -ErrorAction Stop
}
catch {
    $response = $_.Exception.Response
    $reader = New-Object System.IO.StreamReader($response.GetResponseStream())
    $responseBody = $reader.ReadToEnd()
    
    if ($responseBody -match "you have reached the maximum number of requests") {
        Write-Host "✓ Mensagem de bloqueio correta" -ForegroundColor Green
    }
    else {
        Write-Host "✗ Mensagem de bloqueio incorreta" -ForegroundColor Red
        Write-Host "Resposta: $responseBody"
    }
}
Write-Host ""

Write-Host "===================================" -ForegroundColor Cyan
Write-Host "Testes concluídos!" -ForegroundColor Cyan
Write-Host "===================================" -ForegroundColor Cyan
