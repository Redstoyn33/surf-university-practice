import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

final routerProvider = Provider<GoRouter>((ref) {
  final authState = ref.watch(authStateProvider);
  return GoRouter(
    initialLocation: '/splash',
    redirect: (context, state) {
      final isLoggedIn = authState.isAuthenticated;
      final isAuthRoute = state.matchedLocation == '/login' ||
          state.matchedLocation == '/register' ||
          state.matchedLocation == '/splash';

      if (!isLoggedIn && !isAuthRoute) return '/login';
      if (isLoggedIn && isAuthRoute && state.matchedLocation != '/splash') {
        return '/';
      }
      return null;
    },
    routes: [
      GoRoute(
        path: '/splash',
        builder: (_, __) => const _PlaceholderScreen(title: 'Загрузка...'),
      ),
      GoRoute(
        path: '/login',
        builder: (_, __) => const _PlaceholderScreen(title: 'Вход'),
      ),
      GoRoute(
        path: '/register',
        builder: (_, __) => const _PlaceholderScreen(title: 'Регистрация'),
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
            builder: (_, __) => const _PlaceholderScreen(title: 'Расписание'),
            routes: [
              GoRoute(
                path: 'slots/:slotId',
                builder: (_, state) => _PlaceholderScreen(
                  title: 'Слот #${state.pathParameters['slotId']}',
                ),
              ),
              GoRoute(
                path: 'masters/:id',
                builder: (_, state) => _PlaceholderScreen(
                  title: 'Мастер #${state.pathParameters['id']}',
                ),
              ),
              GoRoute(
                path: 'programs/:id',
                builder: (_, state) => _PlaceholderScreen(
                  title: 'Программа #${state.pathParameters['id']}',
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

class _TabScaffold extends StatelessWidget {
  final Widget child;
  const _TabScaffold({required this.child});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: child,
      bottomNavigationBar: NavigationBar(
        destinations: const [
          NavigationDestination(icon: Icon(Icons.calendar_month), label: 'Расписание'),
          NavigationDestination(icon: Icon(Icons.book_online), label: 'Мои записи'),
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
      body: Center(child: Text(title, style: Theme.of(context).textTheme.headlineSmall)),
    );
  }
}

class AuthState {
  final bool isAuthenticated;
  const AuthState({this.isAuthenticated = false});
}

final authStateProvider = StateProvider<AuthState>((_) => const AuthState());
