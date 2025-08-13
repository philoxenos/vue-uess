# Run Backend
Start-Process -NoNewWindow powershell -ArgumentList "-Command", "cd $PSScriptRoot\backend; `$env:CGO_ENABLED=1; go run main.go"

# Run Frontend
Start-Process -NoNewWindow powershell -ArgumentList "-Command", "cd $PSScriptRoot\frontend; npm run dev"

Write-Host "Starting servers..."
Write-Host "Backend API server running on http://localhost:8080"
Write-Host "Frontend development server running on http://localhost:5173"
Write-Host "Press Ctrl+C in each terminal window to stop the servers"
