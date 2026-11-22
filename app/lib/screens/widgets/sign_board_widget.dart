import 'package:flutter/material.dart';

class SignBoardWidget extends StatelessWidget {
  final String message;
  final IconData icon;

  const SignBoardWidget({
    super.key,
    required this.message,
    required this.icon,
  });

  @override
  Widget build(BuildContext context) {
    final colors = Theme.of(context).colorScheme;
    return Center(
      child: Wrap(
        children: [
          Column(
            children: [
              Padding(
                padding: const EdgeInsets.all(8.0),
                child: Icon(
                  icon,
                  color: colors.onSurfaceVariant.withValues(alpha: .75),
                  size: 40,
                ),
              ),
              Text(
                message,
                textAlign: TextAlign.center,
                style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                      color: colors.onSurfaceVariant,
                      fontSize: 14,
                    ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
