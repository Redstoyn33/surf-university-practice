import 'package:dio/dio.dart';
import '../../../core/models/program.dart';

class ProgramRepository {
  final Dio _dio;

  ProgramRepository(this._dio);

  Future<List<Program>> getPrograms() async {
    final response = await _dio.get('/programs');
    return (response.data as List).map((e) => Program.fromJson(e)).toList();
  }

  Future<Program> getProgramById(int id) async {
    final response = await _dio.get('/programs/$id');
    return Program.fromJson(response.data);
  }
}
