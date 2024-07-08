FROM public-cn-beijing.cr.volces.com/public/base:golang-1.17.1-alpine3.14
WORKDIR /app
COPY . /app/
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -gcflags "all=-N -l" -o main -mod=vendor main.go
WORKDIR /opt/application
COPY /app/main /opt/application/
COPY /app/run.sh /opt/application/
USER root
RUN chmod -R 777 /opt/application/run.sh

RUN apk update && \
apk upgrade && \
apk add bash

CMD /opt/application/run.sh