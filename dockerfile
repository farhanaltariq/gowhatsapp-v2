FROM golang:1.21.5-alpine3.18
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main .
EXPOSE 3000
CMD ["/app/main"]
