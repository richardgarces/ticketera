# Dockerfile para Ticketera Web
FROM golang:1.22-alpine as builder
WORKDIR /app
COPY . .
RUN cd cmd/ticketera && go build -o /ticketera

FROM alpine:latest
WORKDIR /app
COPY --from=builder /ticketera /ticketera
COPY web/ web/
EXPOSE 8080
CMD ["/ticketera"]
