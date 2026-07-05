import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../providers/booking_provider.dart';
import '../../../shared/widgets/error_state.dart';
import '../../../shared/widgets/loading_indicator.dart';

class BookingDetailScreen extends ConsumerStatefulWidget {
  final int bookingId;

  const BookingDetailScreen({super.key, required this.bookingId});

  @override
  ConsumerState<BookingDetailScreen> createState() =>
      _BookingDetailScreenState();
}

class _BookingDetailScreenState extends ConsumerState<BookingDetailScreen> {
  late Future _loadFuture;

  @override
  void initState() {
    super.initState();
    _loadFuture = _loadBooking();
  }

  Future _loadBooking() async {
    final repo = ref.read(bookingRepositoryProvider);
    return repo.getBookingById(widget.bookingId);
  }

  Future<void> _cancelBooking() async {
    final navigator = ScaffoldMessenger.of(context);
    try {
      await ref.read(bookingNotifierProvider.notifier).cancelBooking(widget.bookingId);
      navigator.showSnackBar(
        const SnackBar(content: Text('Запись отменена')),
      );
      setState(() => _loadFuture = _loadBooking());
    } on Exception catch (e) {
      final msg = e.toString();
      navigator.showSnackBar(
        SnackBar(
          content: Text(msg.contains('422')
              ? 'Отмена доступна за 4+ часа до начала'
              : 'Произошла ошибка'),
        ),
      );
    }
  }

  void _showCancelDialog() {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Отмена записи'),
        content: const Text('Вы уверены, что хотите отменить запись?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('Нет'),
          ),
          FilledButton(
            onPressed: () {
              Navigator.pop(ctx);
              _cancelBooking();
            },
            child: const Text('Да, отменить'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return FutureBuilder(
      future: _loadFuture,
      builder: (context, snapshot) {
        if (snapshot.connectionState == ConnectionState.waiting) {
          return const Scaffold(body: LoadingIndicator());
        }
        if (snapshot.hasError) {
          return Scaffold(
            body: ErrorState(
              message: 'Не удалось загрузить бронь',
              onRetry: () => setState(() => _loadFuture = _loadBooking()),
            ),
          );
        }

        final booking = snapshot.data!;
        final slot = booking.slot;
        final isActive = booking.isActive;
        final slotTime = DateTime.parse(slot.dateTime);
        final hoursUntil = slotTime.difference(DateTime.now()).inHours;
        final canCancel = isActive && hoursUntil >= 4;

        return Scaffold(
          appBar: AppBar(title: Text('Бронь #${booking.id}')),
          body: SingleChildScrollView(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(slot.program.name,
                    style: theme.textTheme.headlineSmall),
                const SizedBox(height: 8),
                _infoRow(Icons.person, slot.master.name, () {
                  context.push('/schedule/masters/${slot.master.id}');
                }),
                _infoRow(Icons.calendar_today, _formatDate(slot.dateTime), null),
                _infoRow(
                  Icons.access_time,
                  '${slot.dateTime.substring(11, 16)}—${slot.endTime.substring(11, 16)}',
                  null,
                ),
                _infoRow(
                  Icons.build,
                  booking.rentalSelected ? 'Прокат включён' : 'Без проката',
                  null,
                ),
                const SizedBox(height: 8),
                Chip(
                  label: Text(booking.status),
                  backgroundColor: switch (booking.status) {
                    'активна' => Colors.green.withValues(alpha: 0.15),
                    'отменена клиентом' => Colors.orange.withValues(alpha: 0.15),
                    'отменена мастерской' => Colors.red.withValues(alpha: 0.15),
                    _ => null,
                  },
                ),
                if (booking.cancellationReason != null) ...[
                  const SizedBox(height: 8),
                  Text(
                    'Причина отмены: ${booking.cancellationReason}',
                    style: TextStyle(color: theme.colorScheme.error),
                  ),
                ],
                const SizedBox(height: 24),
                if (isActive && !canCancel)
                  Card(
                    color: theme.colorScheme.errorContainer,
                    child: const Padding(
                      padding: EdgeInsets.all(12),
                      child: Text(
                        'Отмена менее чем за 4 часа — обратитесь в поддержку',
                      ),
                    ),
                  ),
                if (canCancel)
                  SizedBox(
                    width: double.infinity,
                    height: 48,
                    child: OutlinedButton.icon(
                      onPressed: _showCancelDialog,
                      icon: const Icon(Icons.cancel),
                      label: const Text('Отменить запись'),
                      style: OutlinedButton.styleFrom(
                        foregroundColor: theme.colorScheme.error,
                        side: BorderSide(color: theme.colorScheme.error),
                      ),
                    ),
                  ),
              ],
            ),
          ),
        );
      },
    );
  }

  Widget _infoRow(IconData icon, String text, VoidCallback? onTap) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: InkWell(
        onTap: onTap,
        child: Row(
          children: [
            Icon(icon, size: 16, color: Theme.of(context).colorScheme.onSurfaceVariant),
            const SizedBox(width: 8),
            Text(text),
            if (onTap != null) ...[
              const Spacer(),
              const Icon(Icons.chevron_right, size: 16),
            ],
          ],
        ),
      ),
    );
  }

  String _formatDate(String dt) {
    final date = DateTime.parse(dt);
    final months = [
      '', 'янв', 'фев', 'мар', 'апр', 'май', 'июн',
      'июл', 'авг', 'сен', 'окт', 'ноя', 'дек',
    ];
    return '${date.day} ${months[date.month]} ${date.year}';
  }
}
