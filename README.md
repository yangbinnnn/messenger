# messenger

消息转发支持
1. 邮箱
2. 企业微信

主要应用场景，报警信息推送等

# 配置
企业微信推送需要申请企业微信以及创建应用

申请企业号参考: https://github.com/Yanjunhui/chat

# 使用
- 支持GET，POST(JSON/Form)
- 参数`tos` 支持批量,使用逗号分隔,如 "tos=abc,xyz"

```
$ curl -i "http://127.0.0.1:4000/sender/mail?tos=abc@xyz.com&subject=hello&content=world"
HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Date: Fri, 13 Apr 2018 08:31:15 GMT
Content-Length: 7

success%                
```
```
$ curl -i "http://127.0.0.1:4000/sender/wechat?tos=xyz&content=hello-world" 
HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Date: Fri, 13 Apr 2018 08:06:05 GMT
Content-Length: 7

success%      
```

