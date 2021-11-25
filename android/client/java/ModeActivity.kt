package com.example.ftpclient

import android.os.Bundle
import android.view.View
import androidx.appcompat.app.AlertDialog
import androidx.appcompat.app.AppCompatActivity
import kotlinx.android.synthetic.main.activity_mode.*

class ModeActivity : AppCompatActivity(), View.OnClickListener {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_mode)
        streamBtn.setOnClickListener(this)
        blockBtn.setOnClickListener(this)
        compressedBtn.setOnClickListener(this)
    }

    override fun onClick(v: View?) {
        when (v?.id) {
            R.id.streamBtn -> {
                try {
                    Connection.getCon()?.mode('S'.code.toByte())
                    AlertDialog.Builder(this).setMessage("Stream mode set")
                        .setPositiveButton("OK", null).create().show()
                } catch (e: Exception) {
                    val error = Connection.exceptionHandle(e)
                    AlertDialog.Builder(this).setMessage(error)
                        .setPositiveButton("OK", null).create().show()
                }
            }
            R.id.blockBtn -> {
                try {
                    Connection.getCon()?.mode('B'.code.toByte())
                    AlertDialog.Builder(this).setMessage("Block mode set")
                        .setPositiveButton("OK", null).create().show()
                } catch (e: Exception) {
                    val error = Connection.exceptionHandle(e)
                    AlertDialog.Builder(this).setMessage(error)
                        .setPositiveButton("OK", null).create().show()
                }
            }
            R.id.compressedBtn -> {
                try {
                    Connection.getCon()?.mode('C'.code.toByte())
                    AlertDialog.Builder(this).setMessage("Compressed mode set")
                        .setPositiveButton("OK", null).create().show()
                } catch (e: Exception) {
                    val error = Connection.exceptionHandle(e)
                    AlertDialog.Builder(this).setMessage(error)
                        .setPositiveButton("OK", null).create().show()
                }
            }
        }
    }
}