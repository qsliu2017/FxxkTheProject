package com.example.ftpclient

import android.content.Intent
import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.view.View
import androidx.appcompat.app.AlertDialog
import client.Client
import kotlinx.android.synthetic.main.activity_conn_mode.*

class ConnModeActivity : AppCompatActivity(), View.OnClickListener {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_conn_mode)
        portBtn.setOnClickListener(this)
        passiveBtn.setOnClickListener(this)
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

        val success = AlertDialog.Builder(this).setPositiveButton(
            "OK"
        ) { _, _ ->
            startActivity(backIntent)
        }
        val fail = AlertDialog.Builder(this).setPositiveButton("OK", null)

        when (v?.id) {
            R.id.portBtn -> {
                try {
                    Connection.getCon()?.connMode(Client.ConnPort)
                    success.setMessage("Port type set").create().show()
                } catch (e: Exception) {
                    val error = Connection.exceptionHandle(e)
                    fail.setMessage(error).create().show()
                }
            }
            R.id.passiveBtn -> {
                try {
                    Connection.getCon()?.connMode(Client.ConnPasv)
                    success.setMessage("Passive type set").create().show()
                } catch (e: Exception) {
                    val error = Connection.exceptionHandle(e)
                    fail.setMessage(error).create().show()
                }
            }
        }
    }
}