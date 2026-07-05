import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../providers/master_provider.dart';
import '../../../shared/widgets/error_state.dart';
import '../../../shared/widgets/empty_state.dart';

class MastersScreen extends ConsumerWidget {
  const MastersScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final mastersAsync = ref.watch(mastersListProvider);
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(title: const Text('Мастера')),
      body: RefreshIndicator(
        onRefresh: () => ref.refresh(mastersListProvider.future),
        child: mastersAsync.when(
          loading: () => const Center(child: CircularProgressIndicator()),
          error: (e, _) => ErrorState(
            message: 'Не удалось загрузить список мастеров',
            onRetry: () => ref.invalidate(mastersListProvider),
          ),
          data: (masters) {
            if (masters.isEmpty) {
              return const EmptyState(
                icon: Icons.people_outline,
                message: 'Нет доступных мастеров',
              );
            }
            return ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: masters.length,
              itemBuilder: (_, i) {
                final master = masters[i];
                return Padding(
                  padding: const EdgeInsets.only(bottom: 12),
                  child: Card(
                    child: InkWell(
                      borderRadius: BorderRadius.circular(12),
                      onTap: () => context.push('/masters/${master.id}'),
                      child: Padding(
                        padding: const EdgeInsets.all(16),
                        child: Row(
                          children: [
                            CircleAvatar(
                              radius: 28,
                              backgroundImage: NetworkImage(master.photo),
                              onBackgroundImageError: (_, __) => Icon(
                                Icons.person,
                                size: 28,
                                color: theme.colorScheme.onSurfaceVariant,
                              ),
                            ),
                            const SizedBox(width: 16),
                            Expanded(
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Text(
                                    master.name,
                                    style: theme.textTheme.titleMedium,
                                  ),
                                  const SizedBox(height: 4),
                                  Row(
                                    children: [
                                      ...List.generate(5, (j) => Icon(
                                        j < master.rating.round()
                                            ? Icons.star
                                            : Icons.star_border,
                                        color: Colors.amber,
                                        size: 16,
                                      )),
                                      const SizedBox(width: 4),
                                      Text(
                                        master.rating.toStringAsFixed(1),
                                        style: theme.textTheme.bodySmall,
                                      ),
                                    ],
                                  ),
                                ],
                              ),
                            ),
                            Chip(
                              label: Text(
                                master.level == 'опытный'
                                    ? 'Опытный'
                                    : 'Новичок',
                                style: theme.textTheme.labelSmall,
                              ),
                              visualDensity: VisualDensity.compact,
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
                );
              },
            );
          },
        ),
      ),
    );
  }
}
