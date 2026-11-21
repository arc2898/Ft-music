import 'package:ft_music/blocs/settings_cubit/cubit/settings_cubit.dart';
import 'package:ft_music/core/theme/app_theme.dart';
import 'package:ft_music/l10n/app_localizations.dart';
import 'package:ft_music/screens/screen/home_views/setting_views/country_setting.dart';
import 'package:ft_music/screens/screen/home_views/setting_views/download_setting.dart';
import 'package:ft_music/screens/screen/home_views/setting_views/player_setting.dart';
import 'package:ft_music/screens/screen/home_views/setting_views/storage_setting.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

class SettingsView extends StatelessWidget {
  const SettingsView({super.key});

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    return BlocBuilder<SettingsCubit, SettingsState>(
      builder: (context, state) {
        return Scaffold(
          backgroundColor: Theme.of(context).scaffoldBackgroundColor,
          appBar: AppBar(
            backgroundColor: Colors.transparent,
            elevation: 0,
            surfaceTintColor: Colors.transparent,
            centerTitle: true,
            title: Text(l10n.settingsTitle),
          ),
          body: Center(
            child: ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 640),
              child: ListView(
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
                children: [
                  _SettingsSection(
                    children: [
                      _SettingsTile(
                        title: 'Appearance',
                        subtitle: 'Choose light, dark or follow system settings.',
                        icon: Icons.palette_outlined,
                        trailing: DropdownButtonHideUnderline(
                          child: DropdownButton<String>(
                            value: state.themeMode,
                            isDense: true,
                            borderRadius: BorderRadius.circular(12),
                            items: const [
                              DropdownMenuItem(value: 'system', child: Text('System')),
                              DropdownMenuItem(value: 'light', child: Text('Light')),
                              DropdownMenuItem(value: 'dark', child: Text('Dark')),
                            ],
                            onChanged: (value) {
                              if (value != null) {
                                context.read<SettingsCubit>().setThemeMode(value);
                              }
                            },
                          ),
                        ),
                        onTap: () {},
                      ),
                    ],
                  ),
                  _SettingsSection(
                    children: [
                      _SettingsTile(
                        title: l10n.settingsPlayer,
                        subtitle: l10n.settingsPlayerSubtitle,
                        icon: Icons.headphones_rounded,
                        onTap: () => _navigate(context, const PlayerSettings()),
                      ),
                      _SettingsTile(
                        title: l10n.settingsDownloads,
                        subtitle: l10n.settingsDownloadsSubtitle,
                        icon: Icons.folder_rounded,
                        onTap: () => _navigate(context, const DownloadSettings()),
                      ),
                    ],
                  ),
                  _SettingsSection(
                    children: [
                      _SettingsTile(
                        title: l10n.settingsLanguageCountry,
                        subtitle: l10n.settingsLanguageCountrySubtitle,
                        icon: Icons.language_rounded,
                        onTap: () => _navigate(context, const CountrySettings()),
                      ),
                      _SettingsTile(
                        title: l10n.settingsStorage,
                        subtitle: l10n.settingsStorageSubtitle,
                        icon: Icons.workspace_premium_rounded,
                        onTap: () => _navigate(context, const BackupSettings()),
                      ),
                    ],
                  ),
                  const SizedBox(height: 32),
                ],
              ),
            ),
          ),
        );
      },
    );
  }

  void _navigate(BuildContext context, Widget screen) {
    Navigator.push(
      context,
      MaterialPageRoute(builder: (context) => screen),
    );
  }
}

class _SettingsSection extends StatelessWidget {
  final List<Widget> children;

  const _SettingsSection({required this.children});

  @override
  Widget build(BuildContext context) {
    final colors = Theme.of(context).colorScheme;
    return Container(
      margin: const EdgeInsets.only(bottom: 18),
      decoration: BoxDecoration(
        color: colors.surfaceContainerLow,
        borderRadius: BorderRadius.circular(20),
      ),
      child: ClipRRect(
        borderRadius: BorderRadius.circular(20),
        child: Column(
          children: [
            for (var index = 0; index < children.length; index++) ...[
              children[index],
              if (index < children.length - 1)
                Divider(
                  height: 1,
                  indent: 66,
                  endIndent: 16,
                  color: colors.outlineVariant.withValues(alpha: .45),
                ),
            ],
          ],
        ),
      ),
    );
  }
}

class _SettingsTile extends StatelessWidget {
  final String title;
  final String subtitle;
  final IconData icon;
  final Widget? trailing;
  final VoidCallback onTap;

  const _SettingsTile({
    required this.title,
    required this.subtitle,
    required this.icon,
    this.trailing,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final colors = Theme.of(context).colorScheme;
    return InkWell(
      onTap: onTap,
      splashColor: colors.primary.withValues(alpha: .08),
      highlightColor: colors.primary.withValues(alpha: .04),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
        child: Row(
          children: [
            Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(
                color: colors.primary.withValues(alpha: .12),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Icon(icon, size: 20, color: colors.primary),
            ),
            const SizedBox(width: 14),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    title,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: Theme.of(context).textTheme.titleSmall?.copyWith(
                          color: colors.onSurface,
                          fontWeight: FontWeight.w700,
                        ),
                  ),
                  const SizedBox(height: 3),
                  Text(
                    subtitle,
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                    style: Theme.of(context).textTheme.bodySmall?.copyWith(
                          color: colors.onSurfaceVariant,
                        ),
                  ),
                ],
              ),
            ),
            const SizedBox(width: 12),
            trailing ??
                Icon(
                  Icons.chevron_right_rounded,
                  color: colors.onSurfaceVariant.withValues(alpha: .55),
                  size: 22,
                ),
          ],
        ),
      ),
    );
  }
}
