import 'slot.dart';

class Booking {
  final int id;
  final int clientId;
  final Slot slot;
  final String status;
  final bool rentalSelected;
  final String createdAt;
  final String? cancellationReason;

  const Booking({
    required this.id,
    required this.clientId,
    required this.slot,
    required this.status,
    required this.rentalSelected,
    required this.createdAt,
    this.cancellationReason,
  });

  bool get isActive => status == 'активна';
  bool get isCancelledByClient => status == 'отменена клиентом';
  bool get isCancelledByWorkshop => status == 'отменена мастерской';

  factory Booking.fromJson(Map<String, dynamic> json) => Booking(
    id: json['id'] as int,
    clientId: json['clientId'] as int,
    slot: Slot.fromJson(json['slot'] as Map<String, dynamic>),
    status: json['status'] as String,
    rentalSelected: json['rentalSelected'] as bool,
    createdAt: json['createdAt'] as String,
    cancellationReason: json['cancellationReason'] as String?,
  );

  Map<String, dynamic> toJson() => {
    'id': id,
    'clientId': clientId,
    'slot': slot.toJson(),
    'status': status,
    'rentalSelected': rentalSelected,
    'createdAt': createdAt,
    'cancellationReason': cancellationReason,
  };
}
