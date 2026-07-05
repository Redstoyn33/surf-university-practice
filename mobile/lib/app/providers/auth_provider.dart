import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../core/storage.dart';
import '../../core/api_client.dart';
import '../../core/models/client.dart';

class AuthNotifier extends StateNotifier<AsyncValue<Client?>> {
  final LocalStorage _storage;
  final ApiClient _api;

  AuthNotifier(this._storage, this._api) : super(const AsyncValue.data(null));

  Future<void> checkSession() async {
    final token = await _storage.getToken();
    if (token != null) {
      final login = await _storage.getClientLogin();
      final id = await _storage.getClientId();
      if (login != null && id != null) {
        state = AsyncValue.data(Client(id: id, login: login, createdAt: ''));
      }
    }
  }

  Future<Client> login(String login, String password) async {
    state = const AsyncValue.loading();
    try {
      final response = await _api.dio.post('/auth/login', data: {
        'login': login,
        'password': password,
      });
      final client = Client.fromJson(response.data['client']);
      final token = response.data['token'] as String;
      await _storage.saveToken(token);
      await _storage.saveClientInfo(id: client.id, login: client.login);
      state = AsyncValue.data(client);
      return client;
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }

  Future<Client> register(String login, String password) async {
    state = const AsyncValue.loading();
    try {
      final response = await _api.dio.post('/auth/register', data: {
        'login': login,
        'password': password,
      });
      final client = Client.fromJson(response.data['client']);
      final token = response.data['token'] as String;
      await _storage.saveToken(token);
      await _storage.saveClientInfo(id: client.id, login: client.login);
      state = AsyncValue.data(client);
      return client;
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }

  Future<void> logout() async {
    await _storage.clearAll();
    state = const AsyncValue.data(null);
  }
}

final authNotifierProvider = StateNotifierProvider<AuthNotifier, AsyncValue<Client?>>((ref) {
  final storage = ref.read(localStorageProvider);
  final api = ref.read(apiClientProvider);
  return AuthNotifier(storage, api);
});
