FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/random-marc .

FROM scratch

COPY --from=builder /out/random-marc /random-marc

ENTRYPOINT ["/random-marc"]
CMD ["-count", "1"]
