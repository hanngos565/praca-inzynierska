FROM golang:alpine as builder
RUN mkdir /build
ADD . /build
WORKDIR /build/cmd
RUN go mod download && go mod tidy && go fmt ../...
RUN apk add gcc musl-dev
RUN go test ../...
RUN go build -o main .

FROM alpine
COPY --from=builder /build/cmd/main /app/
WORKDIR /app
EXPOSE 8081
CMD ./main