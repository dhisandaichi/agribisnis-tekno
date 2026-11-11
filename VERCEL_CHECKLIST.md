# Vercel Deployment Checklist

Frontend Vue dan backend Go berjalan dalam **satu project Vercel** (Go sebagai Vercel Function di `api/index.go`).

1. **Database**
   - [ ] Buat Postgres (Neon via Vercel Marketplace, atau Supabase) dan salin connection string *pooled*.
   - [ ] Jalankan `backend/db/schema.sql`.
2. **Vercel**
   - [ ] Import repo GitHub. Root Directory = root repo. Framework Preset = Other (sudah diatur `vercel.json`).
   - [ ] Environment Variables: `DATABASE_URL`, `JWT_SECRET` (≥ 32 karakter), opsional `BMKG_ADM4`.
   - [ ] Deploy, lalu cek `https://<domain>/api/v1/health` → `{"status":"ok"}`.
3. **Setelah deploy**
   - [ ] Daftar akun pertama, jadikan supervisor lewat SQL (lihat README).
   - [ ] Uji login tiap role, input kunjungan + foto, dashboard supervisor, portal distributor.
