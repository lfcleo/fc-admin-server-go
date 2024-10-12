<h1>FC-Admin-Server-Go</h1>

此项目是[FC-Admin-Server](https://github.com/lfcleo/fc-admin-server)的服务端版本，使用Go语言编写。请配合前端项目[FC-Admin-Server](https://github.com/lfcleo/fc-admin-server)使用。

有编写FC-Admin-Server其它后端语言的大佬们开源的话，可以联系我在项目中展示您的项目地址。（统一下项目名称格式，fc-admin-server-xxx。例如：fc-admin-server-java，fc-admin-server-php，fc-admin-server-python等）

此项目使用框架：gin+gorm

## 准备
此项目运行需要安装`redis`（作者版本:7.2.4）和`mysql`（作者版本:8.3.0）。（安装过程请自行搜索）

mysql创建数据 `fc_admim`（数据库类型：CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci）

本地运行请自行安装Go环境(作者Go版本：1.21.3)，作者本地运行的操作系统是MacOS（Windows系统没有测试过，有问题请联系作者修改Readme，帮助其他小伙伴们避坑。）

## 运行go程序
go环境的基本配置
...

拉取后端代码
```shell
git clone https://github.com/lfcleo/fc-admin-server-go.git
```

修改`config.yaml`中的配置信息，都有注释，根据实际情况修改。配置信息一定要填写正确，尤其是端口号，域名，database（数据库），redis 的信息。注意账号密码不要填写错了。

进入目录
```shell
cd fc-admin-server-go
```

拉取程序所需依赖
```shell
go mod download
```

本地运行
```shell
go run main.go
```
运行报错在终端/控制台有错误信息打印，根据信息修改。

## 前端程序运行。

请查看前端项目[FC-Admin-Server](https://github.com/lfcleo/fc-admin-server)运行。

前端超级管理员账号：admin@admin.com

前端超级管理员手机号：18888888888

前端超级管理员密码：Qwe123