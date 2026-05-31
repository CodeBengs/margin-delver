# OMS Service — Architecture Design Document

**Version:** 1.0  
**Date:** 2026-05-22  
**Module:** `oms-service`  
**Go Version:** 1.24.7

---

## Table of Contents

1. [Overview](#1-overview)
2. [System Context](#2-system-context)
3. [High-Level Architecture](#3-high-level-architecture)
4. [Directory Structure](#4-directory-structure)
5. [Module Anatomy](#5-module-anatomy)
6. [Dependency Injection Strategy](#6-dependency-injection-strategy)
7. [API Layer Design](#7-api-layer-design)
8. [Service Layer Design](#8-service-layer-design)
9. [Repository & Data Layer](#9-repository--data-layer)
10. [Middleware Architecture](#10-middleware-architecture)
11. [Router Architecture](#11-router-architecture)
12. [Configuration & Environment](#12-configuration--environment)
13. [Database Architecture](#13-database-architecture)
14. [Core Libraries (lib/)](#14-core-libraries-lib)
15. [Key Dependencies](#15-key-dependencies)
16. [Naming Conventions](#16-naming-conventions)
17. [Build & Deployment](#17-build--deployment)

---

## 1. Overview

**OMS Service** (Order Management System) adalah backend microservice berbasis Go yang mengelola operasi inti bisnis restoran/F&B, mencakup:

- **Master Data**: Menu, Cabang, Stasiun, Perusahaan, Printer, Meja
- **POS Integration**: Sinkronisasi data, Laporan, Voucher, Member Deposit
- **E-Commerce / Online Order**: Pesanan, Reservasi, Promosi, Pengiriman
- **CMS**: Manajemen konten & versioning
- **ESO (External Service Operations)**: Integrasi layanan pihak ketiga (kurir, payment, dll.)
- **Invoicing**: e-Invoice Malaysia (LHDN)

Arsitektur mengikuti pola **Layered Architecture** dengan elemen **Hexagonal/Ports & Adapters**, mendukung multi-tenant melalui mekanisme *server-code routing*.

---

## 2. System Context

```
┌─────────────────────────────────────────────────────────────────────┐
│                          External Clients                            │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌───────────────────┐   │
│  │ POS App  │  │ CMS App  │  │ ESO App  │  │ Online Order App  │   │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────────┬──────────┘   │
└───────┼─────────────┼─────────────┼──────────────────┼─────────────┘
        │             │             │                  │
        └─────────────┴─────────────┴──────────────────┘
                                │ HTTP REST
                       ┌────────▼────────┐
                       │   OMS Service   │
                       │   (Port 3011)   │
                       └────────┬────────┘
                                │
          ┌─────────────────────┼──────────────────────┐
          │                     │                      │
   ┌──────▼──────┐    ┌────────▼────────┐   ┌────────▼────────┐
   │  MySQL DB   │    │   AWS S3 /      │   │  Third-Party    │
   │  (Main +    │    │   OSS Storage   │   │  APIs (Kurir,   │
   │  Tenants)   │    └─────────────────┘   │  Firebase, ESO) │
   └─────────────┘                          └─────────────────┘
```

---

## 3. High-Level Architecture

```
┌──────────────────────────────────────────────────────────────────┐
│                         main.go (Uber FX)                        │
│                                                                  │
│  lib.NewLib()  →  provider.NewRepo()  →  provider.NewService()  │
│       ↓               ↓                      ↓                  │
│  provider.NewMw()  →  provider.NewHandler()  →  cmd.NewServer() │
└──────────────────────────────────────────────────────────────────┘
         │
         ▼
┌────────────────────────────────────────────────────────────────────┐
│                     HTTP Request Pipeline                          │
│                                                                    │
│  Request → Middleware Chain → Router → Handler → Service → Repo   │
│                                                                    │
│  Middleware:                                                       │
│  - URL Rewrite                                                     │
│  - Sentry Transaction                                              │
│  - Panic Recovery (+ Google Space notif)                          │
│  - JWT Auth / OMS Auth / Basic Auth / ESO Auth                    │
└────────────────────────────────────────────────────────────────────┘
         │
         ▼ (per domain)
┌──────────────────────────┐
│  Handler (Gin — thin)    │   ← validates request, calls service
├──────────────────────────┤
│  Service (Business Logic)│   ← orchestrates, validates rules
├──────────────────────────┤
│  Repository (Data Access)│   ← GORM queries, multi-DB routing
├──────────────────────────┤
│  Entity / Model (GORM)   │   ← struct tags, DB mapping
└──────────────────────────┘
```

---

## 4. Directory Structure

```
oms-service/
├── main.go                     # Aplikasi entry point (Uber FX bootstrap)
├── app.env                     # Environment configuration
├── app.local.env               # Override lokal (tidak di-commit)
├── Makefile                    # Build targets (linux, windows, mac)
├── go.mod / go.sum             # Module dependencies
│
├── cmd/                        # HTTP server lifecycle
│   └── http_server.go          # fx.Lifecycle hooks, graceful shutdown
│
├── lib/                        # Core infrastructure libraries
│   ├── config.go               # AppConfig struct (100+ fields via Viper)
│   ├── database.go             # MySQL multi-connection (DbMain + DbHost[])
│   ├── logger.go               # Zap structured logging
│   ├── storage.go              # AWS S3 / Aliyun OSS client
│   ├── sentry.go               # Sentry error tracking
│   └── lib.go                  # fx.Module provider
│
├── base/                       # Shared infrastructure
│   ├── constants/              # Domain constants (error codes, status, payment, promo types, dll.)
│   ├── helpers/                # 20+ utility packages (array, DB, datetime, encryption, Excel, cURL, dll.)
│   ├── provider/               # Uber FX DI wiring
│   │   ├── repository.go       # Mendaftarkan semua repository
│   │   ├── service.go          # Mendaftarkan semua service
│   │   ├── handler.go          # Mendaftarkan semua handler
│   │   └── middleware.go       # Mendaftarkan semua middleware
│   └── utility/
│       ├── ginx/               # Custom Gin extensions & response wrapper
│       └── erx/                # Custom error types
│
├── middleware/v1/              # HTTP middleware layer
│   ├── jwt/                    # JWT auth
│   ├── oms/                    # OMS authorization & basic auth
│   ├── eso_middleware/         # ESO auth variants (basic, internal, erp, ext, fs, qs, rest)
│   ├── panic_recovery/         # Panic recovery + Google Space notification
│   ├── sentry/                 # Sentry transaction middleware
│   └── url_rewrite/            # URL path rewriting
│
├── router/v1/                  # Route registration
│   ├── router.go               # Central router setup & route group helpers
│   ├── oms_router/             # OMS route groups
│   ├── pos_router/             # POS route groups
│   ├── cms_router/             # CMS versioned route groups
│   ├── eso_router/             # ESO route groups (qsv1, fsv1, extv1, erp)
│   ├── online_order_router/    # Online order route groups
│   └── not_found_router/       # 404 handler
│
├── modules/                    # Core business modules (30+)
│   └── <domain>/               # e.g., station/, branch/, company/, menu/
│       ├── dto/                # Request/Response DTOs
│       ├── handler/            # Gin handler (thin layer)
│       ├── service/            # Business logic
│       └── repository/
│           ├── entity/         # GORM model structs
│           └── mysql/          # MySQL repository implementation
│
├── oms_modules/                # OMS master data modules (Wire-based DI)
│   └── ms_<name>/              # e.g., ms_company/, ms_menu/
│       ├── ms_<name>_dto/
│       ├── ms_<name>_handler/
│       ├── ms_<name>_service/
│       ├── ms_<name>_entity/
│       │   ├── ms_<name>_model/
│       │   └── ms_<name>_repository/
│       └── ms_<name>_provider/ # wire.go + wire_gen.go
│
├── pos_modules/                # POS-specific modules (40+)
│   └── <pos_feature>/          # Sync, Voucher, Report, Member, dll.
│
├── cms_modules/                # CMS modules (15+)
│   └── <cms_feature>/
│
├── eso_modules/                # ESO integration modules (30+)
│   └── <eso_feature>/
│
└── docs/                       # Dokumentasi
    └── architecture-design.md  # Dokumen ini
```

---

## 5. Module Anatomy

Setiap modul mengikuti struktur yang konsisten. Contoh untuk `oms_modules/ms_company/`:

```
ms_company/
├── ms_company_dto/
│   ├── ms_company_request.go       # Struct input binding
│   └── ms_company_response.go      # Struct output serialization
│
├── ms_company_handler/
│   └── ms_company_handler.go       # Gin handler (hanya parsing & delegasi ke service)
│
├── ms_company_service/
│   └── ms_company_service.go       # Business logic & validasi
│
├── ms_company_entity/
│   ├── ms_company_model/
│   │   └── ms_company.go           # GORM entity struct
│   └── ms_company_repository/
│       └── ms_company_repository.go  # GORM query implementation
│
├── ms_company_provider/
│   ├── wire.go                     # Wire dependency declarations
│   └── wire_gen.go                 # Wire generated code
│
└── ms_company_route.go             # RouterGroupFunc registration
```

**Aliran data per request:**

```
HTTP Request
    │
    ▼
Handler.Method(c *gin.Context)
    ├── Bind & validate request DTO
    ├── Call service.Method(ctx, dto)
    │       ├── Business rule validation
    │       ├── Call repository.Query(params)
    │       │       └── GORM → MySQL
    │       └── Transform → response DTO
    └── ginx.Response(c, data, err)
```

---

## 6. Dependency Injection Strategy

OMS Service menggunakan **dua layer DI** yang bekerja bersama:

### 6.1 Uber FX (Application Level)

`main.go` menggunakan `fx.New()` untuk bootstrap seluruh aplikasi:

```go
fx.New(
    lib.NewLib(),               // Config, Logger, DB, Storage, Sentry
    provider.NewRepo(),         // 30+ repositories
    provider.NewService(),      // 30+ services
    provider.NewMw(),           // Middleware set
    provider.NewHandler(),      // 30+ handlers
    ginx.Engine,                // Gin engine
    fx.Invoke(router.Setup),    // Route registration
    fx.Invoke(cmd.NewServer),   // HTTP server start
)
```

### 6.2 Google Wire (Module Level)

Setiap `oms_modules/ms_<name>/ms_<name>_provider/wire.go` mendefinisikan dependensi untuk modul tersebut:

```go
// wire.go
var ProviderSet = wire.NewSet(
    NewRepository,
    NewService,
    NewHandler,
)

// wire_gen.go (generated)
func Initialize(db *lib.MySql, cfg *lib.AppConfig, log *lib.BaseLog) *Handler {
    repo := NewRepository(db, cfg)
    svc  := NewService(repo, log)
    return NewHandler(svc, log)
}
```

### 6.3 Provider Structs (base/provider/)

Empat struct container yang mendaftarkan semua komponen ke FX:

| Struct | File | Isi |
|--------|------|-----|
| `BaseRepository` | `repository.go` | 30+ constructor repository |
| `BaseService` | `service.go` | 30+ constructor service |
| `BaseMiddleware` | `middleware.go` | 6 middleware component |
| `BaseHandler` | `handler.go` | 30+ constructor handler |

---

## 7. API Layer Design

### 7.1 Route Groups & Authentication

| Prefix | Auth | Modul |
|--------|------|-------|
| `/internal/*` | JWT | Operasi internal |
| `/oms/*` | OMS Auth | Core OMS operations |
| `/cms/*` | Basic Auth | CMS management |
| `/v1/auth` | Public | Authentication |
| `/qsv1/*` | ESO QS Auth | ESO Quick Service |
| `/fsv1/*` | ESO FS Auth | ESO Full Service |
| `/extv1/*` | ESO Ext Auth | ESO External v1 |
| `/erp/*` | ERP Auth | ERP integration |
| `/v1/menu` | JWT | Menu API |
| `/v1/order` | JWT | Order API |
| `/v1/setting` | JWT | Settings API |
| `/v1/version` | Basic Auth | CMS Versioning |

### 7.2 Handler Pattern

Handler bersifat **thin** — hanya bertanggung jawab untuk:
1. Membaca & memvalidasi input (binding)
2. Memanggil service
3. Mengembalikan response via `ginx.Response()`

```go
type StationHandler struct {
    service StationServiceInterface
    log     *lib.BaseLog
}

func (h *StationHandler) FindAll(c *gin.Context) {
    var req dto.StationListRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        ginx.BadRequest(c, err)
        return
    }
    data, total, err := h.service.GetAll(c.Request.Context(), req)
    ginx.Response(c, gin.H{"data": data, "total": total}, err)
}
```

### 7.3 Response Format

Response standar menggunakan `ginx` wrapper:

```json
{
  "success": true,
  "message": "OK",
  "data": { ... },
  "meta": {
    "total": 100,
    "page": 1,
    "limit": 10
  }
}
```

---

## 8. Service Layer Design

Service layer mengandung semua **business logic** dan mengkoordinasikan repository:

```go
type CompanyService struct {
    companyRepo  CompanyRepositoryInterface
    branchRepo   BranchRepositoryInterface
    log          *lib.BaseLog
    cfg          *lib.AppConfig
}

func (s *CompanyService) CreateCompany(ctx context.Context, req dto.CreateCompanyRequest) (*dto.CompanyResponse, error) {
    // 1. Validasi business rule
    // 2. Enrich data (generate ID, set audit fields)
    // 3. Koordinasi antar repository
    // 4. Transform ke response DTO
}
```

**Tanggung jawab service:**
- Validasi business rule (bukan hanya input validation)
- Orkestrasi multi-repository
- Transformasi data entity → response DTO
- Logging & error wrapping
- Tidak boleh mengetahui detail HTTP (no `*gin.Context`)

---

## 9. Repository & Data Layer

### 9.1 Repository Pattern

```go
type CompanyRepository struct {
    mysql *lib.MySql
    cfg   *lib.AppConfig
}

// Interface untuk testability
type CompanyRepositoryInterface interface {
    FindByID(serverCode, dbName string, id int64) (*entity.Company, error)
    FindAll(serverCode, dbName string, params dto.ListParams) ([]entity.Company, int64, error)
    Create(serverCode, dbName string, data *entity.Company) error
    Update(serverCode, dbName string, data *entity.Company) error
    SoftDelete(serverCode, dbName string, id int64, deletedBy string) error
}
```

### 9.2 Multi-Tenant Database Routing

OMS mendukung **multi-tenant** melalui `serverCode` routing:

```go
func (r *CompanyRepository) getDB(serverCode string) *gorm.DB {
    if serverCode == "" || serverCode == "MAIN" {
        return r.mysql.DbMain
    }
    // Lookup koneksi tenant dari pool
    for _, db := range r.mysql.DbHost {
        if db.ServerCode == serverCode {
            return db.Conn
        }
    }
    return r.mysql.DbMain // fallback
}
```

### 9.3 Entity / GORM Model

```go
// Contoh entity dengan audit fields standard
type MsCompany struct {
    CompanyID    int64     `gorm:"column:company_id;primaryKey"`
    CompanyCode  string    `gorm:"column:company_code"`
    CompanyName  string    `gorm:"column:company_name"`
    FlagActive   bool      `gorm:"column:flag_active"`
    CreatedBy    string    `gorm:"column:created_by"`
    CreatedDate  time.Time `gorm:"column:created_date;autoCreateTime"`
    EditedBy     string    `gorm:"column:edited_by"`
    EditedDate   time.Time `gorm:"column:edited_date;autoUpdateTime"`
    SyncDate     time.Time `gorm:"column:sync_date"`
}

func (MsCompany) TableName() string { return "ms_company" }
```

**Konvensi entity:**
- Soft delete via `flag_active` (bukan `deleted_at`)
- Audit fields: `created_by`, `created_date`, `edited_by`, `edited_date`, `sync_date`
- Table prefix: `ms_` untuk master data

---

## 10. Middleware Architecture

Middleware diterapkan secara berlapis sesuai route group:

```
Request
  │
  ▼
[1] URL Rewrite Middleware       — normalisasi path
  │
  ▼
[2] Sentry Transaction           — mulai tracking transaksi
  │
  ▼
[3] Panic Recovery               — catch panic, kirim notif ke Google Space
  │
  ▼
[4] Auth Middleware (per group)
  │  ├── JWT Middleware           — validasi Bearer token
  │  ├── OMS Auth Middleware      — OMS-specific auth
  │  ├── Basic Auth Middleware    — username:password base64
  │  └── ESO Auth Variants        — qsv1, fsv1, extv1, erp, rest
  │
  ▼
[5] Route Handler
```

### Middleware Detail

| Middleware | Paket | Fungsi |
|-----------|-------|--------|
| `jwt` | `middleware/v1/jwt` | Decode & validasi JWT Bearer token |
| `oms` | `middleware/v1/oms` | OMS tenant authorization |
| `oms/basic_middleware` | — | Basic auth untuk endpoint OMS |
| `eso_middleware/basic_middleware` | — | Basic auth ESO |
| `eso_middleware/eso_internal_middleware` | — | Auth internal ESO |
| `eso_middleware/erp_middleware` | — | ERP system auth |
| `eso_middleware/extv1_middleware` | — | External partner auth v1 |
| `eso_middleware/fsv1_middleware` | — | Full service auth |
| `eso_middleware/qsv1_middleware` | — | Quick service auth |
| `eso_middleware/rest_auth_middleware` | — | REST auth |
| `panic_recovery` | — | Catch panics, notify Google Space |
| `sentry` | — | Sentry transaction per request |
| `url_rewrite` | — | Rewrite URL paths sebelum routing |

---

## 11. Router Architecture

### 11.1 Route Registration Flow

```
main.go
  └── fx.Invoke(router.RegisterAll)
        ├── RegisterOmsRouter(baseHandler, baseMw, engine)
        ├── RegisterPosRouter(...)
        ├── RegisterCmsRouter(...)
        ├── RegisterEsoRouter(...)
        └── RegisterOnlineOrderRouter(...)
```

### 11.2 Route Group Pattern

Setiap modul mengekspos `RouterGroupFunc`:

```go
// ms_company_route.go
func RouterGroupWithDepsCompany(
    h *ms_company_handler.Handler,
    mw *middleware.BaseMw,
) ginx.RouterGroupFunc {
    return func(rg *gin.RouterGroup) {
        rg.GET("", h.FindAll)
        rg.GET("/:id", h.FindByID)
        rg.POST("", h.Create)
        rg.PUT("/:id", h.Update)
        rg.DELETE("/:id", h.Delete)
    }
}
```

---

## 12. Configuration & Environment

Konfigurasi menggunakan **Viper** membaca dari `app.env`:

### 12.1 Kategori Konfigurasi

| Kategori | Key Prefix | Keterangan |
|----------|-----------|------------|
| App | `APP_*` | Port, environment, timezone |
| Database | `MAIN_DB_*` | Host, port, name, credentials |
| Security | `APP_JWT_SECRET`, `SECURITY_KEY` | JWT & enkripsi |
| Storage | `STORAGE_UPLOAD_*` | S3/OSS credentials & bucket |
| ESO API | `ESO_API_URL` | ESO service endpoint |
| Firebase | `FIREBASE_URL` | Push notification |
| Courier | `GOSEND_*`, `GRAB_*`, `PAXEL_*`, `LALAMOVE_*`, `BLITZ_*` | Delivery services |
| Invoice | `INVOICE_BASE_URL`, `INVOICE_FTP_*` | Malaysia e-Invoice (LHDN) |
| Auth | `OMS_BASIC_AUTH_*` | Basic auth credentials |
| Server | `SERVER_TIME_OUT` | Request timeout (default 60s) |
| POS Single Server | `POS_SINGLE_SERVER_SECRET_KEY` | POS single server auth |

### 12.2 AppConfig Struct

```go
// lib/config.go
type AppConfig struct {
    AppEnv        string `mapstructure:"APP_ENV"`
    AppPort       string `mapstructure:"APP_PORT"`
    TimeZone      string `mapstructure:"TZ"`

    MainDBHost    string `mapstructure:"MAIN_DB_HOST"`
    MainDBPort    string `mapstructure:"MAIN_DB_PORT"`
    MainDBName    string `mapstructure:"MAIN_DB_NAME"`
    MainDBUser    string `mapstructure:"MAIN_DB_USERNAME"`
    MainDBPass    string `mapstructure:"MAIN_DB_PASSWORD"`

    AllowExternalServer bool   `mapstructure:"ALLOW_EXTERNAL_SERVER"`
    JwtSecret           string `mapstructure:"APP_JWT_SECRET"`
    SecurityKey         string `mapstructure:"SECURITY_KEY"`
    // ... 100+ fields lainnya
}
```

---

## 13. Database Architecture

### 13.1 Multi-Connection Setup

```
lib/database.go
  │
  ├── DbMain (*gorm.DB)           — koneksi ke database utama (esb_main)
  │                                  berisi data master tenant / server codes
  │
  └── DbHost []*TenantDB          — pool koneksi ke database tenant
        ├── ServerCode: "OUTLET1"
        │   Conn: *gorm.DB  →  MySQL server tenant 1
        ├── ServerCode: "OUTLET2"
        │   Conn: *gorm.DB  →  MySQL server tenant 2
        └── ...
```

### 13.2 Database Initialization Flow

```
1. Baca config (MAIN_DB_*)
2. Buka koneksi ke DbMain
3. Jika ALLOW_EXTERNAL_SERVER=true:
   a. Query daftar server codes dari DbMain
   b. Buka koneksi ke setiap server tenant
   c. Simpan ke DbHost pool
4. Set connection pool params (MaxIdleConns, MaxOpenConns, ConnMaxLifetime)
```

### 13.3 Table Naming Convention

| Prefix | Domain | Contoh |
|--------|--------|--------|
| `ms_` | Master data | `ms_company`, `ms_branch`, `ms_menu` |
| `tr_` | Transaksi | `tr_order`, `tr_payment` |
| `rpt_` | Report | `rpt_sales_daily` |
| `cfg_` | Konfigurasi | `cfg_setting`, `cfg_printer` |

---

## 14. Core Libraries (lib/)

### 14.1 Logger (`lib.BaseLog`)

Menggunakan **Uber Zap** dengan output JSON:

```go
type BaseLog struct {
    *zap.SugaredLogger
}

// Usage
log.Infow("Company created", "company_id", id, "created_by", user)
log.Errorw("DB error", "error", err, "query", sql)
```

### 14.2 Database (`lib.MySql`)

```go
type MySql struct {
    DbMain *gorm.DB
    DbHost []*TenantDB
    Cfg    *AppConfig
}
```

### 14.3 Storage (`lib.Storage`)

Mendukung AWS S3 dan Aliyun OSS:

```go
type Storage struct {
    Client   *s3.Client
    Bucket   string
    BaseURL  string
}

// Methods: Upload, Delete, GetSignedURL
```

### 14.4 Sentry (`lib.Sentry`)

Error tracking dengan DSN dari konfigurasi:

```go
// Inisialisasi di lib/sentry.go
sentry.Init(sentry.ClientOptions{
    Dsn:         cfg.SentryDSN,
    Environment: cfg.AppEnv,
    Release:     cfg.AppVersion,
})
```

---

## 15. Key Dependencies

| Dependency | Version | Fungsi |
|-----------|---------|--------|
| `gin-gonic/gin` | v1.11.0 | HTTP web framework |
| `gorm.io/gorm` | v1.31.1 | ORM layer |
| `gorm.io/driver/mysql` | latest | MySQL driver |
| `go.uber.org/fx` | v1.24.0 | Application DI & lifecycle |
| `go.uber.org/zap` | v1.27.1 | Structured logging |
| `google/wire` | v0.7.0 | Compile-time DI code generation |
| `spf13/viper` | v1.21.0 | Configuration management |
| `dgrijalva/jwt-go` | v3.2.0 | JWT authentication |
| `getsentry/sentry-go` | v0.40.0 | Error monitoring |
| `aws/aws-sdk-go` | latest | AWS S3 storage |
| `aliyun/aliyun-oss-go-sdk` | latest | Aliyun OSS storage |
| `xuri/excelize/v2` | latest | Excel file generation |
| `pkg/sftp` | latest | SFTP file transfer |
| `shopspring/decimal` | latest | Decimal arithmetic |
| `go-playground/validator` | latest | Struct validation |

---

## 16. Naming Conventions

| Elemen | Konvensi | Contoh |
|--------|---------|--------|
| **Directory** | `snake_case` | `ms_company/`, `pos_modules/` |
| **File** | `snake_case.go` | `ms_company_service.go` |
| **Package** | lowercase, no underscore | `package mscompanyservice` |
| **Struct** | `PascalCase` | `CompanyService`, `MsCompany` |
| **Interface** | `PascalCase` + `Interface` suffix | `CompanyServiceInterface` |
| **Method** | `PascalCase` | `FindByID`, `CreateCompany` |
| **Variable** | `camelCase` | `serverCode`, `dbName` |
| **JSON tag** | `camelCase` | `json:"companyName"` |
| **DB column tag** | `snake_case` | `gorm:"column:company_name"` |
| **Constants** | `UPPER_SNAKE_CASE` | `STATUS_ACTIVE`, `ERR_NOT_FOUND` |
| **Test file** | `*_test.go` | `ms_company_service_test.go` |

---

## 17. Build & Deployment

### 17.1 Makefile Targets

```makefile
make build        # Build untuk Linux AMD64
make build-win    # Build untuk Windows
make build-mac    # Build untuk macOS ARM64
make run          # Build & jalankan aplikasi
```

### 17.2 HTTP Server

- **Port:** `APP_PORT` (default `3011`)
- **Timeout:** `SERVER_TIME_OUT` detik (default `60`)
- **Profiling:** pprof endpoint di port `:6061`
- **Graceful Shutdown:** context-based shutdown via fx lifecycle hooks

### 17.3 Application Startup Sequence

```
1. fx.New() mulai bootstrap
2. lib.NewLib() → init Config (Viper) → Logger (Zap) → DB → Storage → Sentry
3. provider.NewRepo() → konstruksi semua repository dengan injeksi DB & Config
4. provider.NewService() → konstruksi semua service dengan injeksi repo & logger
5. provider.NewMw() → konstruksi middleware set
6. provider.NewHandler() → konstruksi semua handler
7. ginx.Engine → buat Gin engine
8. router.RegisterAll() → daftarkan semua route groups
9. cmd.NewServer() → Start HTTP server (fx.Lifecycle OnStart hook)
10. [Graceful shutdown via fx.Lifecycle OnStop hook]
```

---

## Appendix: Domain Module Inventory

### Core Modules (`modules/`)
Branch, Company, Station, Menu, User, Payment Method, Printer, Table, Voucher, Reservation, Promotion, Order, Product, Category, Modifier, Tax, Discount, Shift, Report, Sync, Member, Loyalty, Delivery, Courier, Setting, License, Device, Permission, Role

### OMS Master Data (`oms_modules/ms_*`)
Company, Branch, Menu, Station, Printer, Table, Payment, Tax, Discount, Modifier, Category, Product, Voucher, Setting, User, Role, Permission, Device

### POS Modules (`pos_modules/`)
Sync data, Voucher management, Member deposit, Daily reports, Stock reports, Settlement, Shift reports, Transaction history, Item reports, Modifier reports

### CMS Modules (`cms_modules/`)
Versioning, Company config, Menu management, Branch config, Master data CMS, Promotions CMS

### ESO Modules (`eso_modules/`)
Online orders, Delivery tracking, Promotions, Reservations, Menu availability, Branch availability, Payment processing, Webhook handlers

---

*Dokumen ini dibuat secara otomatis berdasarkan analisis kode sumber. Perbarui setiap kali ada perubahan arsitektur signifikan.*
