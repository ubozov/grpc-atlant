FROM golang:latest
WORKDIR /
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o grpc-atlant .


FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 ./grpc-atlant .

EXPOSE 10000

CMD ["./grpc-atlant"]