package com.example.ftpclient

import android.content.Intent
import android.os.Bundle
import android.view.Menu
import android.view.MenuItem
import android.view.View
import androidx.appcompat.app.AlertDialog
import androidx.appcompat.app.AppCompatActivity
import androidx.core.content.ContextCompat
import client.Client
import fm.Fm
import kotlinx.android.synthetic.main.activity_main.*

class MainActivity : AppCompatActivity(), View.OnClickListener {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
        downloadBtn.setOnClickListener(this)
        loginBtn.setOnClickListener(this)
        disConnectBtn.setOnClickListener(this)
        Fm.setFileManager(
            Connection.FileManagerImpl(
                ContextCompat.getExternalFilesDirs(
                    this,
                    null
                )[0]
            )
        )
        Client.setBuffer(ByteArray(30 * 1024))
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
                    .setPositiveButton(
                        "Yes"
                    ) { _, _ ->
                        Connection.getCon()?.logout()
                        startActivity(Intent(this, ConnectActivity::class.java))
                    }
                    .setNegativeButton("No", null).create().show()
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
                intent.putExtra("from", "main")
                startActivity(intent)
            }
            R.id.dataMode -> {
                val intent = Intent(this, ModeActivity::class.java)
                intent.putExtra("from", "main")
                startActivity(intent)
            }
            R.id.type -> {
                val intent = Intent(this, TypeActivity::class.java)
                intent.putExtra("from", "main")
                startActivity(intent)
            }
            R.id.structure -> {
                val intent = Intent(this, StructureActivity::class.java)
                intent.putExtra("from", "main")
                startActivity(intent)
            }
        }
        return true
    }
}