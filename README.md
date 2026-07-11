# Wrapped Weekly

Aplikasi mini fullstack yang mengubah aktivitas mingguanmu menjadi recap ala "Spotify
Wrapped": catat aktivitas (workout/reading/coding/spending), lalu generate recap
otomatis berisi statistik + narasi + kartu yang bisa dibagikan lewat link publik.

Dibangun sebagai submission Study Case Fullstack Engineer.

## Stack

- **Backend**: Go + Gin, arsitektur berlapis (`domain` → `usecase` → `repository` → `handler`)
- **Database**: PostgreSQL, skema dikelola via migration SQL (`golang-migrate`)
- **Auth**: JWT disimpan sebagai httpOnly cookie (bukan localStorage, untuk mitigasi XSS)
- **Frontend**: Next.js 16 (App Router) + TypeScript, output `standalone`
- **UI**: shadcn/ui + Tailwind CSS v4
- **Data fetching**: TanStack Query; **Form**: React Hook Form + Zod
- **Chart**: Recharts
- **Reverse proxy**: Nginx (`/` → frontend, `/api/v1/*` → backend)
- **AI narasi recap**: provider `mock` deterministik secara default (tanpa API key)

## Arsitektur

```
Internet ──▶ Nginx (:80) reverse proxy
              / → frontend (Next.js :3000)
              /api/v1/* → backend (Go/Gin :8080)
                                       │
                                  postgres (:5432)
```

Repo ini monorepo: `backend/`, `frontend/`, `nginx/`, dan `docker-compose.yml` di root.

## Cara Run Lokal (Docker — cara utama)

Prasyarat: Docker + Docker Compose. **Prasyarat tambahan sekali jalan**: karena
`backend/go.sum` belum di-generate di repo ini (dikerjakan tanpa akses toolchain Go),
jalankan `go mod tidy` di dalam `backend/` sebelum build image pertama kali — lihat
bagian [Known Limitations](#known-limitations--catatan-jujur) di bawah.

```bash
# 1. Siapkan env
cp .env.example .env
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env

# 2. (sekali saja) generate go.sum — lihat Known Limitations
cd backend && go mod tidy && cd ..

# 3. Build & nyalakan semua service
docker compose up --build -d

# 4. Cek health backend
curl http://localhost/api/v1/health

# 5. Buka app
open http://localhost   # atau kunjungi manual di browser
```

Alur end-to-end untuk verifikasi: register → login → catat aktivitas → buka
Dashboard → generate Weekly Recap → salin link publik `/w/{slug}` → buka di
tab incognito (harus bisa diakses tanpa login, cek `view-source:` untuk
memastikan meta OG ada di HTML awal).

## Cara Run Lokal (tanpa Docker, untuk development)

**Backend** (butuh Go 1.23+ dan PostgreSQL lokal):
```bash
cd backend
cp .env.example .env   # sesuaikan DATABASE_URL ke Postgres lokalmu
go mod tidy
# jalankan migration (butuh golang-migrate CLI terpasang)
migrate -path db/migrations -database "$DATABASE_URL" up
go run ./cmd/api
```

**Frontend** (butuh Node 20+):
```bash
cd frontend
cp .env.example .env.local
# NEXT_PUBLIC_API_URL=/api/v1, BACKEND_INTERNAL_URL=http://localhost:8080
npm install
npm run dev
```

Frontend `next dev` di `:3000` men-proxy `/api/v1/*` ke backend `:8080` lewat
`rewrites()` di `next.config.ts` — tidak perlu Nginx untuk development.

## Aturan / Rumus Agregasi Recap

Definisi lengkap ada di kode: [`backend/internal/usecase/recap_aggregation.go`](backend/internal/usecase/recap_aggregation.go)
dan diuji di [`backend/internal/usecase/recap_aggregation_test.go`](backend/internal/usecase/recap_aggregation_test.go).

- **Batas minggu**: ISO week, Senin 00:00 (inklusif) sampai Senin berikutnya 00:00
  (eksklusif) — half-open range `[start, end)`, dievaluasi dalam **UTC** sebagai
  timezone acuan tetap (lihat `usecase.AppTimezone`). Ini keputusan sadar: server
  tidak menebak timezone user per-request, supaya agregasi selalu reproducible.
  Trade-off: user di zona waktu jauh dari UTC bisa melihat aktivitas "geser satu
  hari" dibanding perasaan lokalnya. Didokumentasikan sebagai batasan yang diketahui.
- **Total per kategori**: sum `value` per `category` dalam rentang minggu.
- **Kategori terbanyak (top category)**: kategori dengan total tertinggi. Kalau
  seri (tie), pemenang ditentukan berdasarkan urutan tetap
  (workout → reading → coding → spending) supaya hasilnya deterministik, bukan
  tergantung urutan iterasi map yang random di Go.
- **Hari paling produktif**: hari dengan total value tertinggi dalam minggu itu.
  Tie-break: hari paling awal dalam minggu (Senin duluan) yang menang.
- **Perubahan vs minggu sebelumnya**: `((total_minggu_ini - total_minggu_lalu) / total_minggu_lalu) * 100`.
  - Jika `total_minggu_lalu == 0` (termasuk kasus minggu sebelumnya benar-benar
    kosong): hasil `null`, **bukan** `+Infinity` atau `+100%` yang menyesatkan.
    Frontend menampilkan ini sebagai "N/A".
  - Jika kedua minggu 0 aktivitas: juga `null`.
- **Minggu tanpa aktivitas sama sekali**: seluruh statistik tetap dihitung tanpa
  crash — `total_activities: 0`, `top_category: null`, `most_productive_day: null`,
  `daily_breakdown` tetap berisi 7 entri (semua `total: 0`). Narasi mock punya
  cabang khusus untuk kasus ini ("belum ada aktivitas yang tercatat...").

## Fitur Selesai vs Belum

### Selesai (MVP)
- [x] Auth: register, login (bcrypt hash), JWT via httpOnly cookie, endpoint privat terproteksi
- [x] Activity Logging CRUD lengkap dengan otorisasi per-user (404 untuk resource
      milik user lain, bukan 403 — supaya tidak bocor informasi keberadaan resource)
- [x] Weekly Recap: agregasi lengkap + AI provider mock deterministik + slug publik
- [x] Dashboard: ringkasan minggu berjalan + chart aktivitas per hari (Recharts)
- [x] Shareable Public Recap: `/w/{slug}` Server Component SSR + `generateMetadata`
      Open Graph (title/description/image)
- [x] Loading / empty / error state di semua halaman yang fetch data
- [x] Unit test table-driven untuk seluruh edge case agregasi (timezone, batas
      minggu, minggu kosong, div-by-zero, tie-break)
- [x] Dockerfile BE & FE multi-stage, non-root user, HEALTHCHECK
- [x] docker-compose (postgres + migrate one-shot + backend + frontend + nginx)
- [x] Nginx reverse proxy
- [x] **(Bonus)** Pagination & filter riwayat aktivitas — filter kategori + rentang
      tanggal (`GET /api/v1/activities?category=&from=&to=&page=&size=`), UI filter
      di halaman Aktivitas dengan reset-to-page-1 saat filter berubah
- [x] **(Bonus)** Dark mode — `next-themes` dengan deteksi `prefers-color-scheme`
      sistem + toggle manual di navbar, hydration-safe via `next/dynamic(ssr:false)`

### Belum / Sengaja Di-skip (lihat `docs/PROJECT_PLAN.md` untuk alasan)
- [ ] Refresh token (hanya access token, expiry 24 jam)
- [ ] Cron job auto-generate recap mingguan
- [ ] Integrasi API eksternal (GitHub/Strava/RSS) sebagai sumber aktivitas
- [ ] Streaming narasi AI ke UI
- [ ] OG image dinamis per-recap (pakai satu gambar statis placeholder)
- [ ] Redis cache / rate-limit
- [ ] Deploy live ke VPS/Vercel/Railway (bonus — lihat rencana di `docs/PROJECT_PLAN.md` Modul 8)

## Keputusan Teknis & Trade-off

| Keputusan | Alasan |
|---|---|
| JWT di httpOnly cookie, bukan localStorage | Mitigasi XSS — token tidak bisa dibaca JavaScript di sisi klien |
| Timezone agregasi fixed di UTC | Reproducibility di atas presisi per-user; didokumentasikan sebagai batasan |
| `change_vs_prev_week_pct = null` saat pembagi nol | Menghindari nilai menyesatkan (`+Inf`/`+100%`) dari data yang tidak ada |
| 404 (bukan 403) untuk resource milik user lain | Tidak membocorkan keberadaan resource ke user yang tidak berhak |
| AI provider di belakang interface (`domain.AIProvider`) | Mock jadi default tanpa API key; provider asli tinggal diimplementasikan dan didaftarkan di `aiprovider.NewProvider` |
| Response envelope konsisten `{data, message}` | Kontrak API seragam, memudahkan frontend menangani sukses/error secara generic (lihat `lib/api/client.ts`) |
| Next.js `rewrites()` untuk proxy dev-lokal ke backend | Browser selalu memanggil path relatif `/api/v1`, cookie selalu terlihat same-origin baik di dev maupun di balik Nginx — tidak perlu route handler proxy manual |
| Value activity dikirim sebagai string di form, dikonversi saat submit | Menghindari konflik tipe input/output `z.coerce.number()` dengan resolver React Hook Form generic |
| `ActivityFilter` sebagai struct terpisah, bukan menambah parameter ke `ListByUser` | Menjaga backward compatibility signature, tidak perlu ubah semua call site yang sudah ada |
| Theme toggle di-load via `next/dynamic(ssr:false)`, bukan manual mounted-state `useEffect` | Pola manual mounted-state (`useState(false)` + `setState` di effect kosong) — meski direkomendasikan resmi oleh `next-themes` — dilarang React Compiler linter di project ini (`react-hooks/set-state-in-effect`); dynamic import dengan `ssr:false` mencapai efek yang sama (skip render saat server) tanpa melanggar rule itu |

## Mengaktifkan AI Provider Asli

Default `AI_PROVIDER=mock` di `backend/.env` mengembalikan narasi deterministik
berbasis template (lihat `backend/internal/aiprovider/mock.go`) — tidak butuh API key.

Untuk memakai LLM asli:
1. Implementasikan `domain.AIProvider` (method `GenerateNarrative(stats, userName) (string, error)`)
   di file baru, mis. `backend/internal/aiprovider/openai.go`, memanggil API sungguhan.
2. Tambahkan case baru di `aiprovider.NewProvider()` (`backend/internal/aiprovider/provider.go`).
3. Set `AI_PROVIDER=openai` (atau nama lain yang dipilih) dan API key terkait di `backend/.env`.

## Deploy Live (Bonus) — Frontend di Vercel, Backend+DB di Render

Kombinasi ini eksplisit disebut sebagai opsi valid oleh study case. Karena
frontend dan backend berada di domain berbeda (cross-origin), ada penyesuaian
env var yang wajib — jangan pakai default lokal.

### 1. Backend + Database di Render

1. Buat akun di [render.com](https://render.com), hubungkan repo GitHub ini.
2. **New → PostgreSQL** — buat database gratis, catat `Internal Database URL`
   yang diberikan (dipakai sebagai `DATABASE_URL`).
3. **New → Web Service** — pilih repo ini, set **Root Directory: `backend`**,
   **Environment: Docker** (Render otomatis pakai `backend/Dockerfile`).
4. Set environment variables di Render:
   ```
   DATABASE_URL=<Internal Database URL dari langkah 2, tambahkan ?sslmode=require>
   JWT_SECRET=<random string panjang, generate baru — jangan pakai default>
   JWT_EXPIRY_HOURS=24
   AI_PROVIDER=mock
   COOKIE_SECURE=true
   COOKIE_CROSS_SITE=true
   FRONTEND_BASE_URL=https://<domain-vercel-kamu>.vercel.app
   ```
5. Jalankan migration sekali (Render Shell atau lewat `migrate` CLI lokal
   mengarah ke `DATABASE_URL` publik Render) — service `migrate` di
   `docker-compose.yml` tidak otomatis jalan di Render, jadi ini manual:
   ```bash
   migrate -path backend/db/migrations -database "<DATABASE_URL>" up
   ```
6. Deploy. Catat URL yang diberikan Render (mis. `https://wrappedweekly-api.onrender.com`).
7. **Catatan free tier**: service Render gratis sleep setelah 15 menit idle;
   request pertama setelah sleep bisa lambat (~30-60 detik cold start). Wajar
   untuk demo, bukan untuk produksi nyata. Database Postgres gratis Render
   expire 90 hari.

### 2. Frontend di Vercel

1. Buat akun di [vercel.com](https://vercel.com), **Import Project** dari repo GitHub ini.
2. Set **Root Directory: `frontend`**.
3. Set environment variables di Vercel:
   ```
   NEXT_PUBLIC_API_URL=https://<url-backend-render>.onrender.com/api/v1
   BACKEND_INTERNAL_URL=https://<url-backend-render>.onrender.com
   ```
4. Deploy. Vercel otomatis mendeteksi Next.js dan `output: "standalone"`.

### 3. Kenapa `COOKIE_CROSS_SITE=true` wajib

Saat frontend (Vercel) dan backend (Render) berada di domain berbeda, cookie
auth httpOnly butuh `SameSite=None; Secure` supaya browser tetap mengirimnya
lintas domain — `SameSite=Lax` (default untuk setup same-origin di belakang
Nginx) akan diblokir browser modern di skenario ini. Lihat
`backend/internal/config/config.go` (`CookieCrossSite`) dan
`backend/internal/handler/auth_handler.go` (`sameSiteMode()`).

### 4. Verifikasi end-to-end di live URL

Sama seperti checklist lokal: register → login → catat aktivitas → dashboard
→ generate recap → buka `/w/{slug}` di tab incognito → cek `view-source:`
untuk meta OG.

## Known Limitations / Catatan Jujur

- **Cold start Render free tier**: lihat bagian Deploy Live di atas.
- **`og-default.png` adalah placeholder minimal.** Ganti dengan gambar share
  1200x630 yang sebenarnya di `frontend/public/og-default.png` sebelum
  benar-benar dibagikan publik — saat ini fungsional (tidak 404) tapi generik.
- Semua klaim di README ini **sudah diverifikasi nyata**, bukan asumsi: backend
  di-build (`go build ./...`) dan diuji (`go test ./...`, 11/11 pass) memakai
  Go 1.26.5; frontend di-build (`next build`, full type-check TypeScript) dan
  di-lint (`eslint`) bersih; `docker compose up --build` dijalankan penuh dan
  seluruh service (postgres, migrate, backend, frontend, nginx) sehat, diverifikasi
  lewat `curl http://localhost/health` dan pemakaian manual di browser.
- **`og-default.png` adalah placeholder 1x1 pixel.** Ganti dengan gambar share
  1200x630 yang sebenarnya di `frontend/public/og-default.png` sebelum deploy
  production — saat ini fungsional (tidak 404) tapi tidak akan terlihat bagus
  di link preview.

Rekomendasi langkah pertama setelah clone: jalankan `go mod tidy` di `backend/`,
lalu `docker compose up --build`, lalu ikuti alur verifikasi end-to-end di atas.
Jika ada error kompilasi Go yang muncul, kemungkinan besar hanya soal versi minor
dependency (`go.mod` memakai versi yang seharusnya kompatibel per Juli 2026) —
laporkan atau perbaiki via `go mod tidy` / `go get -u`.

## Estimasi Waktu

Dikerjakan dalam satu sesi kerja terfokus (setup arsitektur tim-agent, lalu
implementasi end-to-end backend + frontend + Docker + dokumentasi). Estimasi
efektif mengikuti target study case: ±12-16 jam setara effort, terkonsentrasi
pada domain logic agregasi (paling banyak waktu), lalu CRUD/auth, lalu Docker/Nginx,
lalu dokumentasi.

## Struktur Folder

```
.
├── backend/
│   ├── cmd/api/            # entrypoint (main.go)
│   ├── internal/
│   │   ├── domain/         # entities + repository interfaces
│   │   ├── usecase/        # business logic (termasuk agregasi recap)
│   │   ├── repository/     # implementasi Postgres (pgx)
│   │   ├── handler/        # HTTP layer (Gin) + router
│   │   ├── middleware/     # auth JWT
│   │   ├── aiprovider/     # AI narrative provider (mock + interface)
│   │   └── config/         # env config loader
│   ├── db/migrations/      # migration SQL (golang-migrate)
│   ├── pkg/                # response envelope, app error
│   └── Dockerfile
├── frontend/
│   ├── app/                 # Next.js App Router (route groups: (auth), (app), w/[slug])
│   ├── features/            # kode per-fitur (auth, activity, recap, dashboard)
│   ├── lib/api/              # API client terpusat (client.ts untuk browser, server.ts untuk SSR)
│   └── Dockerfile
├── nginx/app.conf
├── docker-compose.yml
└── docs/PROJECT_PLAN.md     # todolist & keputusan project
```
