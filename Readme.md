# Тестовое задание от Medods - Сервис аутентификации

> [!note] 
> **Используемые технологии:**
> - Go
> - JWT
> - PostgreSQL
> - Docker


- [Medods - Technical Assessment](#medods---technical-assessment)
    - [Задание](#задание)
      - [Первый маршрут выдает пару Access, Refresh токенов для пользователя с идентификатором (GUID) указанным в параметре запроса](#первый-маршрут-выдает-пару-access-refresh-токенов-для-пользователя-с-идентификатором-guid-указанным-в-параметре-запроса)
      - [Второй маршрут выполняет Refresh операцию на пару Access, Refresh токенов](#второй-маршрут-выполняет-refresh-операцию-на-пару-access-refresh-токенов)
      - [Требования](#требования)
    - [Running](#running)
    - [Developing](#developing)
    - [Testing](#testing)
    - [Architecture](#architecture)
    - [Endpoints](#endpoints)
      - [`POST /api/v1/auth/register`](#post-apiv1authregister)
        - [Example request 1:](#example-request-1)
        - [Example request 2:](#example-request-2)
        - [Example request 3:](#example-request-3)
        - [Example request 4:](#example-request-4)
        - [Example request 5:](#example-request-5)
        - [Example request 6:](#example-request-6)
      - [`POST /api/v1/auth/login`](#post-apiv1authlogin)
        - [Example request 1:](#example-request-1-1)
        - [Example request 2:](#example-request-2-1)
        - [Example request 3:](#example-request-3-1)
      - [`GET /api/v1/auth/me`](#get-apiv1authme)
        - [Example request 1:](#example-request-1-2)
        - [Example request 2:](#example-request-2-2)
      - [`GET /api/v1/auth/`](#get-apiv1auth)
        - [Example request 1:](#example-request-1-3)
        - [Example request 2:](#example-request-2-3)
        - [Example request 2:](#example-request-2-4)
      - [`POST /api/v1/auth/login/{GUID}`](#post-apiv1authloginguid)
        - [Example request 1:](#example-request-1-4)
        - [Example request 2:](#example-request-2-5)
        - [Example request 3:](#example-request-3-2)
      - [`POST /api/v1/auth/refresh`](#post-apiv1authrefresh)
        - [Example request 1:](#example-request-1-5)
        - [Example request 2:](#example-request-2-6)
        - [Example request 3:](#example-request-3-3)
        - [Example request 4:](#example-request-4-1)
        - [Example request 5:](#example-request-5-1)
        - [Example request 6:](#example-request-6-1)
        - [Example request 7:](#example-request-7)
      - [`GET /api/v1/auth/{GUID}`](#get-apiv1authguid)
        - [Example request 1:](#example-request-1-6)
        - [Example request 2:](#example-request-2-7)
      - [`PATCH /api/v1/auth/{GUID}`](#patch-apiv1authguid)
        - [Example request 1:](#example-request-1-7)
        - [Example request 2:](#example-request-2-8)
        - [Example request 3:](#example-request-3-4)
        - [Example request 4:](#example-request-4-2)
        - [Example request 5:](#example-request-5-2)
      - [`DELETE /api/v1/auth/{GUID}`](#delete-apiv1authguid)
        - [Example request 1:](#example-request-1-8)
        - [Example request 2:](#example-request-2-9)
        - [Example request 3:](#example-request-3-5)
        - [Example request 4:](#example-request-4-3)


### Задание

Написать часть сервиса аутентификации.

Два REST маршрута:

#### Первый маршрут выдает пару Access, Refresh токенов для пользователя с идентификатором (GUID) указанным в параметре запроса
> `POST http://localhost:8080/api/v1/auth/login/{GUID}`

Пример запроса `POST http://localhost:8080/api/v1/auth/login/cec24247-497f-48f2-8a93-4ccdc2fdd65b`


Тело
```
(пусто)
```

Пример ответа
```
{
  "accessToken": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2OTk5MTcsImlhdCI6MTczMzY5OTYxNywiaXAiOiIxNzIuMTguMC4xIiwianRpIjoiMDBlMWFhNGMtN2YxMi00Y2Y4LWIyNGEtYmU1YTM2MjgyNTZjIn0.s9IfgsetfB1HzArmbAz3Rlh2Z4sGA5u2spcB3TT4Q9DFlDFLu9v7R-_kmHfeW1ugEwUnhpOQeja3FDNeHzuMIg",
  "refreshToken": "AOGqTH8STPiySr5aNiglbKwSAAE="
}
```

Реализация [./auth/internal/chi/authcontroller.go#L285](./auth/internal/chi/authcontroller.go#L285)


#### Второй маршрут выполняет Refresh операцию на пару Access, Refresh токенов
> `POST http://localhost:8080/api/v1/auth/refresh`

Пример запроса `POST http://localhost:8080/api/v1/auth/refresh`

Тело
```
{
  "accessToken": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2OTk5MTcsImlhdCI6MTczMzY5OTYxNywiaXAiOiIxNzIuMTguMC4xIiwianRpIjoiMDBlMWFhNGMtN2YxMi00Y2Y4LWIyNGEtYmU1YTM2MjgyNTZjIn0.s9IfgsetfB1HzArmbAz3Rlh2Z4sGA5u2spcB3TT4Q9DFlDFLu9v7R-_kmHfeW1ugEwUnhpOQeja3FDNeHzuMIg",
  "refreshToken": "AOGqTH8STPiySr5aNiglbKwSAAE="
}
```

Пример ответа
```
{
  "accessToken": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2OTk5ODIsImlhdCI6MTczMzY5OTY4MiwiaXAiOiIxNzIuMTguMC4xIiwianRpIjoiMDNjYzA2ZDMtNzc4OS00NjQ2LWIzNzYtOTU1NjU4NmFiNzcxIn0.Xx_9FPHK4nsUyS18F1t5bM56EbLPZ4tE_XWQtvMir6KEqm5GQKRgGoKuflZBJwfsNoKmzxGvEXYZvFGeykkZpQ",
  "refreshToken": "A8wG03eJRkazdpVWWGq3cawSAAE="
}
```

Реализация [./auth/internal/chi/authcontroller.go#L302](./auth/internal/chi/authcontroller.go#L302)



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
  - *Реализация [./auth/internal/jwt/jwtservice.go#L62](./auth/internal/jwt/jwtservice.go#L62)*
- Хранится в базе исключительно в виде bcrypt хеша
  - *Из факта соления через bcrypt следует, что при Refresh операции нельзя найти запись в бд исключительно по Refresh токену*
    - *Значит, нужно дополнительно хранить какой-нибудь идентификатор в Payload*
  - *Создание экземпляра для хранения в бд [./auth/internal/chi/authcontroller.go#L471](./auth/internal/chi/authcontroller.go#L471)*
- Должен быть защищен от изменения на стороне клиента 
  - *Из факта хеширования через bcrypt следует, что подписывать токен (как в JWT) не нужно - проверка целостности осуществляется через bcrypt*
    - *Во время Referesh операции bcrypt хеш передаваемого Refresh токена сравнивается с хешом в базе данных*
- Должен быть защищен от попыток повторного использования
  - *У хранимых в базе данных Refresh токенов есть поле `Active`, на котором висит ограничение "у пользователя может быть только один активный токен"*
    - *Создание нового токена требует отзыва предыдущих*

- Access, Refresh токены обоюдно связаны, Refresh операцию для Access токена можно выполнить только тем Refresh токеном который был выдан вместе с ним
  - *Во время Referesh операции у Access и Refresh токенов проверяется одинаковый ли у них jti*
- Payload токенов должен содержать сведения об ip адресе клиента, которому он был выдан
  - *В обоих токенах есть поле для ip заполняемое по данным из `chi/middleware RealIP`*

- В случае, если ip адрес изменился, при рефреш операции нужно послать email warning на почту юзера (для упрощения можно использовать моковые данные)
  - *Реализация [./auth/internal/smtp/mailservice.go](./auth/internal/smtp/mailservice.go)*

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

Run all tests
```
(cd auth && go test ./...)
```


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




___


### Endpoints


#### `POST /api/v1/auth/register`

##### Example request 1:

Body
```json
{
  "email": "email2@example.com",
  "password": "Hello1234!"
}
```

Example response:
```json
{
  "accessToken": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2OTkyMTQsImlhdCI6MTczMzY5ODkxNCwiaXAiOiIxNzIuMTguMC4xIiwianRpIjoiOTBmOWEwZTQtZDU1MC00MjM1LTg0ZTgtM2Q3YmQzNjAyODUwIn0.hbORB2206F0CIlsscxxtY1Pm1AQ2q4yQKghck7HBAjX7lCMUmjtxMKkXQgiLBrVToKWUmobMGs1wj1LLW13leg",
  "refreshToken": "kPmg5NVQQjWE6D1702AoUKwSAAE="
}
```

##### Example request 2:

Body
```json
{
  "email": "email@example.com",
  "password": "hello1234"
}
```

Example response:
```json
{
  "code": 422,
  "message": {
    "errors": [
      {
        "field": "password",
        "message": "password must contain at least one uppercase letter, one lowercase letter, one number, and one special character"
      }
    ]
  }
}
```

##### Example request 3:

Body
```json
{
  "email": "email@example.com",
  "password": "Hello1234!"
}
```

Example response:
```json
{
  "code": 409,
  "message": "user with this email already exists"
}
```

##### Example request 4:

Body
```json
{
  "email": "email@example.com"
}
```

Example response:
```json
{
  "code": 422,
  "message": {
    "errors": [
      {
        "field": "password",
        "message": "password is required"
      }
    ]
  }
}
```

##### Example request 5:

Body
```json
{
  "hi": "hello"
}
```

Example response:
```json
{
  "code": 400,
  "message": "json: unknown field \"hi\""
}
```

##### Example request 6:

Body
```json
{
  "email": "email@example.com",
  "password": 1234568798
}
```

Example response:
```json
{
  "code": 400,
  "message": "json: cannot unmarshal number into Go struct field CreateUserDto.password of type string"
}
```


___

#### `POST /api/v1/auth/login`


##### Example request 1:

Body
```json
{
  "email": "email@example.com",
  "password": "Hello1234!"
}
```

Example response:
```json
{
  "accessToken": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2Mzc4NTcsImlhdCI6MTczMzYzNzU1NywiaXAiOiIxNzIuMTguMC4xIiwianRpIjoiMDgxNmIyMjctYjk5Ny00ZDQyLTlmY2MtNzUwMmU3MmRjZTRlIiwic3ViIjoiODk4YmU3NjctZjY2Zi00OTRkLWJlOWEtYzFiZTg1NTQ4YmI3In0.wXGy0c_iO1rw4XfZJwJpYg7cu6i1ZGKLqwr8GaRKgGC-V3Mntzan580ZNeurA0SW7LUl2770BaLmhRpy4wNf_A",
  "refreshToken": "CBayJ7mXTUKfzHUC5y3OTqwSAAE="
}
```

##### Example request 2:

Body
```json
{
  "email": "email@example.com",
  "password": "wrongpass"
}
```

Example response:
```json
{
  "code": 403,
  "message": "crypto/bcrypt: hashedPassword is not the hash of the given password"
}
```

##### Example request 3:

Body
```json
{
  "email": "non.existent@example.com",
  "password": "Hello1234!"
}
```

Example response:
```json
{
  "code": 404,
  "message": "user with email non.existent@example.com not found: sql: no rows in result set"
}
```


___

#### `GET /api/v1/auth/me`
- Requires header `Authorization: Bearer eyJhb...`

##### Example request 1:

Headers
```
Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2OTk3NjUsImlhdCI6MTczMzY5OTQ2NSwiaXAiOiIxNzIuMTguMC4xIiwianRpIjoiZDIwODA4ZWYtYWVlNS00MzI0LTg4ZTYtOGU5OTMyYjM3NmQ0In0.k4EnfkiENdBBk0SMvyjLLqgwgTyEc9YVg6dtyvmnX8A9YIk2-7SQPTzwLzPn3dPKUuVu-PxDa1ul9gq8emZ3LA
```


Example response:
```json
{
  "uuid": "239a0eb8-ca36-4930-b087-81a1d6532005",
  "email": "emai1l2@example.com"
}
```

##### Example request 2:

Headers (bad accessToken)
```
Authorization: Bearer BADxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2OTk3NjUsImlhdCI6MTczMzY5OTQ2NSwiaXAiOiIxNzIuMTguMC4xIiwianRpIjoiZDIwODA4ZWYtYWVlNS00MzI0LTg4ZTYtOGU5OTMyYjM3NmQ0In0.k4EnfkiENdBBk0SMvyjLLqgwgTyEc9YVg6dtyvmnX8A9YIk2-7SQPTzwLzPn3dPKUuVu-PxDa1ul9gq8emZ3LA
```


Example response:
```json
{
  "code": 403,
  "message": "error verifying Authorization header: token signature is invalid: signature is invalid"
}
```



___

#### `GET /api/v1/auth/`
- Requires header `Authorization: Bearer eyJhb...`


##### Example request 1:

Headers
```
Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2MzgwMjgsImlhdCI6MTczMzYzNzcyOCwiaXAiOiIxNzIuMTguMC4xIiwianRpIjoiZjk5YjkyNTMtNGIwYi00ZDFlLTg5NmUtOWVjYmMxNzJjNWUwIiwic3ViIjoiODk4YmU3NjctZjY2Zi00OTRkLWJlOWEtYzFiZTg1NTQ4YmI3In0.EYGs2uu7KXIwOJT_xf1bQ1xjxJvMrfhG0gwu67d89mZEJnvV7TWwyp3WmB3UOQSppRALkzxTV9fkdcpp19wEJw
```

Example response:
```json
[
  {
    "email": "email@example.com"
  }
]
```



##### Example request 2:

Headers
```
(empty)
```


Example response:
```json
{
  "code": 403,
  "message": "error verifying Authorization header: Authorization header must be provided"
}
```

##### Example request 2:

Headers (expired access token)
```
Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2MzgwMjgsImlhdCI6MTczMzYzNzcyOCwiaXAiOiIxNzIuMTguMC4xIiwianRpIjoiZjk5YjkyNTMtNGIwYi00ZDFlLTg5NmUtOWVjYmMxNzJjNWUwIiwic3ViIjoiODk4YmU3NjctZjY2Zi00OTRkLWJlOWEtYzFiZTg1NTQ4YmI3In0.EYGs2uu7KXIwOJT_xf1bQ1xjxJvMrfhG0gwu67d89mZEJnvV7TWwyp3WmB3UOQSppRALkzxTV9fkdcpp19wEJw
```

Example response:
```json
{
  "code": 403,
  "message": "error verifying Authorization header: token has invalid claims: token is expired"
}
```



___

#### `POST /api/v1/auth/login/{GUID}`

##### Example request 1:

`POST /api/v1/auth/login/898be767-f66f-494d-be9a-c1be85548bb7`

Body
```json
(empty)
```

Example response:
```json
{
  "accessToken": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2Mzg1NTIsImlhdCI6MTczMzYzODI1MiwiaXAiOiIxNzIuMTguMC4xIiwianRpIjoiNzhkN2EwNTEtYjU0ZS00NDczLWI5MTEtYjE0YzIwZjI5OGJjIiwic3ViIjoiODk4YmU3NjctZjY2Zi00OTRkLWJlOWEtYzFiZTg1NTQ4YmI3In0.9Y7pkxA0iTsp2XX8sbDh3oBtHOPbFTztet-QsMTUkH4mlE3MGjqTekoiqtAphXKEjBN-EIbKJqLeZI5wa6uKtw",
  "refreshToken": "eNegUbVORHO5EbFMIPKYvKwSAAE="
}
```

##### Example request 2:

`POST /api/v1/auth/login/malformed-f66f-494d-be9a-c1be85548bb7`

Body
```json
(empty)
```

Example response:
```json
{
  "code": 400,
  "message": "invalid UUID format: invalid UUID length: 37"
}
```

##### Example request 3:

`POST /api/v1/auth/login/00000000-f66f-494d-be9a-c1be85548bb7`

Body
```json
(empty)
```

Example response:
```json
{
  "code": 404,
  "message": "user not found: sql: no rows in result set"
}
```


___

#### `POST /api/v1/auth/refresh`


##### Example request 1:


Body
```json
{
  "accessToken": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2Mzg3MTUsImlhdCI6MTczMzYzODQxNSwiaXAiOiIxNzIuMTguMC4xIiwianRpIjoiNTI3NTdlNjQtMDlkNi00MWI5LTg4YjgtZjFkMDljNzRhN2IyIiwic3ViIjoiODk4YmU3NjctZjY2Zi00OTRkLWJlOWEtYzFiZTg1NTQ4YmI3In0.g7Zl5_MXMN3IyURgR3VIHjgZZrq1o6YEp3odrSsoQV8YDBj2iqhxEHOkcdRVQsai1xjH53f5pyw1867CQQf4Lg",
  "refreshToken": "UnV+ZAnWQbmIuPHQnHSnsqwSAAE="
}
```

Example response:
```json
{
  "accessToken": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2Mzg3NDUsImlhdCI6MTczMzYzODQ0NSwiaXAiOiIxNzIuMTguMC4xIiwianRpIjoiZmY5N2FlNjMtZjdjMC00NTAyLWIwYmEtNzY4YjAyNmI5N2I3Iiwic3ViIjoiODk4YmU3NjctZjY2Zi00OTRkLWJlOWEtYzFiZTg1NTQ4YmI3In0.4zNEUcQCt0FR03se6x0UN_hYVqRrvBKtjCwGj6OI-k1glfSFdpRfripXFAshNI6RyH07hhZOVNB7oEBJwZH-hg",
  "refreshToken": "/5euY/fARQKwunaLAmuXt6wSAAE="
}
```


##### Example request 2:


Body (missing refreshToken)
```json
{
  "accessToken": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2Mzg3MTUsImlhdCI6MTczMzYzODQxNSwiaXAiOiIxNzIuMTguMC4xIiwianRpIjoiNTI3NTdlNjQtMDlkNi00MWI5LTg4YjgtZjFkMDljNzRhN2IyIiwic3ViIjoiODk4YmU3NjctZjY2Zi00OTRkLWJlOWEtYzFiZTg1NTQ4YmI3In0.g7Zl5_MXMN3IyURgR3VIHjgZZrq1o6YEp3odrSsoQV8YDBj2iqhxEHOkcdRVQsai1xjH53f5pyw1867CQQf4Lg"
}
```

Example response:
```json
{
  "code": 403,
  "message": "crypto/bcrypt: hashedPassword is not the hash of the given password"
}
```


##### Example request 3:


Body (missing accessToken)
```json
{
  "refreshToken": "UnV+ZAnWQbmIuPHQnHSnsqwSAAE="
}
```

Example response:
```json
{
  "code": 400,
  "message": "token is malformed: token contains an invalid number of segments"
}
```

##### Example request 4:


Body (bad accessToken)
```json
{
  "accessToken": "badhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2Mzg4NTMsImlhdCI6MTczMzYzODU1MywiaXAiOiIxNzIuMTguMC4xIiwianRpIjoiN2IxMDFkNDAtZDZhOC00OTNkLWE3NjMtYzNjNjlmYjQ5YTdhIiwic3ViIjoiODk4YmU3NjctZjY2Zi00OTRkLWJlOWEtYzFiZTg1NTQ4YmI3In0.U-TDrcRpRzMoao3gB6D8lLiOKnZf81ZLGUx1NLRNI7TNgyimaYXUftu7NMdPMGqaN3VmCHQsY6LXptoBUASBaQ",
  "refreshToken": "exAdQNaoST2nY8PGn7SaeqwSAAE="
}
```

Example response:
```json
{
  "code": 400,
  "message": "token is malformed: could not JSON decode header: invalid character 'm' looking for beginning of value"
}
```

##### Example request 5:


Body (bad refreshToken)
```json
{
  "accessToken": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2Mzg4NTMsImlhdCI6MTczMzYzODU1MywiaXAiOiIxNzIuMTguMC4xIiwianRpIjoiN2IxMDFkNDAtZDZhOC00OTNkLWE3NjMtYzNjNjlmYjQ5YTdhIiwic3ViIjoiODk4YmU3NjctZjY2Zi00OTRkLWJlOWEtYzFiZTg1NTQ4YmI3In0.U-TDrcRpRzMoao3gB6D8lLiOKnZf81ZLGUx1NLRNI7TNgyimaYXUftu7NMdPMGqaN3VmCHQsY6LXptoBUASBaQ",
  "refreshToken": "baddQNaoST2nY8PGn7SaeqwSAAE="
}
```

Example response:
```json
{
  "code": 403,
  "message": "crypto/bcrypt: hashedPassword is not the hash of the given password"
}
```


##### Example request 6:


Body (revoked refreshToken)
```json
{
  "accessToken": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2Mzg5NDEsImlhdCI6MTczMzYzODY0MSwiaXAiOiIxNzIuMTguMC4xIiwianRpIjoiNjVhMmIwNzQtMjk2NC00OWZkLWE0ZmQtZWUyM2ZhZTNmMWIzIiwic3ViIjoiODk4YmU3NjctZjY2Zi00OTRkLWJlOWEtYzFiZTg1NTQ4YmI3In0.YNJVAB7ICxTrTvxCu06OuJiGZLjhyGdOduoYzWmdRLYpzuoxh1kAt8Y-iIYxoYqVPmFHMqWTTpFB7Bfix_HXdA",
  "refreshToken": "ZaKwdClkSf2k/e4j+uPxs6wSAAE="
}
```

Example response:
```json
{
  "code": 403,
  "message": "crypto/bcrypt: hashedPassword is not the hash of the given password"
}
```

##### Example request 7:


Body (new refreshToken with old accessToken)
```json
{
  "accessToken": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2Mzg5ODMsImlhdCI6MTczMzYzODY4MywiaXAiOiIxNzIuMTguMC4xIiwianRpIjoiNGJjZWU0NjYtZTcxNC00NzE2LWFkODMtNjQyN2RhZmMwZGFmIiwic3ViIjoiODk4YmU3NjctZjY2Zi00OTRkLWJlOWEtYzFiZTg1NTQ4YmI3In0.0uyY8_Zc2oCXQPspeDE88zo2u5AIL0pbroqpT5niTvONnu6xGjt-ZNXXCmBqfvR86WxGcPp2IdAcKj5TsDpGbA",
  "refreshToken": "HGYsNpzcS4iLfmhamWjdqKwSAAE="
}
```

Example response:
```json
{
  "code": 403,
  "message": "jti in accessToken and refreshToken do not match"
}
```

___

#### `GET /api/v1/auth/{GUID}`

##### Example request 1:

`GET /api/v1/auth/898be767-f66f-494d-be9a-c1be85548bb7`

Example response:
```json
{
  "email": "email@example.com"
}
```

##### Example request 2:

`GET /api/v1/auth/badbe767-f66f-494d-be9a-c1be85548bb7`

Example response:
```json
{
  "code": 404,
  "message": "user not found: sql: no rows in result set"
}
```


___

#### `PATCH /api/v1/auth/{GUID}`
- Requires header `Authorization: Bearer eyJhb...`



##### Example request 1:

Body
```json
{
  "email": "NEW_EMAIL@example.com",
  "password": "NewPass123!"
}
```

Example response:
```json
{
  "uuid": "898be767-f66f-494d-be9a-c1be85548bb7",
  "email": "NEW_EMAIL@example.com"
}
```

##### Example request 2:

Body
```json
{
  "email": "third@example.com"
}
```

Example response:
```json
{
  "uuid": "898be767-f66f-494d-be9a-c1be85548bb7",
  "email": "third@example.com"
}
```

##### Example request 3:

Body
```json
{
  "uuid": "000000-f66f-494d-be9a-c1be85548bb7",
  "email": "third@example.com"
}
```

Example response:
```json
{
  "code": 400,
  "message": "json: unknown field \"uuid\""
}
```

##### Example request 4:

Body
```json
{}
```

Example response:
```json
{
  "uuid": "898be767-f66f-494d-be9a-c1be85548bb7",
  "email": "third@example.com"
}
```

##### Example request 5:

Body
```json
{
  "email": "bad_new_email",
  "password": "bad_new_pass"
}
```

Example response:
```json
{
  "code": 422,
  "message": {
    "errors": [
      {
        "field": "email",
        "message": "invalid email format"
      },
      {
        "field": "password",
        "message": "password must contain at least one uppercase letter, one lowercase letter, one number, and one special character"
      }
    ]
  }
}
```


___

#### `DELETE /api/v1/auth/{GUID}`
- Requires header `Authorization: Bearer eyJhb...`

##### Example request 1:

`DELETE http://localhost:8080/api/v1/auth/898be767-f66f-494d-be9a-c1be85548bb7`

Headers
```
Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2MzkzMDUsImlhdCI6MTczMzYzOTAwNSwiaXAiOiIxNzIuMTguMC4xIiwianRpIjoiNjQyNThmZjgtYmFmNC00OTQ4LTliNWMtNTVjNTExOWEyZjI3Iiwic3ViIjoiODk4YmU3NjctZjY2Zi00OTRkLWJlOWEtYzFiZTg1NTQ4YmI3In0.t3GeawGNFsRBiPRzD6OU8ZQAWm66ZgO-kNXdMeistDnng3nAohw1qV_Gmdtjj8Wb4Z_8vgRzy2FdvVBz48NKYg

```

Example response (204):
```json
(empty)
```



##### Example request 2:

`DELETE http://localhost:8080/api/v1/auth/898be767-f66f-494d-be9a-c1be85548bb7`

Headers
```
(empty)
```

Example response:
```json
{
  "code": 403,
  "message": "error verifying Authorization header: token is malformed: token contains an invalid number of segments"
}
```


##### Example request 3:

`DELETE http://localhost:8080/api/v1/auth/898be767-f66f-494d-be9a-c1be85548bb7`


Headers (expired token)
```
Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2MzkzMDUsImlhdCI6MTczMzYzOTAwNSwiaXAiOiIxNzIuMTguMC4xIiwianRpIjoiNjQyNThmZjgtYmFmNC00OTQ4LTliNWMtNTVjNTExOWEyZjI3Iiwic3ViIjoiODk4YmU3NjctZjY2Zi00OTRkLWJlOWEtYzFiZTg1NTQ4YmI3In0.t3GeawGNFsRBiPRzD6OU8ZQAWm66ZgO-kNXdMeistDnng3nAohw1qV_Gmdtjj8Wb4Z_8vgRzy2FdvVBz48NKYg
```

Example response:
```json
{
  "code": 403,
  "message": "error verifying Authorization header: token has invalid claims: token is expired"
}
```

##### Example request 4:

`DELETE http://localhost:8080/api/v1/auth/00000000-f66f-494d-be9a-c1be85548bb7`

Headers (expired token)
```
Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2MzkzMDUsImlhdCI6MTczMzYzOTAwNSwiaXAiOiIxNzIuMTguMC4xIiwianRpIjoiNjQyNThmZjgtYmFmNC00OTQ4LTliNWMtNTVjNTExOWEyZjI3Iiwic3ViIjoiODk4YmU3NjctZjY2Zi00OTRkLWJlOWEtYzFiZTg1NTQ4YmI3In0.t3GeawGNFsRBiPRzD6OU8ZQAWm66ZgO-kNXdMeistDnng3nAohw1qV_Gmdtjj8Wb4Z_8vgRzy2FdvVBz48NKYg
```

Example response:
```json
{
  "code": 404,
  "message": "error deleting user: user not found"
}
```
