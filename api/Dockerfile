FROM golang:alpine as builder 

RUN mkdir /build 

ADD . /build/

WORKDIR /build 

RUN go build -o main .

#stage2 

FROM alpine

RUN adduser -S -D -H -h /app appuser

USER appuser

COPY . /app/

COPY --from=builder /build/main /app/

WORKDIR /app/

CMD ["./main"]
