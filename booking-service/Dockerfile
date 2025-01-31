# Build Stage
FROM golang:1.22.0 AS builder

# ตั้ง working directory
WORKDIR /app

# คัดลอก go.mod และ go.sum ก่อน เพื่อใช้ในการติดตั้ง dependency
COPY go.mod go.sum ./
RUN go mod download

# คัดลอก source code ทั้งหมด
COPY . .

# สร้างแอปพลิเคชันแบบ static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o booking-service ./cmd/main.go

# ขั้นตอนที่ 2: สร้าง final image ด้วย Alpine
FROM alpine:latest

# ติดตั้ง tzdata และ libc6-compat เพื่อรองรับ timezone และ binary dependencies
RUN apk add --no-cache tzdata libc6-compat

# ตั้งค่า timezone เป็น Asia/Bangkok
ENV TZ=Asia/Bangkok


# ติดตั้ง libc6-compat เพื่อรองรับ binary dependencies
RUN apk add --no-cache libc6-compat

# สร้าง working directory
WORKDIR /app

# คัดลอก binary ที่สร้างจากขั้นตอน builder
COPY --from=builder /app/booking-service /app/booking-service

# คัดลอกไฟล์ .env ไปยัง container
COPY .env /app/.env

# กำหนดคำสั่งให้รันแอปพลิเคชัน
CMD ["./booking-service"]

# หากต้องการ expose port ที่แอปพลิเคชันรันอยู่
EXPOSE 50051
