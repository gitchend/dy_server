FROM golang:1.20 as builder
WORKDIR /app
COPY . /app/
RUN GOPROXY=https://goproxy.cn,direct GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o main -mod=vendor cmd/main/main.go

FROM alpine:latest
WORKDIR /opt/application
COPY --from=builder /app/main /opt/application/
COPY --from=builder /app/run.sh /opt/application/
USER root
RUN chmod -R 777 /opt/application/run.sh

CMD /opt/application/run.sh