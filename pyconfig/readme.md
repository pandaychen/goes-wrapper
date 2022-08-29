##  pyconfig包

1、`env.go`：封装了容器环境内以环境变量当做结构体的方法，如`MysqlConfig`结构，定义在环境变量如下：

```text
MYSQLDBNAME=xxx
MYSQLUSER=xxxxx
MYSQLPORT=xxx
MYSQLHOST=xxxx
MYSQLCONNTIMEOUT=20s
MYSQLPASS=xxx
```
