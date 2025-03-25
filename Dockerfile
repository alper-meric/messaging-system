FROM golang:1.22-alpine AS builder

WORKDIR /app

# Bağımlılıkları kopyala ve indir
COPY go.mod go.sum ./
RUN go mod download

# Kaynak kodu kopyala
COPY . .

# Uygulamayı derle
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o messaging-system ./cmd/server

# Çalıştırma aşaması
FROM alpine:latest  

WORKDIR /app

# TLS sertifikaları için gerekli CA sertifikalarını kur
RUN apk --no-cache add ca-certificates

# Çalıştırılabilir dosyayı kopyala
COPY --from=builder /app/messaging-system .

# Konfigürasyon
COPY config.json .

# Uygulama kullanıcısını oluştur ve kullan
RUN adduser -D -g '' appuser
USER appuser

# Uygulamanın çalıştırılacağı portu belirt
EXPOSE 8080

# Uygulamayı çalıştır
CMD ["./messaging-system"] 