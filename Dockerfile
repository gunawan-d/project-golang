#Stage 1 Build
FROM golang:1.23 AS builder

WORKDIR /app
COPY go.mod /app
COPY go.sum /app
RUN go mod download
COPY . /app

#Build Golang
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o service *.go

#Stage 2 Production
FROM alpine:3.20 AS Production
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    nano

RUN apk add --no-cache tzdata 
WORKDIR /app
COPY --from=builder /app/service .
ENV TZ=Asia/Jakarta
# Allow ping command to be executed without sudo
RUN chmod u+s /bin/ping
RUN adduser -D appuser 
USER appuser
EXPOSE 8080
ENTRYPOINT [ "./service" ]
