# margin-delver Codex Context

Project Go backend baru bernama `margin-delver`.

## Stack
- Go
- Gin
- Uber FX
- GORM MySQL
- godotenv
- Zap logger

## Current status
Server sudah jalan di `APP_PORT=3030`.

## Current files
- `main.go`
- `cmd/http_server.go`
- `lib/app_config.go`
- `lib/database.go`
- `lib/base_log.go`
- `doc/api.md`
- `Makefile`
- `router/v1/router.go`
- `middleware/middleware.go`
- `base/response/response.go`
- `base/pagination/pagination.go`
- `provider/repository.go`
- `provider/service.go`
- `provider/handler.go`
- `provider/seeder.go`
- `database/migrations/20260524160000_create_users_table.sql`
- `doc/migration.md`
- `modules/auth/auth_route.go`
- `modules/auth/auth_constant/auth_constant.go`
- `modules/auth/auth_dto/auth_request.go`
- `modules/auth/auth_dto/auth_response.go`
- `modules/auth/auth_handler/new_auth_handler.go`
- `modules/auth/auth_handler/auth_login_handler_impl.go`
- `modules/auth/auth_service/new_auth_service.go`
- `modules/auth/auth_service/auth_login_service_impl.go`
- `modules/auth/auth_entity/auth_model/auth_user.go`
- `modules/auth/auth_entity/auth_repository/new_auth_repository.go`
- `modules/auth/auth_entity/auth_repository/auth_login_repository_impl.go`
- `modules/auth/auth_entity/auth_repository/auth_seed_repository_impl.go`
- `modules/auth/auth_provider/auth_provider.go`
- `.env`

## Architecture target
Ikuti pattern `oms-service`:
- `cmd` untuk bootstrap server
- `lib` untuk config, db, logger
- `router/v1` untuk route registry
- `middleware` untuk CORS, logger, recovery
- `base/response` untuk response helper
- `provider` untuk dependency injection module
- `modules/<module_name>/<module_name>_dto` untuk request dan response DTO
- `modules/<module_name>/<module_name>_constant` untuk constants dan error code
- `modules/<module_name>/<module_name>_handler` untuk Gin handler
- `modules/<module_name>/<module_name>_service` untuk business logic
- `modules/<module_name>/<module_name>_entity` untuk repository interface, model, dan implementasi repository
- `modules/<module_name>/<module_name>_entity/<module_name>_model` untuk GORM model
- `modules/<module_name>/<module_name>_entity/<module_name>_repository` untuk implementasi GORM repository
- `modules/<module_name>/<module_name>_provider` untuk provider manual module

## Coding rules
- Jangan ubah module name `margin-delver`.
- Jangan hardcode port/database config, ambil dari `.env`.
- Jangan taruh route banyak di `cmd/http_server.go`.
- Route baru masuk ke `router/v1`.
- Response API wajib pakai `base/response`.
- Logger pakai `lib.BaseLog`, jangan `fmt.Println` kecuali bootstrap minimal.
- Database schema migration wajib memakai Goose di folder `database/migrations`.
- Jalankan migration lewat Makefile: `make migrate-up`, `make migrate-status`, `make migrate-down`.
- Goose memakai table `goose_db_version` sebagai migration tracking table.
- Jangan buat migration tracking table manual; Goose akan membuat dan mengisinya saat migration dijalankan.
- Jangan pakai GORM `AutoMigrate` untuk schema.
- Seeder default user hanya jalan kalau `DB_SEED_DEFAULT_USER=true`.
- Provider dipisah menjadi repository, service, handler, dan migration.
- Interface layer module pakai suffix `Interface` dan nama spesifik module/use-case, contoh `AuthServiceInterface`, `AuthRepositoryInterface`.
- Subfolder module wajib pakai prefix parent folder, contoh `auth_constant`, `auth_dto`, `auth_service`, `auth_handler`.
- File di dalam subfolder module juga wajib pakai prefix parent/layer, contoh `auth_request.go`, `auth_response.go`, `new_auth_service.go`.
- Kalau menambah dependency, jalankan `go mod tidy`.

## New module rules
Saat membuat module baru, wajib mengikuti struktur dan aturan berikut.

Contoh untuk module `<module_name>`:

```text
modules/<module_name>/
  <module_name>_constant/
    <module_name>_constant.go
  <module_name>_route.go
  <module_name>_dto/
    <module_name>_request.go
    <module_name>_response.go
  <module_name>_handler/
    new_<module_name>_handler.go
    <module_name>_<section>_handler_impl.go
  <module_name>_service/
    new_<module_name>_service.go
    <module_name>_<section>_service_impl.go
  <module_name>_entity/
    <module_name>_model/
      <module_name>_<entity>.go
    <module_name>_repository/
      new_<module_name>_repository.go
      <module_name>_<section>_repository_impl.go
  <module_name>_provider/
    <module_name>_provider.go
```

Rules module baru:
- Nama folder module wajib `snake_case`.
- Semua subfolder wajib memakai prefix nama module.
- Semua file di dalam module wajib memakai prefix nama module.
- Constant dan error code wajib masuk ke `<module_name>_constant`.
- Request dan response struct wajib masuk ke `<module_name>_dto`.
- Route module wajib dipisah di `<module_name>_route.go`.
- Handler wajib tipis: bind request, panggil service, return response.
- Handler constructor wajib dipisah di `new_<module_name>_handler.go`.
- Implementasi handler per section/action wajib dipisah ke file sendiri, contoh `auth_login_handler_impl.go`.
- Service interface, struct, dan constructor wajib dipisah di `new_<module_name>_service.go`.
- Implementasi service per section/action wajib dipisah ke file sendiri, contoh `auth_login_service_impl.go`.
- Implementasi service dibuat tipis: panggil repository, handle error seperlunya, lalu return response.
- Repository interface, struct, dan constructor wajib dipisah di `<module_name>_entity/<module_name>_repository/new_<module_name>_repository.go`.
- Implementasi query repository per section/action wajib dipisah ke file sendiri, contoh `auth_login_repository_impl.go`.
- Business logic wajib di `<module_name>_service`.
- Implementasi GORM repository wajib di `<module_name>_entity/<module_name>_repository`.
- GORM model wajib di `<module_name>_entity/<module_name>_model`.
- Interface layer wajib memakai suffix `Interface` dan nama spesifik module/use-case.
- Struct service dan repository wajib memakai nama spesifik, contoh `AuthService`, `AuthRepository`.
- Constructor tetap memakai nama standar `NewRepository`, `NewService`, dan `NewHandler`.
- Register constructor baru ke `provider/repository.go`, `provider/service.go`, dan `provider/handler.go`.
- Module provider manual wajib ada di `<module_name>_provider/<module_name>_provider.go`.
- Global provider hanya boleh register constructor dari module provider manual, bukan langsung dari service/repository/handler.
- Route baru wajib dibuat di `<module_name>_route.go`, lalu didaftarkan lewat `router/v1`, bukan di `cmd/http_server.go`.
- Response API wajib memakai `base/response`.
- Pagination/list response wajib memakai helper reusable dari `base/pagination` jika dibutuhkan.
- Migration schema baru wajib dibuat sebagai file Goose SQL di `database/migrations`.
- Migration yang sudah tercatat di `goose_db_version` tidak boleh diedit untuk perubahan baru; buat file migration baru.
- Cara menjalankan Goose migration wajib didokumentasikan di `doc/migration.md`.

Layer split rules:
- File `new_<module_name>_handler.go` hanya boleh berisi struct handler dan constructor.
- File `<module_name>_route.go` hanya boleh berisi route registration module.
- File `<module_name>_<section>_handler_impl.go` hanya boleh berisi method handler untuk section/action tersebut.
- File `new_<module_name>_service.go` hanya boleh berisi service interface, struct service, dan constructor.
- File `<module_name>_<section>_service_impl.go` hanya boleh berisi method service untuk section/action tersebut.
- File `new_<module_name>_repository.go` hanya boleh berisi repository interface, struct repository, dan constructor.
- File `<module_name>_<section>_repository_impl.go` hanya boleh berisi query/mapping repository untuk section/action tersebut.
- File `<module_name>_provider.go` hanya boleh berisi function provider manual seperti `NewRepositoryProvider`, `NewServiceProvider`, dan `NewHandlerProvider`.
- Module provider boleh punya initializer manual seperti `InitializeAuthHandler(log, cfg, db)` untuk dipakai route module.
- Jangan gabungkan constructor dengan implementation method dalam satu file.
- Jangan pakai nama generic seperti `ServiceInterface`, `RepositoryInterface`, `service`, atau `repository`; pakai nama spesifik seperti `AuthServiceInterface`, `AuthRepositoryInterface`, `AuthService`, `AuthRepository`.
- Service method mengikuti pola delegasi:

```go
func (service *AuthService) Login(
	ctx context.Context,
	request *authdto.LoginRequest,
) (*authdto.LoginResponse, error) {
	response, err := service.authRepository.Login(ctx, request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
```

## Done
- [x] Bootstrap server dengan Gin dan Uber FX.
- [x] Load config dari `.env` lewat `lib.AppConfig`.
- [x] Setup database GORM MySQL lewat `lib.NewDatabase`.
- [x] Setup logger Zap lewat `lib.BaseLog`.
- [x] Setup middleware CORS, request logger, dan recovery.
- [x] Setup base response helper di `base/response`.
- [x] Buat provider layer di `provider/repository.go`, `provider/service.go`, dan `provider/handler.go`.
- [x] Rapikan nama file `cmd/htpp_server.go` menjadi `cmd/http_server.go`.
- [x] Buat module `auth` mengikuti pattern folder `architecture-design.md`.
- [x] Rename subfolder `auth` agar memakai prefix parent folder.
- [x] Rename file module `auth` agar memakai prefix parent folder.
- [x] Buat model database `User`.
- [x] Tambahkan query repository berbasis GORM MySQL untuk `auth`.
- [x] Tambahkan request dan response DTO untuk login.
- [x] Tambahkan service login dengan validasi bcrypt.
- [x] Tambahkan handler login.
- [x] Tambahkan route auth lewat `router/v1`.
- [x] Tambahkan migration schema Goose di `database/migrations`.
- [x] Pisahkan seed default user ke `provider/seeder.go`.
- [x] Tambahkan optional seed default user dari `.env`.
- [x] Tambahkan dokumentasi endpoint dasar di `doc/api.md`.

## Current endpoint
- `POST /internal/v1/auth/login`

## TODO
- [x] Tambahkan response helper untuk created, bad request, unauthorized, not found, dan internal error.
- [x] Tambahkan struktur pagination reusable di `base`.
- [x] Sesuaikan provider DI dengan pola `NewRepository`, `NewService`, `NewHandler`, `NewSeeder`.
- [x] Sesuaikan naming interface module dengan suffix `Interface`.
- [ ] Tambahkan middleware auth untuk validasi Bearer token.
- [ ] Simpan token/session ke database atau cache.
- [ ] Tambahkan endpoint logout dan me/profile.
- [ ] Jalankan Goose migration di database lokal setelah diizinkan.
- [ ] Jaga unit test auth login tetap mencakup success, invalid request, invalid credentials, dan internal error.
- [ ] Jalankan `go mod tidy` hanya jika ada dependency baru.
- [ ] Jalankan build/test setelah diizinkan.

## Next task
Tambahkan middleware auth dan persistence token/session supaya hasil login bisa dipakai untuk endpoint private.
