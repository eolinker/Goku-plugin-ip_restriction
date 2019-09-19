# Goku Plugin：IP Restriction

| 插件名称  | 文件名.so |  插件类型  | 错误处理方式 | 作用范围 |  优先级  |
| ------------ | ------------ | ------------ | ------------ | ------------ | ------------ |
| IP黑白名单  | goku-ip_restriction | 访问策略 | 继续后续操作 | 转发前  | 990 |

IP黑名单指除黑名单外的IP均可访问，IP白名单指除白名单外的IP不能访问，网关通过会 **X-Real-IP** 头判断客户端真实IP。

##### 一、IP配置支持以下写法：
（1）192.168.0.1
（2）192.168.0.1/26
（3）192.168.0.*

注：\*仅支持放最后一位，如：192.\*、192.168.\*、192.168.0.\*

##### 二、配合nginx的X-Real-IP使用：
（1）若客户端与网关之间不存在代理服务器，此时从请求中解析出的IP地址就是实际客户端的IP，网关会把该地址设为X-Real-IP头的值；
（2）若客户端与网关之间存在多层代理，则需在 **第一层代理** 中设置X-Real-IP请求头，此时网关会把代理传来的X-Real-IP转发到服务器。

**代理配置示例**：
```
location / {
	root   html;
	index  index.html index.htm index.php;
	proxy_set_header X-Real-IP $remote_addr;
	proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
	proxy_pass http://192.168.247.132;
	}
```

# 目录

- [编译教程](#编译教程 "编译教程")
- [安装教程](#安装教程 "安装教程")
- [使用教程](#使用教程 "使用教程")
- [更新日志](#更新日志 "更新日志")

# 编译教程

#### 环境要求
* 系统：基于 Linux 内核（2.6.23+）的系统，CentOS、RedHat 等均可；

* golang版本号：12.x及其以上

* 环境变量设置：
	* GO111MODULE：on
	
	* GOPROXY：https://goproxy.io


#### 编译步骤

1.clone项目

2.进入项目文件夹，执行**build.sh**
```
cd goku-ip_restriction && chmod +x build.sh && ./build.sh
```

###### 注：build.sh为通用的插件编译脚本，自定义插件时可以拷贝直接使用。

3.执行第2步将会生成文件： **{插件名}.so**

将该文件上传到**节点服务器运行目录**下的**plugin**文件夹，然后在控制台安装插件即可使用。

# 安装教程
前往 Goku API Gateway 官方网站查看：[插件安装教程](url "https://help.eolinker.com/#/tutorial/?groupID=c-341&productID=19")

# 使用教程

#### 配置页面

进入控制台 >> 策略管理 >> 某策略 >> 策略插件 >> IP黑白名单插件：

![](http://data.eolinker.com/course/rrwbADA8d86ccce880127198916f0b0306250ecacb002e0)

#### 配置参数

| 参数名 | 说明   |  值可能性
| ------------ | ------------ |  ------------ |  
|  ipListType | IP名单类型| white/black/none  | 
| ipWhiteList  | IP白名单列表 |   |
| ipBlackList  | IP黑名单列表 |  |  |

#### 配置示例

```
{
    "ipListType":"black",
    "ipWhiteList":["127.0.0.1"],
    "ipBlackList":["127.0.0.2"]
}
```

#### ### 返回示例

```
	HTTP/1.1 403 Forbidden
	Content-Type: text/plain
	Content-Length: 18

	[ERROR]Illegal IP!
```
