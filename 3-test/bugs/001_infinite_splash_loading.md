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
- `app.dart` — использует `ref.read(routerProvider)` (не watch!), поэтому GoRouter создаётся **один раз** и не пересоздаётся при изменении `authNotifierProvider`
- `router.dart` — GoRouter + SplashScreen, **redirect guard — корень бага**
- `auth_provider.dart` — `AuthNotifier`, начальное состояние `AsyncValue.data(null)`
- **Первоисточник:** redirect guard на строках 22–31 не выпускает с `/splash`:
  - Для авторизованных: условие `state.matchedLocation != '/splash'` блокирует редирект
  - Для неавторизованных: `isAuthRoute=true` блокирует редирект на `/login`
- **Дополнительно:** SplashScreen вызывает `checkSession()`, но не навигирует после его завершения — полагается на redirect guard, который не срабатывает

### Промт: исправить баг (первая попытка)
Прочитай router.dart и auth_provider.dart, исправь бесконечную загрузку на сплэше.

### Ответ (провалилось):
**Две правки:**
1. `auth_provider.dart:13` — начальное состояние `AsyncValue.loading()`
2. `router.dart:21-31` — redirect: `if (AsyncLoading) return null; if (/splash) → /schedule or /login`

**Почему не сработало:** `app.dart` использует `ref.read(routerProvider)`, а не `ref.watch()`. GoRouter создаётся один раз с начальным состоянием. Даже если auth-состояние меняется, `routerProvider` не пересоздаётся, redirect не переоценивается.

### Промт: исправить баг (финальное решение)
Первая попытка не помогла — UI показывает тот же сплэш. Найди настоящую причину и исправь.

### Ответ:
**Корень:** SplashScreen не навигирует сам после `checkSession()`. Надо чтобы он вызывал `context.go()` явно.

**Три правки (все в `router.dart`):**
1. Redirect guard больше не пытается обработать `/splash` — **SplashScreen сам решает, куда идти**
2. `SplashScreen.initState` вызывает `_checkAndNavigate()` (async)
3. `_checkAndNavigate()` ждёт `checkSession()`, читает свежий `authState`, и вызывает `context.go('/schedule')` или `context.go('/login')`

```dart
// router.dart:21-25 — redirect guard пропускает splash
redirect: (context, state) {
  if (state.matchedLocation == '/splash') return null;
  // ... защита остальных роутов
}

// router.dart:117-133 — SplashScreen сам навигирует
Future<void> _checkAndNavigate() async {
  await ref.read(authNotifierProvider.notifier).checkSession();
  final authState = ref.read(authNotifierProvider);
  if (authState is AsyncData && authState.value != null) {
    context.go('/schedule');
  } else {
    context.go('/login');
  }
}
```

## Файлы изменений
- `mobile/lib/app/router.dart`:
  - Строки 21–25 — redirect guard: `/splash` никогда не редиректится
  - Строки 117–133 — `SplashScreen._checkAndNavigate()` с явным `context.go()`

## Верификация
После установки APK на устройство `uiautomator dump` показывает экран логина:
- content-desc="Глини"
- content-desc="Войдите в аккаунт"
- content-desc="Войти"
- content-desc="Нет аккаунта? Зарегистрироваться"
