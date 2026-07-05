import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'storage.dart';
import 'models/api_error.dart';

class AuthInterceptor extends Interceptor {
  final LocalStorage _storage;

  AuthInterceptor(this._storage);

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) async {
    final token = await _storage.getToken();
    if (token != null) {
      options.headers['Authorization'] = 'Bearer $token';
    }
    handler.next(options);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    if (err.response?.statusCode == 401) {
      await _storage.clearAll();
    }
    handler.next(err);
  }
}

class ErrorInterceptor extends Interceptor {
  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    final apiError = _parseError(err);
    handler.reject(DioException(
      requestOptions: err.requestOptions,
      response: err.response,
      type: err.type,
      error: apiError,
      message: apiError.message,
    ));
  }

  ApiError _parseError(DioException err) {
    if (err.response?.data is Map) {
      final data = err.response!.data as Map;
      if (data['error'] is String) {
        return ApiError(message: data['error'] as String);
      }
    }
    switch (err.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.receiveTimeout:
        return ApiError(message: 'Сервер не отвечает. Попробуйте позже.');
      case DioExceptionType.connectionError:
        return ApiError(message: 'Нет соединения с интернетом.');
      default:
        return ApiError(message: 'Произошла ошибка. Попробуйте позже.');
    }
  }
}

class ApiClient {
  late final Dio dio;

  ApiClient(LocalStorage storage) {
    dio = Dio(BaseOptions(
      baseUrl: 'http://10.0.2.2:8080',
      connectTimeout: const Duration(seconds: 10),
      receiveTimeout: const Duration(seconds: 10),
      headers: {'Content-Type': 'application/json'},
    ));
    dio.interceptors.addAll([
      AuthInterceptor(storage),
      ErrorInterceptor(),
    ]);
  }
}

final localStorageProvider = Provider<LocalStorage>((_) => LocalStorage());

final apiClientProvider = Provider<ApiClient>((ref) {
  final storage = ref.read(localStorageProvider);
  return ApiClient(storage);
});
