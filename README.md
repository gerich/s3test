### Запуск тестов
* Для начала надо установить [HURL](https://hurl.dev/docs/installation.html). Для macos: `brew install hurl`.
* Затем запустить `go run cmd/api/main.go`
* Затем `make tests`

### Добавление новго сервера
* Раскомментить internal/config/config.go:34
* Запустить/Перезапустить `go run cmd/api/main.go`
* `make test-after-add-new-node`

### Получение JWT
Список доступных пользователей и их JWT:
* foo  |  eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiZm9vIn0.6ee7NeQqcFzmPymq-5N3h52P98_JCOKwsrwye_SVfZA
* bar  |  eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiYmFyIn0.mtdZLBruufzEP2TdwnHaSosfaBI-clTWIMjD9Isr8JA
* baz  |  eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiYmF6In0.0YfojECty3dJgSIZTXcG7sIjW2sbtPVYV7TdpDUZ1vc

Можно получить список спомощью `go run cmd/jwt/main.go`
Для тестов используется **foo**


