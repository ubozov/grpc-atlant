FROM golang:latest

WORKDIR /

COPY . .

RUN go build -o grpc-atlant .

EXPOSE 10000

CMD ["./grpc-atlant"]