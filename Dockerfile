FROM golang:1.21
LABEL \
    image.author="zhangjay200@gmail.com" \
    appName="solar-metrics" \
    version="V1.0.0"

ENV WORK_DIR="/data/app"
ENV PROFILE="prod"

WORKDIR $WORK_DIR

# 将编译后的Go二进制文件和配置文件复制到容器中
COPY solar-metrics $WORK_DIR
COPY config/* $WORK_DIR/config

EXPOSE 9999

ENTRYPOINT ["./solar-metrics"]
