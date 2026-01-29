# Архитектура приложения

## Общая идея

Приложение построено по принципам **слоистой архитектуры** с явным разделением ответственности и строгим направлением зависимостей.

Ключевая цель архитектуры:

- изоляция бизнес-логики от инфраструктуры;
- возможность тестирования слоёв по отдельности;
- отсутствие циклических зависимостей между пакетами;
- явная сборка приложения в одном месте (composition root).

---

## Обзор слоёв

Архитектура логически делится на следующие слои (снизу вверх):

```text
Domain → Repository → Service → Transport (HTTP)

  ↑          ↑           ↑          ↑
App (composition root)
```

### Кратко по ролям

| Слой                 | Назначение                        |
|----------------------|-----------------------------------|
| **Domain**           | Бизнес-контракты и сущности       |
| **Repository**       | Доступ к хранилищам (PostgreSQL)  |
| **Service**          | Бизнес-логика и оркестрация       |
| **Transport (HTTP)** | Приём и обработка HTTP-запросов   |
| **App**              | Сборка и инициализация приложения |

### Направление зависимостей (ключевое правило)

```text
domain     ← repository
domain     ← service
service    ← http
http       ← app
repository ← app
service    ← app
```

❗ **Ни один нижний слой не знает о верхнем.**

## Domain layer

**Пакет:** `internal/domain`

### Назначение Domain слоя

Domain-слой содержит:

- бизнес-интерфейсы;
- бизнес-типы (DTO / entity);
- контракты, не зависящие от инфраструктуры.

Этот слой **не знает**:

- о БД;
- о HTTP;
- о логгере;
- о конфигурации.

### Пример Domain слоя

```go
// internal/domain/user/person.go
type PersonRepository interface {
    Create(ctx context.Context) error
}
```

Domain задаёт **что должно быть сделано**, но не **как**.

## Repository layer

**Пакет:** `internal/repository/postgres`

### Назначение Repository слоя

Repository-слой:

- реализует интерфейсы из domain;
- инкапсулирует SQL и работу с БД;
- не содержит бизнес-логики.

### Зависимости Repository слоя

```text
repository → domain
```

Этот слой **знает**:

- о database/sql;
- о PostgreSQL;
- о схемах и таблицах.

Этот слой **не знает**:

- о сервисах;
- о HTTP;
- о приложении целиком

### Пример Repository слоя

```go
type PersonRepository struct {
    db *sql.DB
}

func NewPersonRepository(db *sql.DB) *PersonRepository {
    return &PersonRepository{db: db}
}
```

## Service layer

**Пакет:** `internal/service`

### Назначение Service слоя

Service-слой содержит:

- бизнес-логику;
- оркестрацию нескольких репозиториев;
- правила и сценарии (use cases).

Service — это **центр бизнес-логики приложения**.

### Зависимости Service слоя

```text
service → domain
```

Service зависит **только от интерфейсов**, а не от реализаций.

### Пример

```go
type personService struct {
    personRepo         user.PersonRepository
    userCredentialRepo user.UserCredentialRepository
}
```

Service:

- не знает, какая БД используется;
- не знает, через какой транспорт пришёл запрос.

## Transport layer (HTTP)

**Пакет:** `internal/transport/http`

### Назначение Transport слоя

HTTP-слой:

- принимает HTTP-запросы;
- валидирует входные данные;
- вызывает сервисы;
- формирует HTTP-ответы.

Этот слой не содержит бизнес-логики.

### Структура Transport слоя

```go
transport/http/
  ├── handler/        // HTTP handlers
  ├── middleware/     // middleware (logging, recover, auth)
  ├── response/       // HTTP response helpers
  ├── deps.go         // зависимости HTTP слоя
```

### Зависимости Transport слоя

```text
http → service
```

HTTP-слой знает **только те сервисы, которые ему нужны**, через интерфейсы.

## Deps (HTTP dependencies)

**Файл:** `internal/transport/http/deps.go`

```go
type Deps struct {
   Person service.PersonService
   //...
}
```

### Зачем нужен `Deps`

Deps — это контракт зависимостей HTTP-слоя.

Он:

- описывает, какие сервисы нужны HTTP;
- не знает, как они создаются;
- не зависит от app.

⚠️ HTTP-слой **не знает** про агрегат `Services` и не должен его знать.

## App layer (composition root)

**Пакет:** `internal/app`

### Назначение App слоя

app — это **composition root**:

- инициализация конфигурации;
- подключение к БД;
- применение миграций;
- сборка репозиториев;
- сборка сервисов;
- запуск HTTP-сервера.

Это **единственное место**, где разрешено:

- знать обо всех слоях;
- связывать их между собой.

### Пример инициализации App слоя

```go
a.repos = initRepositories(db)
a.services = initServices(a.repos)
a.httpServer = initHTTPServer(a.cfg, a.logger, a.services)
```
