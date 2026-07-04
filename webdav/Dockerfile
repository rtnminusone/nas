FROM httpd:2.4

# 1. 락DB 폴더 및 데이터 폴더 생성
# 데비안 아파치의 기본 실행 유저도 'www-data'입니다.
RUN mkdir -p /var/lib/dav /data && \
    chown -R www-data:www-data /var/lib/dav /data

# 2. 엔트리포인트 주입
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 80

ENTRYPOINT ["/entrypoint.sh"]