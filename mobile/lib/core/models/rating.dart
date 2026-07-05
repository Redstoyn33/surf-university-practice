class Rating {
  final int id;
  final int clientId;
  final int masterId;
  final int slotId;
  final int score;
  final String createdAt;

  const Rating({
    required this.id,
    required this.clientId,
    required this.masterId,
    required this.slotId,
    required this.score,
    required this.createdAt,
  });

  factory Rating.fromJson(Map<String, dynamic> json) => Rating(
    id: json['id'] as int,
    clientId: json['clientId'] as int,
    masterId: json['masterId'] as int,
    slotId: json['slotId'] as int,
    score: json['score'] as int,
    createdAt: json['createdAt'] as String,
  );

  Map<String, dynamic> toJson() => {
    'id': id,
    'clientId': clientId,
    'masterId': masterId,
    'slotId': slotId,
    'score': score,
    'createdAt': createdAt,
  };
}
