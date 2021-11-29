package com.example.ftpserver

import android.content.Intent
import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import android.view.View
import androidx.appcompat.app.AlertDialog
import androidx.core.content.ContextCompat
import kotlinx.android.synthetic.main.activity_start.*
import server.Server

class StartActivity : AppCompatActivity(), View.OnClickListener {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_start)
        startBtn.setOnClickListener(this)
    }

    override fun onClick(v: View?) {
        when (v?.id) {
            R.id.startBtn -> {
                // Launch the server
                val portNum = port.text.toString().toLong()
                if (portNum < 1024) {
                    AlertDialog.Builder(this).setMessage("Port should be larger than 1024")
                        .setPositiveButton("OK", null).create().show()
                } else {
                    val intent = Intent(this, FtpService::class.java)
                    intent.putExtra("port", portNum)
                    startService(intent)
                    startActivity(Intent(this, StopActivity::class.java))
                }
            }
        }
    }
}