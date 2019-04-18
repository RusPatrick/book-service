FROM golang:alpine AS builder
RUN apk add git
WORKDIR $GOPATH/src/github.com/ruspatrick/book-service
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o books-service .
FROM scratch
WORKDIR /app/bin
COPY --from=builder /go/src/github.com/ruspatrick/book-service/books-service .
COPY config.json .
CMD [ "/app/bin/books-service" ]