FROM golang:1.20 as builder
# 指定构建过程中的工作目录
WORKDIR /app
# 将当前目录（dockerfile所在目录）下所有文件都拷贝到工作目录下（.dockerignore中文件除外）
COPY . /app/

# 执行代码编译命令。操作系统参数为linux，编译后的二进制产物命名为main，并存放在当前目录下。
RUN GOPROXY=https://goproxy.cn,direct GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o main -mod=vendor cmd/main/main.go

FROM ubuntu:20.04

WORKDIR /opt/application

COPY --from=builder /app/main /opt/application/
COPY --from=builder /app/run.sh /opt/application/

USER root

RUN chmod -R 777 /opt/application/run.sh

CMD /opt/application/run.sh
