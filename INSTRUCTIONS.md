# Инструкции

## Бэкенд (Go)

```bash
cd backend

# Запуск сервера (с миграциями + сид-данными)
make run
# или: nix-shell -p go --run "go run ./cmd/server"

# Прогнать тесты
make test
# или: nix-shell -p go --run "go test ./... -count=1"
```

Сервер слушает `:8080`. Сид-данные (3 мастера, 3 программы, слоты на неделю) заливаются при первом запуске. Для тестового логина: `test@test.com` / `123456`.

## Мобильное приложение (Flutter)

```bash
cd mobile

# 1. Подключить андроид-устройство по ADB
adb connect <IP:PORT>          # пример: adb connect 192.168.1.219:44949

# 2. Собрать и запустить на устройстве
flutter run -d '<DEVICE_NAME>'  # пример: flutter run -d 'WP19 Pro'

# 3. Проверить код на ошибки
dart analyze lib/
```

### Если не собирается (Android SDK)

Настроить Android SDK через nix:

```bash
# Временная ANDROID_HOME создаётся в /tmp/android-home (уже настроена)
# Если нужно пересоздать:
flutter config --android-sdk /tmp/android-home
export ANDROID_HOME=/tmp/android-home
export ANDROID_SDK_ROOT=/tmp/android-home
flutter run -d '<DEVICE_NAME>'
```

### API URL

По умолчанию мобилка стучится на `http://192.168.1.101:8080` (хост-машина в локальной сети).  
Переопределить можно в `lib/core/api_client.dart` — константа `apiHostOverride`.
