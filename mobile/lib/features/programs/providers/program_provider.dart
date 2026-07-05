import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/api_client.dart';
import '../../../core/models/program.dart';
import '../repository/program_repository.dart';

final programRepositoryProvider = Provider<ProgramRepository>((ref) {
  final api = ref.read(apiClientProvider);
  return ProgramRepository(api.dio);
});

final programsListProvider = FutureProvider<List<Program>>((ref) {
  return ref.read(programRepositoryProvider).getPrograms();
});

final programByIdProvider =
    FutureProvider.family<Program, int>((ref, id) {
  return ref.read(programRepositoryProvider).getProgramById(id);
});
