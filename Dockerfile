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
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Ho_Chi_Minh

# Copy file chạy (đã ngầm chứa toàn bộ HTML nhờ lệnh go:embed)
COPY --from=builder /app/server .

# Chỉ cần copy duy nhất thư mục static (vì CSS/JS cần đọc từ bên ngoài)
COPY --from=builder /app/static ./static

EXPOSE 8080
CMD ["./server"]
