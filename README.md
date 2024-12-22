### Запустить проект можно этой командой)
  go run ./cmd/calc_service/...

###Первый тестик
curl -Uri 'http://localhost:8080/api/v1/calculate' `
     -Method POST `
     -ContentType 'application/json' `
     -Body '{"expression": "2+2*2"}'
