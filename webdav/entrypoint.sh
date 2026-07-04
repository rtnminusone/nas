#!/bin/bash

# 1. 환경변수 매핑
AUTH_USER="${WEBDAV_USERNAME:-admin}"
AUTH_PASS="${WEBDAV_PASSWORD:-admin}"

# 2. 🌟 htpasswd 정확한 공식 경로로 수정하여 인증 파일 생성
/usr/local/apache2/bin/htpasswd -cb /usr/local/apache2/conf/webdav.htpasswd "$AUTH_USER" "$AUTH_PASS"

HTTPD_CONF="/usr/local/apache2/conf/httpd.conf"

# 3. 필수 모듈 활성화
sed -i 's/#LoadModule dav_module/LoadModule dav_module/' "$HTTPD_CONF"
sed -i 's/#LoadModule dav_fs_module/LoadModule dav_fs_module/' "$HTTPD_CONF"
sed -i 's/#LoadModule headers_module/LoadModule headers_module/' "$HTTPD_CONF"

# 4. WebDAV 및 윈도우 탐색기 최적화 설정 주입
cat <<EOF >> "$HTTPD_CONF"

DavLockDB /var/lib/dav/DavLock

Alias / /data/
<Directory /data/>
    Dav On
    Options Indexes FollowSymLinks
    AllowOverride None
    Require valid-user

    # 인증 설정
    AuthType Basic
    AuthName "WebDAV Storage"
    AuthUserFile /usr/local/apache2/conf/webdav.htpasswd

    Header always set Access-Control-Allow-Origin "*"
    Header always set Access-Control-Allow-Methods "GET, POST, OPTIONS, HEAD, PUT, DELETE, PROPFIND, PROPPATCH, COPY, MOVE, LOCK, UNLOCK"
    Header always set Access-Control-Allow-Headers "User-Agent, Authorization, Content-Type, Depth, If-Modified-Since, Cache-Control, X-Requested-With"
</Directory>
EOF

# 5. 아파치 엔진 구동
exec httpd -D FOREGROUND