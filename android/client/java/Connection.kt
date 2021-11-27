package com.example.ftpclient

import android.util.Log
import client.Client
import client.FtpClient
import fm.Fm.readEOF
import fm.MyFile
import fm.MyFileManager
import java.io.*

// 单例类，在activity之间实现数据传递
object Connection {
    // 一个private的FtpClient对象作为其属性，可以为null值
    private var connection: FtpClient? = null
    fun setCon(address: String): Int {
        try {
            connection = Client.newFtpClient(address)
        } catch (e: java.lang.Exception) {
            return -1
        }
        return 0
    }

    fun getCon(): FtpClient? {
        return connection
    }

    fun exceptionHandle(e: Exception): String {
        return when (e) {
            Client.getErrModeNotSupported() -> {
                "Error: mode is not supported"
            }
            Client.getErrPasswordNotMatch() -> {
                "Error: wrong password"
            }
            Client.getErrUsernameNotExist() -> {
                "Error: username not exist"
            }
            else -> {
                "Error: " + e.message
            }
        }
    }

    class FileManagerImpl(private val context: File) : MyFileManager {
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
            Log.d("ClientWrite", String(content))
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
            return fileBufferedInputStream?.read(buffer)?.toLong() ?: 0
        }
    }
}