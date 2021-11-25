package com.example.ftpclient

import android.util.Log
import client.Client
import client.FtpClient
import fm.Fm
import fm.MyFile
import fm.MyFileManager
import java.io.File

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
        override fun getFile(path: String): MyFile {
            // 返回一个FIleImpl
            val file = File(context, path)
            Log.d("File", file.path)
            if (!file.exists()) {
                file.createNewFile()
            }
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
                return Fm.readEOF()
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