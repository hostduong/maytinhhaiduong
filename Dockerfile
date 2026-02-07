# Bước 1: Build
FROM golang:1.20-alpine as builder
WORKDIR /app
COPY . .
ENV CGO_ENABLED=0 
RUN go mod tidy
RUN go build -o server main.go

# Bước 2: Run
FROM alpine:latest
WORKDIR /root/
RUN apk add --no-cache ca-certificates

# Copy file chạy (Đã có)
COPY --from=builder /app/server .

# --- [THÊM DÒNG NÀY ĐỂ COPY GIAO DIỆN] ---
COPY --from=builder /app/giao_dien ./giao_dien
# ------------------------------------------

EXPOSE 8080
CMD ["./server"]
