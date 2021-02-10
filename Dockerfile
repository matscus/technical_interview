# Builder
FROM golang:1.15.2-alpine3.12 as builder

RUN apk update && apk upgrade && \
    apk --update add git make

WORKDIR /application

COPY . .

RUN make engine

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /application 

WORKDIR /application 

EXPOSE 9443

COPY --from=builder /application/engine /application

CMD /application/engine