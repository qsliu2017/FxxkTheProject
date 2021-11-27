package com.example.ftpclient

import android.content.Intent
import android.os.Bundle
import android.view.View
import androidx.appcompat.app.AlertDialog
import androidx.appcompat.app.AppCompatActivity
import client.Client
import kotlinx.android.synthetic.main.activity_type.*

class TypeActivity : AppCompatActivity(), View.OnClickListener {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_type)
        ascBtn.setOnClickListener(this)
        biBtn.setOnClickListener(this)
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
            R.id.ascBtn -> {
                try {
                    Connection.getCon()?.type(Client.TypeAscii)
                    success.setMessage("Ascii type set").create().show()
                } catch (e: Exception) {
                    val error = Connection.exceptionHandle(e)
                    fail.setMessage(error).create().show()
                }
            }
            R.id.biBtn -> {
                try {
                    Connection.getCon()?.type(Client.TypeBinary)
                    success.setMessage("Binary type set").create().show()
                } catch (e: Exception) {
                    val error = Connection.exceptionHandle(e)
                    fail.setMessage(error).create().show()
                }
            }
        }
    }
}