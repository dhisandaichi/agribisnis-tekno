# Frontend Agritec - SFA Progressive Web App (Vue 3 + Vite)

Frontend aplikasi Sales Force Automation Agritec, dibangun dengan Vue 3, Vite, dan Vuetify, berjalan sebagai PWA.
Untuk PoC, frontend terhubung **langsung ke Supabase** (Auth, Postgres + RLS, Storage). Backend Go di `../backend` tidak dipakai.

## Environment Variables
Salin `.env.example` ke `.env.local` untuk development lokal:

```bash
cp .env.example .env.local
```

| Variable | Keterangan |
| --- | --- |
| `VITE_SUPABASE_URL` | URL project Supabase, mis. `https://xxxx.supabase.co` |
| `VITE_SUPABASE_ANON_KEY` | Publishable key (`sb_publishable_...`). Aman di browser karena akses data dibatasi RLS. |

Aplikasi sengaja gagal start (error di console) jika salah satu variabel kosong.

## Deploy ke Vercel
1. Import repo di Vercel, set **Root Directory** = `frontend`. Framework terdeteksi sebagai Vite dari `vercel.json`.
2. Isi kedua environment variable di atas untuk Production dan Preview.
3. Deploy. `vercel.json` sudah mengatur rewrite SPA, header keamanan, dan cache untuk service worker.
4. Di Supabase > Authentication > URL Configuration, isi **Site URL** dan **Redirect URLs** dengan domain Vercel.

## Fitur
- **Auth & role**: Supabase Auth. Role disimpan di tabel `profiles`. Pendaftaran mandiri hanya untuk `sales` / `farmer`; `supervisor` diset admin.
- **Kunjungan sales**: GPS + foto, dikompres di browser ke JPEG < 500 KB, lalu diunggah ke bucket privat `visit-photos`. Data disimpan lewat RPC `submit_visit` yang atomik dan idempoten.
- **Offline**: app shell di-precache oleh service worker. Kunjungan yang dibuat saat offline disimpan di IndexedDB dan dikirim otomatis saat online.
- **Dashboard supervisor**: ringkasan (`supervisor_summary`), kunjungan terbaru dengan foto (signed URL), dan saldo hutang distributor (`distributor_balances`).
- Halaman petani masih konten statis (di luar cakupan PoC).

## Local Development

```bash
npm install
npm run dev      # development server
npm run build    # build production ke dist/
npm run preview  # jalankan hasil build
```

Ikon PWA dibuat dari `public/logo.svg`:

```bash
npx pwa-assets-generator --preset minimal-2023 public/logo.svg
```
