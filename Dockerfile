FROM public-cn-beijing.cr.volces.com/public/base:golang-1.17.1-alpine3.14 as builder
WORKDIR /app
COPY . /app/
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -gcflags "all=-N -l" -o main -mod=vendor main.go

FROM public-cn-beijing.cr.volces.com/public/base:alpine-3.13
WORKDIR /opt/application
COPY --from=builder /app/main /opt/application/
COPY --from=builder /app/run.sh /opt/application/
USER root
RUN chmod -R 777 /opt/application/run.sh
CMD /opt/application/run.sh