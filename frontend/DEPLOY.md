# SFA Agritec - Frontend Vercel Deployment Guide

## Prerequisites

- Node.js >= 18
- Vercel CLI: `npm i -g vercel`

## Environment Variables (set via Vercel Dashboard or CLI)

| Variable | Value |
|---|---|
| `VITE_API_URL` | `https://<your-backend-domain>/api/v1` |
| `VITE_SUPABASE_URL` | Your Supabase project URL |
| `VITE_SUPABASE_ANON_KEY` | Your Supabase anon key |

## Deploy to Vercel (custom domain)

```bash
# 1. Install dependencies
npm install

# 2. Build production
npm run build

# 3. Deploy (first time - interactive, links to project)
vercel --prod

# 4. Set custom domain via Vercel Dashboard:
#    Project -> Settings -> Domains -> Add domain
#    e.g: app.agritec.id
```

## Firebase Hosting (Alternative)

```bash
npm install -g firebase-tools
firebase login
npm run build
firebase deploy
```
