import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../providers/master_provider.dart';
import '../../../shared/widgets/error_state.dart';
import '../../../shared/widgets/loading_indicator.dart';

class MasterProfileScreen extends ConsumerWidget {
  final int id;

  const MasterProfileScreen({super.key, required this.id});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final masterAsync = ref.watch(masterByIdProvider(id));
    final theme = Theme.of(context);

    return masterAsync.when(
      loading: () => const Scaffold(body: LoadingIndicator()),
      error: (e, _) => Scaffold(
        body: ErrorState(
          message: 'Не удалось загрузить профиль мастера',
          onRetry: () => ref.invalidate(masterByIdProvider(id)),
        ),
      ),
      data: (master) {
        return Scaffold(
          appBar: AppBar(
            title: Text(master.name),
            actions: [
              IconButton(
                icon: Icon(master.level == 'опытный'
                    ? Icons.verified
                    : Icons.auto_awesome),
                tooltip: master.level,
                onPressed: null,
              ),
            ],
          ),
          body: SingleChildScrollView(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Center(
                  child: CircleAvatar(
                    radius: 48,
                    backgroundImage: NetworkImage(master.photo),
                    onBackgroundImageError: (_, __) => Icon(
                      Icons.person,
                      size: 48,
                      color: theme.colorScheme.onSurfaceVariant,
                    ),
                  ),
                ),
                const SizedBox(height: 16),
                Center(
                  child: Text(master.name, style: theme.textTheme.headlineSmall),
                ),
                const SizedBox(height: 8),
                Center(
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      ...List.generate(5, (i) => Icon(
                        i < master.rating.round()
                            ? Icons.star
                            : Icons.star_border,
                        color: Colors.amber,
                        size: 20,
                      )),
                      const SizedBox(width: 8),
                      Text(
                        master.rating.toStringAsFixed(1),
                        style: theme.textTheme.bodyLarge,
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 8),
                Center(
                  child: Chip(
                    label: Text(
                      master.level == 'опытный' ? 'Опытный мастер' : 'Новичок',
                    ),
                  ),
                ),
                const SizedBox(height: 24),
                Text('Программы', style: theme.textTheme.titleMedium),
                const SizedBox(height: 8),
                Wrap(
                  spacing: 8,
                  runSpacing: 4,
                  children: master.programIds.map((pid) => ActionChip(
                    label: Text('Программа #$pid'),
                    onPressed: () => context.push('/schedule/programs/$pid'),
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
