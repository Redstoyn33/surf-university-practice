import 'program.dart';
import 'master.dart';

class Slot {
  final int id;
  final String dateTime;
  final String endTime;
  final Program program;
  final Master master;
  final int totalSpots;
  final int availableSpots;
  final bool rentalAvailable;
  final double rentalPrice;

  const Slot({
    required this.id,
    required this.dateTime,
    required this.endTime,
    required this.program,
    required this.master,
    required this.totalSpots,
    required this.availableSpots,
    required this.rentalAvailable,
    required this.rentalPrice,
  });

  bool get isAvailable => availableSpots > 0;

  factory Slot.fromJson(Map<String, dynamic> json) => Slot(
    id: json['id'] as int,
    dateTime: json['dateTime'] as String,
    endTime: json['endTime'] as String,
    program: Program.fromJson(json['program'] as Map<String, dynamic>),
    master: Master.fromJson(json['master'] as Map<String, dynamic>),
    totalSpots: json['totalSpots'] as int,
    availableSpots: json['availableSpots'] as int,
    rentalAvailable: json['rentalAvailable'] as bool,
    rentalPrice: (json['rentalPrice'] as num).toDouble(),
  );

  Map<String, dynamic> toJson() => {
    'id': id,
    'dateTime': dateTime,
    'endTime': endTime,
    'program': program.toJson(),
    'master': master.toJson(),
    'totalSpots': totalSpots,
    'availableSpots': availableSpots,
    'rentalAvailable': rentalAvailable,
    'rentalPrice': rentalPrice,
  };
}
