FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/server ./cmd/server

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /out/server /app/server
COPY internal/templates /app/internal/templates
COPY static /app/static
COPY migrations /app/migrations
ENV PORT=8080
EXPOSE 8080
CMD ["/app/server"]
