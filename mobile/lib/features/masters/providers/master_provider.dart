import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/api_client.dart';
import '../../../core/models/master.dart';
import '../repository/master_repository.dart';

final masterRepositoryProvider = Provider<MasterRepository>((ref) {
  final api = ref.read(apiClientProvider);
  return MasterRepository(api.dio);
});

final mastersListProvider = FutureProvider<List<Master>>((ref) {
  return ref.read(masterRepositoryProvider).getMasters();
});

final masterByIdProvider =
    FutureProvider.family<Master, int>((ref, id) {
  return ref.read(masterRepositoryProvider).getMasterById(id);
});
