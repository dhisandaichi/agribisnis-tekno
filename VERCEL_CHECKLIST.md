# Vercel + Supabase Deployment Checklist (PoC)

Arsitektur PoC: frontend Vue (Vercel) langsung ke Supabase. Backend Go **tidak** di-deploy.

## Environment Variables (Vercel > Settings > Environment Variables)
| Variable | Nilai | Environment |
| --- | --- | --- |
| `VITE_SUPABASE_URL` | `https://<project-ref>.supabase.co` | Production, Preview |
| `VITE_SUPABASE_ANON_KEY` | Publishable key `sb_publishable_...` | Production, Preview |

Jangan pernah memasukkan `service_role` / secret key ke variabel `VITE_*`, karena semua variabel `VITE_*` ikut ter-bundle ke browser.

## Checklist
1. **Database**
   - [ ] Migrasi `supabase/migrations/*.sql` sudah diterapkan (bukan `migrations/` atau `supabase_queries.sql` yang legacy).
   - [ ] Security Advisor di dashboard Supabase tanpa ERROR.
   - [ ] Minimal satu akun supervisor sudah dibuat dan role-nya diset via SQL.
2. **Auth**
   - [ ] Site URL & Redirect URLs berisi domain Vercel production.
   - [ ] Confirm email: dimatikan untuk PoC, atau SMTP custom sudah dipasang.
3. **Vercel**
   - [ ] Root Directory = `frontend`.
   - [ ] Kedua env var terisi, lalu redeploy setelah mengubahnya (variabel Vite dibaca saat build).
4. **Smoke test production**
   - [ ] Daftar/login sales → input kunjungan dengan GPS + foto → muncul di dashboard supervisor.
   - [ ] Mode pesawat → input kunjungan → online lagi → kunjungan tersinkron otomatis.
   - [ ] Login sales tidak bisa membuka `/supervisor`.
   - [ ] App bisa di-install (Add to Home Screen) di Android.
