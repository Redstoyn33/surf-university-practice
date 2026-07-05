import 'package:dio/dio.dart';
import '../../../core/models/rating.dart';

class RatingRepository {
  final Dio _dio;

  RatingRepository(this._dio);

  Future<Rating> createRating({
    required int masterId,
    required int slotId,
    required int score,
  }) async {
    final response = await _dio.post('/ratings', data: {
      'masterId': masterId,
      'slotId': slotId,
      'score': score,
    });
    return Rating.fromJson(response.data);
  }
}
