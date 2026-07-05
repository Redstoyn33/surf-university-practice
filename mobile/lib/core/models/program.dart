class Program {
  final int id;
  final String name;
  final String description;
  final int maxCapacity;
  final List<int> masterIds;

  const Program({
    required this.id,
    required this.name,
    required this.description,
    required this.maxCapacity,
    required this.masterIds,
  });

  factory Program.fromJson(Map<String, dynamic> json) => Program(
    id: json['id'] as int,
    name: json['name'] as String,
    description: json['description'] as String,
    maxCapacity: json['maxCapacity'] as int,
    masterIds: (json['masterIds'] as List).cast<int>(),
  );

  Map<String, dynamic> toJson() => {
    'id': id,
    'name': name,
    'description': description,
    'maxCapacity': maxCapacity,
    'masterIds': masterIds,
  };
}
