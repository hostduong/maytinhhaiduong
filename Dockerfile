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

# Copy file chạy
COPY --from=builder /app/server .

# --- COPY ĐẦY ĐỦ CÁC THƯ MỤC HỆ THỐNG ---
COPY --from=builder /app/giao_dien ./giao_dien
COPY --from=builder /app/giao_dien_admin ./giao_dien_admin
COPY --from=builder /app/giao_dien_master ./giao_dien_master

# [THÊM DÒNG NÀY]: Copy thư mục chứa file CSS/JS dùng chung
COPY --from=builder /app/static ./static
# ------------------------------------------

EXPOSE 8080
CMD ["./server"]
