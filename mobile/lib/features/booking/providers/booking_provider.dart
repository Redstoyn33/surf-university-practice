import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/api_client.dart';
import '../../../core/models/booking.dart';
import '../../../core/models/slot.dart';
import '../repository/booking_repository.dart';

final bookingRepositoryProvider = Provider<BookingRepository>((ref) {
  final api = ref.read(apiClientProvider);
  return BookingRepository(api.dio);
});

final slotByIdProvider = FutureProvider.family<Slot, int>((ref, id) {
  return ref.read(bookingRepositoryProvider).getSlotById(id);
});

class BookingNotifier extends StateNotifier<AsyncValue<List<Booking>>> {
  final BookingRepository _repository;

  BookingNotifier(this._repository) : super(const AsyncValue.data([]));

  Future<void> loadMyBookings({String? status}) async {
    state = const AsyncValue.loading();
    try {
      final bookings = await _repository.getMyBookings(status: status);
      state = AsyncValue.data(bookings);
    } catch (e, st) {
      state = AsyncValue.error(e, st);
    }
  }

  Future<Booking> createBooking({
    required int slotId,
    required bool rentalSelected,
  }) async {
    final booking = await _repository.createBooking(
      slotId: slotId,
      rentalSelected: rentalSelected,
    );
    return booking;
  }

  Future<Booking> cancelBooking(int bookingId) async {
    final booking = await _repository.cancelBooking(bookingId);
    return booking;
  }
}

final bookingNotifierProvider =
    StateNotifierProvider<BookingNotifier, AsyncValue<List<Booking>>>((ref) {
  return BookingNotifier(ref.read(bookingRepositoryProvider));
});
