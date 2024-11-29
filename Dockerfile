FROM golang:21
LABEL \
        image.author="zhangjay200@gmail.com" \
        appName="solar-metrics" \
        version="V1.0.0"
ENV WORK_DIR="/data/app"
WORKDIR /data/app
#获取打包的最新文件
COPY solar-metrics $WORK_DIR
COPY config/* $WORK_DIR/config

EXPOSE 8080

ENTRYPOINT["./solar-metrics"]

