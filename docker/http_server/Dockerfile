FROM golang:1.21
ENV TZ Europe/Moscow

WORKDIR /var/app/src

COPY  . /var/app/src
RUN go build -o /var/app/bin/server ./cmd/api

WORKDIR /var/app/bin

CMD ["./server"]
