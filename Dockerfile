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

# Cài đặt chứng chỉ SSL và Timezone (Cần thiết cho Go)
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Ho_Chi_Minh

# Copy file thực thi
COPY --from=builder /app/server .

# CHỈ CẦN COPY THƯ MỤC STATIC (CSS/JS/IMG)
# (HTML đã được embed trực tiếp vào file chạy 'server' nhờ tính năng go:embed)
COPY --from=builder /app/static ./static

EXPOSE 8080
CMD ["./server"]
