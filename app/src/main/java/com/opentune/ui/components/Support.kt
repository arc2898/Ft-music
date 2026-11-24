package com.opentune.ui.components

import android.content.Context
import android.content.Intent
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.rounded.Star
import androidx.compose.material3.MaterialTheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.unit.dp
import androidx.core.content.edit
import androidx.core.net.toUri
import java.util.concurrent.TimeUnit

/** GitHub project link and the optional once-a-month source-star invitation. */
object Support {
    const val REPO = "https://github.com/arc2898/Ft-music"
    enum class Invite(internal val daysBefore: Int, internal val askedKey: String, internal val doneKey: String) {
        STAR(5, "askedAt", "done"),
    }
    private val MIN_GAP_MS = TimeUnit.DAYS.toMillis(30)
    private fun prefs(context: Context) = context.getSharedPreferences("support", Context.MODE_PRIVATE)

    fun usedToday(context: Context, now: Long = System.currentTimeMillis()) {
        val p = prefs(context)
        val day = TimeUnit.MILLISECONDS.toDays(now)
        if (p.getLong("lastDay", -1) == day) return
        p.edit { putLong("lastDay", day); putInt("days", p.getInt("days", 0) + 1) }
    }

    fun due(context: Context, now: Long = System.currentTimeMillis()): Invite? {
        if (!com.opentune.data.config.RemoteConfig.current.communityPrompts) return null
        val p = prefs(context)
        val invite = Invite.STAR
        return invite.takeIf {
            shouldAsk(p.getBoolean(it.doneKey, false), p.getInt("days", 0), p.getLong(it.askedKey, 0), now, it.daysBefore)
        }
    }

    internal fun shouldAsk(done: Boolean, days: Int, askedAt: Long, now: Long, daysBefore: Int = Invite.STAR.daysBefore): Boolean =
        !done && days >= daysBefore && now - askedAt >= MIN_GAP_MS

    fun answered(context: Context, invite: Invite, never: Boolean, now: Long = System.currentTimeMillis()) = prefs(context).edit {
        putLong(invite.askedKey, now)
        if (never) putBoolean(invite.doneKey, true)
    }
}

@Composable
fun InviteCard(invite: Support.Invite, onDismiss: () -> Unit) {
    val context = LocalContext.current
    androidx.compose.runtime.LaunchedEffect(invite) { Support.answered(context, invite, never = false) }
    FloatingDialog(
        onDismissRequest = onDismiss,
        icon = { com.opentune.ui.theme.Icon(Icons.Rounded.Star, null) },
        title = { com.opentune.ui.theme.Text("Made for music. Built by arc2898.") },
        text = {
            Column {
                com.opentune.ui.theme.Text("FT-music is a free, open-source fork built on OpenTune. Visit its GitHub project and follow development there.")
                androidx.compose.material3.FilledTonalButton(
                    onClick = {
                        runCatching { context.startActivity(Intent(Intent.ACTION_VIEW, Support.REPO.toUri())) }
                        onDismiss()
                    },
                    modifier = Modifier.fillMaxWidth().padding(top = 8.dp),
                ) {
                    com.opentune.ui.theme.Icon(Icons.Rounded.Star, null)
                    com.opentune.ui.theme.Text("FT-music on GitHub", Modifier.padding(start = 8.dp))
                }
                androidx.compose.material3.TextButton(onClick = { Support.answered(context, invite, never = true); onDismiss() }) {
                    com.opentune.ui.theme.Text("Don't show again")
                }
            }
        },
        confirmButton = { SheetButton("Maybe later", onClick = onDismiss, closes = true) },
    )
}
