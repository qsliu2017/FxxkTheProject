package com.example.ftpclient

import android.content.Intent
import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.view.View
import androidx.appcompat.app.AlertDialog
import kotlinx.android.synthetic.main.activity_main.*

class MainActivity : AppCompatActivity(), View.OnClickListener {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
        downloadBtn.setOnClickListener(this)
        loginBtn.setOnClickListener(this)
        disConnectBtn.setOnClickListener(this)
    }

    override fun onClick(v: View?) {
        when (v?.id) {
            R.id.downloadBtn -> {
                val intent = Intent(this, DownloadActivity::class.java)
                intent.putExtra("from", "main")
                startActivity(intent)
            }
            R.id.loginBtn -> {
                startActivity(Intent(this, LoginActivity::class.java))
            }
            R.id.disConnectBtn -> {
                AlertDialog.Builder(this).setMessage("Disconnect?")
                    .setPositiveButton("Yes"
                    ) { _, _ ->
                        Connection.getCon()?.logout()
                        startActivity(Intent(this, ConnectActivity::class.java))
                    }
                    .setNegativeButton("No", null).create().show()
            }
        }
    }
}