package com.example.ftpserver

import android.content.Intent
import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.view.View
import kotlinx.android.synthetic.main.activity_start.*

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
                val intent = Intent(this, FtpService::class.java)
                intent.putExtra("port", portNum)
                startService(intent)
                startActivity(Intent(this, StopActivity::class.java))
            }
        }
    }
}