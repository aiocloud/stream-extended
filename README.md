# stream-extended
```
iptables -t nat -A OUTPUT -d 1.1.1.1/32 -j RETURN
iptables -t nat -A OUTPUT -m owner --pid-owner 1000 -j RETURN
iptables -t nat -A OUTPUT -p tcp --dport 80 -j DNAT --to-destination 127.0.0.1:80
iptables -t nat -A OUTPUT -p tcp --dport 443 -j DNAT --to-destination 127.0.0.1:443
```
