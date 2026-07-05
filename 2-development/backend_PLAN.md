# Backend Implementation Checklist

## Архитектура

```
Handler (HTTP) → Service (Business Logic) → Repository (SQLite via sqlx)
```

- **Handler**: парсинг запроса, валидация, вызов service, JSON-ответ
- **Service**: бизнес-правила (окно отмены 4ч, окно оценки 1–48ч, проверка мест)
- **Repository**: sqlx-запросы к SQLite, маппинг строк в struct

**Точка входа**: `cmd/server/main.go` — config → DB (sqlx + goose) → repository → service → handler → router → listen

**Токен**: JWT (`golang-jwt/jwt/v5`), claims: `client_id`, `exp`, `iat`. Секрет из `AUTH_SECRET` env.

---

## Стек

| Компонент | Выбор |
|---|---|
| HTTP | `chi` |
| DB | `modernc.org/sqlite` (pure Go, без CGO) |
| Query | `sqlx` |
| Миграции | `pressly/goose` (SQL-файлы, встраивание `embed`) |
| JWT | `golang-jwt/jwt/v5` |
| Валидация | `go-playground/validator` |
| Логи | `log/slog` (stdlib) |
| Тесты | `testing` + `testify` + `httptest` |

---

## Итерация 0 — Скелет

- [ ] `internal/config/config.go` — чтение env: `AUTH_SECRET`, `DB_PATH`, `ADDR`
- [ ] `internal/domain/client.go` — `Client` struct + `ClientRepository` interface
- [ ] `internal/domain/master.go` — `Master` struct + `MasterRepository` interface
- [ ] `internal/domain/program.go` — `Program` struct + `ProgramRepository` interface
- [ ] `internal/domain/slot.go` — `Slot` struct, `SlotFilter` + `SlotRepository` interface
- [ ] `internal/domain/booking.go` — `Booking` struct + `BookingRepository` interface
- [ ] `internal/domain/rating.go` — `Rating` struct + `RatingRepository` interface
- [ ] `internal/repository/db.go` — подключение SQLite, PRAGMA foreign_keys, запуск goose миграций через `embed`
- [ ] `internal/middleware/auth.go` — JWT middleware: Bearer → parse → validate → inject `client_id` in context
- [ ] `internal/router.go` — chi router, группировка routes, подключение middleware
- [ ] `cmd/server/main.go` — bootstrap: config → DB → repository → service → handler → router → listen
- [ ] `Makefile` — цели `build`, `run`, `migrate-up`, `migrate-down`, `test`
- [ ] **Verify**: `go build ./...` проходит, сервер стартует и отвечает 404 на неизвестный route

---

## Итерация 1 — Auth (FT-01, FT-02)

### Repository
- [ ] `internal/repository/client.go: InsertClient(client) (Client, error)` — INSERT + RETURNING id
- [ ] `internal/repository/client.go: GetClientByLogin(login) (Client, error)` — SELECT by login

### Service
- [ ] `internal/service/auth.go: Register(login, password)` — валидация (login не пуст, password ≥ 6), bcrypt → hash, InsertClient, 409 на duplicate
- [ ] `internal/service/auth.go: Login(login, password)` — GetClientByLogin, bcrypt.Compare, JWT (sub=client_id, exp=24h)

### Handler
- [ ] `internal/handler/auth.go: Register` — POST /auth/register → 201 + Client | 400 | 409
- [ ] `internal/handler/auth.go: Login` — POST /auth/login → 200 + {token, client} | 401

### Tests
- [ ] Repository: InsertClient + GetClientByLogin (in-memory SQLite)
- [ ] Repository: duplicate login → error
- [ ] Service: Register — success case
- [ ] Service: Register — duplicate login → error
- [ ] Service: Login — wrong password → 401
- [ ] Handler: POST /auth/register — 201, 400, 409 (httptest)
- [ ] Handler: POST /auth/login — 200, 401 (httptest)
- [ ] **Verify**: `curl POST /auth/register + /auth/login` работает end-to-end

---

## Итерация 2 — Schedule (FT-03, FT-09, FT-10)

### Seed
- [ ] `internal/repository/seed.go` — заполнить `masters`, `programs`, `masters_programs`, `slots` тестовыми данными (5–10 слотов на ближайшие дни)

### Repository
- [ ] `internal/repository/slot.go: QuerySlots(filter) ([]Slot, error)` — SELECT с JOIN program + master, фильтры dateFrom/dateTo/masterId/programId
- [ ] `internal/repository/slot.go: GetSlotByID(id) (Slot, error)` — SELECT с JOIN
- [ ] `internal/repository/master.go: QueryMasters() ([]Master, error)` — SELECT + список programIds
- [ ] `internal/repository/master.go: GetMasterByID(id) (Master, error)` — SELECT + programIds
- [ ] `internal/repository/program.go: QueryPrograms() ([]Program, error)` — SELECT + список masterIds
- [ ] `internal/repository/program.go: GetProgramByID(id) (Program, error)` — SELECT + masterIds

### Handler
- [ ] `internal/handler/slot.go: ListSlots` — GET /slots?dateFrom=&dateTo=&masterId=&programId= → 200 + []Slot
- [ ] `internal/handler/slot.go: GetSlot` — GET /slots/{id} → 200 + Slot | 404
- [ ] `internal/handler/master.go: List` — GET /masters → 200 + []Master
- [ ] `internal/handler/master.go: Get` — GET /masters/{id} → 200 + Master | 404
- [ ] `internal/handler/program.go: List` — GET /programs → 200 + []Program
- [ ] `internal/handler/program.go: Get` — GET /programs/{id} → 200 + Program | 404

### Tests
- [ ] Repository: QuerySlots с разными фильтрами
- [ ] Repository: GetSlotByID — found + not found
- [ ] Handler: GET /slots — 200 с пагинацией по датам
- [ ] Handler: GET /slots — фильтр masterId/programId
- [ ] Handler: GET /masters, /masters/{id} — 200
- [ ] Handler: GET /programs, /programs/{id} — 200
- [ ] **Verify**: полный цикл register → login → GET /slots с токеном

---

## Итерация 3 — Booking (FT-04, FT-05, FT-06)

### Repository
- [ ] `internal/repository/booking.go: InsertBooking(clientID, slotID, rentalSelected) (Booking, error)` — в составе транзакции
- [ ] `internal/repository/booking.go: QueryBookingsByClient(clientID, statusFilter) ([]Booking, error)` — SELECT с JOIN slot + program + master
- [ ] `internal/repository/booking.go: GetBookingByID(id) (Booking, error)` — SELECT с JOIN
- [ ] `internal/repository/booking.go: UpdateBookingStatus(id, status, reason) error`
- [ ] `internal/repository/slot.go: DecrementSpots(tx, slotID) error` — UPDATE available_spots - 1 WHERE available_spots > 0
- [ ] `internal/repository/slot.go: IncrementSpots(tx, slotID) error` — UPDATE available_spots + 1

### Service
- [ ] `internal/service/booking.go: CreateBooking(clientID, slotID, rentalSelected)` — транзакция: проверка мест + нет активной брони → DecrementSpots → InsertBooking → COMMIT. 409 при конфликте
- [ ] `internal/service/booking.go: CancelBooking(bookingID, clientID)` — проверка принадлежности, проверка `slot.date_time - now() ≥ 4h` (иначе 422), транзакция: UpdateBookingStatus → IncrementSpots
- [ ] `internal/service/booking.go: GetMyBookings(clientID, statusFilter)` — QueryBookingsByClient

### Handler
- [ ] `internal/handler/booking.go: Create` — POST /bookings (auth) → 201 + Booking | 400 | 409 | 401
- [ ] `internal/handler/booking.go: ListMy` — GET /bookings (auth) → 200 + []Booking | 401
- [ ] `internal/handler/booking.go: Get` — GET /bookings/{id} (auth) → 200 + Booking | 404 | 401
- [ ] `internal/handler/booking.go: Cancel` — PATCH /bookings/{id}/cancel (auth) → 200 | 422 | 404 | 401

### Tests
- [ ] Repository: InsertBooking — success
- [ ] Repository: DecrementSpots при available=0 → error
- [ ] Service: CreateBooking — полный success flow
- [ ] Service: CreateBooking — нет мест → 409
- [ ] Service: CreateBooking — двойная бронь → 409
- [ ] Service: CancelBooking — ≥4h → success
- [ ] Service: CancelBooking — <4h → 422
- [ ] Service: CancelBooking — чужой bookingID → error
- [ ] Handler: POST /bookings — 201, 409, 401 (httptest)
- [ ] Handler: PATCH /bookings/{id}/cancel — 200, 422 (httptest)
- [ ] **Verify**: register → login → book slot → GET /bookings → cancel → GET /bookings (status changed)

---

## Итерация 4 — Rating (FT-07)

### Repository
- [ ] `internal/repository/rating.go: InsertRating(clientID, masterID, slotID, score) (Rating, error)`
- [ ] `internal/repository/rating.go: GetRatingByClientAndSlot(clientID, slotID) (Rating, error)` — для проверки 409 на дубликат

### Service
- [ ] `internal/service/rating.go: CreateRating(clientID, masterID, slotID, score)` — проверка: есть активная бронь на slot → проверка окна `end_time + 1h ≤ now ≤ end_time + 48h` → 422 если вне окна → InsertRating → 409 на дубликат

### Handler
- [ ] `internal/handler/rating.go: Create` — POST /ratings (auth) → 201 + Rating | 400 | 401 | 422 | 409

### Tests
- [ ] Service: CreateRating — success
- [ ] Service: CreateRating — вне окна (рано/поздно) → 422
- [ ] Service: CreateRating — нет активной брони → 422
- [ ] Service: CreateRating — повторная оценка → 409
- [ ] Handler: POST /ratings — 201, 422, 409 (httptest)
- [ ] **Verify**: полный E2E: register → login → book → проходим время (mock) → rate → проверяем рейтинг мастера

---

## Итерация 5 — Прочее и полировка

- [ ] `internal/middleware/logger.go` — request logging: method, path, duration, status
- [ ] CORS middleware (chi built-in) — разрешить origin для мобильного клиента
- [ ] Error handling: единый формат `{"error": "message"}` для всех 4xx/5xx
- [ ] Panic recovery middleware
- [ ] Graceful shutdown (`signal.NotifyContext`)
- [ ] HTTP-таймауты (ReadTimeout, WriteTimeout, IdleTimeout)
- [ ] **Verify**: полный интеграционный тест register → login → get slots → book → get bookings → cancel → book again → rate (сквозной через httptest)

---

## Реализация миграций (файлы уже созданы)

- [x] `migrations/001_create_clients.sql`
- [x] `migrations/002_create_masters.sql`
- [x] `migrations/003_create_programs.sql`
- [x] `migrations/004_create_masters_programs.sql`
- [x] `migrations/005_create_slots.sql`
- [x] `migrations/006_create_bookings.sql`
- [x] `migrations/007_create_ratings.sql`

---

## Структура файлов после всех итераций

```
cmd/server/main.go
internal/
├── config/config.go
├── domain/
│   ├── client.go
│   ├── master.go
│   ├── program.go
│   ├── slot.go
│   ├── booking.go
│   └── rating.go
├── repository/
│   ├── db.go
│   ├── seed.go
│   ├── client.go
│   ├── master.go
│   ├── program.go
│   ├── slot.go
│   ├── booking.go
│   └── rating.go
├── service/
│   ├── auth.go
│   ├── booking.go
│   └── rating.go
├── handler/
│   ├── auth.go
│   ├── slot.go
│   ├── master.go
│   ├── program.go
│   ├── booking.go
│   └── rating.go
├── middleware/
│   ├── auth.go
│   └── logger.go
└── router.go
migrations/
├── 001_create_clients.sql
├── ...
└── 007_create_ratings.sql
go.mod
Makefile
```
