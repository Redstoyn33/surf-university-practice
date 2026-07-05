# Фича: Вкладка «Мастера» со списком всех мастеров

## Описание
Добавлена третья вкладка в нижнюю навигацию — «Мастера». Отображает список всех мастеров, загружаемый с `GET /masters`.

## Реализация

### Новые файлы
- `mobile/lib/features/masters/screens/masters_screen.dart` — экран списка мастеров

### Изменённые файлы
- `mobile/lib/app/router.dart`:
  - Добавлен импорт `masters_screen.dart`
  - Добавлен `StatefulShellBranch` с путём `/masters` и вложенным `:id` для профиля
  - `NavigationBar` расширен до 3 вкладок (Расписание / Мастера / Мои записи)

### Детали
- **Провайдер:** существующий `mastersListProvider` (`FutureProvider<List<Master>>`) из `master_provider.dart`
- **Экран:** `ConsumerWidget` с `RefreshIndicator`, `ListView.builder`, `Card` + `InkWell` на каждый элемент
- **Карточка мастера:** `CircleAvatar` (фото), имя, рейтинг (звёзды), уровень (`Опытный` / `Новичок`)
- **Навигация:** тап по карточке → `context.push('/masters/:id')` → `MasterProfileScreen` (внутри той же вкладки)
- **Pull-to-refresh:** `ref.refresh(mastersListProvider.future)`
- **Пустой список / ошибка:** `EmptyState` / `ErrorState` с `onRetry`

## Промт
> в 3-test/changes/ по похожей схеме с багами, но только для новых фич. реализуем вкладку с списком всех мастеров

## Верификация
После установки APK: логин → таб «Мастера» → отображаются 3 мастера (Анна Кузнецова, Иван Петров, Елена Соколова) с рейтингом и уровнем. Тап по карточке открывает профиль мастера.
