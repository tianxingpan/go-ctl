# go-ctl
go-ctl是一款SSH的工具客户端辅助工具，支持同时操作多台远程机器执行命令，或者批量上传文件等等，是一款不错的辅助工具。

## 子命令
### cli
控制远程机器执行命令，并返还结果。
#### 参数说明
##### command
远程执行的命令，支持短命名*-c* 和长命名 *--command*
```shell script
./go-ctl cli -c='ls /'
```
##### hosts
目标机器IP，如果是多台机器，IP用英文逗号分割，支持短命名*-H* 和长命名 *--hosts*
##### port
目标机器端口，支持短命名*-P* 和长命名 *--port*
##### user
目标机器登陆用户名，支持短命名*-u* 和长命名 *--user*
##### password
目标机器登陆密码，支持短命名*-p* 和长命名 *--password*
#### 用例
```shell script
./go-ctl cli -c='ls /' -H=127.0.0.1 -P=36000 -P="test@test"
```

### push
批量上传文件或者目录到目标机器
#### 参数说明
##### source
本地待传文件或目录，支持短命名*-s* 和长命名 *--source*
##### destination
远程服务器接受目录，支持短命名*-d* 和长命名 *--destination*
##### port
目标机器端口，支持短命名*-P* 和长命名 *--port*
##### user
目标机器登陆用户名，支持短命名*-u* 和长命名 *--user*
##### password
目标机器登陆密码，支持短命名*-p* 和长命名 *--password*
#### 用例
```shell script
# 可以带上目标机器IP
./go-ctl push -s=./test.txt -d=/data/test -H=127.0.0.1 -P=36000 -P="test@test"
# 可以不带目标机器IP，如果没带，则为上一次操作的IP地址
./go-ctl push -s=./test.txt -d=/data/test
```
### pull
下载目标机器的文件或者目录，支持多台机器。
#### 参数说明
##### source
远程服务器待传文件或目录，支持短命名*-s* 和长命名 *--source*
##### destination
本地服务器接受目录，支持短命名*-d* 和长命名 *--destination*
##### port
目标机器端口，支持短命名*-P* 和长命名 *--port*
##### user
目标机器登陆用户名，支持短命名*-u* 和长命名 *--user*
##### password
目标机器登陆密码，支持短命名*-p* 和长命名 *--password*
#### 用例
```shell script
# 可以带上目标机器IP
./go-ctl pull -s=/data/test -d=/data/test/test.tx -H=127.0.0.1 -P=36000 -P="test@test"
# 可以不带目标机器IP，如果没带，则为上一次操作的IP地址
./go-ctl pull -s=/data/test -d=/data/test/test.tx
```

### env
设置go-ctl配置和环境变量

#### 参数说明

#### 用例

