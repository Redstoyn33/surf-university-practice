import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../providers/my_bookings_provider.dart';
import '../widgets/booking_card.dart';
import '../../../shared/widgets/error_state.dart';
import '../../../shared/widgets/empty_state.dart';

class MyBookingsScreen extends ConsumerStatefulWidget {
  const MyBookingsScreen({super.key});

  @override
  ConsumerState<MyBookingsScreen> createState() => _MyBookingsScreenState();
}

class _MyBookingsScreenState extends ConsumerState<MyBookingsScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  static const _tabs = ['Активные', 'Прошедшие', 'Отменённые'];
  static const _statusFilters = <String?>['активна', null, null];

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: _tabs.length, vsync: this);
    _tabController.addListener(() {
      if (!_tabController.indexIsChanging) _loadBookings();
    });
    WidgetsBinding.instance.addPostFrameCallback((_) => _loadBookings());
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  void _loadBookings() {
    final status = _statusFilters[_tabController.index];
    ref.read(myBookingsProvider.notifier).load(status: status);
  }

  @override
  Widget build(BuildContext context) {
    final bookingsAsync = ref.watch(myBookingsProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Мои записи'),
        bottom: TabBar(
          controller: _tabController,
          tabs: _tabs.map((t) => Tab(text: t)).toList(),
        ),
      ),
      body: RefreshIndicator(
        onRefresh: () async => _loadBookings(),
        child: bookingsAsync.when(
          loading: () => const Center(child: CircularProgressIndicator()),
          error: (e, _) => ErrorState(
            message: 'Не удалось загрузить записи',
            onRetry: _loadBookings,
          ),
          data: (bookings) {
            if (bookings.isEmpty) {
              return EmptyState(
                icon: Icons.event_busy,
                message: _tabController.index == 0
                    ? 'Нет активных записей'
                    : _tabController.index == 1
                        ? 'Нет прошедших записей'
                        : 'Нет отменённых записей',
                actionLabel: 'К расписанию',
                onAction: () => context.go('/schedule'),
              );
            }
            return ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: bookings.length,
              itemBuilder: (_, i) => Padding(
                padding: const EdgeInsets.only(bottom: 8),
                child: BookingCard(booking: bookings[i]),
              ),
            );
          },
        ),
      ),
    );
  }
}
