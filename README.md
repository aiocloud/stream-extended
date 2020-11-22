# stream-extended
```
# 绕过解锁机器
iptables -t nat -A OUTPUT -d 1.1.1.1/32 -j RETURN

# 绕过 stream-extended 程序
iptables -t nat -A OUTPUT -m owner --uid-owner 1000 -j RETURN

# 劫持本机 TCP 80 443 连接
iptables -t nat -A OUTPUT -p tcp --dport 80 -j DNAT --to-destination 127.0.0.1:80
iptables -t nat -A OUTPUT -p tcp --dport 443 -j DNAT --to-destination 127.0.0.1:443

# 丢弃本机 UDP 443 流量（防止 QUIC 连接）
iptables -t nat -A OUTPUT -p udp --dport 443 -J DROP

# 保存 iptables 规则
netfilter-persistent save
```
