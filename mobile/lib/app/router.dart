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

final routerProvider = Provider<GoRouter>((ref) {
  final authState = ref.watch(authNotifierProvider);

  return GoRouter(
    initialLocation: '/splash',
    redirect: (context, state) {
      final isLoggedIn = authState is AsyncData && authState.value != null;
      final isAuthRoute = state.matchedLocation == '/login' ||
          state.matchedLocation == '/register' ||
          state.matchedLocation == '/splash';

      if (!isLoggedIn && !isAuthRoute) return '/login';
      if (isLoggedIn && isAuthRoute && state.matchedLocation != '/splash') {
        return '/schedule';
      }
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
      ShellRoute(
        builder: (_, __, child) => _TabScaffold(child: child),
        routes: [
          GoRoute(
            path: '/',
            redirect: (_, __) => '/schedule',
          ),
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
          GoRoute(
            path: '/my-bookings',
            builder: (_, __) => const _PlaceholderScreen(title: 'Мои записи'),
            routes: [
              GoRoute(
                path: 'bookings/:id',
                builder: (_, state) => _PlaceholderScreen(
                  title: 'Бронь #${state.pathParameters['id']}',
                ),
              ),
              GoRoute(
                path: 'rate',
                builder: (_, __) => const _PlaceholderScreen(title: 'Оценка'),
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
    Future.microtask(() {
      ref.read(authNotifierProvider.notifier).checkSession();
    });
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
  final Widget child;
  const _TabScaffold({required this.child});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: child,
      bottomNavigationBar: NavigationBar(
        destinations: const [
          NavigationDestination(
              icon: Icon(Icons.calendar_month), label: 'Расписание'),
          NavigationDestination(
              icon: Icon(Icons.book_online), label: 'Мои записи'),
        ],
        selectedIndex: 0,
        onDestinationSelected: (i) {
          if (i == 0) context.go('/schedule');
          if (i == 1) context.go('/my-bookings');
        },
      ),
    );
  }
}

class _PlaceholderScreen extends StatelessWidget {
  final String title;
  const _PlaceholderScreen({required this.title});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(title)),
      body: Center(
          child: Text(title, style: Theme.of(context).textTheme.headlineSmall)),
    );
  }
}
