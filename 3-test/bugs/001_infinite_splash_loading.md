# Баг: Бесконечная загрузка на сплэш-экране

## Симптом
Приложение открывается, показывает сплэш-экран с логотипом «Глини» и крутящимся `CircularProgressIndicator`. Загрузка не заканчивается — экран висит бесконечно, независимо от того, авторизован пользователь или нет.

## Цель
После завершения проверки сессии (`checkSession`) приложение должно:
- перейти на `/schedule`, если пользователь авторизован;
- перейти на `/login`, если не авторизован.

## Расследование

### Промт: исследовать код стартапа мобильного приложения
Найди точку входа, сплэш, роутер, auth-провайдер, API-клиент — всё, что участвует в startup'е. Верни пути файлов, номера строк и содержимое ключевых файлов.

### Ответ:
- `main.dart` → `GliniApp` → `MaterialApp.router`
- `app.dart` — `ConnectivityBanner` оборачивает child
- `router.dart` — GoRouter + SplashScreen, **redirect guard — корень бага**
- `auth_provider.dart` — `AuthNotifier`, начальное состояние `AsyncValue.data(null)`
- `api_client.dart` — `apiHostOverride = 'http://192.168.1.101:8080'`
- **Первоисточник:** redirect guard на строках 22–31 не выпускает с `/splash`:
  - Для авторизованных: условие `state.matchedLocation != '/splash'` блокирует редирект
  - Для неавторизованных: `isAuthRoute=true` блокирует редирект на `/login`

### Промт: исправить баг
Прочитай router.dart и auth_provider.dart, исправь бесконечную загрузку на сплэше.

### Ответ:
**Две правки:**

1. **`auth_provider.dart:13`** — изменить начальное состояние с `AsyncValue.data(null)` на `AsyncValue.loading()`, чтобы redirect guard мог отличить «ещё проверяем» от «не авторизован»:
   - Было: `: super(const AsyncValue.data(null));`
   - Стало: `: super(const AsyncValue.loading());`

2. **`router.dart:21-31`** — переписать redirect guard:
   - Если `AsyncLoading` → `return null` (ждём checkSession)
   - Если на `/splash` и авторизован → `/schedule`
   - Если на `/splash` и не авторизован → `/login`
   - Стандартные проверки для защиты остальных роутов

## Файлы изменений
- `mobile/lib/features/auth/providers/auth_provider.dart` — строка 13
- `mobile/lib/app/router.dart` — строки 21–33
