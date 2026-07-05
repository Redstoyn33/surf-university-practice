import 'package:flutter/material.dart';

class RentalSwitch extends StatelessWidget {
  final bool value;
  final double price;
  final ValueChanged<bool> onChanged;

  const RentalSwitch({
    super.key,
    required this.value,
    required this.price,
    required this.onChanged,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Card(
      child: SwitchListTile(
        title: const Text('Прокат оборудования'),
        subtitle: Text('+${price.toStringAsFixed(0)} ₽'),
        secondary: Icon(
          Icons.build,
          color: value ? theme.colorScheme.primary : null,
        ),
        value: value,
        onChanged: onChanged,
      ),
    );
  }
}
