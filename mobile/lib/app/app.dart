import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../core/theme.dart';
import '../shared/widgets/connectivity_banner.dart';
import 'router.dart';

class GliniApp extends ConsumerWidget {
  const GliniApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = ref.read(routerProvider);
    return MaterialApp.router(
      title: 'Глини',
      theme: AppTheme.light,
      routerConfig: router,
      debugShowCheckedModeBanner: false,
      builder: (context, child) {
        return GestureDetector(
          onTap: () => FocusScope.of(context).unfocus(),
          child: ConnectivityBanner(child: child ?? const SizedBox.shrink()),
        );
      },
    );
  }
}
