# SFA Agritec - Backend Vercel Deployment Guide

## Prerequisites

- Go >= 1.22
- Vercel CLI: `npm i -g vercel`

## Environment Variables (WAJIB diset di Vercel Dashboard)

| Variable | Keterangan |
|---|---|
| `SUPABASE_DB_URL` | PostgreSQL connection string dari Supabase (Settings → Database → Connection String → URI). Pastikan `?sslmode=require` ada di akhir URL. |
| `JWT_SECRET` | Secret key untuk signing JWT. Gunakan string random panjang (min 32 karakter). |

> **Penting:** Jangan commit nilai ini ke Git. Set hanya melalui Vercel Dashboard atau `vercel env add`.

## Deploy Backend ke Vercel

```bash
# Dari folder /backend
cd backend

# Pertama kali - setup project
vercel

# Deploy production
vercel --prod

# Atau set env via CLI:
vercel env add SUPABASE_DB_URL production
vercel env add JWT_SECRET production
```

## Custom Domain

1. Masuk ke **Vercel Dashboard** → Project Backend → **Settings** → **Domains**
2. Tambahkan domain, misalnya: `api.agritec.id`
3. Update DNS record di domain registrar:
   - Tambahkan **CNAME** record: `api` → `cname.vercel-dns.com`
4. Update `VITE_API_URL` di frontend `.env` menjadi `https://api.agritec.id/api/v1`

## Verifikasi API

Setelah deploy, test endpoint:

```bash
# Health check
curl https://api.agritec.id/

# Login
curl -X POST https://api.agritec.id/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"sales@agritec.com","password":"password123"}'
```
