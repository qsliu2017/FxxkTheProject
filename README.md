# FxxkTheProject

## FTP协议概述

FTP协议是建立在TCP连接上的应用层协议。

### 工作流程
```ascii
                                      -------------
                                      |/---------\|
                                      ||   User  ||    --------
                                      ||Interface|<--->| User |
                                      |\----^----/|    --------
            ----------                |     |     |
            |/------\|  FTP Commands  |/----V----\|
            ||Server|<---------------->|   User  ||
            ||  PI  ||   FTP Replies  ||    PI   ||
            |\--^---/|                |\----^----/|
            |   |    |                |     |     |
--------    |/--V---\|      Data      |/----V----\|    --------
| File |<--->|Server|<---------------->|  User   |<--->| File |
|System|    || DTP  ||   Connection   ||   DTP   ||    |System|
--------    |\------/|                |\---------/|    --------
            ----------                -------------

            Server-FTP                   USER-FTP
```

FTP协议工作的流程大致如下：

1. 服务端协议解释进程(Server PI, server protocal interpreter)监听一个端口（FTP知名端口为21），客户端协议解释进程(User PI, user protocol interpreter)创建一个TCP连接到Server PI，这个连接用作FTP控制，传输FTP控制指令和响应（合称FTP控制流）。
2. 文件数据通过FTP数据流传输。需要建立数据流连接时，客户端监听一个任意端口，通过控制流`PORT`指令把这个端口号通知服务端；服务端新建一个TCP连接到这个端口，完成数据流连接。
3. 需要断开FTP连接时，User PI通过控制发送一个`QUIT`指令，服务端收到指令后响应并断开连接。

### 最小实现
rfc959规定FTP的最小实现应该包括：
|                     |                                                      |
| ------------------- | ---------------------------------------------------- |
| 数据格式(TYPE)      | ASCII Non-print                                      |
| 传输模式(MODE)      | Stream                                               |
| 数据结构(STRUCTURE) | File, Record                                         |
| 控制命令(COMMANDS)  | USER, QUIT, PORT, TYPE, MODE, STRU, RETR, STOR, NOOP |

- 传输模式(MODE)
  - 流式传输(Stream) 数据按照字节流传输。

- 数据结构(STRUCTURE)
  - 文件结构(File) 没有内部结构，文件视为字节流。
  - 记录结构(Record) 

- 控制命令(COMMANDS)
  | 权限控制      | 参数设置        | 服务命令      |
  | ------------- | --------------- | ------------- |
  | USER 用户登入 | PORT 数据流端口 | RETR 读取文件 |
  | QUIT 用户登出 | TYPE 表示格式   | STOR 保存文件 |
  |               | MODE 传输模式   | NOOP 空指令   |
  |               | STRU 数据结构   |               |
