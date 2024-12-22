### Запустить проект можно этой командой)
  go run ./cmd/calc_service/...

### Тестики 😋😋😋
### correct 😋
curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2"
}'

curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'

curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "(2+2)*2"
}'

curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "(1 + 2) * (3 - 4) / 2"
}'

curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "  3 + 5 * 2  "
}'

curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "10 + 2 * 6 - 4 / 2"
}'

### no correct


### Выражение с буквой
curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+x"
}'

### Деление на ноль
curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "9/0"
}'

### Не фулл выражение 
curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "(9+1"
}'

curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "5 🥲 4"
}'

🥲🥲🥲🥲🥲🥲🥲🥲🥲🥲🥲### Надеюсь все ###🥲🥲🥲🥲🥲🥲🥲🥲🥲🥲🥲🥲





























