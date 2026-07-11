# Wrapped Weekly

Aplikasi yang mengubah aktivitas mingguanmu jadi recap ala "Spotify Wrapped": catat
aktivitas (workout/reading/coding/spending), lalu generate recap otomatis berisi
statistik + narasi + kartu yang bisa dibagikan lewat link publik.

Submission Study Case Fullstack Engineer.

## Stack

- **Backend**: Go + Gin, layered (`domain` → `usecase` → `repository` → `handler`), PostgreSQL + `golang-migrate`
- **Auth**: JWT via httpOnly cookie
- **Frontend**: Next.js 16 (App Router) + TypeScript (`output: standalone`), shadcn/ui + Tailwind v4
- **Data fetching**: TanStack Query; **Form**: React Hook Form + Zod; **Chart**: Recharts
- **Reverse proxy**: Nginx (`/` → frontend, `/api/v1/*` → backend)
- **AI narasi**: provider `mock` deterministik by default (tanpa API key)

## Arsitektur

```
Internet ──▶ Nginx (:80)
              / → frontend (Next.js :3000)
              /api/v1/* → backend (Go/Gin :8080) → postgres (:5432)
```

Monorepo: `backend/`, `frontend/`, `nginx/`, `docker-compose.yml` di root.

## Cara Run (Docker)

```bash
cp .env.example .env
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env

docker compose up --build -d
curl http://localhost/health
```

Buka `http://localhost`. Verifikasi end-to-end: register → login → catat aktivitas →
Dashboard → generate Weekly Recap → salin link publik `/w/{slug}` → buka di tab
incognito (harus bisa diakses tanpa login; `view-source:` untuk cek meta OG ada
di HTML awal).

## Cara Run (tanpa Docker, development)

**Backend** (Go 1.23+, PostgreSQL lokal):
```bash
cd backend
cp .env.example .env   # sesuaikan DATABASE_URL
migrate -path db/migrations -database "$DATABASE_URL" up
go run ./cmd/api
```

**Frontend** (Node 20+):
```bash
cd frontend
cp .env.example .env.local
npm install && npm run dev
```

`next dev` men-proxy `/api/v1/*` ke backend `:8080` lewat `rewrites()` di
`next.config.ts` — tidak perlu Nginx untuk development.

## Rumus Agregasi Recap

Implementasi + test: [`recap_aggregation.go`](backend/internal/usecase/recap_aggregation.go) /
[`recap_aggregation_test.go`](backend/internal/usecase/recap_aggregation_test.go).

- **Batas minggu**: ISO week, Senin 00:00–Senin berikutnya 00:00 (half-open), dievaluasi
  di **UTC** tetap (`usecase.AppTimezone`) — bukan timezone per-user, supaya agregasi
  reproducible. Trade-off: user jauh dari UTC bisa lihat aktivitas "geser satu hari".
- **Top category / hari produktif**: total tertinggi; tie-break deterministik
  (urutan kategori tetap; hari paling awal menang).
- **Perubahan vs minggu lalu**: `null` (bukan `+Inf`/`+100%`) kalau minggu lalu totalnya 0.
- **Minggu kosong**: statistik tetap dihitung tanpa crash (`total: 0`, `top_category: null`,
  `daily_breakdown` tetap 7 entri).

## Fitur Selesai vs Belum

### Selesai
- [x] Auth (register/login, bcrypt, JWT httpOnly cookie)
- [x] Activity CRUD dengan otorisasi per-user (404, bukan 403, untuk resource user lain)
- [x] Weekly Recap: agregasi + AI mock + slug publik
- [x] Dashboard: ringkasan minggu + chart (Recharts)
- [x] Public Recap `/w/{slug}`: SSR + `generateMetadata` Open Graph
- [x] Loading/empty/error state di semua halaman fetch
- [x] Unit test table-driven untuk semua edge case agregasi
- [x] Docker multi-stage BE+FE, docker-compose (postgres+migrate+backend+frontend+nginx), Nginx
- [x] **(Bonus)** Filter kategori + rentang tanggal di riwayat aktivitas
- [x] **(Bonus)** Dark mode (`next-themes`)

### Sengaja tidak dikerjakan
Refresh token, cron auto-recap, integrasi API eksternal, streaming narasi AI, OG image
dinamis, Redis, deploy live (lihat bagian Deploy di bawah). Detail alasan: `docs/PROJECT_PLAN.md`.

## Keputusan Teknis Penting

- **JWT di httpOnly cookie** (bukan localStorage) — mitigasi XSS.
- **AI provider di belakang interface** `domain.AIProvider` — mock default, provider asli
  tinggal didaftarkan di `aiprovider.NewProvider()`.
- **Response envelope seragam** `{data, message}` di semua endpoint.
- **`ActivityFilter` sebagai struct terpisah** — jaga backward compatibility signature repository.
- Detail lain (kenapa `rewrites()` bukan route handler, kenapa value form pakai string, dll)
  ada sebagai komentar inline di kode terkait.

## Mengaktifkan AI Provider Asli

Default `AI_PROVIDER=mock` (deterministik, tanpa API key). Untuk LLM asli: implementasikan
`domain.AIProvider` di file baru (mis. `aiprovider/openai.go`), daftarkan di
`aiprovider.NewProvider()`, lalu set `AI_PROVIDER=<nama>` + API key di `backend/.env`.

## Deploy Live (Bonus — tidak dieksekusi)

PDF study case: deploy live adalah nilai tambah (+15%), bukan syarat lulus. Diputuskan
untuk skip karena platform gratis (Render, Fly.io) sekarang mewajibkan verifikasi kartu.
Kode sudah siap untuk deploy cross-domain (Vercel + Render/Railway/Fly.io) lewat env var
`COOKIE_CROSS_SITE=true` (lihat `backend/internal/config/config.go`) yang mengubah cookie
jadi `SameSite=None; Secure` — wajib saat frontend dan backend beda domain.

## Known Limitations

- **`og-default.png` adalah placeholder minimal** — ganti dengan gambar 1200x630 asli
  sebelum dibagikan publik sungguhan.
- Semua klaim di README ini sudah diverifikasi nyata: `go build`/`go test` (11/11 pass,
  Go 1.26.5), `next build`/`eslint` bersih, `docker compose up --build` jalan penuh dan
  seluruh service sehat.

## Estimasi Waktu

±12-16 jam efektif: paling banyak di domain logic agregasi, lalu CRUD/auth, Docker/Nginx,
dokumentasi.

## Struktur Folder

```
.
├── backend/
│   ├── cmd/api/            # entrypoint
│   ├── internal/{domain,usecase,repository,handler,middleware,aiprovider,config}/
│   ├── db/migrations/
│   └── Dockerfile
├── frontend/
│   ├── app/                # App Router: (auth), (app), w/[slug]
│   ├── features/           # kode per-fitur
│   ├── lib/api/             # client.ts (browser) + server.ts (SSR)
│   └── Dockerfile
├── nginx/app.conf
├── docker-compose.yml
└── docs/PROJECT_PLAN.md
```
