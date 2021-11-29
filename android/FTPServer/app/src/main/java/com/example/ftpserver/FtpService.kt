package com.example.ftpserver

import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.PendingIntent
import android.app.PendingIntent.FLAG_MUTABLE
import android.app.Service
import android.content.Context
import android.content.Intent
import android.os.Build
import android.os.IBinder
import androidx.appcompat.app.AlertDialog
import androidx.core.app.NotificationCompat
import androidx.core.content.ContextCompat
import java.io.File

class FtpService : Service() {

    override fun onBind(intent: Intent): IBinder? {
        return null
    }

    override fun onCreate() {
        super.onCreate()
        val manager = getSystemService(Context.NOTIFICATION_SERVICE) as NotificationManager
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val channel = NotificationChannel(
                "ftp_service", "service notification",
                NotificationManager.IMPORTANCE_DEFAULT
            )
            manager.createNotificationChannel(channel)
        }
        val intent = Intent(this, StartActivity::class.java)
        val pi = PendingIntent.getActivity(this, 0, intent, FLAG_MUTABLE)
        val notification = NotificationCompat.Builder(this, "ftp_service")
            .setContentIntent(pi).build()
        startForeground(1, notification)
        val path = ContextCompat.getExternalFilesDirs(this, null)[0].toString()
        val file = File(path)
        if (!file.exists())
            file.createNewFile()
    }

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        val port = intent?.getLongExtra("port", 5000)
        val result = port?.let {
            MyServer.startServer(
                it,
                ContextCompat.getExternalFilesDirs(this, null)[0].toString()
            )
        }
        if (result == -1) {
            AlertDialog.Builder(this).setMessage("Fail to connect to the server")
                .setPositiveButton("OK", null).create().show()
        }
        return super.onStartCommand(intent, flags, startId)
    }

    override fun onDestroy() {
        MyServer.stopServer()
        super.onDestroy()
    }
}
