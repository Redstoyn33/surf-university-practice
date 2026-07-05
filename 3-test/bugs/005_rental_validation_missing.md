# BUG-002: Отсутствует валидация rentalSelected на слотах без проката

## Симптом
Бронь создаётся с `rentalSelected=true` на слоте, где `rentalAvailable=false`.
Сервер не проверяет, что запрос проката возможен только на слотах,
где прокат действительно доступен.

## Шаги воспроизведения
1. Найти слот с `rentalAvailable=false` и `availableSpots > 0`
2. POST /bookings с `{"slotId":<id>, "rentalSelected":true}`
3. В ответе 201, `rentalSelected: true`

**Ожидаемый результат:** HTTP 400/422 с сообщением
`"прокат недоступен для данного слота"` или `rentalSelected` принудительно
выставляется в `false`

**Фактический результат:** HTTP 201, `rentalSelected: true`,
бронь создана с неподдерживаемой опцией

## Причина
В `internal/service/booking.go` метод `CreateBooking` не проверяет
`slot.RentalAvailable` перед установкой `rentalSelected`.
В `internal/repository/booking.go:InsertBooking` нет соответствующей проверки.

## Предлагаемое исправление
Добавить проверку в `CreateBooking`:
```go
if rentalSelected && !slot.RentalAvailable {
    return domain.Booking{}, fmt.Errorf("rental not available for this slot: %w", domain.ErrValidation)
}
```
