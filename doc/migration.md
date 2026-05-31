# Migration

Project ini memakai Goose untuk database migration.

Folder migration:

```text
migrations
```

## Install Goose

Jika Goose belum tersedia di mesin lokal:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Pastikan folder binary Go ada di `PATH`.

## Run Migration

Dengan nilai `.env` saat ini:

```text
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=
DB_NAME=margin_delver
```

Untuk menjalankan migrasi otomatis saat aplikasi start, tambahkan:

```text
DB_RUN_MIGRATIONS=true
```

Run migration via Makefile:

```bash
make migrate-up
```

Run migration:

```bash
goose -dir migrations mysql "root:@tcp(localhost:3306)/margin_delver?parseTime=true" up
```

Jika password database tidak kosong, format DSN:

```bash
goose -dir migrations mysql "username:password@tcp(host:port)/database_name?parseTime=true" up
```

## Check Status

Via Makefile:

```bash
make migrate-status
```

Command Goose langsung:

```bash
goose -dir migrations mysql "root:@tcp(localhost:3306)/margin_delver?parseTime=true" status
```

## Rollback

Rollback migration terakhir:

```bash
make migrate-down
```

Command Goose langsung:

```bash
goose -dir migrations mysql "root:@tcp(localhost:3306)/margin_delver?parseTime=true" down
```

Rollback ke version tertentu:

```bash
goose -dir migrations mysql "root:@tcp(localhost:3306)/margin_delver?parseTime=true" down-to 0
```

## Create Migration

Contoh membuat file migration baru:

```bash
make migrate-create name=create_example_table
```

Command Goose langsung:

```bash
goose -dir migrations create create_example_table sql
```

Catatan:

- Migration schema tidak dijalankan otomatis ketika server start.
- Goose otomatis membuat table `goose_db_version` untuk mencatat migration yang sudah sukses.
- File migration yang version-nya sudah tercatat di `goose_db_version` tidak akan dijalankan ulang saat `goose up`.
- Jika `DB_RUN_MIGRATIONS=true`, aplikasi akan menjalankan goose migration secara otomatis saat start.
- Jangan pakai GORM `AutoMigrate` untuk schema.
- Seed default user dipisah dari migration dan dikontrol lewat `.env`.
