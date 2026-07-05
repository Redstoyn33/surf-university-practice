import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../features/auth/screens/login_screen.dart';
import '../features/auth/screens/registration_screen.dart';
import '../features/auth/providers/auth_provider.dart';
import '../features/schedule/screens/schedule_screen.dart';
import '../features/masters/screens/master_profile_screen.dart';
import '../features/programs/screens/program_detail_screen.dart';
import '../features/booking/screens/slot_detail_screen.dart';
import '../features/booking/screens/my_bookings_screen.dart';
import '../features/booking/screens/booking_detail_screen.dart';
import '../features/ratings/screens/rate_master_screen.dart';

final routerProvider = Provider<GoRouter>((ref) {
  return GoRouter(
    initialLocation: '/splash',
    redirect: (context, state) {
      final authState = ref.read(authNotifierProvider);
      final isLoggedIn = authState is AsyncData && authState.value != null;

      // SplashScreen управляет навигацией сам
      if (state.matchedLocation == '/splash') return null;

      final isAuthRoute = state.matchedLocation == '/login' ||
          state.matchedLocation == '/register';

      if (!isLoggedIn && !isAuthRoute) return '/login';
      if (isLoggedIn && isAuthRoute) return '/schedule';

      return null;
    },
    routes: [
      GoRoute(
        path: '/splash',
        builder: (_, __) => const SplashScreen(),
      ),
      GoRoute(
        path: '/login',
        builder: (_, __) => const LoginScreen(),
      ),
      GoRoute(
        path: '/register',
        builder: (_, __) => const RegistrationScreen(),
      ),
      StatefulShellRoute.indexedStack(
        builder: (_, __, navigationShell) =>
            _TabScaffold(navigationShell: navigationShell),
        branches: [
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: '/schedule',
                builder: (_, __) => const ScheduleScreen(),
                routes: [
                  GoRoute(
                    path: 'slots/:slotId',
                    builder: (_, state) => SlotDetailScreen(
                      slotId: int.parse(state.pathParameters['slotId']!),
                    ),
                  ),
                  GoRoute(
                    path: 'masters/:id',
                    builder: (_, state) => MasterProfileScreen(
                      id: int.parse(state.pathParameters['id']!),
                    ),
                  ),
                  GoRoute(
                    path: 'programs/:id',
                    builder: (_, state) => ProgramDetailScreen(
                      id: int.parse(state.pathParameters['id']!),
                    ),
                  ),
                ],
              ),
            ],
          ),
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: '/my-bookings',
                builder: (_, __) => const MyBookingsScreen(),
                routes: [
                  GoRoute(
                    path: 'bookings/:id',
                    builder: (_, state) => BookingDetailScreen(
                      bookingId: int.parse(state.pathParameters['id']!),
                    ),
                  ),
                  GoRoute(
                    path: 'rate',
                    builder: (_, state) {
                      final q = state.uri.queryParameters;
                      return RateMasterScreen(
                        masterId: int.parse(q['masterId'] ?? '0'),
                        slotId: int.parse(q['slotId'] ?? '0'),
                        masterName: q['masterName'] ?? '',
                        programName: q['programName'] ?? '',
                        slotDate: q['slotDate'] ?? '',
                      );
                    },
                  ),
                ],
              ),
            ],
          ),
        ],
      ),
    ],
  );
});

class SplashScreen extends ConsumerStatefulWidget {
  const SplashScreen({super.key});

  @override
  ConsumerState<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends ConsumerState<SplashScreen> {
  @override
  void initState() {
    super.initState();
    _checkAndNavigate();
  }

  Future<void> _checkAndNavigate() async {
    await ref.read(authNotifierProvider.notifier).checkSession();
    final authState = ref.read(authNotifierProvider);
    final isLoggedIn = authState is AsyncData && authState.value != null;
    if (isLoggedIn) {
      context.go('/schedule');
    } else {
      context.go('/login');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.weekend,
                size: 80, color: Theme.of(context).colorScheme.primary),
            const SizedBox(height: 16),
            Text('Глини', style: Theme.of(context).textTheme.headlineMedium),
            const SizedBox(height: 24),
            const CircularProgressIndicator(),
          ],
        ),
      ),
    );
  }
}

class _TabScaffold extends StatelessWidget {
  final StatefulNavigationShell navigationShell;
  const _TabScaffold({required this.navigationShell});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: navigationShell,
      bottomNavigationBar: NavigationBar(
        destinations: const [
          NavigationDestination(
              icon: Icon(Icons.calendar_month), label: 'Расписание'),
          NavigationDestination(
              icon: Icon(Icons.book_online), label: 'Мои записи'),
        ],
        selectedIndex: navigationShell.currentIndex,
        onDestinationSelected: navigationShell.goBranch,
      ),
    );
  }
}


