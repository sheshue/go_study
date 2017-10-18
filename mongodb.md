# window10安装mongodb
## 下载mongodb
[mongoDB官网](http://www.mongodb.org/)  
[直接下载zip包](http://dl.mongodb.org/dl/win32/x86_64)
## 安装mongodb
1. 创建文件路径：D:\mongodb
2. 将刚刚下载的zip解压在mongodb，修改文件夹名为bin
3. 然后在D:\mongodb下新建文件夹data
4. 在data下建立文件夹db,log
5. 在log文件夹下创建日志文件MongoDB.log
6. 至此mongodb有一下文件：  
    * D:\mongodb\bin
    * D:\mongodb\data
    * D:\mongodb\data\db
    * D:\mongodb\data\log
    * D:\mongodb\data\log\MongoDB.log
7. 命令行安装  
    `D:\mongodb\bin>mongod -dbpath "D:\mongodb\data\db"`
8. 出现以下内容说明已经安装完成  
    `waiting for connections on port 27017`
9. 安装完成后打开 *http://127.0.0.1:27017*  
    `It looks like you are trying to access MongoDB over HTTP on the native driver port. `  
    `说明已经安装成功了。`
## 开启mongodb服务
1. 在mongodb文件下建立mongodb.conf配置文件  
    `dbpath=D:\program\mongodb\data\db #数据库路径`  
    `logpath=D:\program\mongodb\data\log\MongoDB.log #日志输出文件路径`  
    `logappend=true #错误日志采用追加模式，配置这个选项后mongodb的日志会追加到  现有的日志文件，而不是从新创建一个新文件`  
    `journal=true #启用日志文件，默认启用`  
    `quiet=true #这个选项可以过滤掉一些无用的日志信息，若需要调试使用请设置为false`  
    `port=27017 #端口号 默认为27017`  
2. 在命令行转到 *D:\mongodb\bin*,运行  
    `mongod --config D:\program\mongodb\mongo.config`
3. 打开另一个命令行窗口,转到*D:\mongodb\bin*,运行  
    `mongo`  
    `此操作打开mongodb控制台`
    
