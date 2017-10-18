# redis Linux安装及配置
## 下载redix安装redis
[redis官网](https://redis.io/download)  
1. redis官网有安装命令，以下是我安装redis的命令  
    `wget http://download.redis.io/releases/redis-4.0.2.tar.gz`  
    `tar xzf redis-4.0.2.tar.gz`  
    `cd redis-4.0.2`  
    `make test`  
    `make install`  
2. redis编译安装完成后，在安装目录下有个配置文件redis.cong，src目录下有3个可执行文件redis-server\redis-cli\redis-benchmark。将这四个文件拷贝到一个目录下/usr/redis/  
    `mkdir /usr/redis`  
    `cp redis.conf  /usr/redis`  
    `cd src`  
    `cp redis-server  /usr/redis`  
    `cp redis-benchmark /usr/redis`  
    `cp redis-cli  /usr/redis`  
3. 配置redis.conf,新建reids日志文件/usr/redis/log/redis.log,并设置允许后台运行,允许其他ip能够访问  
    `vi /usr/redis/redis.conf`  
    *将daemonize的属性改为 yes*  
    *将 logfile 的值设为 /usr/redis/log/redis.log*  
    *将protected-mode的属性改为 no*  
    *将bind的属性设为 0.0.0.0 允许任何ip访问*   
    `wq`保存并退出  
    *新建reids日志文件*  
    `mkdir /usr/redis/log`  
    `vi /usr/redis/log/redis.log`  
4. 运行redis
    `redis-server /usr/redis/redis.conf`
