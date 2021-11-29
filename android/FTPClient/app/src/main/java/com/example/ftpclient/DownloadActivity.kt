package com.example.ftpclient

import android.content.Intent
import android.os.Bundle
import android.util.Log
import android.view.Menu
import android.view.MenuItem
import android.view.View
import androidx.appcompat.app.AlertDialog
import androidx.appcompat.app.AppCompatActivity
import androidx.core.content.ContextCompat
import kotlinx.android.synthetic.main.activity_download.*
import kotlinx.android.synthetic.main.activity_download.localName
import kotlinx.android.synthetic.main.activity_download.remoteName
import kotlinx.android.synthetic.main.activity_test.*
import java.io.File

class DownloadActivity : AppCompatActivity(), View.OnClickListener {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_download)
        downloadBtn.setOnClickListener(this)
    }

    override fun onClick(v: View?) {
        when (v?.id) {
            R.id.downloadBtn -> {
                val local = localName.text.toString().trim()
                val remote = remoteName.text.toString().trim().split(";")

                // Download a file
                try {
                    if (remote.size > 1) {
                        val dir = File(
                            ContextCompat.getExternalFilesDirs(
                                this,
                                null
                            )[0], local
                        )
                        dir.mkdir()

                        for (file in remote) {
                            val localFile = File(dir, file)
                            localFile.createNewFile()
                            Connection.getCon()?.retrieve(localFile.toString(), file)
                        }
                    } else {
                        Connection.getCon()?.retrieve(local, remote[0])
                    }
                    val from = intent.getStringExtra("from").toString()
                    val dialog = AlertDialog.Builder(this).setMessage("Download successfully!")
                    if (from == "main") {
                        dialog.setPositiveButton(
                            "OK"
                        ) { _, _ ->
                            startActivity(Intent(this, MainActivity::class.java))
                        }
                    } else {
                        dialog.setPositiveButton(
                            "OK"
                        ) { _, _ ->
                            startActivity(Intent(this, UserActivity::class.java))
                        }
                    }
                    dialog.create().show()
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
                intent.putExtra("from", "download")
                startActivity(intent)
            }
            R.id.dataMode -> {
                val intent = Intent(this, ModeActivity::class.java)
                intent.putExtra("from", "download")
                startActivity(intent)
            }
            R.id.type -> {
                val intent = Intent(this, TypeActivity::class.java)
                intent.putExtra("from", "download")
                startActivity(intent)
            }
            R.id.structure -> {
                val intent = Intent(this, StructureActivity::class.java)
                intent.putExtra("from", "download")
                startActivity(intent)
            }
            R.id.information -> {
                val intent = Intent(this, InformationActivity::class.java)
                intent.putExtra("from", "download")
                startActivity(intent)
            }
        }
        return true
    }
}