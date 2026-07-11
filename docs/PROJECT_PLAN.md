# Wrapped Weekly — Project Plan

Sumber requirement: `Study Case Fullstack Engineer.pdf` (root project).
Dikelola oleh tim agent: `orchestrator` → `ba` → `backend` + `frontend` → `devops` → `qa`.

## Keputusan yang Sudah Difinalkan

| Keputusan | Pilihan | Alasan |
|---|---|---|
| JWT storage di frontend | httpOnly cookie | Aman dari XSS, best-practice; butuh Next.js route handler sebagai proxy set-cookie |
| Refresh token | **Skip di MVP** | PDF: nilai plus, bukan wajib. Access token expiry panjang (mis. 24 jam) untuk demo |
| Deploy live (bonus +15%) | Dikerjakan **paling akhir**, setelah MVP + Docker lokal solid | Prioritaskan nilai inti dulu |
| Nice-to-have | **Tidak ada** kecuali unit test agregasi (sudah wajib secara implisit) | PDF: "jangan kejar semua nice-to-have" |
| Repo | Monorepo: `backend/`, `frontend/`, `nginx/`, root `docker-compose.yml` | Sesuai anjuran PDF §4 |

## Urutan Modul (tiap modul = 1 siklus ba → backend/frontend → devops → qa)

### Modul 0 — Project Scaffolding
- [ ] Init `backend/` (Go module, struktur `cmd/api`, `domain/`, `usecase/`,
      `repository/`, `handler/`), `go.mod`
- [ ] Init `frontend/` (Next.js App Router + TS + Tailwind + shadcn/ui init)
- [ ] `.env.example` root + `backend/.env.example`
- [ ] `docker-compose.yml` skeleton (postgres saja dulu) + `nginx/app.conf` skeleton
- [ ] Endpoint `/health` di backend
- Agent: backend, frontend, devops (BA tidak perlu untuk scaffolding murni)

### Modul 1 — Auth (Register & Login)
- [ ] BA: user story + acceptance criteria + edge case (email dipakai ganda, password lemah, dll)
- [ ] Backend: `users` table (migration), bcrypt hash, `POST /api/v1/auth/register`,
      `POST /api/v1/auth/login`, JWT issue, middleware Bearer/cookie auth
- [ ] Frontend: halaman register & login (RHF+Zod), Next.js route handler untuk
      set httpOnly cookie, auth guard untuk route privat
- [ ] QA: test register/login termasuk edge case dari BA

### Modul 2 — Activity Logging (CRUD)
- [ ] BA: user story + business rules (kategori apa saja, batas value, dll) + edge case
- [ ] Backend: `activities` table (migration, FK ke user), CRUD endpoint dengan
      ownership check ketat (404 bukan 403 untuk resource milik user lain — putuskan
      dan dokumentasikan konsisten)
- [ ] Frontend: form catat aktivitas + list + edit + hapus, loading/empty/error state
- [ ] QA: test CRUD + ownership (user A tidak bisa akses data user B)

### Modul 3 — Weekly Recap (Domain Logic Inti — paling kritis)
- [ ] BA: definisi rumus agregasi eksplisit (total per kategori, kategori terbanyak,
      hari paling produktif, formula "perubahan vs minggu lalu"), edge case wajib:
      timezone, batas Senin–Minggu, minggu kosong, div-by-zero
- [ ] Backend: usecase agregasi + `AIProvider` interface (mock deterministik default),
      `recaps` table dengan `slug` unik, `POST /api/v1/recaps/generate`,
      `GET /api/v1/recaps` (list), `GET /api/v1/recaps/{id}`
- [ ] Frontend: trigger generate recap, tampilkan statistik + narasi
- [ ] QA: **unit test table-driven wajib** untuk semua edge case agregasi (lihat qa.md)

### Modul 4 — Dashboard
- [ ] BA: apa saja yang wajib tampil di ringkasan minggu berjalan
- [ ] Backend: endpoint ringkasan + data chart (aktivitas per hari / komposisi kategori)
- [ ] Frontend: dashboard page + minimal 1 chart (Recharts), loading/empty/error state
- [ ] QA: test endpoint dashboard

### Modul 5 — Shareable Public Recap
- [ ] BA: aturan akses publik (apa yang boleh terlihat tanpa login)
- [ ] Backend: `GET /api/v1/recaps/public/{slug}` tanpa auth
- [ ] Frontend: halaman `/w/{slug}` — **Server Component + SSR wajib**,
      `generateMetadata()` dengan Open Graph title/description/image
- [ ] QA: test endpoint publik, verifikasi tidak bocor data privat, cek meta tag OG
      muncul di response HTML (view-source, bukan cuma di browser setelah hydrate)

### Modul 6 — Dockerization + Nginx (deliverable inti, bobot 20%)
- [ ] Dockerfile `backend/` (multi-stage, non-root, HEALTHCHECK)
- [ ] Dockerfile `frontend/` (multi-stage, `output: standalone`, non-root)
- [ ] `docker-compose.yml` lengkap: postgres + migrate (one-shot) + backend + frontend + nginx
- [ ] `nginx/app.conf` reverse proxy final (`/` → frontend, `/api/v1/*` → backend)
- [ ] Verifikasi nyata: `docker compose up --build -d` sukses, semua service healthy,
      alur end-to-end (register → catat aktivitas → generate recap → buka link
      publik → dashboard) jalan lewat `http://localhost`
- Agent: devops (setelah modul 1–5 minimal MVP-nya ada)

### Modul 7 — Dokumentasi & Kerapian (bobot 15%)
- [ ] `README.md`: cara run lokal, arsitektur singkat, rumus agregasi recap,
      fitur selesai vs belum, keputusan teknis & trade-off (termasuk tabel di atas),
      estimasi waktu yang dihabiskan
- [ ] Review commit history — pastikan bercerita, bukan satu commit raksasa
- [ ] Pastikan tidak ada secret ter-commit (cek `.gitignore` mencakup `.env`)

### Modul 8 — Bonus: Deploy Live (opsional, paling akhir)
- [ ] Pilih target: VPS (DigitalOcean/Lightsail) dengan Docker penuh, ATAU
      Frontend Vercel + Backend/DB Railway/Render/Fly.io
- [ ] Domain/subdomain → server, TLS via Certbot jika VPS
- [ ] Verifikasi end-to-end di URL live, tulis URL + langkah di README

## Yang Sengaja TIDAK Dikerjakan (out of scope MVP ini)

- Refresh token flow
- Cron job auto-generate recap
- Integrasi API eksternal (GitHub/Strava/RSS)
- Streaming narasi AI ke UI
- OG image dinamis
- Redis cache/rate-limit
- Pagination & filter riwayat aktivitas
- Dark mode eksplisit (di luar loading/empty/error state yang tetap wajib)

Alasan: PDF eksplisit meminta scope kecil + eksekusi rapi lebih baik dari fitur
banyak setengah jadi. Bisa direvisit setelah Modul 0–7 selesai dan waktu masih ada.

## Cara Menjalankan Tim Agent per Modul

Panggil `orchestrator` dengan instruksi: "kerjakan Modul N — [nama]". Orchestrator
akan scan state project, delegasi ke ba → backend/frontend → devops (jika relevan)
→ qa sesuai `.claude/agents/orchestrator.md`, dan lapor balik sebelum lanjut modul
berikutnya.
