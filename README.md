# stream-extended
```
# 安装程序
curl -fsSL https://git.io/JkMmc | bash

# 编辑配置文件
nano /etc/stream-extended.json

# 绕过解锁机器（注意修改这里的 IP 地址）
iptables -t nat -A OUTPUT -d 1.1.1.1/32 -j RETURN

# 绕过 stream-extended 程序
iptables -t nat -A OUTPUT -m owner --uid-owner 1234 -j RETURN

# 劫持本机 TCP 80 443 连接
iptables -t nat -A OUTPUT -p tcp --dport 80 -j DNAT --to-destination 127.0.0.1:60080
iptables -t nat -A OUTPUT -p tcp --dport 443 -j DNAT --to-destination 127.0.0.1:60443

# 丢弃本机 UDP 443 流量（防止 QUIC 连接）
iptables -A OUTPUT -p udp --dport 443 -j DROP

# 保存 iptables 规则
netfilter-persistent save
```
