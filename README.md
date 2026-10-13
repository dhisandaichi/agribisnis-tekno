# SFA Agritec - PoC

Sales Force Automation untuk PT. Agribisnis Technology Indonesia.

## Arsitektur PoC
```
Browser (PWA Vue 3)  ──►  Supabase (Auth · Postgres + PostGIS + RLS · Storage)
      ▲
   Vercel (static hosting, root dir: frontend/)
```

| Folder | Status | Isi |
| --- | --- | --- |
| `frontend/` | **Aktif** | PWA Vue 3 + Vuetify. Lihat `frontend/README.md`. |
| `supabase/migrations/` | **Aktif** | Skema, RLS, RPC, dan bucket Storage. Sumber kebenaran database. |
| `supabase/seed.sql` | **Aktif** | Data contoh distributor/order (tanpa user). |
| `backend/` | Legacy | API Go/Gin. Tidak dipakai di PoC; tidak bisa dijalankan di Vercel apa adanya. |
| `migrations/`, `supabase_queries.sql` | Legacy | Skema lama untuk backend Go (tabel `users` + password hash). **Jangan dijalankan** di Supabase. |
| `docker-compose.yml`, `analytics/` | Legacy | Stack lokal VPS (Postgres, Redis, Grafana). |

## Setup Supabase
1. Terapkan migrasi `supabase/migrations/*.sql` (Supabase CLI `supabase db push`, atau SQL Editor), lalu `supabase/seed.sql` bila perlu data contoh.
2. Authentication > URL Configuration: isi Site URL & Redirect URLs dengan domain Vercel.
3. Authentication > Sign In / Providers > Email: untuk PoC, matikan **Confirm email** atau siapkan SMTP sendiri. SMTP bawaan Supabase dibatasi sekitar 2 email/jam.
4. Buat akun supervisor: daftar biasa (atau Authentication > Users > Add user), lalu jalankan:
   ```sql
   update public.profiles set role = 'supervisor' where email = 'supervisor@contoh.com';
   ```

## Model akses (RLS)
| Role | profiles | farmers / visits | distributors / orders | foto kunjungan |
| --- | --- | --- | --- | --- |
| sales | baca milik sendiri | tambah & baca milik sendiri | - | unggah & baca folder sendiri |
| supervisor | baca semua | baca semua | kelola semua | baca semua |
| farmer | baca milik sendiri | - | - | - |
| anon | - | - | - | - |

Tidak ada user yang bisa mengubah role-nya sendiri; perubahan role hanya lewat SQL/admin.
