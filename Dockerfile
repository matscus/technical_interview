# Builder
FROM golang:1.15.2-alpine3.12 as builder

WORKDIR /application

RUN apk update && apk upgrade && \
    apk --update add git make bash openssl &&\
    openssl req -newkey rsa:2048 -x509 -nodes -keyout /application/server.key -new -out /application/server.pem -subj "/C=RU/ST=interview/L=docker/O=matscus/OU=IT/CN=localhost" -sha256 -days 3650

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

COPY --from=builder /application/server.key /application

COPY --from=builder /application/server.pem /application

CMD /application/engine