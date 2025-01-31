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
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o user-management-service ./cmd/main.go

# ขั้นตอนที่ 2: สร้าง final image ด้วย Alpine
FROM alpine:latest

# ติดตั้ง libc6-compat เพื่อรองรับ binary dependencies
RUN apk add --no-cache libc6-compat

# สร้าง working directory
WORKDIR /app

# คัดลอก binary ที่สร้างจากขั้นตอน builder
COPY --from=builder /app/user-management-service /app/user-management-service

# คัดลอกไฟล์ .env ไปยัง container
COPY .env /app/.env

# กำหนดคำสั่งให้รันแอปพลิเคชัน
CMD ["./user-management-service"]

# หากต้องการ expose port ที่แอปพลิเคชันรันอยู่
EXPOSE 50051


