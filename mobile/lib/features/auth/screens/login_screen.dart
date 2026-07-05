import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../providers/auth_provider.dart';
import '../widgets/auth_form.dart';

class LoginScreen extends ConsumerStatefulWidget {
  const LoginScreen({super.key});

  @override
  ConsumerState<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends ConsumerState<LoginScreen> {
  final _formKey = GlobalKey<FormState>();
  final _loginController = TextEditingController();
  final _passwordController = TextEditingController();
  bool _obscurePassword = true;

  @override
  void dispose() {
    _loginController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;
    ref.read(authNotifierProvider.notifier).login(
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
      subtitle: 'Войдите в аккаунт',
      buttonText: 'Войти',
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
                textInputAction: TextInputAction.done,
                onFieldSubmitted: (_) => _submit(),
                validator: (v) =>
                    v == null || v.isEmpty ? 'Введите пароль' : null,
              ),
            ],
          ),
        ),
      ],
      footer: TextButton(
        onPressed: () => context.go('/register'),
        child: const Text('Нет аккаунта? Зарегистрироваться'),
      ),
    );
  }

  String _errorMessage(Object error) {
    final msg = error.toString();
    if (msg.contains('401') || msg.contains('неверный логин')) {
      return 'Неверный логин или пароль';
    }
    if (msg.contains('400')) return 'Проверьте введённые данные';
    return 'Произошла ошибка. Попробуйте позже.';
  }
}
