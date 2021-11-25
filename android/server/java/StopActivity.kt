package com.example.ftpserver

import android.content.Intent
import android.os.Bundle
import android.view.View
import androidx.appcompat.app.AlertDialog
import androidx.appcompat.app.AppCompatActivity
import kotlinx.android.synthetic.main.activity_stop.*
import server.Server

class StopActivity : AppCompatActivity(), View.OnClickListener {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_stop)
        stopBtn.setOnClickListener(this)
        Server.setLogger(MyServer.Logger(logText))
    }

    override fun onClick(v: View?) {
        when (v?.id) {
            R.id.stopBtn -> {
                AlertDialog.Builder(this).setMessage("Stop the server?")
                    .setPositiveButton(
                        "Yes"
                    ) { _, _ ->
                        // Stop the server
                        stopService(Intent(this, FtpService::class.java))
                        startActivity(Intent(this, StartActivity::class.java))
                    }
                    .setNegativeButton("No", null).create().show()
            }
        }
    }
}