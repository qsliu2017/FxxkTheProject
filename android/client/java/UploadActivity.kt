package com.example.ftpclient

import android.content.Intent
import android.os.Bundle
import android.view.View
import androidx.appcompat.app.AlertDialog
import androidx.appcompat.app.AppCompatActivity
import androidx.core.content.ContextCompat
import fm.Fm
import kotlinx.android.synthetic.main.activity_upload.*

class UploadActivity : AppCompatActivity(), View.OnClickListener {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_upload)
        uploadBtn.setOnClickListener(this)
        Fm.setFileManager(
            Connection.FileManagerImpl(
                ContextCompat.getExternalFilesDirs(
                    this,
                    null
                )[0]
            )
        )
    }

    override fun onClick(v: View?) {
        when (v?.id) {
            R.id.uploadBtn -> {
                val local = localName.text.toString().trim()
                val remote = remoteName.text.toString().trim()

                // Upload a file
                try {
                    Connection.getCon()?.store(local, remote)
                    AlertDialog.Builder(this).setMessage("Upload successfully!")
                        .setPositiveButton(
                            "OK"
                        ) { _, _ ->
                            startActivity(Intent(this, UserActivity::class.java))
                        }.create().show()
                } catch (e: Exception) {
                    val error = Connection.exceptionHandle(e)
                    AlertDialog.Builder(this).setMessage(error)
                        .setPositiveButton("OK", null).create().show()
                }
            }
        }
    }
}