package com.example.ftpclient

import android.content.Intent
import android.os.Bundle
import android.view.View
import androidx.appcompat.app.AlertDialog
import androidx.appcompat.app.AppCompatActivity
import kotlinx.android.synthetic.main.activity_login.*
import java.lang.Exception

class LoginActivity : AppCompatActivity(), View.OnClickListener {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_login)
        loginBtn.setOnClickListener(this)
    }

    override fun onClick(v: View?) {
        when (v?.id) {
            R.id.loginBtn -> {
                val user = username.text.toString()
                val pwd = password.text.toString()
                try {
                    Connection.getCon()?.login(user, pwd)
                    startActivity(Intent(this, UserActivity::class.java))
                } catch (e: Exception) {
                    val error = Connection.exceptionHandle(e)
                    AlertDialog.Builder(this).setMessage(error)
                        .setPositiveButton("OK", null).create().show()
                }
            }
        }
    }
}