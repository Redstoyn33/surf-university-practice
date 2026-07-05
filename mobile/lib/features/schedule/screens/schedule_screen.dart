import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../providers/schedule_provider.dart';
import '../widgets/date_bar.dart';
import '../widgets/slot_card.dart';
import '../../../shared/widgets/error_state.dart';
import '../../../shared/widgets/empty_state.dart';

class ScheduleScreen extends ConsumerStatefulWidget {
  const ScheduleScreen({super.key});

  @override
  ConsumerState<ScheduleScreen> createState() => _ScheduleScreenState();
}

class _ScheduleScreenState extends ConsumerState<ScheduleScreen> {
  DateTime _selectedDate =
      DateTime(DateTime.now().year, DateTime.now().month, DateTime.now().day);

  void _loadSlots() {
    final from = _formatDate(_selectedDate);
    final to = _formatDate(_selectedDate);
    ref.read(slotsNotifierProvider.notifier).loadSlots(dateFrom: from, dateTo: to);
  }

  String _formatDate(DateTime date) =>
      '${date.year}-${date.month.toString().padLeft(2, '0')}-${date.day.toString().padLeft(2, '0')}';

  @override
  Widget build(BuildContext context) {
    final slotsAsync = ref.watch(slotsNotifierProvider);

    if (slotsAsync is AsyncData) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (!mounted) return;
        final current = ref.read(slotsNotifierProvider);
        if (current is AsyncData && current.requireValue.isEmpty) {
          _loadSlots();
        }
      });
    }

    return RefreshIndicator(
      onRefresh: () async => _loadSlots(),
      child: CustomScrollView(
        slivers: [
          SliverAppBar(
            title: const Text('Расписание'),
            pinned: true,
          ),
          SliverToBoxAdapter(
            child: DateBar(
              selectedDate: _selectedDate,
              onDateSelected: (date) {
                setState(() => _selectedDate = date);
                _loadSlots();
              },
            ),
          ),
          const SliverToBoxAdapter(child: SizedBox(height: 8)),
          slotsAsync.when(
            loading: () => const SliverFillRemaining(
              child: Center(child: CircularProgressIndicator()),
            ),
            error: (e, _) => SliverFillRemaining(
              child: ErrorState(
                message: 'Не удалось загрузить расписание',
                onRetry: _loadSlots,
              ),
            ),
            data: (slots) {
              if (slots.isEmpty) {
                return const SliverFillRemaining(
                  child: EmptyState(
                    icon: Icons.event_busy,
                    message: 'Нет слотов на эту дату',
                  ),
                );
              }
              return SliverList(
                delegate: SliverChildBuilderDelegate(
                  (_, i) => Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
                    child: SlotCard(slot: slots[i]),
                  ),
                  childCount: slots.length,
                ),
              );
            },
          ),
        ],
      ),
    );
  }
}
