class Client {
  final int id;
  final String login;
  final String createdAt;

  const Client({
    required this.id,
    required this.login,
    required this.createdAt,
  });

  factory Client.fromJson(Map<String, dynamic> json) => Client(
    id: json['id'] as int,
    login: json['login'] as String,
    createdAt: json['createdAt'] as String,
  );

  Map<String, dynamic> toJson() => {
    'id': id,
    'login': login,
    'createdAt': createdAt,
  };
}
