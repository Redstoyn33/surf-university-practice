import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/models/booking.dart';
import '../repository/booking_repository.dart';
import 'booking_provider.dart';

final myBookingsProvider =
    StateNotifierProvider<MyBookingsNotifier, AsyncValue<List<Booking>>>(
  (ref) => MyBookingsNotifier(ref.read(bookingRepositoryProvider)),
);

class MyBookingsNotifier extends StateNotifier<AsyncValue<List<Booking>>> {
  final BookingRepository _repository;

  MyBookingsNotifier(this._repository) : super(const AsyncValue.data([]));

  Future<void> load({String? status}) async {
    state = const AsyncValue.loading();
    try {
      final bookings = await _repository.getMyBookings(status: status);
      state = AsyncValue.data(bookings);
    } catch (e, st) {
      state = AsyncValue.error(e, st);
    }
  }

  Future<void> refresh() async {
    await load();
  }
}
