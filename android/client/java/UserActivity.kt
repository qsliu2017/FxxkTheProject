package com.example.ftpclient

import android.content.Intent
import android.os.Bundle
import android.view.View
import androidx.appcompat.app.AlertDialog
import androidx.appcompat.app.AppCompatActivity
import kotlinx.android.synthetic.main.activity_user.*

class UserActivity : AppCompatActivity(), View.OnClickListener {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_user)
        downloadBtn.setOnClickListener(this)
        uploadBtn.setOnClickListener(this)
        modeBtn.setOnClickListener(this)
        logoutBtn.setOnClickListener(this)
    }

    override fun onClick(v: View?) {
        when (v?.id) {
            R.id.downloadBtn -> {
                val intent = Intent(this, DownloadActivity::class.java)
                intent.putExtra("from", "user")
                startActivity(intent)
            }
            R.id.uploadBtn -> {
                startActivity(Intent(this, UploadActivity::class.java))
            }
            R.id.modeBtn -> {
                startActivity(Intent(this, ModeActivity::class.java))
            }
            R.id.logoutBtn -> {
                AlertDialog.Builder(this).setMessage("Log out?")
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