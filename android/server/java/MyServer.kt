package com.example.ftpserver

import android.util.Log
import android.widget.TextView
import server.*
import server.Server.readEOF
import java.io.File

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
        override fun getFile(path: String): MyFile? {
            // 返回一个FIleImpl
            val file = File(context, path)
            Log.d("File", file.path)
            if (!file.exists())
                return null
            return FileImpl(file)
        }
    }

    class FileImpl(private val file: File) : MyFile {
        private var readOffset = 0
        private val fileContent = file.readBytes()
        private val fileLen = fileContent.size

        override fun write(content: ByteArray): Long {
            // 返回写入长度
            file.writeBytes(content)
            return content.size.toLong()
        }

        override fun close() {
            // 关闭文件
            return
        }

        override fun read(buffer: ByteArray): Long {
            // 写到ByteArray，返回写的长度
            if (fileLen == readOffset)
                return readEOF()
            var len = buffer.size
            if (buffer.size > fileLen - readOffset)
                len = fileLen - readOffset
            val buf = file.readBytes()
            buf.copyInto(buffer, 0, readOffset, readOffset + len)
            readOffset += len
            return buf.size.toLong()
        }
    }
}