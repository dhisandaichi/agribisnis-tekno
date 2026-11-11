<template>
  <div class="auth-shell">
    <!-- Panel merek: tampil di kiri pada desktop, menjadi header di HP -->
    <aside class="auth-brand">
      <div class="auth-brand__bg" />
      <div class="auth-brand__content">
        <router-link to="/" class="auth-brand__home">
          <span class="brand-logo-badge"><img src="/brand/logo.png" alt="Logo Agrina 99" width="48" height="48" /></span>
          <span>
            <span class="auth-brand__kicker">Pupuk Organik <b>Premium</b></span>
            <span class="auth-brand__name">AGRINA 99</span>
          </span>
        </router-link>

        <div class="auth-brand__hero">
          <div class="auth-brand__pill">SFA Agritec</div>
          <h1>Tanah Subur, <br />Panen Melimpah, <br /><span>Petani Sejahtera.</span></h1>
          <p>Satu aplikasi untuk tim sales lapangan, distributor dan kios, serta petani mitra PT. Agribisnis Technology Indonesia.</p>
          <ul class="auth-brand__points">
            <li><v-icon size="18">mdi-map-marker-check</v-icon> Kunjungan sales dengan GPS & foto</li>
            <li><v-icon size="18">mdi-store</v-icon> Order kios & distributor lebih cepat</li>
            <li><v-icon size="18">mdi-sprout</v-icon> Edukasi budidaya untuk petani</li>
          </ul>
        </div>

        <img class="auth-brand__packs" src="/brand/products.jpg" alt="Produk Agrina 99, Dolomite Super Premium dan Top Super Booster" />
      </div>
    </aside>

    <main class="auth-main">
      <div class="auth-card">
        <slot />
      </div>
      <div class="auth-foot">
        <router-link to="/"><v-icon size="16">mdi-arrow-left</v-icon> Lihat produk Agrina 99</router-link>
        <span>© {{ year }} PT. Agribisnis Technology Indonesia</span>
      </div>
    </main>
  </div>
</template>

<script>
export default {
  data() {
    return { year: new Date().getFullYear() }
  },
}
</script>

<style scoped>
.auth-shell {
  min-height: 100vh;
  display: grid;
  grid-template-columns: minmax(0, 1.05fr) minmax(0, 1fr);
  background: #f6f8f5;
}
.auth-brand {
  position: relative;
  overflow: hidden;
  color: #fff;
  position: sticky;
  top: 0;
  height: 100vh;
  min-height: 560px;
  background: linear-gradient(150deg, var(--brand-navy-dark) 0%, var(--brand-navy) 38%, var(--brand-green-dark) 100%);
}
.auth-brand__bg {
  position: absolute;
  inset: 0;
  background: url('/brand/field.jpg') center / cover no-repeat;
  opacity: .22;
  mix-blend-mode: luminosity;
}
.auth-brand__bg::after {
  content: '';
  position: absolute;
  inset: 0;
  background: radial-gradient(circle at 85% 10%, rgba(0, 154, 74, .55), transparent 55%);
}
.auth-brand__content {
  position: relative;
  height: 100%;
  display: flex;
  flex-direction: column;
  gap: 28px;
  padding: 40px clamp(24px, 5vw, 64px);
}
.auth-brand__home {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  color: #fff;
  text-decoration: none;
}
.auth-brand__kicker { display: block; font-size: 11px; letter-spacing: .08em; text-transform: uppercase; opacity: .9; }
.auth-brand__kicker b { color: #ffb37f; }
.auth-brand__name { display: block; font-size: 22px; font-weight: 800; letter-spacing: .02em; }
.auth-brand__hero { margin-top: auto; }
.auth-brand__pill {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 999px;
  background: rgba(255, 255, 255, .14);
  border: 1px solid rgba(255, 255, 255, .25);
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 14px;
}
.auth-brand h1 {
  font-size: clamp(30px, 3.4vw, 44px);
  line-height: 1.12;
  font-weight: 800;
  margin: 0 0 14px;
}
.auth-brand h1 span { color: #7ee2a8; }
.auth-brand p { max-width: 460px; opacity: .88; margin: 0 0 18px; line-height: 1.6; }
.auth-brand__points { list-style: none; padding: 0; margin: 0; display: grid; gap: 8px; }
.auth-brand__points li { display: flex; align-items: center; gap: 10px; font-size: 15px; }
.auth-brand__points .v-icon { color: #7ee2a8; }
.auth-brand__packs {
  width: 100%;
  max-width: 560px;
  max-height: 32vh;
  object-fit: cover;
  object-position: left center;
  border-radius: 16px;
  box-shadow: 0 18px 40px rgba(0, 0, 0, .35);
  margin-top: auto;
}

.auth-main {
  position: relative;
  z-index: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px 24px;
}
.auth-card {
  width: 100%;
  max-width: 440px;
  background: #fff;
  border-radius: 20px;
  padding: 32px 28px;
  box-shadow: 0 1px 2px rgba(16, 24, 40, .06), 0 12px 32px rgba(16, 24, 40, .08);
  border-top: 5px solid var(--brand-green);
}
.auth-foot {
  margin-top: 20px;
  width: 100%;
  max-width: 440px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  justify-content: space-between;
  font-size: 13px;
  color: #667085;
}
.auth-foot a { color: var(--brand-navy); text-decoration: none; font-weight: 500; }

@media (max-width: 959px) {
  .auth-shell { grid-template-columns: minmax(0, 1fr); }
  .auth-brand { position: relative; height: auto; min-height: 0; }
  .auth-brand__content { padding: 22px 20px 64px; gap: 16px; }
  .auth-brand__hero p, .auth-brand__points, .auth-brand__packs { display: none; }
  .auth-brand__pill { margin-bottom: 8px; }
  .auth-brand h1 { font-size: 24px; margin: 0; }
  .auth-brand h1 br { display: none; }
  .auth-main { justify-content: flex-start; padding: 0 16px 24px; margin-top: -44px; }
  .auth-card { padding: 24px 20px; }
}
</style>
