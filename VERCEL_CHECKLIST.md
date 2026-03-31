# Vercel Deployment Checklist & API Reference

Berdasarkan pengecekan kode pada repository ini, berikut adalah hasil analisis kesiapan deployment untuk platform Vercel, serta daftar API dan API keys yang Anda butuhkan.

## 🏗️ Status Kesiapan Deployment

### 1. Frontend (Vue 3 + Vite)
- **Status:** ✅ **Siap Deploy di Vercel**
- **Keterangan:** Konfigurasi Vercel sudah benar melalui file `frontend/vercel.json`! Setting tersebut mengaktifkan routing SPA (Single Page Application) yang melakukan *rewrite* traffic ke `index.html`. Output `dist` dan build command `npm run build` juga standar dan bisa langsung dijalankan oleh Vercel.

### 2. Backend (Go + Gin)
- **Status:** 🚧 **Perlu Penyesuaian Jika Dideploy di Vercel**
- **Keterangan:** Backend saat ini menggunakan Gin dan `router.Run(":" + port)` untuk berjalan terus-menerus (*long-running server*). Vercel adalah platform *Serverless Functions* sehingga mengekspektasikan file Go mengekspor `http.HandlerFunc` di folder `/api/`.
- **Rekomendasi Deployment Backend:**
  - **Saran A (Termudah):** Deploy backend secara independen di platform lain seperti **Railway**, **Render**, atau VPS. Repository Anda sudah menggunakan `docker-compose.yml` sehingga platform-platform tersebut (termasuk Supabase/PostgreSQL) sangat cocok.
  - **Saran B:** Jika tetap ngotot ingin menggunakan Vercel untuk backend, maka `main.go` harus direfactor agar handler-handler gin dibungkus dalam `http.Handler` di dalam folder baru bernama `api/`.

---

## 🔑 Daftar Environment Variables dan API Key

Ini adalah daftar yang **WAJIB** dimasukkan pada menu **Settings > Environment Variables** saat proses deployment di Vercel.

### Frontend (`frontend/.env`)
Ini adalah API kunci yang dipanggil di sisi client/frontend.

| Variable | Keterangan | Wajib Diisi? |
| :--- | :--- | :--- |
| `VITE_SUPABASE_URL` | URL dari DB & Auth Supabase (contoh: `https://xxxx.supabase.co`). | 🔴 Wajib |
| `VITE_SUPABASE_ANON_KEY` | Public/Anon Key project Supabase Anda. | 🔴 Wajib |
| `VITE_BMKG_API_URL` | Endpoint ramalan cuaca digital BMKG. | 🔴 Wajib |
| `VITE_MAPBOX_API_KEY` | Key untuk provider Peta Mapbox. | 🟡 Opsional (jika peta jalan) |

### Backend (`backend/.env`)
Jika backend ikut dideploy, ia membutuhkan koneksi DB yang diatur di variabel ini.

| Variable | Keterangan | Wajib Diisi? |
| :--- | :--- | :--- |
| `DB_HOST` | Host dari endpoint Postgres Supabase (Contoh: `aws-ap-southeast-1.pooler.supabase.com`) | 🔴 Wajib |
| `DB_PORT` | Port untuk Koneksi. Default dari pooling biasanya `6543`. | 🔴 Wajib |
| `DB_USER` | Username user (`postgres.xxxxxx`). | 🔴 Wajib |
| `DB_PASSWORD` | Password akun Supabase. | 🔴 Wajib |
| `DB_NAME` | Nama tabel DB (Biasanya `postgres`). | 🔴 Wajib |
| `JWT_SECRET` | Token Rahasia untuk men-generate Authorization/Login token user. Silakan buat teks random. | 🔴 Wajib |
| `PORT` | Auto diisi oleh cloud provider. | 🔵 Auto |

---

## 🔌 Daftar API Internal (Backend)
Berikut ini route API yang digunakan oleh aplikasi:
- `POST  /api/v1/auth/register` ▶️ Mendaftarkan Akun Baru
- `POST  /api/v1/auth/login` ▶️ Fitur Login (Akan mengembalikan token JWT untuk autentikasi)
- `POST  /api/v1/sales/visits` ▶️ Mendaftarkan kunjungan Sales baru (Butuh JWT Token dengan peran "sales")

---

## ✅ Checklist Vercel Deployment

1. **Persiapan Database Base**
   - [ ] Import/jalankan query di file `migrations/` & `supabase_queries.sql` pada portal SQL Editor di Supabase.
   - [ ] Catat Supabase URL & Anon Key (dari Settings > API Supabase).
   - [ ] Catat Host, DB, User, dan Password (dari Settings > Database Supabase).

2. **Eksekusi di Vercel Panel**
   - [ ] Sambungkan Vercel dengan Repositori Github Anda.
   - [ ] Di bagian pengaturan Configure Project, ubah **Root Directory** menjadi `frontend`.
   - [ ] Framework Preset biarkan Vite (atau Vue.js 3).
   - [ ] Masukkan ke-4 variabel dari tabel Frontend di atas ke form **Environment Variables** di Vercel.
   - [ ] Klik **Deploy**.

**Catatan:** PWA Plugin sudah disetting. Jika konfigurasi env var valid, website Anda otomatis akan berjalan sempurna sebagai PWA!
