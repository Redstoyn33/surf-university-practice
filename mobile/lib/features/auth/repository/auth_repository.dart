import 'package:dio/dio.dart';
import '../../../core/models/client.dart';

class AuthRepository {
  final Dio _dio;

  AuthRepository(this._dio);

  Future<({Client client, String token})> login(
    String login,
    String password,
  ) async {
    final response = await _dio.post('/auth/login', data: {
      'login': login,
      'password': password,
    });
    return (
      client: Client.fromJson(response.data['client']),
      token: response.data['token'] as String,
    );
  }

  Future<({Client client, String token})> register(
    String login,
    String password,
  ) async {
    final response = await _dio.post('/auth/register', data: {
      'login': login,
      'password': password,
    });
    return (
      client: Client.fromJson(response.data['client']),
      token: response.data['token'] as String,
    );
  }
}
