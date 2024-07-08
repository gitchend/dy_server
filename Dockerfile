FROM public-cn-beijing.cr.volces.com/public/base:golang-1.17.1-alpine3.14
USER root
WORKDIR /app
COPY . /app/
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -gcflags "all=-N -l" -o main -mod=vendor main.go
RUN chmod -R 777 /app/run.sh

RUN apk update && \
apk upgrade && \
apk add bash

CMD /app/run.sh