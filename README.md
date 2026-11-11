# SFA Agritec — PT. Agribisnis Technology Indonesia

Aplikasi Sales Force Automation (PWA) untuk sales lapangan, supervisor, distributor/kios, dan petani.

## Arsitektur production (satu project Vercel)

```
Browser (PWA Vue) ──► /api/v1/*  ──► Go (Gin) Vercel Function ──► PostgreSQL (Neon / Supabase)
                └──► /*          ──► frontend/dist (static)
```

| Bagian | Teknologi | Lokasi |
| --- | --- | --- |
| Frontend (PWA) | Vue 3 + Vite + Vuetify | `frontend/` |
| Backend API | Go + Gin + GORM, JWT + bcrypt | `backend/`, entry Vercel `api/index.go` |
| Skema database | PostgreSQL | `backend/db/schema.sql` |

## Fitur per role

- **Sales**: input kunjungan petani (data petani, komoditas, luas lahan, GPS wajib, foto dikompres di HP), simpan offline + kirim otomatis saat online, riwayat kunjungan.
- **Supervisor**: ringkasan (petani, kunjungan, distributor, piutang), tabel kunjungan + foto + peta + ekspor CSV, kelola distributor & order (tandai lunas), kelola pengguna (role, nonaktifkan, password sementara), harga komoditas & konten edukasi.
- **Distributor/Kios**: tagihan & jatuh tempo, riwayat order, ajukan order baru.
- **Petani**: konten edukasi, prakiraan cuaca BMKG, harga komoditas.

Role *supervisor* tidak bisa dipilih saat daftar; hanya supervisor lain (atau SQL) yang bisa memberikannya.

## Environment variables (Vercel → Settings → Environment Variables)

| Variable | Keterangan | Wajib |
| --- | --- | --- |
| `DATABASE_URL` | Connection string Postgres (pakai URL *pooled* untuk Neon/Supabase, `sslmode=require`) | Ya |
| `JWT_SECRET` | String acak ≥ 32 karakter (`openssl rand -base64 48`) | Ya |
| `BMKG_ADM4` | Kode kelurahan untuk cuaca, contoh `31.71.03.1001` | Opsional |

## Deploy

1. Buat database Postgres lalu jalankan `backend/db/schema.sql` (SQL editor, atau `DATABASE_URL=... go run ./backend/cmd/server -migrate`).
2. Import repo ke Vercel dengan **Root Directory = root repo** (bukan `frontend`). `vercel.json` sudah mengatur build frontend, fungsi Go, dan routing.
3. Isi environment variables di atas, lalu deploy.
4. Daftar akun pertama lewat aplikasi, lalu jadikan supervisor:
   ```sql
   update users set role = 'supervisor' where email = 'email-anda@contoh.com';
   ```

## Pengembangan lokal

```bash
docker compose up -d db                       # Postgres lokal (skema otomatis dibuat)
cp backend/.env.example backend/.env          # isi DATABASE_URL / DB_* dan JWT_SECRET
go run ./backend/cmd/server                   # API di http://localhost:8080/api/v1
cd frontend && npm install && npm run dev     # http://localhost:5173 (proxy /api ke :8080)
```

## Daftar API

Publik: `POST /auth/register`, `POST /auth/login`, `GET /health`, `GET /weather`
Login: `GET /me`, `PUT /me/password`, `GET /content/prices`, `GET /content/education`, `GET /visits/:id/photo`
Sales: `GET|POST /sales/visits`
Distributor: `GET /distributor`, `POST /distributor/orders`
Supervisor: `/supervisor/summary`, `/supervisor/visits`, `/supervisor/distributors`, `/supervisor/orders` (+`/:id/status`), `/supervisor/users` (+`/:id`, `/:id/reset-password`), `/supervisor/prices`, `/supervisor/education`

## Deploy ke VPS (Docker + Caddy, HTTPS otomatis)

Prasyarat: Docker + Docker Compose di VPS, port 80/443 terbuka, dan domain (A record) mengarah ke IP VPS.

```bash
git clone https://github.com/dhisandaichi/agribisnis-tekno && cd agribisnis-tekno
cp .env.production.example .env.production      # isi DOMAIN, DATABASE_URL (Neon pooled), JWT_SECRET
docker compose -f docker-compose.prod.yml --env-file .env.production build
docker compose -f docker-compose.prod.yml --env-file .env.production run --rm app -migrate
docker compose -f docker-compose.prod.yml --env-file .env.production up -d
curl https://$DOMAIN/api/v1/health              # {"status":"ok"}
```

Update versi baru: `git pull && docker compose -f docker-compose.prod.yml --env-file .env.production up -d --build`.
