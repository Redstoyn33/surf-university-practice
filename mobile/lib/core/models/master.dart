class Master {
  final int id;
  final String name;
  final String photo;
  final double rating;
  final String level;
  final List<int> programIds;

  const Master({
    required this.id,
    required this.name,
    required this.photo,
    required this.rating,
    required this.level,
    required this.programIds,
  });

  factory Master.fromJson(Map<String, dynamic> json) => Master(
    id: json['id'] as int,
    name: json['name'] as String,
    photo: json['photo'] as String,
    rating: (json['rating'] as num).toDouble(),
    level: json['level'] as String,
    programIds: (json['programIds'] as List).cast<int>(),
  );

  Map<String, dynamic> toJson() => {
    'id': id,
    'name': name,
    'photo': photo,
    'rating': rating,
    'level': level,
    'programIds': programIds,
  };
}
