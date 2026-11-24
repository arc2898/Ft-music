# FT-music fork branding and build notes

FT-music is built by [arc2898](https://github.com/arc2898) on top of [OpenTune](https://github.com/justshivv/Opentune). It is a community fork, not an official OpenTune release and not affiliated with or endorsed by the upstream maintainers.

The Android app uses FT-music branding and the application ID `com.ftmusic.player.preview`. Its launcher chooser includes FT-music colorways plus circle, diamond, and hexagon variants. Downloads may use mobile data by default; the Wi-Fi-only preference remains available. Discord presence and upstream self-update prompts are removed from the Android and Desktop builds.

Build the universal Android preview with `./gradlew :app:assembleDebug` and run Android unit tests with `./gradlew :app:testDebugUnitTest` (JDK 21). Build Debian packages with `./gradlew :desktop:packageDeb` or `./gradlew :desktop:packageReleaseDeb`.

The OpenTune `LICENSE` (GPL-3.0), `THIRD_PARTY_NOTICES.md`, and upstream copyright notices are retained. Redistribution must comply with those licenses and applicable corresponding-source obligations.
