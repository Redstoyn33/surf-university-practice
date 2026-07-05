import 'package:dio/dio.dart';
import '../../../core/models/slot.dart';

class ScheduleRepository {
  final Dio _dio;

  ScheduleRepository(this._dio);

  Future<List<Slot>> getSlots({
    String? dateFrom,
    String? dateTo,
    int? masterId,
    int? programId,
  }) async {
    final params = <String, dynamic>{};
    if (dateFrom != null) params['dateFrom'] = dateFrom;
    if (dateTo != null) params['dateTo'] = dateTo;
    if (masterId != null) params['masterId'] = masterId;
    if (programId != null) params['programId'] = programId;

    final response = await _dio.get('/slots', queryParameters: params);
    return (response.data as List).map((e) => Slot.fromJson(e)).toList();
  }
}
