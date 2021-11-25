package com.example.ftpclient

import android.content.Intent
import android.os.Bundle
import android.view.View
import androidx.appcompat.app.AlertDialog
import androidx.appcompat.app.AppCompatActivity
import kotlinx.android.synthetic.main.activity_connect.*

class ConnectActivity : AppCompatActivity(), View.OnClickListener {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_connect)
        connectBtn.setOnClickListener(this)
    }

    override fun onClick(v: View?) {
        when (v?.id) {
            R.id.connectBtn -> {
                val address = serverAddr.text.toString()
                // Try to connect the server
                val result = Connection.setCon(address)

                if (result == 0) { // Establish a connection
                    startActivity(Intent(this, MainActivity::class.java))
                } else { // Invalid server address, or the server is not running
                    AlertDialog.Builder(this).setMessage("Fail to connect to the server")
                        .setPositiveButton("OK", null).create().show()
                }
            }
        }
    }
}