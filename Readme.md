# 使用
    远端 go run srv.go
    本地 go run clt.go


# 配置 conf.ini
    [common]
    ;加密方式
    cipher = caesar
    ;代理协议
    proto = socks5
    ;token
    token = zcsk18

    [srv]
    ;服务端 域名/ip
    host = 127.0.0.1
    ;服务端 端口
    port = 2081

    [clt]
    ;本地 端口
    port = 1081

    [caesar]
    ;凯撒加密 偏移量(0 - 255)
    dis = 5

#测试
    curl --socks5-hostname 127.0.0.1:1081 https://google.com
