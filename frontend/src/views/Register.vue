<template>
  <AuthShell>
    <div class="mb-5">
      <div class="text-overline text-primary font-weight-bold" style="letter-spacing: .12em">Bergabung bersama Agrina 99</div>
      <h2 class="auth-title">Buat akun baru</h2>
      <div class="text-body-medium text-medium-emphasis">Pilih peran Anda, lalu lengkapi data di bawah.</div>
    </div>

    <v-alert v-if="error" type="error" variant="tonal" class="mb-4" density="compact">{{ error }}</v-alert>
    <v-alert v-if="success" type="success" variant="tonal" class="mb-4" density="compact">{{ success }}</v-alert>

    <v-form ref="form" @submit.prevent="handleRegister" :disabled="!!success">
      <div class="text-body-medium font-weight-medium mb-2">Daftar sebagai</div>
      <div class="role-grid mb-5" role="radiogroup" aria-label="Daftar sebagai">
        <button v-for="r in roles" :key="r.value" type="button" role="radio" :aria-checked="role === r.value"
          class="role-card" :class="{ 'role-card--active': role === r.value }" @click="role = r.value">
          <v-icon size="26">{{ r.icon }}</v-icon>
          <span class="role-card__label">{{ r.label }}</span>
          <span class="role-card__hint">{{ r.hint }}</span>
        </button>
      </div>

      <v-text-field v-model.trim="fullName" label="Nama lengkap" prepend-inner-icon="mdi-account-outline" color="primary"
        variant="outlined" :rules="[rules.required]" autocomplete="name" class="mb-1" />
      <v-text-field v-model.trim="phone" label="No. HP / WhatsApp (opsional)" prepend-inner-icon="mdi-whatsapp" color="primary"
        variant="outlined" type="tel" autocomplete="tel" class="mb-1" />
      <v-text-field v-model.trim="email" label="Email" type="email" prepend-inner-icon="mdi-email-outline" color="primary"
        variant="outlined" :rules="[rules.required, rules.email]" autocomplete="email" class="mb-1" />
      <v-text-field v-model="password" label="Password" :type="show ? 'text' : 'password'" prepend-inner-icon="mdi-lock-outline"
        :append-inner-icon="show ? 'mdi-eye-off' : 'mdi-eye'" @click:append-inner="show = !show" color="primary"
        variant="outlined" :rules="[rules.required, rules.min]" autocomplete="new-password" hint="Minimal 8 karakter" />

      <v-btn type="submit" color="primary" block size="x-large" class="auth-submit mt-4" :loading="loading" append-icon="mdi-arrow-right">
        Daftar Sekarang
      </v-btn>
      <div class="text-body-small text-medium-emphasis mt-3 d-flex ga-2">
        <v-icon size="16" color="secondary">mdi-shield-account-outline</v-icon>
        <span>Akun Supervisor dibuat oleh admin. Hubungi admin Agritec untuk akses supervisor.</span>
      </div>
    </v-form>

    <div class="auth-divider"><span>Sudah punya akun?</span></div>
    <v-btn to="/login" block variant="outlined" color="secondary" size="large" prepend-icon="mdi-login">Masuk</v-btn>
  </AuthShell>
</template>

<script>
import AuthShell from '../components/AuthShell.vue'

export default {
  components: { AuthShell },
  data() {
    return {
      fullName: '',
      phone: '',
      email: '',
      password: '',
      show: false,
      role: ['farmer', 'sales', 'distributor'].includes(this.$route.query.role) ? this.$route.query.role : 'farmer',
      roles: [
        { value: 'farmer', label: 'Petani', hint: 'Tips budidaya', icon: 'mdi-sprout' },
        { value: 'distributor', label: 'Kios / Distributor', hint: 'Pesan stok', icon: 'mdi-storefront-outline' },
        { value: 'sales', label: 'Sales', hint: 'Kunjungan lapangan', icon: 'mdi-briefcase-outline' },
      ],
      error: null,
      success: null,
      loading: false,
      rules: {
        required: (v) => !!v || 'Wajib diisi',
        email: (v) => /.+@.+\..+/.test(v) || 'Format email tidak valid',
        min: (v) => (v && v.length >= 8) || 'Minimal 8 karakter',
      },
    }
  },
  methods: {
    async handleRegister() {
      const { valid } = await this.$refs.form.validate()
      if (!valid) return
      this.error = null
      this.loading = true
      try {
        await this.$store.dispatch('register', {
          email: this.email,
          password: this.password,
          role: this.role,
          fullName: this.fullName,
          phone: this.phone,
        })
        this.success = 'Pendaftaran berhasil. Mengarahkan ke halaman login...'
        setTimeout(() => this.$router.push('/login'), 1500)
      } catch (err) {
        this.error = err.message || 'Gagal mendaftar.'
      } finally {
        this.loading = false
      }
    },
  },
}
</script>

<style scoped>
.auth-title { font-size: 26px; font-weight: 800; color: var(--brand-navy-dark); line-height: 1.2; margin: 4px 0 6px; }
.auth-submit { font-weight: 700; letter-spacing: .02em; }
.auth-divider { display: flex; align-items: center; gap: 12px; color: #667085; font-size: 13px; margin: 22px 0 14px; }
.auth-divider::before, .auth-divider::after { content: ''; flex: 1; height: 1px; background: #e4e7ec; }
.role-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 10px; }
.role-card {
  display: flex; flex-direction: column; align-items: center; gap: 4px; text-align: center;
  padding: 14px 6px 12px; border-radius: 14px; border: 1.5px solid #e4e7ec; background: #fff;
  color: #475467; cursor: pointer; transition: all .15s ease; font: inherit;
}
.role-card:hover { border-color: var(--brand-green); }
.role-card:focus-visible { outline: 3px solid rgba(0, 154, 74, .35); outline-offset: 2px; }
.role-card--active { border-color: var(--brand-green); background: #ecf8f1; color: var(--brand-green-dark); box-shadow: 0 0 0 1px var(--brand-green) inset; }
.role-card__label { font-weight: 700; font-size: 13px; line-height: 1.2; color: var(--brand-navy-dark); }
.role-card__hint { font-size: 11px; color: #667085; line-height: 1.2; }
</style>
