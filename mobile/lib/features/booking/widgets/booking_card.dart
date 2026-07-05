import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../../core/models/booking.dart';

class BookingCard extends StatelessWidget {
  final Booking booking;
  final VoidCallback? onRate;

  const BookingCard({super.key, required this.booking, this.onRate});

  bool get _isPast {
    try {
      return DateTime.parse(booking.slot.dateTime).isBefore(DateTime.now());
    } catch (_) {
      return false;
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final slot = booking.slot;
    final time = slot.dateTime.length >= 16
        ? slot.dateTime.substring(11, 16)
        : slot.dateTime;
    final date = slot.dateTime.length >= 10
        ? slot.dateTime.substring(0, 10)
        : slot.dateTime;

    final statusColor = switch (booking.status) {
      'активна' => Colors.green,
      'отменена клиентом' => Colors.orange,
      'отменена мастерской' => Colors.red,
      _ => Colors.grey,
    };

    return Card(
      child: InkWell(
        borderRadius: BorderRadius.circular(12),
        onTap: () => context.push('/my-bookings/bookings/${booking.id}'),
        child: Padding(
          padding: const EdgeInsets.all(12),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Expanded(
                    child: Text(
                      slot.program.name,
                      style: theme.textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                    decoration: BoxDecoration(
                      color: statusColor.withValues(alpha: 0.15),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Text(
                      booking.status,
                      style: theme.textTheme.labelSmall?.copyWith(
                        color: statusColor,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 8),
              Row(
                children: [
                  Icon(Icons.person, size: 14, color: theme.colorScheme.onSurfaceVariant),
                  const SizedBox(width: 4),
                  Text(slot.master.name, style: theme.textTheme.bodySmall),
                  const SizedBox(width: 16),
                  Icon(Icons.calendar_today, size: 14, color: theme.colorScheme.onSurfaceVariant),
                  const SizedBox(width: 4),
                  Text(date, style: theme.textTheme.bodySmall),
                  const SizedBox(width: 16),
                  Icon(Icons.access_time, size: 14, color: theme.colorScheme.onSurfaceVariant),
                  const SizedBox(width: 4),
                  Text(time, style: theme.textTheme.bodySmall),
                ],
              ),
              if (booking.cancellationReason != null) ...[
                const SizedBox(height: 8),
                Text(
                  'Причина: ${booking.cancellationReason}',
                  style: theme.textTheme.bodySmall?.copyWith(color: Colors.red),
                ),
              ],
              if (_isPast && booking.isActive) ...[
                const SizedBox(height: 8),
                SizedBox(
                  width: double.infinity,
                  child: OutlinedButton.icon(
                    onPressed: onRate,
                    icon: const Icon(Icons.star_outline, size: 16),
                    label: const Text('Оценить мастера'),
                    style: OutlinedButton.styleFrom(
                      foregroundColor: theme.colorScheme.primary,
                    ),
                  ),
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }
}
