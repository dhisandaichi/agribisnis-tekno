<template>
  <AuthShell>
    <div class="mb-6">
      <div class="text-overline text-primary font-weight-bold" style="letter-spacing: .12em">Selamat datang kembali</div>
      <h2 class="auth-title">Masuk ke SFA Agritec</h2>
      <div class="text-body-medium text-medium-emphasis">Gunakan email dan password akun Anda.</div>
    </div>

    <v-alert v-if="error" type="error" variant="tonal" class="mb-4" density="compact">{{ error }}</v-alert>
    <v-alert v-if="info" type="info" variant="tonal" class="mb-4" density="compact">{{ info }}</v-alert>

    <v-form ref="form" @submit.prevent="handleLogin">
      <v-text-field v-model.trim="email" label="Email" type="email" autocomplete="email" color="primary"
        prepend-inner-icon="mdi-email-outline" variant="outlined" :rules="[rules.required, rules.email]" class="mb-1" />
      <v-text-field v-model="password" label="Password" :type="show ? 'text' : 'password'" autocomplete="current-password"
        color="primary" prepend-inner-icon="mdi-lock-outline" :append-inner-icon="show ? 'mdi-eye-off' : 'mdi-eye'"
        @click:append-inner="show = !show" variant="outlined" :rules="[rules.required]" />

      <div class="d-flex justify-end mb-4 mt-n2">
        <a href="#" class="auth-link text-body-medium" @click.prevent="forgot">Lupa password?</a>
      </div>

      <v-btn type="submit" color="primary" block size="x-large" class="auth-submit" :loading="loading" append-icon="mdi-arrow-right">
        Masuk
      </v-btn>
    </v-form>

    <div class="auth-divider"><span>Belum punya akun?</span></div>
    <v-btn to="/register" block variant="outlined" color="secondary" size="large" prepend-icon="mdi-account-plus-outline">
      Daftar akun baru
    </v-btn>
  </AuthShell>
</template>

<script>
import { homeForRole } from '../services/format'
import AuthShell from '../components/AuthShell.vue'

export default {
  components: { AuthShell },
  data() {
    return {
      email: '',
      password: '',
      show: false,
      error: null,
      info: null,
      loading: false,
      rules: {
        required: (v) => !!v || 'Wajib diisi',
        email: (v) => /.+@.+\..+/.test(v) || 'Format email tidak valid',
      },
    }
  },
  methods: {
    async handleLogin() {
      const { valid } = await this.$refs.form.validate()
      if (!valid) return
      this.error = null
      this.info = null
      this.loading = true
      try {
        const profile = await this.$store.dispatch('login', { email: this.email, password: this.password })
        this.$router.push(homeForRole(profile.role))
      } catch (err) {
        this.error = err.message || 'Login gagal'
      } finally {
        this.loading = false
      }
    },
    forgot() {
      this.error = null
      this.info = 'Hubungi supervisor Agritec untuk mendapatkan password sementara, lalu ganti password setelah login (ikon kunci di kanan atas).'
    },
  },
}
</script>

<style scoped>
.auth-title { font-size: 26px; font-weight: 800; color: var(--brand-navy-dark); line-height: 1.2; margin: 4px 0 6px; }
.auth-link { color: var(--brand-navy); font-weight: 500; text-decoration: none; }
.auth-link:hover { text-decoration: underline; }
.auth-submit { font-weight: 700; letter-spacing: .02em; }
.auth-divider { display: flex; align-items: center; gap: 12px; color: #667085; font-size: 13px; margin: 22px 0 14px; }
.auth-divider::before, .auth-divider::after { content: ''; flex: 1; height: 1px; background: #e4e7ec; }
</style>
