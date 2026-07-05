# Баг: Вкладка расписания не грузится каждый второй раз

## Симптом
При переключении на вкладку «Расписание» каждый второй раз данные не загружаются — не отправляется запрос на бэкенд. По умолчанию не выбрана никакая дата (хотя данные на сегодня отображаются). Если кликнуть на дату — данные загружаются, и после этого переключение вкладок работает стабильно.

## Цель
При открытии вкладки «Расписание» данные всегда должны загружаться автоматически, а сегодняшняя дата должна быть подсвечена как выбранная.

## Расследование

### Промт
> новый баг: каждый второй раз при открытии вкладки расписание она не грузится, при этом не отправляется никакого запроса на бэкенд, так же по умолчанию не выбрана никакая дата, но когда страница загружается отображаются данные на сегодня, хотя дата выбранной не отображается, и если страница не грузится и я кликаю на дату всё рагружается, при этом не происходит каждого второго раза и грузится всегда

### Ответ:
**Две независимые проблемы:**

1. **`router.dart`** — `ShellRoute` (не `StatefulShellRoute`) не сохраняет состояние вкладок. GoRouter при переключении мог пересоздавать или кешировать виджет в alternating-режиме. `initState` срабатывал через раз.

2. **`date_bar.dart:28-29`** — `date == selectedDate` сравнивает полный `DateTime` с микросекундами. `_selectedDate = DateTime.now()` включает время, а `DateBar` генерирует даты как `DateTime(year, month, day)` (полночь). Сравнение всегда `false` — выделение даты не работало.

### Финальное решение

**1. `router.dart`** — `ShellRoute` → `StatefulShellRoute.indexedStack`:
```dart
StatefulShellRoute.indexedStack(
  builder: (_, __, navigationShell) =>
      _TabScaffold(navigationShell: navigationShell),
  branches: [
    StatefulShellBranch(routes: [
      GoRoute(path: '/schedule', builder: ...),
    ]),
    StatefulShellBranch(routes: [
      GoRoute(path: '/my-bookings', builder: ...),
    ]),
  ],
)
```

**2. `date_bar.dart:28-29`** — сравнение по году/месяцу/дню:
```dart
final isSelected = date.year == selectedDate.year &&
    date.month == selectedDate.month &&
    date.day == selectedDate.day;
final isToday = date.year == today.year &&
    date.month == today.month &&
    date.day == today.day;
```

**3. `schedule_screen.dart`** — загрузка данных из `build()` через `addPostFrameCallback` при `AsyncData` с пустым списком.

## Файлы изменений
- `mobile/lib/app/router.dart` — `ShellRoute` → `StatefulShellRoute.indexedStack`, `_TabScaffold` принимает `StatefulNavigationShell`
- `mobile/lib/features/schedule/widgets/date_bar.dart` — сравнение дат
- `mobile/lib/features/schedule/screens/schedule_screen.dart` — загрузка из `build()`

## Верификация
После установки APK: логин → вкладка «Расписание» загружается сразу, сегодняшняя дата подсвечена, переключение вкладок работает стабильно.
