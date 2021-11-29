package com.example.ftpclient

import android.content.Intent
import android.os.Bundle
import android.view.Menu
import android.view.MenuItem
import android.view.View
import androidx.appcompat.app.AlertDialog
import androidx.appcompat.app.AppCompatActivity
import androidx.core.content.ContextCompat
import kotlinx.android.synthetic.main.activity_main.*
import java.io.File

class MainActivity : AppCompatActivity(), View.OnClickListener {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
        downloadBtn.setOnClickListener(this)
        loginBtn.setOnClickListener(this)
        disConnectBtn.setOnClickListener(this)
        anonymousBtn.setOnClickListener(this)
        val path = ContextCompat.getExternalFilesDirs(this, null)[0].toString()
        val file = File(path)
        if (!file.exists())
            file.createNewFile()
        Connection.getCon()?.setRootDir(path)
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
            R.id.anonymousBtn -> {
                try {
                    Connection.getCon()?.login("anonymous", "anonymous")
                    startActivity(Intent(this, UserActivity::class.java))
                } catch (e: Exception) {
                    val error = Connection.exceptionHandle(e)
                    AlertDialog.Builder(this).setMessage(error)
                        .setPositiveButton("OK", null).create().show()
                }
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
            R.id.information -> {
                val intent = Intent(this, InformationActivity::class.java)
                intent.putExtra("from", "main")
                startActivity(intent)
            }
        }
        return true
    }
}