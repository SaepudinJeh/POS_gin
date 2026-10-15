# POS Backend

REST API backend untuk aplikasi Point of Sale (POS), dibangun dengan Go, Gin, dan GORM.

---

## Persyaratan

- Go 1.21+
- PostgreSQL 13+
- Git
- Air (opsional, untuk hot reload) — `go install github.com/air-verse/air@latest`
- gotestsum (opsional, untuk output test cantik) — `go install gotest.tools/gotestsum@latest`

Pastikan `$GOPATH/bin` sudah ada di `PATH`.

---

## Setup

### 1. Clone repository

```bash
git clone <repository-url>
cd POS
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Setup environment

Copy file contoh, lalu sesuaikan nilainya:

```bash
cp .env.example .env
```

Lihat `.env.example` untuk daftar variabel yang dibutuhkan.

### 4. Setup database

Buat database PostgreSQL:

```sql
CREATE DATABASE pos_db;
```

---

## Menjalankan Aplikasi

### Development (hot reload)

```bash
air
```

### Development (tanpa hot reload)

```bash
go run cmd/server/main.go
```

### Production

Build binary:

```bash
go build -o bin/pos-api cmd/server/main.go
```

Jalankan:

```bash
export APP_ENV=production
export GIN_MODE=release
./bin/pos-api
```

---

## Testing

### Test semua package

```bash
go test ./...
```

### Output lebih cantik (opsional, butuh gotestsum)

```bash
gotestsum --format pkgname
```

### Test package tertentu

```bash
go test ./internal/services/
```

### Test satu fungsi

```bash
gotestsum --format pkgname -- -run TestCategoryService_Create ./internal/services/
```

### Force rerun (bypass cache)

```bash
go test -count=1 ./...
```

---

## Coverage

Generate:

```bash
gotestsum --format pkgname -- -coverprofile coverage.out ./...
```

Lihat ringkas:

```bash
go tool cover -func coverage.out
```

Lihat visual di browser:

```bash
go tool cover -html coverage.out
```