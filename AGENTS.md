# AGENTS.md

## Проект

Монолитное Go-приложение (`backend/`) + Flutter-клиент (`mobile/`). Проект только начат — есть скелет и планы, реализации кода нет.

**Источник истины для требований:** `1-analysis/` (Markdown-спеки на русском). Там же OpenAPI-спецификация (`1-analysis/api/openapi.yaml`).

**Язык проекта:** русский (UI, ошибки, статусы, названия сущностей).

## Структура

```
1-analysis/         # Документация (не закоммичена)
backend/            # Go backend
  cmd/server/main.go
  internal/         # domain/ handler/ service/ repository/ middleware/ config/
  migrations/       # 7 goose-миграций (SQLite)
  PLAN.md           # Чек-лист реализации
  Makefile
mobile/             # Flutter клиент
  lib/              # app/ core/ features/ (пусто)
  pubspec.yaml
  PLAN.md           # Чек-лист реализации
```

## Backend (Go)

- **Go 1.22**, module `github.com/glini/backend`
- **Стек (план, в go.mod пока нет зависимостей):** chi + sqlx + modernc.org/sqlite (pure Go, без CGO) + goose + golang-jwt/v5 + go-playground/validator + log/slog
- **Архитектура:** Handler → Service → Repository, слои через интерфейсы из `internal/domain/`
- **Точка входа:** `cmd/server/main.go`
- **Команды:**
  - `make build` — `go build -o bin/server ./cmd/server`
  - `make run` — `go run ./cmd/server`
  - `make migrate-up` — `goose -dir migrations sqlite3 ./data/glini.db up`
  - `make migrate-down` — `goose -dir migrations sqlite3 ./data/glini.db down`
- **Тесты:** in-memory SQLite (`:memory:`), `testing` + `testify` + `httptest`
- **Миграции:** 7 файлов (001–007), FK-зависимости соблюдены. После подключения к БД обязательно `PRAGMA foreign_keys = ON`
- **Бизнес-правила (ключевые):**
  - Отмена брони — только ≥4ч до начала, иначе 422
  - Оценка мастера — окно от 1ч до 48ч после окончания слота
  - Клиент записывает только себя
  - Двойная бронь — 409 (уникальный partial index `idx_active_booking WHERE status = 'активна'`)
- **Endpoints (12):** `/auth/register`, `/auth/login`, `/slots` (GET list + GET by id), `/bookings` (GET list + POST + GET by id + PATCH cancel), `/ratings` (POST), `/masters` (GET list + GET by id), `/programs` (GET list + GET by id)
- **JWT:** claims `client_id`, `exp=24h`, секрет из env `AUTH_SECRET`

## Mobile (Flutter)

- **Flutter SDK:** `>=3.4.0 <4.0.0`
- **Стек (план, часть зависимостей ещё не добавлена):** flutter_riverpod + go_router + dio + flutter_secure_storage + freezed + json_serializable
- **Архитектура:** Screen → Provider/Notifier → Repository (Dio)
- **Маршрутизация (GoRouter):** `/splash` → `/login` / `/register` → TabBar (`/schedule` ⟷ `/my-bookings`) с push-экранами `/slots/:id`, `/masters/:id`, `/programs/:id`, `/bookings/:id`, `/rate`
- **9 экранов:** SCR-001 (Login) → SCR-009 (Program) — все спецификации в `1-analysis/5-mobile-app-spec/`
- **5 бизнес-логик:** LOGIC-001 (Auth) → LOGIC-005 (Schedule) — в `1-analysis/5-mobile-app-spec/09_Логики/`

## Общее

- **Git:** 5 коммитов, `1-analysis/` не в `.gitignore`, не закоммичена (не добавлена)
- **CI/CD:** нет. **Тесты:** нет (предстоит написать). **Docker:** нет.
- **Планы реализации:** `backend/PLAN.md` (6 итераций), `mobile/PLAN.md` (7 итераций)
- **Seed-данные:** `internal/repository/seed.go` — masters, programs, slots (не реализован)
