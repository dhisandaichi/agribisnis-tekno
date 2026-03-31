# Frontend Agritec - SFA Progressive Web App (Vue 3 + Vite)

This is the frontend component of the Agritec Sales Force Automation application, built with Vue 3, Vite, and Vuetify. It is configured to run as a Progressive Web Application (PWA).

## Infrastructure & Deployment
The frontend is configured for deployment to **Vercel**, and temporarily will use **Supabase** as its backend database/auth (as an MVP substitute or complement for the Golang Microservices).

### Environment Variables
Before deploying, make sure you configure your Vercel project with the required API keys. Copy the `.env.example` file to a new file named `.env.local` to configure them for local development.

```bash
# In the `frontend` directory:
cp .env.example .env.local
```

Required keys for Supabase:
- `VITE_SUPABASE_URL`
- `VITE_SUPABASE_ANON_KEY`

### Deployment to Vercel

1. **Import the Project:** Go to the Vercel Dashboard and click "Add New Project" > "Import from Git".
2. **Select Framework:** Once imported, Vercel should automatically detect **Vite**.
3. **Configure Root Directory:** Make sure to set the Root Directory to `frontend`.
4. **Environment Variables:** In the Vercel deploy configuration under "Environment Variables", paste the values from your Supabase Project Settings (URL and anon key).
5. **Deploy:** Click deploy. Vercel will process the `vercel.json` file automatically to correctly route requests for your Single Page Application (SPA).

## Local Development

```bash
# Install dependencies
npm install

# Run local development server
npm run dev

# Build for production
npm run build
```

## Features
- **PWA Ready**: Offline caching, service workers, installable.
- **Supabase Connected**: Pre-configured Supabase SDK client snippet at `src/services/supabase.js`.
- **SPA Routing configured**: Includes fallbacks configured in Vercel.
