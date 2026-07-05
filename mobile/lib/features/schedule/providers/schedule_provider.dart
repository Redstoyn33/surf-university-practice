import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/api_client.dart';
import '../../../core/models/slot.dart';
import '../repository/schedule_repository.dart';

final scheduleRepositoryProvider = Provider<ScheduleRepository>((ref) {
  final api = ref.read(apiClientProvider);
  return ScheduleRepository(api.dio);
});

class SlotsNotifier extends StateNotifier<AsyncValue<List<Slot>>> {
  final ScheduleRepository _repository;

  SlotsNotifier(this._repository) : super(const AsyncValue.data([]));

  Future<void> loadSlots({
    String? dateFrom,
    String? dateTo,
    int? masterId,
    int? programId,
  }) async {
    state = const AsyncValue.loading();
    try {
      final slots = await _repository.getSlots(
        dateFrom: dateFrom,
        dateTo: dateTo,
        masterId: masterId,
        programId: programId,
      );
      state = AsyncValue.data(slots);
    } catch (e, st) {
      state = AsyncValue.error(e, st);
    }
  }
}

final slotsNotifierProvider =
    StateNotifierProvider<SlotsNotifier, AsyncValue<List<Slot>>>((ref) {
  return SlotsNotifier(ref.read(scheduleRepositoryProvider));
});
