# Medods - Technical Assessment

> [!note] 
> **Используемые технологии:**
> - Go
> - JWT
> - PostgreSQL
> - Docker

### Задание

Написать часть сервиса аутентификации.

Два REST маршрута:

- Первый маршрут выдает пару Access, Refresh токенов для пользователя с идентификатором (GUID) указанным в параметре запроса
  - *TODO*
- Второй маршрут выполняет Refresh операцию на пару Access, Refresh токенов

#### Требования

Access токен
- Тип JWT
- Aлгоритм SHA512
  - *Для подписывания JWT токена используется алгоритм HMAC-SHA512*
- Хранить в базе строго запрещено

Refresh токен
- Тип произвольный
- Формат передачи base64
  - *С учетом лимита на создание bcrypt хеша в 72 символа, максимальное количество байт которое сможет хранить в себе base64 строка (72 * 3/4) 54 байта*
    - *На идентификатор (UUID) уходит 16 байт*
    - *На IP адрес (`netip.Addr`) уходит используется 4 или 16 байт в зависимости от версии IPv4 или IPv6*
- Хранится в базе исключительно в виде bcrypt хеша
  - *Из факта соления через bcrypt следует, что при Refresh операции нельзя найти запись в бд исключительно по Refresh токену*
    - *Нужно хранить GUID пользователя в Access токене*
- Должен быть защищен от изменения на стороне клиента 
  - *Из факта хеширования через bcrypt следует, что подписывать токен (как в JWT) не нужно - проверка целостности осуществляется через bcrypt*
    - Во время Referesh операции bcrypt хеш передаваемого Refresh токена сравнивается с хешом в базе данных
- Должен быть защищен от попыток повторного использования
  - *У хранимых в базе данных Refresh токенов есть поле `Active`, на котором висит ограничение "у пользователя может быть только один активный токен"*
    - *Создание нового токена требует отзыва предыдущих*

- Access, Refresh токены обоюдно связаны, Refresh операцию для Access токена можно выполнить только тем Refresh токеном который был выдан вместе с ним
  - *Во время Referesh операции у Access и Refresh токенов проверяется одинаковый ли у них jti*
- Payload токенов должен содержать сведения об ip адресе клиента, которому он был выдан
  - *В обоих токенах есть поле для ip заполняемое по данным из `chi/middleware RealIP`*

- В случае, если ip адрес изменился, при рефреш операции нужно послать email warning на почту юзера (для упрощения можно использовать моковые данные)
  - *TODO*

Будет плюсом, если получится использовать Docker и покрыть код тестами.


### Running

1. Set up environment variables
```bash
cp .env.example .env
```

2. Run `docker-compose`
```
docker-compose up --build
```

### Developing

Installing uninstalled (but imported) dependencies
```bash
(cd auth && go mod tidy)
```

### Testing
> ref: https://go.dev/doc/code#Testing



### Architecture

Patterns:
- **DDD** - Domain Driven Design
  - "Ensure that you solve valid problem in the optimal way. After that implement the solution in a way that your business will understand without any extra translation from technical language needed"
  - Applied to:
    - Root package [auth.go](./auth/auth.go)
- **DIP** - Dependency Inversion Principle
  - "D" in SOLID
    - "High-level modules should not depend on low-level modules. Both should depend on abstractions"
    - "Abstractions should not depend upon details. Details should depend upon abstractions"
  - Applied to:
    - Internal modules implementing DIP
      - Auth Controller [authcontroller.go](./auth/internal/chi/authcontroller.go)
      - Auth Service [authservice.go](./auth/internal/postgres/authservice.go)
- **CQRS** - Command and Query Responsibility Segregation
  - "Every method should either be a command that performs an action, or a query that returns data to the caller, but not both"
  - Applied to:
    - Auth Service [authservice.go](./auth/internal/postgres/authservice.go)

Reference:
- https://go.dev/tour/list
- https://www.youtube.com/watch?v=8uiZC0l4Ajw
- https://www.reddit.com/r/golang/comments/1310xxl/comment/jhymmry/
- https://www.gobeyond.dev/standard-package-layout/
- https://www.reddit.com/r/golang/comments/wbawx5/comment/ii5m2ox/
- https://threedots.tech/post/ddd-lite-in-go-introduction/
- https://security.stackexchange.com/questions/79577/whats-the-difference-between-hmac-sha256key-data-and-sha256key-data
- https://stackoverflow.com/a/54378384
- https://datatracker.ietf.org/doc/html/rfc7519#section-4.1
- https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries/
- https://security.stackexchange.com/questions/39849/does-bcrypt-have-a-maximum-password-length