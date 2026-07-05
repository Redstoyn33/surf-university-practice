class ApiError {
  final String message;

  const ApiError({required this.message});

  @override
  String toString() => message;
}
