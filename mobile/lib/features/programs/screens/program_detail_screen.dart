import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../providers/program_provider.dart';
import '../../../shared/widgets/error_state.dart';
import '../../../shared/widgets/loading_indicator.dart';

class ProgramDetailScreen extends ConsumerWidget {
  final int id;

  const ProgramDetailScreen({super.key, required this.id});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final programAsync = ref.watch(programByIdProvider(id));
    final theme = Theme.of(context);

    return programAsync.when(
      loading: () => const Scaffold(body: LoadingIndicator()),
      error: (e, _) => Scaffold(
        body: ErrorState(
          message: 'Не удалось загрузить программу',
          onRetry: () => ref.invalidate(programByIdProvider(id)),
        ),
      ),
      data: (program) {
        return Scaffold(
          appBar: AppBar(title: Text(program.name)),
          body: SingleChildScrollView(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(program.name, style: theme.textTheme.headlineSmall),
                const SizedBox(height: 12),
                Chip(
                  avatar: const Icon(Icons.group, size: 16),
                  label: Text('До ${program.maxCapacity} участников'),
                ),
                const SizedBox(height: 16),
                Text('Описание', style: theme.textTheme.titleMedium),
                const SizedBox(height: 8),
                Text(program.description, style: theme.textTheme.bodyLarge),
                const SizedBox(height: 24),
                Text('Мастера', style: theme.textTheme.titleMedium),
                const SizedBox(height: 8),
                Wrap(
                  spacing: 8,
                  runSpacing: 4,
                  children: program.masterIds.map((mid) => ActionChip(
                    label: Text('Мастер #$mid'),
                    onPressed: () => context.push('/schedule/masters/$mid'),
                  )).toList(),
                ),
              ],
            ),
          ),
        );
      },
    );
  }
}
