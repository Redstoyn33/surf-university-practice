import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/models/client.dart';
import '../../../core/storage.dart';
import '../../../core/api_client.dart';
import '../repository/auth_repository.dart';

class AuthNotifier extends StateNotifier<AsyncValue<Client?>> {
  final AuthRepository _repository;
  final LocalStorage _storage;

  AuthNotifier(this._repository, this._storage)
    : super(const AsyncValue.data(null));

  Future<void> checkSession() async {
    final token = await _storage.getToken();
    if (token != null) {
      final login = await _storage.getClientLogin();
      final id = await _storage.getClientId();
      if (login != null && id != null) {
        state = AsyncValue.data(Client(id: id, login: login, createdAt: ''));
        return;
      }
    }
    state = const AsyncValue.data(null);
  }

  Future<void> login(String login, String password) async {
    state = const AsyncValue.loading();
    try {
      final result = await _repository.login(login, password);
      await _storage.saveToken(result.token);
      await _storage.saveClientInfo(
        id: result.client.id,
        login: result.client.login,
      );
      state = AsyncValue.data(result.client);
    } catch (e, st) {
      state = AsyncValue.error(e, st);
    }
  }

  Future<void> register(String login, String password) async {
    state = const AsyncValue.loading();
    try {
      final result = await _repository.register(login, password);
      await _storage.saveToken(result.token);
      await _storage.saveClientInfo(
        id: result.client.id,
        login: result.client.login,
      );
      state = AsyncValue.data(result.client);
    } catch (e, st) {
      state = AsyncValue.error(e, st);
    }
  }

  Future<void> logout() async {
    await _storage.clearAll();
    state = const AsyncValue.data(null);
  }
}

final authNotifierProvider =
    StateNotifierProvider<AuthNotifier, AsyncValue<Client?>>((ref) {
  final storage = ref.read(localStorageProvider);
  final api = ref.read(apiClientProvider);
  return AuthNotifier(AuthRepository(api.dio), storage);
});
