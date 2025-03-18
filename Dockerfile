FROM golang:latest AS builder
WORKDIR /app
COPY src .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
RUN adduser -D cook
RUN mkdir /myapp && chown -R cook /myapp
WORKDIR /myapp
COPY --from=builder /app/main .
EXPOSE 8080
USER cook
ENTRYPOINT ["./main"]
CMD ["serve"]

