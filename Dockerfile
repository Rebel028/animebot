FROM golang:alpine AS builder
WORKDIR /build
ADD go.mod .
COPY . .
RUN go build -o animebot main.go
FROM alpine
WORKDIR /build
COPY --from=builder /build/animebot /build/animebot
RUN ls -all
ENTRYPOINT [ "/build/animebot" ]