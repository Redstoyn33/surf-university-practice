# Баг: Бесконечная загрузка на сплэш-экране

## Симптом
Приложение открывается, показывает сплэш-экран с логотипом «Глини» и крутящимся `CircularProgressIndicator`. Загрузка не заканчивается — экран висит бесконечно, независимо от того, авторизован пользователь или нет.

## Цель
После завершения проверки сессии (`checkSession`) приложение должно:
- перейти на `/schedule`, если пользователь авторизован;
- перейти на `/login`, если не авторизован.

## Расследование

### Промт: исследование
> надо разрешить баг, для бага надо будет создать свой md файл в директории 3-test/bugs/ , файл должен содержать симптом, цель, а также прикреплены мои промты с твоими сжатыми ответами при работе с решением бага. бекенд поднят, проверяй логи с него. телефон для отладки с приложением запущен на 192.168.1.219:39937, текущий баг: приложение открывается но бесконечно грузится

### Промт: попытка исправления
> можешь ли ты проверить что на экране? потому что проблема не решена и всёё езё показывается загрузка

### Первая попытка (провалилась)
**Две правки:**
1. `auth_provider.dart:13` — начальное состояние `AsyncValue.loading()`
2. `router.dart:21-31` — redirect: `if (AsyncLoading) return null; if (/splash) → /schedule or /login`

**Почему не сработало:** `app.dart` использует `ref.read(routerProvider)`, а не `ref.watch()`. GoRouter создаётся один раз с начальным состоянием. Даже если auth-состояние меняется, `routerProvider` не пересоздаётся, redirect не переоценивается.

### Финальное решение
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
