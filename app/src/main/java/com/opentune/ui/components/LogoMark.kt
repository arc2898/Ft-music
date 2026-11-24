package com.opentune.ui.components

import androidx.compose.ui.geometry.CornerRadius
import androidx.compose.ui.geometry.RoundRect
import androidx.compose.ui.graphics.Path
import androidx.compose.ui.graphics.vector.PathParser

/** FT-music's logo mark, using the 512 by 512 viewBox of brand/icon.svg. */
internal const val MARK_VIEWPORT = 512f

private val FT_MARK_PATHS = listOf(
    "M168 150q3-15 18-15h151l-30 42H207l-7 40h101l-28 40h-80l-21 111-55 34 44-252z",
    "M284 214h128l-30 42h-68l-30 121-57 36 57-199z",
    "M213 346c-18 0-34 14-34 32s16 31 34 31c19 0 34-15 34-34v-29z",
    "M324 102c39-36 92-43 124-17-4 32-25 61-49 82 14-26 20-50 14-67-22-12-56-10-89 2z",
    "M190 405c-37 29-91 34-123 8 5-28 21-51 43-70-12 25-16 45-10 59 21 12 55 12 90 3z",
)

private val FT_WAVE_BARS = listOf(
    floatArrayOf(79f, 245f, 94f, 259f),
    floatArrayOf(99f, 228f, 113f, 276f),
    floatArrayOf(119f, 208f, 133f, 296f),
    floatArrayOf(379f, 245f, 394f, 259f),
    floatArrayOf(399f, 228f, 413f, 276f),
    floatArrayOf(419f, 208f, 433f, 296f),
)

/** The FT monogram, orbit and rounded sound-wave bars as a single filled path. */
internal fun markPath(): Path = Path().apply {
    FT_MARK_PATHS.forEach { data ->
        addPath(PathParser().parsePathString(data).toPath())
    }
    FT_WAVE_BARS.forEach { (left, top, right, bottom) ->
        addRoundRect(RoundRect(left, top, right, bottom, CornerRadius(7f, 7f)))
    }
}
