# ДЗ 6. Тесты для Ledger и Gateway

Продолжение ДЗ 5, к коду добавились тесты. Gateway это REST API на порту 8080. Он импортирует модуль `ledger` через `replace ledger => ../ledger` в `gateway/go.mod`, данные хранятся в памяти. Чтобы импорт работал, `ledger` стал обычным пакетом, а сценарии из прошлых ДЗ переехали в `ledger/cmd`.

В `gateway/internal/api` лежат DTO с функциями перевода (`dto.go`), обработчики (`handlers.go`) и middleware, которое пишет метод, путь и время запроса (`middleware.go`).

## Запуск

```bash
cd gateway
go run .
```

Все ответы приходят в JSON, ошибки в виде `{"error":"..."}`.

## Эндпоинты

`POST /api/transactions` создает транзакцию: 201 при успехе, 400 если данные не прошли проверку, 409 если превышен бюджет, 500 при других ошибках.

`GET /api/transactions` отдает список транзакций, пустой `[]`, если их нет.

`POST /api/budgets` создает или меняет бюджет: 201 при успехе, 400 при неправильных данных.

`GET /api/budgets` отдает список бюджетов.

## Примеры cURL

```bash
curl -X POST http://localhost:8080/api/budgets \
  -H "Content-Type: application/json" \
  -d '{"category":"еда","limit":1000,"period":"месяц"}'
```

```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{"amount":450,"category":"еда","description":"ланч","date":"2025-09-10"}'
```

Еще 700 на еду уже не влезут в лимит 1000, придет 409:

```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{"amount":700,"category":"еда","description":"ужин","date":"2025-09-10"}'
```

Списки можно открыть прямо в браузере: http://localhost:8080/api/transactions и http://localhost:8080/api/budgets.

## Тесты

Тесты лежат в `ledger/ledger_test.go` и `gateway/internal/api/handlers_test.go`. Модулей два, поэтому запускать в каждой папке отдельно:

```bash
cd ledger
go test ./... -cover -coverprofile=cover.out
go tool cover -html=cover.out -o cover.html
```

```bash
cd gateway
go test ./... -cover -coverprofile=cover.out
go tool cover -html=cover.out -o cover.html
```

Покрытие тестами: пакет `ledger` ~64%, пакет `internal/api` ~79%.
