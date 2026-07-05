import 'package:dio/dio.dart';
import '../../../core/models/slot.dart';
import '../../../core/models/booking.dart';

class BookingRepository {
  final Dio _dio;

  BookingRepository(this._dio);

  Future<Slot> getSlotById(int id) async {
    final response = await _dio.get('/slots/$id');
    return Slot.fromJson(response.data);
  }

  Future<Booking> createBooking({
    required int slotId,
    required bool rentalSelected,
  }) async {
    final response = await _dio.post('/bookings', data: {
      'slotId': slotId,
      'rentalSelected': rentalSelected,
    });
    return Booking.fromJson(response.data);
  }

  Future<List<Booking>> getMyBookings({String? status}) async {
    final params = <String, dynamic>{};
    if (status != null) params['status'] = status;
    final response = await _dio.get('/bookings', queryParameters: params);
    return (response.data as List).map((e) => Booking.fromJson(e)).toList();
  }

  Future<Booking> getBookingById(int id) async {
    final response = await _dio.get('/bookings/$id');
    return Booking.fromJson(response.data);
  }

  Future<Booking> cancelBooking(int bookingId) async {
    final response = await _dio.patch('/bookings/$bookingId/cancel');
    return Booking.fromJson(response.data);
  }
}
