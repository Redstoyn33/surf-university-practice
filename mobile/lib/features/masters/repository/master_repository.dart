import 'package:dio/dio.dart';
import '../../../core/models/master.dart';

class MasterRepository {
  final Dio _dio;

  MasterRepository(this._dio);

  Future<List<Master>> getMasters() async {
    final response = await _dio.get('/masters');
    return (response.data as List).map((e) => Master.fromJson(e)).toList();
  }

  Future<Master> getMasterById(int id) async {
    final response = await _dio.get('/masters/$id');
    return Master.fromJson(response.data);
  }
}
