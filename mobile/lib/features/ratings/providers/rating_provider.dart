import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/api_client.dart';
import '../../../core/models/rating.dart';
import '../repository/rating_repository.dart';

final ratingRepositoryProvider = Provider<RatingRepository>((ref) {
  final api = ref.read(apiClientProvider);
  return RatingRepository(api.dio);
});

class RatingNotifier extends StateNotifier<AsyncValue<Rating?>> {
  final RatingRepository _repository;

  RatingNotifier(this._repository) : super(const AsyncValue.data(null));

  Future<Rating> createRating({
    required int masterId,
    required int slotId,
    required int score,
  }) async {
    state = const AsyncValue.loading();
    try {
      final rating = await _repository.createRating(
        masterId: masterId,
        slotId: slotId,
        score: score,
      );
      state = AsyncValue.data(rating);
      return rating;
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }
}

final ratingNotifierProvider =
    StateNotifierProvider<RatingNotifier, AsyncValue<Rating?>>((ref) {
  return RatingNotifier(ref.read(ratingRepositoryProvider));
});
