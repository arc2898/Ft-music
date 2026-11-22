import 'package:ft_music/screens/widgets/sign_board_widget.dart';
import 'package:flutter/material.dart';
import 'package:icons_plus/icons_plus.dart';

class DownloadsView extends StatelessWidget {
  const DownloadsView({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Theme.of(context).scaffoldBackgroundColor,
      appBar: AppBar(
        centerTitle: true,
        title: Text(
          'Downloads',
          style: Theme.of(context).textTheme.titleLarge?.copyWith(
            color: Theme.of(context).colorScheme.onSurface,
            fontWeight: FontWeight.bold,
          ),
        ),
      ),
      body: const Center(
        child: SignBoardWidget(
            message: "No Downloads Yet", icon: MingCute.download_2_fill),
      ),
    );
  }
}
