package com.example.ftpclient

import android.content.Intent
import android.os.Bundle
import android.view.Menu
import android.view.MenuItem
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

    override fun onCreateOptionsMenu(menu: Menu?): Boolean {
        menuInflater.inflate(R.menu.main, menu)
        return true
    }

    override fun onOptionsItemSelected(item: MenuItem): Boolean {
        when (item.itemId) {
            R.id.conMode -> {
                val intent = Intent(this, ConnModeActivity::class.java)
                intent.putExtra("from", "upload")
                startActivity(intent)
            }
            R.id.dataMode -> {
                val intent = Intent(this, ModeActivity::class.java)
                intent.putExtra("from", "upload")
                startActivity(intent)
            }
            R.id.type -> {
                val intent = Intent(this, TypeActivity::class.java)
                intent.putExtra("from", "upload")
                startActivity(intent)
            }
            R.id.structure -> {
                val intent = Intent(this, StructureActivity::class.java)
                intent.putExtra("from", "upload")
                startActivity(intent)
            }
        }
        return true
    }
}