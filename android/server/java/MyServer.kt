package com.example.ftpserver

import android.util.Log
import android.widget.TextView
import server.*
import server.OutputStream
import server.Server.readEOF
import java.io.*

object MyServer {
    private lateinit var myServer: FtpServer

    fun startServer(port: Long): Int {
        myServer = Server.newFtpServer()
        try {
            myServer.listen(port)
        } catch (e: Exception) {
            return -1
        }
        return 0
    }

    fun stopServer() {
        myServer.close()
    }

    class Logger(private val logText: TextView) : OutputStream {
        override fun write(log: ByteArray?): Long {
            log?.let { String(it) }?.let { Log.d("Server", it) }
            val text = logText.text.toString() + "\n" + String(log!!)
            logText.text = text
            return log.toString().length.toLong()
        }
    }

    class FileManagerImpl(private val context: File) : MyFileManager {
//        fun getFile(path: String): MyFile? {
//            // 返回一个FIleImpl
//            val file = File(context, path)
//            Log.d("File", file.path)
//            if (!file.exists())
//                return null
//            return FileImpl(file)
//        }

        override fun create(path: String): MyFile {
            val file = File(context, path)
            val pf = file.parentFile
            if (!pf.exists()) {
                pf.mkdir()
            }
            file.createNewFile()
            return FileImpl(file)
        }

        override fun open(path: String): MyFile {
            val file = File(context, path)
            return FileImpl(file)
        }
    }

    class FileImpl(private val file: File) : MyFile {
        private var fileBufferedInputStream: BufferedInputStream? = null
        private var fileBufferedOutputStream: BufferedOutputStream? = null

        override fun write(content: ByteArray): Long {
            if (fileBufferedOutputStream == null)
                fileBufferedOutputStream = BufferedOutputStream(FileOutputStream(file))
            // 返回写入长度
            fileBufferedOutputStream?.write(content)
            return content.size.toLong()
        }

        override fun close() {
            // 关闭文件
            return
        }

        override fun read(buffer: ByteArray): Long {
            if (fileBufferedInputStream == null)
                fileBufferedInputStream = BufferedInputStream(FileInputStream(file))
            // 写到ByteArray，返回写的长度
            if (fileBufferedInputStream?.available() == 0)
                return readEOF()
            val content = fileBufferedInputStream?.read(buffer)
            return content?.toLong() ?: 0
        }
    }
}