FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY src .
RUN go build -o cook-bin ./main.go

FROM alpine:3.21.3
RUN adduser -D cook
COPY --from=builder /app/cook-bin /cook-bin
EXPOSE 8080
USER cook
ENTRYPOINT ["/cook-bin"]
CMD ["serve"]

