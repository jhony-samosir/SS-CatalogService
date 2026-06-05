# SS-CatalogService

## Overview

`SS-CatalogService` adalah microservice inti yang mengelola katalog produk dan data taksonomi untuk platform e-commerce SamStore. Dibangun dengan **Go 1.26**, service ini dirancang untuk performa tinggi dalam menangani pencarian produk, indexing, penentuan harga (pricing), update inventaris, dan varian produk multi-level.

Service ini terintegrasi dengan **Meilisearch** untuk pencarian teks penuh (full-text search) yang sangat cepat, serta menggunakan **Bigcache** sebagai in-memory cache lokal untuk mengurangi beban query read ke database PostgreSQL.

Pola **Hexagonal Architecture** (atau Clean Architecture) digunakan secara ketat untuk memisahkan domain logic, usecase, delivery (HTTP handlers), dan infrastruktur (database, cache, message broker).

---

## Tech Stack

| Kategori       | Teknologi                              |
| -------------- | -------------------------------------- |
| Runtime        | Go (version 1.26+)                     |
| Web Framework  | Gin-Gonic                              |
| Database / ORM | PostgreSQL / GORM                      |
| Search Engine  | Meilisearch (meilisearch-go)           |
| Caching        | Bigcache (allegro/bigcache)            |
| Message Broker | RabbitMQ (amqp091-go)                  |
| Telemetry      | OpenTelemetry (dengan Gin instrumentation) |

---

## Arsitektur

Service ini menerapkan **Clean Hexagonal Architecture**:

```text
SS-CatalogService/
├── cmd/
│   ├── api/                  # Entry point untuk HTTP REST API server
│   └── seed/                 # Utility untuk inisialisasi data awal (seeding)
├── config/                   # Parsing environment variables (JWT, DB, RabbitMQ, Meilisearch)
├── db/
│   └── migrations/           # Skrip migrasi SQL berversi (golang-migrate)
├── internal/
│   ├── delivery/             # HTTP Handlers (Gin) dan Middleware
│   │   └── http/
│   │       ├── router.go     # Definisi route API
│   │       ├── middleware/   # Auth, CorrelationID, Logger
│   │       └── v1/           # Handlers: Product, Variant, Inventory, Category, Brand, dll
│   ├── domain/               # Core entities, constants, dan port interfaces (repository & usecase)
│   ├── infrastructure/       # Implementasi driver eksternal (DB, Cache, Meilisearch, RabbitMQ)
│   ├── repository/           # Implementasi GORM dan Meilisearch repository
│   ├── usecase/              # Core business logic
│   └── worker/               # Background worker (misal: Outbox publisher)
└── pkg/                      # Utilities (logger, response formatter)
```

---

## Fitur Utama

- **Taksonomi Katalog Lengkap**: Mendukung produk, brand, kategori, atribut, opsi, dan varian produk multi-level.
- **Harga & Inventaris**: Melacak harga produk, penyesuaian stok, dan riwayat harga.
- **Pencarian & Indexing Berkinerja Tinggi**: Sinkronisasi otomatis entitas produk dengan Meilisearch. Mendukung *faceted search* (filter berdasarkan harga, kategori, brand, atribut).
- **In-Memory Caching**: Penggunaan Bigcache untuk mengurangi latensi database pada path baca yang sering diakses (high-traffic read paths).
- **Event-Driven Outbox**: Publikasi event transaksional (seperti stok berubah atau harga berubah) ke RabbitMQ tanpa risiko kegagalan parsial (*dual write problem*).
- **Review & Rating**: Mengelola ulasan pengguna terhadap produk dan rangkuman rating.
- **Bundel Produk & Digital License**: Dukungan untuk produk paket (bundle) dan produk digital (pengunggahan file/lisensi).

---

## API Endpoints

Semua endpoint API diawali dengan `/api/catalog/v1` dan diproteksi (kecuali yang ditandai khusus).

### Products (`/products`)
| Method | Endpoint                        | Deskripsi                                  |
| ------ | ------------------------------- | ------------------------------------------ |
| POST   | `/products`                     | Buat produk baru (Auth)                    |
| PUT    | `/products/:id`                 | Update produk (Auth)                       |
| GET    | `/products`                     | Ambil daftar produk                        |
| GET    | `/products/search`              | Cari produk (via Meilisearch)              |
| GET    | `/products/faceted-search`      | Cari dengan facet filter                   |
| GET    | `/products/:id`                 | Detail produk beserta variannya            |
| GET    | `/products/:id/price-history`   | Riwayat perubahan harga produk             |

### Inventory (`/inventory`)
| Method | Endpoint               | Deskripsi                                    |
| ------ | ---------------------- | -------------------------------------------- |
| GET    | `/inventory`           | Cek status stok                              |
| POST   | `/inventory/adjust`    | Sesuaikan stok manual atau karena pembelian  |

### Categories, Brands, Attributes, Tags
Endpoints CRUD lengkap untuk tiap entitas pendukung katalog (dibutuhkan JWT untuk operasi tulis).
- `/categories`
- `/brands`
- `/attributes`
- `/tags`

### Reviews (`/reviews`)
| Method | Endpoint                        | Deskripsi                                  |
| ------ | ------------------------------- | ------------------------------------------ |
| POST   | `/reviews`                      | Submit ulasan baru (Auth)                  |
| GET    | `/reviews/product/:id`          | Ambil ulasan untuk satu produk             |
| GET    | `/reviews/product/:id/summary`  | Dapatkan agregat rating (1-5 bintang)      |
| PATCH  | `/reviews/:id/status`           | Update status ulasan (approve/reject) (Auth)|

### Lainnya
- `/variants`: Mengelola varian produk (SKU).
- `/bundles`: Mengelola bundel produk.
- `/warehouses`: Mengelola gudang/lokasi penyimpanan.
- `/sellers`: Registrasi/detail penjual (seller/toko).
- `/imports`: Impor massal data katalog.
- `/digital`: Pengelolaan lisensi dan file digital.
- `/audit-logs`: Log audit operasi katalog.
- `/health`: Cek kesehatan sistem (digunakan oleh API Gateway).

---

## Environment Variables

| Variable                      | Deskripsi                                      | Wajib |
| ----------------------------- | ---------------------------------------------- | ----- |
| `APP_PORT`                    | Port untuk Gin HTTP server (default: 8081)     | Tidak |
| `ENVIRONMENT`                 | Mode runtime (`development`, `production`)     | ✅    |
| `DB_DSN`                      | DSN PostgreSQL                                 | ✅    |
| `RABBITMQ_URL`                | URL koneksi ke RabbitMQ                        | ✅    |
| `MEILI_URL`                   | Endpoint host Meilisearch                      | ✅    |
| `MEILI_MASTER_KEY`            | Master key Meilisearch API                     | ✅    |
| `JWT_SECRET`                  | Fallback key lokal, tapi validasi utama oleh Gateway | ✅ |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | Endpoint OTel Collector                        | ✅    |

---

## Instalasi & Menjalankan

### Prasyarat

- Go 1.26+
- PostgreSQL
- Meilisearch
- RabbitMQ

### Setup

```bash
git clone <repository>
cd SamStore/SS-CatalogService
go mod download
```

### Menjalankan Server API

```bash
go run cmd/api/main.go
```

### Inisialisasi Database (Seeding)

Jika database masih kosong dan butuh data awal:

```bash
go run cmd/seed/main.go
```

### Build & Test

```bash
# Build executable binary
go build -o bin/catalog-service cmd/api/main.go

# Jalankan seluruh unit tests
go test ./...
```

---

## Integrasi & Pola Pesan (Messaging)

Service ini bertindak sebagai **Producer** untuk event katalog.
Menggunakan **Outbox Pattern**, setiap operasi bisnis yang mutasi data (misal: stok produk berkurang) akan mencatat `OutboxEvent` di database yang sama dalam satu transaksi SQL. Worker terpisah akan *poll* tabel outbox dan mempublikasikan event tersebut ke `samstore.events` topic exchange di RabbitMQ.

**Routing Key yang Diterbitkan:**
- `catalog.product.created`
- `catalog.product.updated`
- `catalog.inventory.adjusted`

---

## Known Issues

Tidak ada issue yang teridentifikasi dari source code.

## Future Improvements

- Menambahkan auto-migration run saat startup aplikasi (otomatis memicu `golang-migrate`).
- Menerapkan distributed cache eksternal (seperti Redis) berdampingan dengan Bigcache jika instance service perlu di-scale out besar-besaran secara stateful.
