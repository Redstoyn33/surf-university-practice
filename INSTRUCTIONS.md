# Инструкции

## Бэкенд (Go)

### NixOS

```bash
cd backend

# Запуск сервера (с миграциями + сид-данными)
make run
# или: nix-shell -p go --run "go run ./cmd/server"

# Прогнать тесты
make test
# или: nix-shell -p go --run "go test ./... -count=1"
```

### Ubuntu (и другие дистрибутивы)

```bash
cd backend

# Установка Go (если ещё не установлен)
sudo apt update && sudo apt install -y golang-go

# Запуск сервера (с миграциями + сид-данными)
make run
# или: go run ./cmd/server

# Прогнать тесты
make test
# или: go test ./... -count=1
```

Сервер слушает `:8080`. Сид-данные (3 мастера, 3 программы, слоты на неделю) заливаются при первом запуске. Для тестового логина: `test@test.com` / `123456`.

## Мобильное приложение (Flutter)

### NixOS

```bash
cd mobile

# 1. Подключить андроид-устройство по ADB
adb connect <IP:PORT>          # пример: adb connect 192.168.1.219:44949

# 2. Собрать и запустить на устройстве
flutter run -d '<DEVICE_NAME>'  # пример: flutter run -d 'WP19 Pro'

# 3. Проверить код на ошибки
dart analyze lib/
```

### Ubuntu (и другие дистрибутивы)

```bash
cd mobile

# Установка Flutter (если ещё не установлен)
# См. https://docs.flutter.dev/get-started/install/linux
sudo snap install flutter --classic

# 1. Подключить андроид-устройство по ADB
adb connect <IP:PORT>          # пример: adb connect 192.168.1.219:44949

# 2. Собрать и запустить на устройстве
flutter run -d '<DEVICE_NAME>'  # пример: flutter run -d 'WP19 Pro'

# 3. Проверить код на ошибки
dart analyze lib/
```

### Android SDK (Ubuntu)

```bash
# Установка Android SDK через командную строку
sudo apt install -y android-sdk
export ANDROID_HOME=/usr/lib/android-sdk
export ANDROID_SDK_ROOT=/usr/lib/android-sdk

# Или вручную через Android Studio:
# 1. Установить Android Studio (snap install android-studio)
# 2. Открыть, пройти настройку SDK (меню: Configure → SDK Manager)
# 3. ANDROID_HOME = ~/Android/Sdk
echo 'export ANDROID_HOME=$HOME/Android/Sdk' >> ~/.bashrc
echo 'export ANDROID_SDK_ROOT=$HOME/Android/Sdk' >> ~/.bashrc
source ~/.bashrc

# Принять лицензии
yes | $ANDROID_HOME/cmdline-tools/latest/bin/sdkmanager --licenses

# Проверить Flutter
flutter doctor -v

# Запуск
flutter run -d '<DEVICE_NAME>'
```

## Сборка APK

Собрать APK без запуска на устройство:

```bash
cd mobile

# Debug APK (быстрая сборка, подписано debug-ключом)
flutter build apk --debug

# Release APK (оптимизировано, требует keystore)
flutter build apk --release

# APK будет в: build/app/outputs/flutter-apk/app-debug.apk
#                        build/app/outputs/flutter-apk/app-release.apk

# Установить готовый APK на устройство:
adb connect <IP:PORT>
adb install build/app/outputs/flutter-apk/app-debug.apk
```

### API URL

По умолчанию мобилка стучится на `http://192.168.1.101:8080` (хост-машина в локальной сети).  
Для эмулятора Android (Android Studio AVD) используется `http://10.0.2.2:8080`.  
Переопределить можно в `lib/core/api_client.dart` — константа `apiHostOverride`.
