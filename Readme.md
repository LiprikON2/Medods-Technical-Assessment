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
- Второй маршрут выполняет Refresh операцию на пару Access, Refresh токенов

#### Требования

Access токен
- Тип JWT
- Aлгоритм SHA512
  - *Для подписывания JWT токена используется алгоритм HMAC-SHA512*
- Хранить в базе строго запрещено

Refresh токен
- Тип произвольный
  - *Используется так же JWT*
- Формат передачи base64
  - *Подписанные JWT токены возвращаются в виде base64* 
- Хранится в базе исключительно в виде bcrypt хеша
  - *TODO*
- Должен быть защищен от изменения на стороне клиента 
  - *Защита реализована JWT подписью при помощи переменной из среды (`JWT_REFRESH_SECRET`)*
- Должен быть защищен от попыток повторного использования
  - *Защита реализована заданием срока годности `exp` в Payload (8 часов для Refresh токена и 5 минут для Access токена)*
  - *Так же, токен может быть отозван исходя из значения поля `Revoked` в базе данных*
Access, Refresh токены обоюдно связаны, Refresh операцию для Access токена можно выполнить только тем Refresh токеном который был выдан вместе с ним
- *Генерация токенов производится вместе*

Payload токенов должен содержать сведения об ip адресе клиента, которому он был выдан
- *Поле `ip` в Payload, заполняемое при помощи `chi/middleware RealIP`*

В случае, если ip адрес изменился, при рефреш операции нужно послать email warning на почту юзера (для упрощения можно использовать моковые данные)
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



### Structure

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