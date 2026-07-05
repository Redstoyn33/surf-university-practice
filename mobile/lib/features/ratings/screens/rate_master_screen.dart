import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../providers/rating_provider.dart';

class RateMasterScreen extends ConsumerStatefulWidget {
  final int masterId;
  final int slotId;
  final String masterName;
  final String programName;
  final String slotDate;

  const RateMasterScreen({
    super.key,
    required this.masterId,
    required this.slotId,
    required this.masterName,
    required this.programName,
    required this.slotDate,
  });

  @override
  ConsumerState<RateMasterScreen> createState() => _RateMasterScreenState();
}

class _RateMasterScreenState extends ConsumerState<RateMasterScreen> {
  int _score = 0;
  bool _isSubmitting = false;

  Future<void> _submit() async {
    if (_score == 0) return;
    setState(() => _isSubmitting = true);
    try {
      await ref.read(ratingNotifierProvider.notifier).createRating(
        masterId: widget.masterId,
        slotId: widget.slotId,
        score: _score,
      );
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Спасибо за оценку!')),
      );
      context.pop();
    } on Exception catch (e) {
      if (!mounted) return;
      final msg = e.toString();
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(
            msg.contains('422')
                ? 'Оценка доступна в окне 1–48 ч после завершения'
                : msg.contains('409')
                    ? 'Вы уже оценили мастера'
                    : 'Произошла ошибка',
          ),
        ),
      );
    } finally {
      if (mounted) setState(() => _isSubmitting = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(title: const Text('Оценка мастера')),
      body: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          children: [
            CircleAvatar(
              radius: 40,
              child: Icon(Icons.person, size: 40, color: theme.colorScheme.onPrimaryContainer),
            ),
            const SizedBox(height: 16),
            Text(widget.masterName, style: theme.textTheme.headlineSmall),
            const SizedBox(height: 4),
            Text(
              '${widget.programName} • ${widget.slotDate}',
              style: theme.textTheme.bodyMedium?.copyWith(
                color: theme.colorScheme.onSurfaceVariant,
              ),
            ),
            const SizedBox(height: 32),
            Text('Оцените мастера', style: theme.textTheme.titleLarge),
            const SizedBox(height: 16),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: List.generate(5, (i) {
                final starIdx = i + 1;
                return IconButton(
                  iconSize: 48,
                  icon: Icon(
                    starIdx <= _score ? Icons.star : Icons.star_border,
                    color: Colors.amber,
                  ),
                  onPressed: _isSubmitting
                      ? null
                      : () => setState(() => _score = starIdx),
                );
              }),
            ),
            const SizedBox(height: 32),
            SizedBox(
              width: double.infinity,
              height: 48,
              child: FilledButton(
                onPressed: _score == 0 || _isSubmitting ? null : _submit,
                child: _isSubmitting
                    ? const SizedBox(
                        width: 24,
                        height: 24,
                        child: CircularProgressIndicator(strokeWidth: 2),
                      )
                    : const Text('Отправить оценку'),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
