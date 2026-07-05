import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../providers/auth_provider.dart';
import '../widgets/auth_form.dart';

class RegistrationScreen extends ConsumerStatefulWidget {
  const RegistrationScreen({super.key});

  @override
  ConsumerState<RegistrationScreen> createState() => _RegistrationScreenState();
}

class _RegistrationScreenState extends ConsumerState<RegistrationScreen> {
  final _formKey = GlobalKey<FormState>();
  final _loginController = TextEditingController();
  final _passwordController = TextEditingController();
  final _confirmController = TextEditingController();
  bool _obscurePassword = true;

  @override
  void dispose() {
    _loginController.dispose();
    _passwordController.dispose();
    _confirmController.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;
    ref.read(authNotifierProvider.notifier).register(
      _loginController.text.trim(),
      _passwordController.text,
    );
  }

  @override
  Widget build(BuildContext context) {
    final authState = ref.watch(authNotifierProvider);

    ref.listen(authNotifierProvider, (prev, next) {
      next.mapOrNull(
        error: (e) => ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(_errorMessage(e.error))),
        ),
      );
    });

    final isLoading = authState is AsyncLoading;

    return AuthForm(
      title: 'Глини',
      subtitle: 'Создайте аккаунт',
      buttonText: 'Зарегистрироваться',
      isLoading: isLoading,
      onSubmitted: _submit,
      fields: [
        Form(
          key: _formKey,
          child: Column(
            children: [
              TextFormField(
                controller: _loginController,
                decoration: const InputDecoration(
                  labelText: 'Логин',
                  prefixIcon: Icon(Icons.person),
                ),
                textInputAction: TextInputAction.next,
                validator: (v) => v == null || v.trim().isEmpty
                    ? 'Введите логин'
                    : null,
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _passwordController,
                decoration: InputDecoration(
                  labelText: 'Пароль',
                  prefixIcon: const Icon(Icons.lock),
                  suffixIcon: IconButton(
                    icon: Icon(
                      _obscurePassword
                          ? Icons.visibility_off
                          : Icons.visibility,
                    ),
                    onPressed: () =>
                        setState(() => _obscurePassword = !_obscurePassword),
                  ),
                ),
                obscureText: _obscurePassword,
                textInputAction: TextInputAction.next,
                validator: (v) {
                  if (v == null || v.isEmpty) return 'Введите пароль';
                  if (v.length < 6) return 'Пароль должен быть ≥ 6 символов';
                  return null;
                },
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _confirmController,
                decoration: const InputDecoration(
                  labelText: 'Подтверждение пароля',
                  prefixIcon: Icon(Icons.lock_outline),
                ),
                obscureText: true,
                textInputAction: TextInputAction.done,
                onFieldSubmitted: (_) => _submit(),
                validator: (v) {
                  if (v != _passwordController.text) {
                    return 'Пароли не совпадают';
                  }
                  return null;
                },
              ),
            ],
          ),
        ),
      ],
      footer: TextButton(
        onPressed: () => context.go('/login'),
        child: const Text('Уже есть аккаунт? Войти'),
      ),
    );
  }

  String _errorMessage(Object error) {
    final msg = error.toString();
    if (msg.contains('409') || msg.contains('уже занят')) {
      return 'Этот логин уже занят';
    }
    if (msg.contains('400')) return 'Проверьте введённые данные';
    return 'Произошла ошибка. Попробуйте позже.';
  }
}
