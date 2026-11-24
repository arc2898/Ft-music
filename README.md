# FT-music

**Built by [arc2898](https://github.com/arc2898), on top of [OpenTune](https://github.com/justshivv/Opentune).**

FT-music is a community fork of the GPL-3.0 OpenTune music player. The customized native builds use FT-music branding and retain the upstream project’s copyright, license, and third-party notices. This is not an official OpenTune release.

## What changed in this fork

- Android and Desktop app branding, About text, and package metadata identify FT-music and arc2898.
- Android includes selectable launcher-icon styles, including circle, diamond, and hexagon.
- Android downloads are allowed to use mobile data by default; Wi-Fi-only remains available in settings. Existing installs migrate once to this default, and queued downloads are re-enqueued.
- Discord presence/login and the in-app self-updater/update-lock flow have been removed from the Android app; Discord presence/settings have been removed from Desktop.
- Android package ID: `com.ftmusic.player.preview`; current preview version: `1.1.1` (code `48`).

The browser player/extension is retained from upstream and is not part of the native-app cleanup described above.

## Build Android

Use JDK 21 and an Android SDK with platform 37 installed. From the repository root:

```sh
./gradlew :app:testDebugUnitTest
./gradlew :app:assembleDebug
```

The universal debug APK is written under `app/build/outputs/apk/`. See [docs/building.md](docs/building.md) for the upstream build overview.

## Build Desktop Debian packages

```sh
./gradlew :desktop:packageDeb
./gradlew :desktop:packageReleaseDeb
```

The generated packages are under `desktop/build/compose/binaries/main/deb/`.

## Source and notices

The source archive includes the FT-music changes, upstream `LICENSE`, and `THIRD_PARTY_NOTICES.md`. Redistribution must follow the GNU General Public License v3.0 and the licenses of bundled third-party components.
