package com.example.ftpclient

import android.content.Intent
import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.view.View
import androidx.appcompat.app.AlertDialog
import kotlinx.android.synthetic.main.activity_information.*

class InformationActivity : AppCompatActivity(), View.OnClickListener {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_information)
        okBtn.setOnClickListener(this)

        try {
            connModeInfo.text = Connection.getCon()?.connMode.toString()
            modeInfo.text = Connection.getCon()?.mode.toString()
            typeInfo.text = Connection.getCon()?.type.toString()
            struInfo.text = Connection.getCon()?.structure.toString()
            userInfo.text = Connection.getCon()?.username.toString()
        } catch (e: Exception) {
            val error = Connection.exceptionHandle(e)
            AlertDialog.Builder(this).setMessage(error)
                .setPositiveButton("OK", null).create().show()
        }
    }

    override fun onClick(v: View?) {
        lateinit var backIntent: Intent
        when (intent.getStringExtra("from")) {
            "main" -> {
                backIntent = Intent(this, MainActivity::class.java)
            }
            "user" -> {
                backIntent = Intent(this, UserActivity::class.java)
            }
            "download" -> {
                backIntent = Intent(this, DownloadActivity::class.java)
            }
            "upload" -> {
                backIntent = Intent(this, UploadActivity::class.java)
            }
            "login" -> {
                backIntent = Intent(this, LoginActivity::class.java)
            }
        }

        when (v?.id) {
            R.id.okBtn -> {
                startActivity(backIntent)
            }
        }
    }
}