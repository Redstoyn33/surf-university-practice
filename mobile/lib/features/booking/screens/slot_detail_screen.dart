import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../providers/booking_provider.dart';
import '../widgets/rental_switch.dart';
import '../../../shared/widgets/error_state.dart';
import '../../../shared/widgets/loading_indicator.dart';

class SlotDetailScreen extends ConsumerStatefulWidget {
  final int slotId;

  const SlotDetailScreen({super.key, required this.slotId});

  @override
  ConsumerState<SlotDetailScreen> createState() => _SlotDetailScreenState();
}

class _SlotDetailScreenState extends ConsumerState<SlotDetailScreen> {
  bool _rentalSelected = false;

  Future<void> _book() async {
    final navigator = ScaffoldMessenger.of(context);
    try {
      await ref.read(bookingNotifierProvider.notifier).createBooking(
        slotId: widget.slotId,
        rentalSelected: _rentalSelected,
      );
      navigator.showSnackBar(const SnackBar(content: Text('Вы записаны!')));
      context.go('/my-bookings');
    } on Exception catch (e) {
      final msg = e.toString();
      if (msg.contains('409')) {
        navigator.showSnackBar(
          const SnackBar(content: Text('Вы уже записаны на этот слот')),
        );
      } else {
        navigator.showSnackBar(
          SnackBar(content: Text(msg.contains('400')
              ? 'Проверьте данные'
              : 'Произошла ошибка. Попробуйте позже.')),
        );
      }
    }
  }

  void _showBookingDialog() {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Подтверждение записи'),
        content: const Text('Вы уверены, что хотите записаться?'),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('Отмена')),
          FilledButton(onPressed: () {
            Navigator.pop(ctx);
            _book();
          }, child: const Text('Записаться')),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final slotAsync = ref.watch(slotByIdProvider(widget.slotId));
    final theme = Theme.of(context);

    return slotAsync.when(
      loading: () => const Scaffold(body: LoadingIndicator()),
      error: (e, _) => Scaffold(
        body: ErrorState(
          message: 'Не удалось загрузить слот',
          onRetry: () => ref.invalidate(slotByIdProvider(widget.slotId)),
        ),
      ),
      data: (slot) {
        final isPast = DateTime.parse(slot.dateTime).isBefore(DateTime.now());
        final isSoldOut = slot.availableSpots <= 0;
        final canBook = !isPast && !isSoldOut;

        return Scaffold(
          appBar: AppBar(title: Text(slot.program.name)),
          body: SingleChildScrollView(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(slot.program.name, style: theme.textTheme.headlineSmall),
                const SizedBox(height: 8),
                Row(
                  children: [
                    const Icon(Icons.calendar_today, size: 16),
                    const SizedBox(width: 4),
                    Text(
                      _formatDateTime(slot.dateTime),
                      style: theme.textTheme.bodyLarge,
                    ),
                  ],
                ),
                const SizedBox(height: 4),
                Row(
                  children: [
                    const Icon(Icons.access_time, size: 16),
                    const SizedBox(width: 4),
                    Text(
                      '${slot.dateTime.substring(11, 16)}—${slot.endTime.substring(11, 16)}',
                      style: theme.textTheme.bodyLarge,
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                InkWell(
                  onTap: () => context.push('/schedule/masters/${slot.master.id}'),
                  child: Row(
                    children: [
                      CircleAvatar(
                        radius: 24,
                        backgroundImage: NetworkImage(slot.master.photo),
                        onBackgroundImageError: (_, __) => const Icon(Icons.person),
                      ),
                      const SizedBox(width: 12),
                      Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(slot.master.name, style: theme.textTheme.titleMedium),
                          Text(slot.master.level, style: theme.textTheme.bodySmall),
                        ],
                      ),
                      const Spacer(),
                      const Icon(Icons.chevron_right),
                    ],
                  ),
                ),
                const SizedBox(height: 16),
                InkWell(
                  onTap: () => context.push('/schedule/programs/${slot.program.id}'),
                  child: Row(
                    children: [
                      const Icon(Icons.info_outline, size: 16),
                      const SizedBox(width: 4),
                      Text('О программе', style: theme.textTheme.titleSmall),
                      const Spacer(),
                      const Icon(Icons.chevron_right),
                    ],
                  ),
                ),
                const SizedBox(height: 16),
                Row(
                  children: [
                    Icon(Icons.people, color: isSoldOut ? theme.colorScheme.error : theme.colorScheme.primary),
                    const SizedBox(width: 4),
                    Text(
                      isSoldOut
                          ? 'Нет свободных мест'
                          : '${slot.availableSpots} из ${slot.totalSpots} мест свободно',
                      style: TextStyle(
                        color: isSoldOut ? theme.colorScheme.error : null,
                      ),
                    ),
                  ],
                ),
                if (slot.rentalAvailable) ...[
                  const SizedBox(height: 16),
                  RentalSwitch(
                    value: _rentalSelected,
                    price: slot.rentalPrice,
                    onChanged: (v) => setState(() => _rentalSelected = v),
                  ),
                ],
                if (isPast)
                  Padding(
                    padding: const EdgeInsets.only(top: 16),
                    child: Chip(
                      avatar: const Icon(Icons.block),
                      label: const Text('Слот уже прошёл'),
                      backgroundColor: theme.colorScheme.errorContainer,
                    ),
                  ),
                const SizedBox(height: 24),
                SizedBox(
                  width: double.infinity,
                  height: 48,
                  child: FilledButton(
                    onPressed: canBook ? _showBookingDialog : null,
                    child: Text(isSoldOut ? 'Нет мест' : 'Записаться'),
                  ),
                ),
              ],
            ),
          ),
        );
      },
    );
  }

  String _formatDateTime(String dt) {
    final date = DateTime.parse(dt);
    final months = [
      '', 'янв', 'фев', 'мар', 'апр', 'май', 'июн',
      'июл', 'авг', 'сен', 'окт', 'ноя', 'дек',
    ];
    return '${date.day} ${months[date.month]} ${date.year}';
  }
}
