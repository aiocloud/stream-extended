#!/usr/bin/env bash
echo=echo
for cmd in echo /bin/echo; do
    $cmd >/dev/null 2>&1 || continue

    if ! $cmd -e "" | grep -qE '^-e'; then
        echo=$cmd
        break
    fi
done

CSI=$($echo -e "\033[")
CEND="${CSI}0m"
CDGREEN="${CSI}32m"
CRED="${CSI}1;31m"
CGREEN="${CSI}1;32m"
CYELLOW="${CSI}1;33m"
CBLUE="${CSI}1;34m"
CMAGENTA="${CSI}1;35m"
CCYAN="${CSI}1;36m"

OUT_ALERT() {
    echo -e "${CYELLOW}$1${CEND}"
}

OUT_ERROR() {
    echo -e "${CRED}$1${CEND}"

    exit 1
}

OUT_INFO() {
    echo -e "${CCYAN}$1${CEND}"
}

if [[ -f /etc/redhat-release ]]; then
    release="centos"
elif cat /etc/issue | grep -q -E -i "debian"; then
    release="debian"
elif cat /etc/issue | grep -q -E -i "ubuntu"; then
    release="ubuntu"
elif cat /etc/issue | grep -q -E -i "centos|red hat|redhat"; then
    release="centos"
elif cat /proc/version | grep -q -E -i "raspbian|debian"; then
    release="debian"
elif cat /proc/version | grep -q -E -i "ubuntu"; then
    release="ubuntu"
elif cat /proc/version | grep -q -E -i "centos|red hat|redhat"; then
    release="centos"
else
    OUT_ERROR "[错误] 不支持的操作系统！"
fi

OUT_ALERT "[信息] 下载程序中"
cd ~ && rm -fr release
wget -O release.zip https://github.com/aiocloud/stream-extended/releases/latest/download/release.zip || exit 1

OUT_ALERT "[信息] 解压程序中"
unzip release.zip && rm -f release.zip && cd release

OUT_ALERT "[信息] 设置权限中"
chmod +x stream-extended

OUT_ALERT "[提示] 复制配置中"
cp -f default.json /etc/stream-extended.json

OUT_ALERT "[提示] 复制程序中"
cp -f stream-extended /usr/bin

OUT_ALERT "[提示] 创建用户中"
userdel -r -f stream-extended
groupdel stream-extended
groupadd -g 1234 stream-extended
useradd -M -s /bin/false -u 1234 -g 1234 stream-extended

OUT_ALERT "[提示] 配置服务中"
cat >/etc/systemd/system/stream-extended.service <<EOF
[Unit]
Description=Stream Unlock Service [Extended]
After=network.target

[Service]
Type=simple
User=stream-extended
Group=stream-extended
ExecStart=/usr/bin/stream-extended -c /etc/stream-extended.json
Restart=always
RestartSec=4

[Install]
WantedBy=multi-user.target
EOF

OUT_ALERT "[提示] 重载服务中"
systemctl daemon-reload

OUT_INFO "[信息] 部署完毕！"
cd ~ && rm -fr release
exit 0
