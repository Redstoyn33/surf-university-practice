import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../../core/models/slot.dart';

class SlotCard extends StatelessWidget {
  final Slot slot;

  const SlotCard({super.key, required this.slot});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isAvailable = slot.availableSpots > 0;
    final time = slot.dateTime.length >= 16
        ? slot.dateTime.substring(11, 16)
        : slot.dateTime;
    final endTime = slot.endTime.length >= 16
        ? slot.endTime.substring(11, 16)
        : slot.endTime;

    return Card(
      color: isAvailable ? null : theme.colorScheme.surfaceContainerHighest,
      child: InkWell(
        borderRadius: BorderRadius.circular(12),
        onTap: isAvailable ? () => context.push('/schedule/slots/${slot.id}') : null,
        child: Padding(
          padding: const EdgeInsets.all(12),
          child: Row(
            children: [
              Container(
                width: 60,
                height: 60,
                decoration: BoxDecoration(
                  color: theme.colorScheme.primaryContainer,
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Center(
                  child: Text(
                    time,
                    style: theme.textTheme.titleMedium?.copyWith(
                      fontWeight: FontWeight.bold,
                      color: theme.colorScheme.onPrimaryContainer,
                    ),
                  ),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      slot.program.name,
                      style: theme.textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Row(
                      children: [
                        Icon(Icons.person, size: 14, color: theme.colorScheme.onSurfaceVariant),
                        const SizedBox(width: 4),
                        Text(slot.master.name, style: theme.textTheme.bodySmall),
                        const SizedBox(width: 12),
                        Icon(Icons.access_time, size: 14, color: theme.colorScheme.onSurfaceVariant),
                        const SizedBox(width: 4),
                        Text('$time—$endTime', style: theme.textTheme.bodySmall),
                      ],
                    ),
                    const SizedBox(height: 4),
                    Row(
                      children: [
                        if (slot.rentalAvailable)
                          Padding(
                            padding: const EdgeInsets.only(right: 8),
                            child: Icon(Icons.build, size: 14, color: theme.colorScheme.primary),
                          ),
                        Text(
                          isAvailable
                              ? '${slot.availableSpots}/${slot.totalSpots} мест'
                              : 'Нет мест',
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: isAvailable
                                ? theme.colorScheme.primary
                                : theme.colorScheme.error,
                          ),
                        ),
                        const Spacer(),
                        if (slot.rentalAvailable)
                          Text(
                            '+${slot.rentalPrice.toStringAsFixed(0)} ₽',
                            style: theme.textTheme.bodySmall?.copyWith(
                              color: theme.colorScheme.onSurfaceVariant,
                            ),
                          ),
                      ],
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
