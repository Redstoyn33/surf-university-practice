# Mobile (Flutter) Implementation Checklist

## Архитектура

```
Screen (UI) → Provider/Notifier (State) → Repository (API via Dio)
```

- **Screen**: Widget-дерево экрана, состояния `AsyncValue` (loading/error/data)
- **Provider/Notifier**: Riverpod — бизнес-логика + вызов API + локальное состояние
- **Repository**: Dio-клиент, маппинг JSON ↔ модель, перехватчики JWT + ошибок

**Маршрутизация**: GoRouter — стек аутентификации (SCR-001/002) + TabBar (SCR-003 ⟷ SCR-005) + модальные детали (SCR-004, 006, 007, 008, 009)

---

## Стек

| Компонент | Выбор |
|---|---|
| State management | `flutter_riverpod` + `riverpod_annotation` |
| Routing | `go_router` |
| HTTP | `dio` (JWT interceptor, error interceptor) |
| Хранилище токена | `flutter_secure_storage` |
| Генерация | `freezed` (data classes) + `json_serializable` |
| UI | Material 3 |
| Тесты | `flutter_test` + `mocktail` |

---

## Итерация 0 — Скелет

- [ ] Инициализировать Flutter-проект (`flutter create --org com.glini glini`)
- [ ] `pubspec.yaml` — зависимости: `flutter_riverpod`, `go_router`, `dio`, `flutter_secure_storage`, `freezed`, `json_serializable`, `build_runner`
- [ ] `lib/core/theme.dart` — Material 3 theme (colors: тёплая гамма глины/керамики)
- [ ] `lib/core/api_client.dart` — Dio-клиент: base URL, JWT-перехватчик (читает токен из SecureStorage), error-перехватчик
- [ ] `lib/core/api_client.dart: AuthInterceptor` — при 401 очищает токен и редиректит на экран входа
- [ ] `lib/core/api_client.dart: ErrorInterceptor` — парсит ошибки в единый `ApiError` класс
- [ ] `lib/core/storage.dart` — обёртка над `flutter_secure_storage` (get/set/delete token + client)
- [ ] `lib/core/models/` — `freezed` data-классы: `Client`, `Master`, `Program`, `Slot`, `Booking`, `Rating` c `json_serializable`
- [ ] `lib/core/models/api_error.dart` — `ApiError` с message-ом
- [ ] `lib/app/app.dart` — `MaterialApp.router` с GoRouter
- [ ] `lib/app/router.dart` — GoRouter: splash (проверка токена), auth stack, tab bar, detail screens
- [ ] `lib/app/providers/auth_provider.dart` — `AuthNotifier`: `isAuthenticated`, `client`, `token`, методы `login()`, `register()`, `logout()`, `checkSession()`
- [ ] **Verify**: `flutter run` — видим splash, проверка токена → auth экран или главная

---

## Итерация 1 — Auth (SCR-001, SCR-002 · LOGIC-001)

### Repository
- [ ] `lib/features/auth/repository/auth_repository.dart` — `login(login, password) → {token, client}`, `register(login, password) → {token, client}`, `validateToken() → bool`

### Screens
- [ ] `lib/features/auth/screens/login_screen.dart` — SCR-001: поле логина, поле пароля (с visibility toggle), кнопка "Войти"
- [ ] `lib/features/auth/screens/registration_screen.dart` — SCR-002: поля логин/пароль/подтверждение, кнопка "Зарегистрироваться"
- [ ] `lib/features/auth/widgets/auth_form.dart` — общий каркас формы (логотип, поля, ссылки)

### Providers
- [ ] `lib/features/auth/providers/auth_provider.dart` — `AuthNotifier`: вызывает repository, сохраняет token в SecureStorage, устанавливает `isAuthenticated`, переключает route

### States & Validation
- [ ] Login: loading (лоадер на кнопке), error 401 ("Неверный логин или пароль"), error 400/5xx (снек)
- [ ] Registration: loading, error 409 ("Этот логин уже занят"), error 400 ("Пароль должен быть ≥ 6 символов")
- [ ] Валидация на клиенте: login не пуст, password ≥ 6, password == confirmPassword (для регистрации)

### UI States
- [ ] Loading: `CircularProgressIndicator` на кнопке
- [ ] Error: `SnackBar` с текстом ошибки
- [ ] Success: переход на главный экран

### Tests
- [ ] AuthRepository: login/register — успех и ошибки (mock Dio)
- [ ] AuthNotifier: правильные переходы состояний
- [ ] LoginScreen: ввод данных, тап на "Войти", валидация
- [ ] RegistrationScreen: password mismatch, success flow
- [ ] **Verify**: register → auto-login → main screen

---

## Итерация 2 — Schedule + Master + Program (SCR-003, SCR-008, SCR-009 · LOGIC-005)

### Repository
- [ ] `lib/features/schedule/repository/schedule_repository.dart` — `getSlots(dateFrom, dateTo, masterId?, programId?) → List<Slot>`
- [ ] `lib/features/masters/repository/master_repository.dart` — `getMasters() → List<Master>`, `getMasterById(id) → Master`
- [ ] `lib/features/programs/repository/program_repository.dart` — `getPrograms() → List<Program>`, `getProgramById(id) → Program`

### Screens
- [ ] `lib/features/schedule/screens/schedule_screen.dart` — SCR-003: горизонтальный календарь + список слотов на выбранную дату
- [ ] `lib/features/schedule/widgets/date_bar.dart` — горизонтальный скролл дат на 7 дней, today выделен
- [ ] `lib/features/schedule/widgets/slot_card.dart` — карточка слота: время, программа, мастер, свободные места, иконка проката
- [ ] `lib/features/masters/screens/master_profile_screen.dart` — SCR-008: фото, имя, уровень, рейтинг, список программ
- [ ] `lib/features/programs/screens/program_detail_screen.dart` — SCR-009: название, описание, макс. вместимость, список мастеров

### Navigation
- [ ] TabBar: Schedule ⇔ MyBookings (GoRouter StatefulShellRoute)
- [ ] SCR-003 → SCR-008 (тап по мастеру)
- [ ] SCR-003 → SCR-009 (тап по программе)

### UI States
- [ ] Schedule: loading (скелетон списка), empty ("Нет слотов на эту дату"), error (снек + retry)
- [ ] SlotCard: full (серая, "Нет мест"), available (активная, visible)
- [ ] MasterProfile: loading, error
- [ ] ProgramDetail: loading, error
- [ ] Pull-to-refresh на расписании

### Providers
- [ ] `schedule_provider.dart` — `SlotsNotifier`: загрузка слотов по фильтрам, переключение даты, фильтры по мастеру/программе
- [ ] `master_provider.dart` — `getMasterById` future provider
- [ ] `program_provider.dart` — `getProgramById` future provider

### Tests
- [ ] SlotCard: full vs available отображение
- [ ] ScheduleScreen: загрузка и отображение списка
- [ ] DateBar: переключение дат обновляет список
- [ ] **Verify**: login → расписание → слоты видны → тап по мастеру → профиль → тап по программе → описание

---

## Итерация 3 — Slot Details + Booking (SCR-004 · LOGIC-004)

### Repository
- [ ] `lib/features/booking/repository/booking_repository.dart` — `getSlotById(id) → Slot`, `createBooking(slotId, rentalSelected) → Booking`, `getMyBookings(status?) → List<Booking>`, `getBookingById(id) → Booking`, `cancelBooking(bookingId) → Booking`

### Screen
- [ ] `lib/features/booking/screens/slot_detail_screen.dart` — SCR-004: программа, мастер (тап → SCR-008), дата/время, описание, свободные места, переключатель проката, кнопка "Записаться"
- [ ] `lib/features/booking/widgets/rental_switch.dart` — переключатель аренды инструментов

### Logic
- [ ] LOGIC-004: тап "Записаться" → проверка авторизации → диалог подтверждения → POST /bookings → 201 → переход в SCR-005 MyBookings | 409 → снек "Вы уже записаны" | 400 → снек ошибки

### Navigation
- [ ] SCR-004 → SCR-005 (успех брони)
- [ ] SCR-004 → SCR-008 (тап по мастеру)
- [ ] SCR-004 → SCR-009 (тап по программе)

### UI States
- [ ] Loading: шиммер деталей слота
- [ ] Slot unavailable (past / sold out): кнопка "Записаться" неактивна + пояснение
- [ ] Client already booked: индикатор "Вы записаны", кнопка неактивна
- [ ] Booking: loading на кнопке, error (снек), success (снек + переход)
- [ ] Confirmation dialog перед бронированием

### Provider
- [ ] `booking_provider.dart` — `BookingNotifier`: createBooking, getMyBookings, cancelBooking

### Tests
- [ ] SlotDetailScreen: отображение всех полей
- [ ] Booking flow: create → 201 → redirect
- [ ] Booking flow: 409 → снек
- [ ] **Verify**: login → slot details → book → redirected to MyBookings

---

## Итерация 4 — My Bookings + Booking Details + Cancel (SCR-005, SCR-006 · LOGIC-002)

### Screens
- [ ] `lib/features/booking/screens/my_bookings_screen.dart` — SCR-005: табы "Активные" / "Прошедшие" / "Отменённые", список броней с программой, датой, мастером, статусом
- [ ] `lib/features/booking/widgets/booking_card.dart` — карточка брони: программа, дата/время, мастер, статус, кнопка "Оценить" (для прошедших неоценённых)
- [ ] `lib/features/booking/screens/booking_detail_screen.dart` — SCR-006: программа, дата, мастер, статус, прокат, блок действий

### Cancel Logic (LOGIC-002)
- [ ] Активная бронь ≥ 4ч: кнопка "Отменить запись" → диалог подтверждения → PATCH /bookings/{id}/cancel → 200 → статус обновлён | 422 → снек "Отмена доступна за 4+ часа до начала"
- [ ] Активная бронь < 4ч: текст "Отмена менее чем за 4 часа — обратитесь в поддержку", кнопка скрыта
- [ ] Статус "отменена мастерской": показать причину отмены
- [ ] Прошедшая / отменена клиентом: никаких действий

### UI States
- [ ] MyBookings: loading, empty ("Нет записей" + кнопка "К расписанию"), error
- [ ] BookingDetail: loading, error
- [ ] Cancel: dialog → loading на кнопке → success (обновление статуса) / error (снек)

### Navigation
- [ ] TabBar иконка "Мои записи"
- [ ] SCR-005 → SCR-006 (тап по брони)
- [ ] SCR-005 → SCR-007 (тап "Оценить")
- [ ] SCR-005 → SCR-003 (empty state "К расписанию")
- [ ] SCR-006 → SCR-008 (тап по мастеру)

### Provider
- [ ] `my_bookings_provider.dart` — загрузка списка броней по фильтру статуса

### Tests
- [ ] MyBookingsScreen: табы переключают фильтр
- [ ] BookingCard: правильный статус бейджа
- [ ] Cancel flow: confirm → PATCH → status changed
- [ ] Cancel flow: <4h → 422 → снек
- [ ] **Verify**: book → see in MyBookings → cancel → status changed → reopen MyBookings

---

## Итерация 5 — Rating (SCR-007 · LOGIC-003)

### Repository
- [ ] `lib/features/ratings/repository/rating_repository.dart` — `createRating(masterId, slotId, score) → Rating`

### Screen
- [ ] `lib/features/ratings/screens/rate_master_screen.dart` — SCR-007: фото + имя мастера, программа + дата, 5 звёзд (интерактив), кнопка "Отправить оценку"

### Logic (LOGIC-003)
- [ ] Проверка окна: 1–48ч после end_time (сервер валидирует, клиент блокирует UI)
- [ ] Кнопка "Отправить" активна только при выбранной оценке
- [ ] POST /ratings → 201 → снек "Спасибо за оценку!" → возврат SCR-005 | 422 → снек "Оценка доступна в окне 1–48ч" | 409 → снек "Вы уже оценили мастера"

### UI States
- [ ] Initial: кнопка неактивна (нет оценки)
- [ ] Loading: лоадер на кнопке
- [ ] Success: снек + возврат на SCR-005
- [ ] Error: снек с ошибкой

### Navigation
- [ ] SCR-007 → SCR-005 (успех оценки)

### Provider
- [ ] `rate_master_provider.dart` — `RatingNotifier`: проверка окна, отправка

### Tests
- [ ] RateMasterScreen: выбор звёзд, активность кнопки
- [ ] CreateRating: 201 → success flow
- [ ] CreateRating: 422 → error snack
- [ ] **Verify**: book → session passes (mock time) → rate → "Спасибо!" → back to MyBookings

---

## Итерация 6 — Полировка

- [ ] Pull-to-refresh на всех списках (Schedule, MyBookings)
- [ ] Error states: кастомные виджеты с retry-кнопкой для каждого экрана
- [ ] Empty states: заглушки с понятным текстом и действием
- [ ] SnackBar: единый стиль (цвета success/error/info)
- [ ] Loading: shimmer/skeleton для списков, `CircularProgressIndicator` для кнопок
- [ ] Keyboard handling: dismiss при тапе вне поля, scroll on keyboard show
- [ ] Deep links / push notifications: placeholder-обработчики
- [ ] Splash screen: проверка токена при запуске (1–2 сек, логотип)
- [ ] Internet connectivity check: `connectivity_plus` — показать баннер "Нет соединения"
- [ ] **Verify**: полный E2E: splash → register → schedule → slot detail → book → my bookings → cancel → re-book → rate

---

## Навигация (GoRouter)

```
/splash              → SplashScreen (проверка токена → auth или main)
/login               → LoginScreen (SCR-001)
/register            → RegistrationScreen (SCR-002)
/                    → TabBar (StatefulShellRoute)
  /schedule          → ScheduleScaffold (SCR-003)
    /slots/:slotId   → SlotDetailScreen (SCR-004)  — push
    /masters/:id     → MasterProfileScreen (SCR-008) — push
    /programs/:id    → ProgramDetailScreen (SCR-009) — push
  /my-bookings       → MyBookingsScaffold (SCR-005)
    /bookings/:id    → BookingDetailScreen (SCR-006) — push
    /rate            → RateMasterScreen (SCR-007) — push
      ?masterId&slotId&masterName&programName&slotDate
```

**Guard**: redirect в `/login` при `!isAuthenticated && route требует auth`.

---

## Структура файлов после всех итераций

```
lib/
├── main.dart
├── app/
│   ├── app.dart              # MaterialApp.router
│   └── router.dart           # GoRouter + guards
├── core/
│   ├── theme.dart
│   ├── api_client.dart       # Dio + interceptors
│   ├── storage.dart          # flutter_secure_storage
│   └── models/               # freezed data classes
│       ├── client.dart
│       ├── master.dart
│       ├── program.dart
│       ├── slot.dart
│       ├── booking.dart
│       ├── rating.dart
│       └── api_error.dart
├── features/
│   ├── auth/
│   │   ├── screens/
│   │   │   ├── login_screen.dart
│   │   │   └── registration_screen.dart
│   │   ├── widgets/
│   │   │   └── auth_form.dart
│   │   ├── repository/
│   │   │   └── auth_repository.dart
│   │   └── providers/
│   │       └── auth_provider.dart
│   ├── schedule/
│   │   ├── screens/
│   │   │   └── schedule_screen.dart
│   │   ├── widgets/
│   │   │   ├── date_bar.dart
│   │   │   └── slot_card.dart
│   │   ├── repository/
│   │   │   └── schedule_repository.dart
│   │   └── providers/
│   │       └── schedule_provider.dart
│   ├── booking/
│   │   ├── screens/
│   │   │   ├── slot_detail_screen.dart
│   │   │   ├── my_bookings_screen.dart
│   │   │   └── booking_detail_screen.dart
│   │   ├── widgets/
│   │   │   ├── rental_switch.dart
│   │   │   └── booking_card.dart
│   │   ├── repository/
│   │   │   └── booking_repository.dart
│   │   └── providers/
│   │       └── booking_provider.dart
│   ├── ratings/
│   │   ├── screens/
│   │   │   └── rate_master_screen.dart
│   │   ├── repository/
│   │   │   └── rating_repository.dart
│   │   └── providers/
│   │       └── rating_provider.dart
│   ├── masters/
│   │   ├── screens/
│   │   │   └── master_profile_screen.dart
│   │   ├── repository/
│   │   │   └── master_repository.dart
│   │   └── providers/
│   │       └── master_provider.dart
│   └── programs/
│       ├── screens/
│       │   └── program_detail_screen.dart
│       ├── repository/
│       │   └── program_repository.dart
│       └── providers/
│           └── program_provider.dart
└── shared/
    ├── widgets/
    │   ├── error_state.dart
    │   ├── empty_state.dart
    │   └── loading_indicator.dart
    └── extensions/
        └── context_extensions.dart
```
