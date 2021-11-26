package com.example.ftpclient

import android.content.Intent
import android.os.Bundle
import android.view.Menu
import android.view.MenuItem
import android.view.View
import androidx.appcompat.app.AlertDialog
import androidx.appcompat.app.AppCompatActivity
import kotlinx.android.synthetic.main.activity_user.*

class UserActivity : AppCompatActivity(), View.OnClickListener {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_user)
        downloadBtn.setOnClickListener(this)
        uploadBtn.setOnClickListener(this)
        logoutBtn.setOnClickListener(this)
    }

    override fun onClick(v: View?) {
        when (v?.id) {
            R.id.downloadBtn -> {
                val intent = Intent(this, DownloadActivity::class.java)
                intent.putExtra("from", "user")
                startActivity(intent)
            }
            R.id.uploadBtn -> {
                startActivity(Intent(this, UploadActivity::class.java))
            }
            R.id.logoutBtn -> {
                AlertDialog.Builder(this).setMessage("Log out?")
                    .setPositiveButton("Yes"
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
                intent.putExtra("from", "user")
                startActivity(intent)
            }
            R.id.dataMode -> {
                val intent = Intent(this, ModeActivity::class.java)
                intent.putExtra("from", "user")
                startActivity(intent)
            }
            R.id.type -> {
                val intent = Intent(this, TypeActivity::class.java)
                intent.putExtra("from", "user")
                startActivity(intent)
            }
            R.id.structure -> {
                val intent = Intent(this, StructureActivity::class.java)
                intent.putExtra("from", "user")
                startActivity(intent)
            }
        }
        return true
    }
}