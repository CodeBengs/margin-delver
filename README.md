# Margin Delver

`margin-delver` adalah layanan backend Go dengan arsitektur modular, database migration menggunakan Goose, dan lapisan otentikasi internal.

## Ringkasan

- Bahasa: Go
- Versi Go: `1.23.0`
- Migration: Goose
- Database: MySQL
- Folder migrasi: `migrations/`
- Endpoint auth login: `POST /internal/v1/auth/login`

## Struktur proyek

- `cmd/` - entrypoint aplikasi dan CLI migrasi
- `lib/` - konfigurasi, logger, database, migrasi
- `migrations/` - file Goose migration
- `modules/` - domain logic, termasuk auth
- `provider/` - wiring repository/service/handler
- `doc/` - dokumentasi API dan setup

## Dependensi utama

Dependensi yang digunakan di aplikasi ini:

- `github.com/gin-gonic/gin` v1.11.0
- `github.com/go-sql-driver/mysql` v1.8.1
- `github.com/joho/godotenv` v1.5.1
- `github.com/pressly/goose/v3` v3.16.0
- `go.uber.org/fx` v1.24.0
- `go.uber.org/zap` v1.28.0
- `gorm.io/gorm` v1.31.1
- `gorm.io/driver/mysql` v1.6.0

Dependensi lainnya diatur melalui `go.mod` dan akan diunduh otomatis oleh Go.

## Setup lingkungan

1. Pastikan Go `1.23.0` sudah terinstall.
2. Buat file `.env` di root project. Contoh:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=
DB_NAME=margin_delver
DB_RUN_MIGRATIONS=true
DB_SEED_DEFAULT_USER=true
AUTH_DEFAULT_USERNAME=DELVERADMIN1
AUTH_DEFAULT_PASSWORD=delverAdmin1
AUTH_DEFAULT_NAME="Delver Administrator"
```

3. Jangan commit `.env` ke Git.
4. Jika ingin menggunakan Goose CLI global:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Build

Bangun aplikasi dan helper migrasi:

```bash
make build
```

Build untuk Windows:

```bash
make build-win
```

Build untuk macOS (ARM):

```bash
make build-mac
```

## Menjalankan aplikasi

Setelah build, jalankan:

```bash
./bin/margin-delver
```

Atau jalankan langsung tanpa build:

```bash
go run .
```

## Migrasi database

### Dengan binary migrasi

Setelah `make build`:

```bash
./bin/migrate status
./bin/migrate up
```

### Dengan `go run`

```bash
go run ./cmd/migrate status
go run ./cmd/migrate up
```

### Dengan Goose CLI langsung

```bash
goose -dir migrations mysql "root:@tcp(localhost:3306)/margin_delver?parseTime=true" status
goose -dir migrations mysql "root:@tcp(localhost:3306)/margin_delver?parseTime=true" up
```
## Setup tambahan

- `DB_RUN_MIGRATIONS=true` akan menjalankan migrasi otomatis jika aplikasi memanggil migrasi pada startup.
- Gunakan `DB_SEED_DEFAULT_USER=true` untuk seed default user jika tabel `users` kosong.
- `migrations/` sekarang adalah folder migrasi utama.

## Catatan

- Pastikan `go mod tidy` dijalankan saat menambahkan dependency baru.
- Jika butuh helper script untuk migrasi singkat, `./bin/migrate` sudah tersedia setelah build.
