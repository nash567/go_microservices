FROM alpine:latest

RUN mkdir /app
WORKDIR /app
COPY authApp /app
COPY config.yml /app

CMD ["/app/authApp"]