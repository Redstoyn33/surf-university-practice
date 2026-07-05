# Баг: Успешный логин/регистрация не вызывает переход на следующий экран

## Симптом
После ввода корректных данных и нажатия «Войти» (или «Зарегистрироваться») API возвращает 200 с токеном и данными клиента, но экран не меняется — пользователь остаётся на форме авторизации.

## Цель
После успешной авторизации приложение должно перейти на `/schedule`.

## Расследование

### Промт
> окей, тот же подход, только теперь не коммить без подтверждения о рабочем фиксе. новый баг: регистрация работает успешно, и авторизация тоже возвращяет правильные коды, но при успешной авторизации ничего не меняется и переход на следующий экран не происходит

### Ответ:
**Проблема в связке двух механизмов:**

1. **`login_screen.dart` / `registration_screen.dart`** — `_submit()` вызывает `login()` без `await`. Даже если добавить `ref.listen` на `data`-состояние, навигация не срабатывает из-за проблемы №2.

2. **`router.dart:17`** — `routerProvider` использует `final authState = ref.watch(authNotifierProvider)` и захватывает это значение в замыкание redirect-колбэка. GoRouter создаётся **один раз** (т.к. `app.dart` читает через `ref.read`). Когда `authNotifierProvider` меняется после логина, `routerProvider` не пересоздаётся, и redirect guard видит **старое** `authState` (= `AsyncData(null)`). При попытке `context.go('/schedule')` redirect проверяет `isLoggedIn = false` и редиректит обратно на `/login`.

### Промт (уточнение)
> давай я тебе помогу, ты явно кликаешь не ту кнопку, что вводить?

> и да, при вводе открывается клавиатура и кнопки сдвигаются, скорей всего промазывал из за этого, я проверил и ничего не переключается

### Финальное решение
**Три правки:**

1. **`router.dart:16-22`** — убрать `ref.watch` из тела `Provider`, перенести `ref.read(authNotifierProvider)` внутрь redirect-колбэка:

```dart
final routerProvider = Provider<GoRouter>((ref) {
  return GoRouter(
    initialLocation: '/splash',
    redirect: (context, state) {
      final authState = ref.read(authNotifierProvider); // свежее состояние
      final isLoggedIn = authState is AsyncData && authState.value != null;
      // ...
    },
  );
});
```

2. **`login_screen.dart:27-36`** — `_submit()` теперь `await` логин и явно навигирует:

```dart
Future<void> _submit() async {
  if (!_formKey.currentState!.validate()) return;
  await ref.read(authNotifierProvider.notifier).login(...);
  if (!mounted) return;
  final authState = ref.read(authNotifierProvider);
  if (authState is AsyncData && authState.value != null) {
    context.go('/schedule');
  }
}
```

3. **`registration_screen.dart:29-38`** — то же самое для регистрации.

## Файлы изменений
- `mobile/lib/app/router.dart` — строка 16–22 (redirect читает `ref.read`)
- `mobile/lib/features/auth/screens/login_screen.dart` — строки 27–36 (`_submit` await + navigate)
- `mobile/lib/features/auth/screens/registration_screen.dart` — строки 29–38 (`_submit` await + navigate)

## Верификация
После установки APK и логина `bugtest` / `123456` приложение переходит на экран расписания. При повторном открытии сплэш-экран распознаёт сохранённый токен и сразу ведёт в меню.
