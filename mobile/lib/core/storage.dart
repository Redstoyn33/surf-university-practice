import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class LocalStorage {
  static const _tokenKey = 'auth_token';
  static const _clientIdKey = 'client_id';
  static const _clientLoginKey = 'client_login';

  final FlutterSecureStorage _storage;

  LocalStorage() : _storage = const FlutterSecureStorage();

  Future<void> saveToken(String token) =>
    _storage.write(key: _tokenKey, value: token);

  Future<String?> getToken() =>
    _storage.read(key: _tokenKey);

  Future<void> deleteToken() =>
    _storage.delete(key: _tokenKey);

  Future<void> saveClientInfo({required int id, required String login}) async {
    await _storage.write(key: _clientIdKey, value: id.toString());
    await _storage.write(key: _clientLoginKey, value: login);
  }

  Future<int?> getClientId() async {
    final v = await _storage.read(key: _clientIdKey);
    return v != null ? int.tryParse(v) : null;
  }

  Future<String?> getClientLogin() =>
    _storage.read(key: _clientLoginKey);

  Future<void> clearAll() async {
    await _storage.deleteAll();
  }
}
