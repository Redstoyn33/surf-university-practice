# Plan реализации бэкенда

## Архитектура

```
Handler (HTTP) → Service (Business Logic) → Repository (SQLite via sqlx)
```

- **Handler**: приём запросов, парсинг, валидация, вызов service, возврат JSON
- **Service**: бизнес-правила (окно 4ч для отмены, 1–48ч для оценки, проверка мест)
- **Repository**: sqlx-запросы к SQLite, маппинг строк в domain-структуры

**Точка входа**: `cmd/server/main.go` — инициализация БД, миграции, DI, запуск `http.Server`

---

## Phase 1 — Critical (Auth + Schedule)

### 1.1 Seed-данные
Файл `internal/repository/seed.go` — заполняет таблицы `masters`, `programs`, `masters_programs`, `slots` тестовыми данными.

### 1.2 Auth — FT-01, FT-02
| Endpoint | Handler | Service | Repository |
|---|---|---|---|
| `POST /auth/register` | `handler/auth.go: Register` | валидация логина/пароля, bcrypt, создание клиента | `InsertClient` |
| `POST /auth/login` | `handler/auth.go: Login` | проверка пароля, генерация JWT | `GetClientByLogin` |

*JWT* — `golang-jwt/jwt/v5`, claims: `client_id`, `exp`, `iat`. Секрет из `AUTH_SECRET` env.

### 1.3 Schedule — FT-03
| Endpoint | Handler | Service | Repository |
|---|---|---|---|
| `GET /slots` | `handler/slot.go: ListSlots` | фильтры dateFrom/dateTo/masterId/programId | `QuerySlots` |
| `GET /slots/{id}` | `handler/slot.go: GetSlot` | — | `GetSlotByID` |

*Ответ включает вложенные `program` и `master`* (JOIN-ы).

### 1.4 Masters & Programs (read-only) — FT-09, FT-10
| Endpoint | Handler | Service | Repository |
|---|---|---|---|
| `GET /masters` | `handler/master.go: List` | — | `QueryMasters` |
| `GET /masters/{id}` | `handler/master.go: Get` | — | `GetMasterByID` |
| `GET /programs` | `handler/program.go: List` | — | `QueryPrograms` |
| `GET /programs/{id}` | `handler/program.go: Get` | — | `GetProgramByID` |

*Ответы включают списки ID программ/мастеров (из masters_programs).*

---

## Phase 2 — Critical (Booking)

### 2.1 Booking CRUD — FT-04, FT-05, FT-06

| Endpoint | Handler | Service | Repository |
|---|---|---|---|
| `POST /bookings` | `handler/booking.go: Create` | проверить `available_spots > 0`, атомарный декремент + вставка (транзакция), проверка уникальности активной брони | `InsertBooking`, `DecrementSpots` |
| `GET /bookings` | `handler/booking.go: ListMy` | фильтр по status | `QueryBookingsByClient` |
| `GET /bookings/{id}` | `handler/booking.go: Get` | проверка: бронь принадлежит текущему клиенту | `GetBookingByID` |
| `PATCH /bookings/{id}/cancel` | `handler/booking.go: Cancel` | проверить `≥4ч до date_time`, иначе 422. Статус → `отменена клиентом`, `available_spots++` (транзакция) | `UpdateBookingStatus`, `IncrementSpots` |

### 2.2 Booking бизнес-правила

- **Транзакция в Create**: `BEGIN → SELECT available_spots FOR UPDATE → проверка >0 → INSERT booking → UPDATE available_spots-1 → COMMIT`
- **Отмена ≥4ч**: в service вычисляем разницу `slot.date_time - now()`, если `< 4h → return 422`
- **Атомарность**: SQLite с `BEGIN IMMEDIATE` — гарантия от двойных бронь на уровне БД (уникальный индекс `idx_active_booking` — вторая линия защиты)

---

## Phase 3 — Medium (Rating)

### 3.1 Rating — FT-07

| Endpoint | Handler | Service | Repository |
|---|---|---|---|
| `POST /ratings` | `handler/rating.go: Create` | проверить окно 1–48ч после end_time, проверить активную бронь клиента на этот слот, 409 если уже оценено | `InsertRating` |

*Проверка окна*: `now() BETWEEN slot.end_time + 1h AND slot.end_time + 48h`.

---

## Cross-cutting

| Слой | Файл | Что делает |
|---|---|---|
| middleware/auth.go | JWT-проверка | Извлекает Bearer-токен, валидирует, кладёт `client_id` в контекст |
| middleware/logger.go | request logging | slog + duration, request ID |
| internal/config/config.go | env → struct | `AUTH_SECRET`, `DB_PATH`, `ADDR` |
| internal/domain/*.go | entity structs + interfaces | domain-модели + контракты репозиториев |
| cmd/server/main.go | bootstrap | config → DB (sqlx + goose) → repository → service → handler → router → listen |

---

## Стек

| Компонент | Выбор |
|---|---|
| HTTP | `chi` — stdlib-совместимый роутер, middleware chaining |
| DB | `modernc.org/sqlite` — pure Go SQLite (без CGO) |
| Query | `sqlx` — маппинг rows → struct, именованные параметры |
| Миграции | `pressly/goose` — SQL-файлы с `-- +goose Up/Down`, встраивание через `embed` |
| JWT | `golang-jwt/jwt/v5` |
| Валидация | `go-playground/validator` — теги в request DTO |
| Логи | `log/slog` (stdlib) |
| Тесты | `testing` + `testify` + `httptest` |

---

## Структура реализуемых файлов

```
internal/
├── domain/
│   ├── client.go        # Client struct + ClientRepository interface
│   ├── master.go        # Master struct + MasterRepository interface
│   ├── program.go       # Program struct + ProgramRepository interface
│   ├── slot.go          # Slot struct, SlotFilter + SlotRepository interface
│   ├── booking.go       # Booking struct + BookingRepository interface
│   └── rating.go        # Rating struct + RatingRepository interface
├── handler/
│   ├── auth.go          # Register, Login
│   ├── slot.go          # ListSlots, GetSlot
│   ├── booking.go       # Create, ListMy, Get, Cancel
│   ├── master.go        # List, Get
│   ├── program.go       # List, Get
│   └── rating.go        # Create
├── service/
│   ├── auth.go          # Register, Login
│   ├── booking.go       # CreateBooking, CancelBooking
│   └── rating.go        # CreateRating
├── repository/
│   ├── client.go        # sqlx implementation
│   ├── booking.go       # sqlx implementation
│   ├── slot.go          # sqlx implementation
│   ├── master.go        # sqlx implementation
│   ├── program.go       # sqlx implementation
│   ├── rating.go        # sqlx implementation
│   └── db.go            # goose migrations runner, PRAGMA foreign_keys
├── middleware/
│   └── auth.go          # JWT middleware
├── config/
│   └── config.go        # env → struct
└── router.go            # chi router setup with all routes
```

---

## Бизнес-правила (алгоритмы)

### Register
1. Валидация: login не пуст, password ≥ 6 символов
2. bcrypt(password) → hash
3. INSERT clients → если `UNIQUE constraint` → 409
4. 201 + Client

### Login
1. SELECT client by login → не найден → 401
2. bcrypt.CompareHashAndPassword → не совпал → 401
3. JWT с `sub=client_id` + `exp=24h`
4. 200 + `{token, client}`

### CreateBooking (транзакция)
```
BEGIN IMMEDIATE
  SELECT available_spots FROM slots WHERE id=? → 0? → 409
  SELECT id FROM bookings WHERE client_id=? AND slot_id=? AND status='активна' → found? → 409
  INSERT INTO bookings (client_id, slot_id, rental_selected) VALUES (?,?,?)
  UPDATE slots SET available_spots = available_spots - 1 WHERE id=?
COMMIT
```

### CancelBooking (транзакция)
```
SELECT s.date_time, s.end_time FROM bookings b JOIN slots s ON b.slot_id=s.id WHERE b.id=?
→ now + 4h > date_time? → 422
BEGIN IMMEDIATE
  UPDATE bookings SET status='отменена клиентом' WHERE id=?
  UPDATE slots SET available_spots = available_spots + 1 WHERE id=?
COMMIT
```

### CreateRating
```
SELECT s.end_time FROM bookings b JOIN slots s ON b.slot_id=s.id
  WHERE b.client_id=? AND b.slot_id=? AND b.status='активна'
→ no active booking? → 422

end_time + 1h ≤ now ≤ end_time + 48h? → нет? → 422

INSERT INTO ratings → UNIQUE(client_id, slot_id)? → 409
```

---

## Тесты

- **Repository**: in-memory SQLite (`:memory:`) с goose UP, тесты транзакций
- **Handler**: httptest.Server + mock service (через интерфейсы domain)
- **Service**: юнит-тесты бизнес-правил (окно отмены, окно оценки)
- **Интеграционные**: полный flow (register → login → book → cancel → rate)
