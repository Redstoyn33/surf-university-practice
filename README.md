# glini — Гончарная мастерская: запись онлайн

Монолитное Go-приложение (бэкенд) + Flutter-клиент для записи в гончарную мастерскую.  
Клиенты регистрируются, просматривают расписание мастеров и программ, бронируют слоты, отменяют брони (с бизнес-правилами) и оставляют оценки.

## Стек

| Слой | Технологии |
|---|---|
| **Бэкенд** | Go 1.22+, chi, sqlx, modernc.org/sqlite (pure Go, без CGO), goose, golang-jwt/jwt/v5, bcrypt, log/slog |
| **Мобильный клиент** | Flutter 3.41+, Dart 3.11+, flutter_riverpod, go_router, Dio, flutter_secure_storage, Material 3 |
| **Тесты (бэкенд)** | testing + testify + httptest (in-memory SQLite) |

## API (12 endpoints)

| Метод | Путь | Описание |
|---|---|---|
| POST | `/auth/register` | Регистрация клиента |
| POST | `/auth/login` | Вход, получение JWT |
| GET | `/slots` | Список слотов |
| GET | `/slots/:id` | Детали слота |
| GET | `/bookings` | Брони клиента |
| POST | `/bookings` | Создать бронь |
| GET | `/bookings/:id` | Детали брони |
| PATCH | `/bookings/:id/cancel` | Отменить бронь (≥4ч до начала) |
| POST | `/ratings` | Оценить мастера (1ч–48ч после слота) |
| GET | `/masters` | Список мастеров |
| GET | `/masters/:id` | Детали мастера |
| GET | `/programs` | Список программ |
| GET | `/programs/:id` | Детали программы |

## Структура

```
1-analysis/          # Спеки и OpenAPI (не закоммичено)
2-development/       # Планы реализации (копии)
backend/             # Go-сервер
  cmd/server/main.go # Точка входа
  internal/
    domain/          # Интерфейсы и модели
    handler/         # HTTP-обработчики (chi)
    service/         # Бизнес-логика
    repository/      # SQLite (sqlx) + seed
    middleware/      # JWT-аутентификация
    config/          # Конфигурация из env
  migrations/        # 7 goose-миграций
  Makefile
mobile/              # Flutter-клиент
  lib/
    app/             # Маршрутизация (GoRouter)
    core/            # API-клиент (Dio), тема, хранилище
    features/        # Экран → Notifier → Repository
      auth/          # Логин/регистрация
      schedule/      # Расписание, слоты, мастера, программы
      booking/       # Брони, детали, отмена
      ratings/       # Оценки мастеров
```

## Быстрый старт

### Бэкенд

```bash
cd backend
make run  # миграции + сид-данные (3 мастера, слоты на неделю)
```

Тестовый логин: `test@test.com` / `123456`

### Мобильное приложение

```bash
cd mobile
flutter run -d '<DEVICE>'
# API URL настраивается в lib/core/api_client.dart
```

## Бизнес-правила

- **Отмена брони** — только ≥4ч до начала, иначе 422
- **Оценка мастера** — окно 1ч–48ч после окончания слота
- **Двойная бронь** — 409 Conflict (unique partial index)
- **Клиент записывает только себя** (client_id из JWT)
