package com.example.ftpserver

import android.util.Log
import android.widget.TextView
import server.*
import server.OutputStream
import server.Server.readEOF
import java.io.*

object MyServer {
    private var myServer: FtpServer? = null

    fun startServer(port: Long, path: String): Int {
        myServer = Server.newFtpServer()
        try {
            myServer?.setRootDir(path)
            myServer?.listen(port)
        } catch (e: Exception) {
            return -1
        }
        return 0
    }

    fun stopServer() {
        myServer?.close()
    }

    class Logger(private val logText: TextView) : OutputStream {
        override fun write(log: ByteArray?): Long {
            log?.let { String(it) }?.let { Log.d("Server", it) }
            val text = logText.text.toString() + "\n" + String(log!!)
            logText.text = text
            return log.toString().length.toLong()
        }
    }
}