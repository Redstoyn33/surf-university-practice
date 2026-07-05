import 'package:flutter/material.dart';

class AuthForm extends StatelessWidget {
  final String title;
  final String subtitle;
  final List<Widget> fields;
  final String buttonText;
  final bool isLoading;
  final VoidCallback? onSubmitted;
  final Widget? footer;

  const AuthForm({
    super.key,
    required this.title,
    required this.subtitle,
    required this.fields,
    required this.buttonText,
    this.isLoading = false,
    this.onSubmitted,
    this.footer,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Scaffold(
      body: SafeArea(
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(24),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Icon(Icons.weekend, size: 80, color: theme.colorScheme.primary),
                const SizedBox(height: 16),
                Text(title, style: theme.textTheme.headlineMedium),
                const SizedBox(height: 4),
                Text(subtitle, style: theme.textTheme.bodyMedium?.copyWith(
                  color: theme.colorScheme.onSurfaceVariant,
                )),
                const SizedBox(height: 32),
                ...fields,
                const SizedBox(height: 16),
                SizedBox(
                  width: double.infinity,
                  height: 48,
                  child: FilledButton(
                    onPressed: isLoading ? null : onSubmitted,
                    child: isLoading
                        ? const SizedBox(
                            width: 24,
                            height: 24,
                            child: CircularProgressIndicator(strokeWidth: 2),
                          )
                        : Text(buttonText),
                  ),
                ),
                if (footer != null) ...[
                  const SizedBox(height: 16),
                  footer!,
                ],
              ],
            ),
          ),
        ),
      ),
    );
  }
}
