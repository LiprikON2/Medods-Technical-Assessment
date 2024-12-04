# Medods - Technical Assessment

### Running

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
- https://www.youtube.com/watch?v=8uiZC0l4Ajw
- https://www.reddit.com/r/golang/comments/1310xxl/comment/jhymmry/
- https://www.gobeyond.dev/standard-package-layout/
- https://www.reddit.com/r/golang/comments/wbawx5/comment/ii5m2ox/
- https://threedots.tech/post/ddd-lite-in-go-introduction/

