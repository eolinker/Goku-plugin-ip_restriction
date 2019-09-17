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
- [安装教程](#安装教程 "安装教程")
- [使用教程](#使用教程 "使用教程")
- [更新日志](#更新日志 "更新日志")

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