package com.example.ftpclient

import client.Client
import client.FtpClient

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
}