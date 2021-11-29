package com.example.ftpclient

import android.os.Bundle
import android.view.View
import androidx.appcompat.app.AlertDialog
import androidx.appcompat.app.AppCompatActivity
import androidx.core.content.ContextCompat
import client.Client
import kotlinx.android.synthetic.main.activity_test.*
import java.io.File
import java.util.*

class TestActivity : AppCompatActivity(), View.OnClickListener {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_test)
        testBtn.setOnClickListener(this)
    }

    override fun onClick(v: View?) {
        val local = localName.text.toString().trim()
        val remote = remoteName.text.toString().trim()
        val result = StringBuffer()

        when (v?.id) {
            R.id.testBtn -> {
                try {
                    var start: Long?
                    var end: Long?
                    var blockSize: Long
                    var blockTime: Long
                    var rate: Double?
                    val fileTree = File(ContextCompat.getExternalFilesDirs(this, null)[0], local).walk()
                    var size: Double = 0.0
                    var rateStr:String
                    fileTree.forEach { size += it.length().toDouble() }

                    Connection.getCon()?.mode(Client.ModeBlock)
                    var i = 100
                    while (i < 200) {
                        blockSize = i.toLong()
                        Connection.getCon()?.setBlockSize(blockSize)
                        start = Date().time
                        Connection.getCon()?.store(local, remote)
                        end = Date().time
                        blockTime = end - start
                        rate = size / blockTime
                        rateStr = String.format("%.4f", rate)

                        result.append("block size: $i; block time: $blockTime; rate: $rateStr Byte/ms\n")
                        i += 10
                    }

                    Connection.getCon()?.mode(Client.ModeStream)
                    start = Date().time
                    Connection.getCon()?.store(local, remote)
                    end = Date().time
                    val streamTime = end - start
                    rate = size / streamTime
                    rateStr = String.format("%.4f", rate)
                    result.append("stream time: $streamTime; rate: $rateStr Byte/ms\n")

                    AlertDialog.Builder(this)
                        .setMessage(
                            result
                        )
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