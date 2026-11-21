import 'dart:developer';

import 'package:ft_music/core/di/service_locator.dart';
import 'package:ft_music/plugins/blocs/chart/chart_bloc.dart';
import 'package:ft_music/plugins/blocs/chart/chart_event.dart';
import 'package:ft_music/plugins/blocs/chart/chart_state.dart';
import 'package:ft_music/plugins/blocs/plugin/plugin_bloc.dart';
import 'package:ft_music/plugins/blocs/plugin/plugin_state.dart';
import 'package:ft_music/screens/screen/chart/chart_view.dart';
import 'package:ft_music/screens/screen/chart/chart_widget.dart';
import 'package:carousel_slider/carousel_slider.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:responsive_framework/responsive_framework.dart';

/// Home carousel for chart-provider artwork. Chart details are fetched only
/// when a card is opened, keeping the initial home load lightweight.
class CaraouselWidget extends StatefulWidget {
  const CaraouselWidget({super.key});

  @override
  State<CaraouselWidget> createState() => _CaraouselWidgetState();
}

class _CaraouselWidgetState extends State<CaraouselWidget> {
  late final ChartBloc _chartBloc;

  @override
  void initState() {
    super.initState();
    _chartBloc = ChartBloc(pluginService: ServiceLocator.pluginService);
    _loadChartsFromPlugin();
  }

  void _loadChartsFromPlugin() {
    final chartProviders = context.read<PluginBloc>().state.loadedChartProviders;
    if (chartProviders.isNotEmpty) {
      final pluginId = chartProviders.first.manifest.id;
      log('Loading charts from plugin: $pluginId', name: 'ChartCarousel');
      _chartBloc.add(LoadCharts(pluginId: pluginId));
    } else {
      log('No chart provider plugins loaded — charts unavailable',
          name: 'ChartCarousel');
    }
  }

  @override
  void dispose() {
    _chartBloc.close();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final colors = Theme.of(context).colorScheme;
    return BlocBuilder<PluginBloc, PluginState>(
      builder: (context, pluginState) {
        if (pluginState.loadedChartProviders.isEmpty) {
          _chartBloc.add(const ClearCharts());
        } else if (_chartBloc.state.chartsStatus == ChartStatus.initial) {
          _loadChartsFromPlugin();
        }

        return BlocBuilder<ChartBloc, ChartState>(
          bloc: _chartBloc,
          builder: (context, chartState) {
            final visibleCharts = chartState.charts;
            if (visibleCharts.isEmpty) return const SizedBox.shrink();

            final isCompact = ResponsiveBreakpoints.of(context).isMobile ||
                ResponsiveBreakpoints.of(context).isTablet;
            final compactHeight =
                (MediaQuery.sizeOf(context).height * .32).clamp(220.0, 300.0);

            return Padding(
              padding: const EdgeInsets.fromLTRB(18, 20, 18, 0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    'Charts for you',
                    style: Theme.of(context).textTheme.titleLarge?.copyWith(
                          color: colors.onSurface,
                          fontWeight: FontWeight.w800,
                          letterSpacing: -.3,
                        ),
                  ),
                  const SizedBox(height: 3),
                  Text(
                    'Fresh artwork from your chart sources',
                    style: Theme.of(context).textTheme.bodySmall?.copyWith(
                          color: colors.onSurfaceVariant,
                        ),
                  ),
                  const SizedBox(height: 10),
                  CarouselSlider.builder(
                    itemCount: visibleCharts.length,
                    itemBuilder: (context, index, realIndex) {
                      final chart = visibleCharts[index];
                      return Semantics(
                        button: true,
                        label: 'Open ${chart.title} chart',
                        child: InkWell(
                          borderRadius: BorderRadius.circular(18),
                          onTap: () {
                            Navigator.push(
                              context,
                              MaterialPageRoute(
                                builder: (_) => ChartScreen(
                                  pluginId: _chartBloc.state.activePluginId ?? '',
                                  chartId: chart.id,
                                  chartTitle: chart.title,
                                ),
                              ),
                            );
                          },
                          child: ChartWidget(
                            chart: chart,
                            pluginId: _chartBloc.state.activePluginId ?? '',
                          ),
                        ),
                      );
                    },
                    options: CarouselOptions(
                      height: isCompact ? compactHeight.toDouble() : 250,
                      viewportFraction: ResponsiveBreakpoints.of(context).isMobile
                          ? .68
                          : ResponsiveBreakpoints.of(context).isTablet
                              ? .42
                              : .30,
                      autoPlay: visibleCharts.length > 1,
                      autoPlayInterval: const Duration(seconds: 4),
                      enableInfiniteScroll: visibleCharts.length > 1,
                      enlargeFactor: .14,
                      initialPage: 0,
                      pauseAutoPlayOnTouch: true,
                      padEnds: true,
                      enlargeCenterPage: ResponsiveBreakpoints.of(context).isMobile,
                    ),
                  ),
                ],
              ),
            );
          },
        );
      },
    );
  }
}
