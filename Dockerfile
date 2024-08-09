FROM amd64/centos:7
LABEL maintainer="zhoushoujianwork@163.com"
COPY bin/goosefs-cli2api-linux-amd64 /usr/bin/goosefs-cli2api
COPY entrypoint.sh /root/entrypoint.sh
RUN chmod +x /root/entrypoint.sh 
WORKDIR /root
EXPOSE 8080
ENTRYPOINT /root/entrypoint.sh
