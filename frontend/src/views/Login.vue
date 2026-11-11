<template>
  <v-container class="fill-height d-flex align-center justify-center" style="min-height: 100vh">
    <v-card width="420" max-width="100%" elevation="4" class="pa-2">
      <v-card-title class="text-center pt-4">
        <v-icon color="primary" size="48">mdi-sprout</v-icon>
        <div class="text-headline-small mt-2">Login Agritec</div>
        <div class="text-body-medium text-medium-emphasis">Sales Force Automation</div>
      </v-card-title>
      <v-card-text>
        <v-alert v-if="error" type="error" variant="tonal" class="mb-4" density="compact">{{ error }}</v-alert>
        <v-alert v-if="info" type="success" variant="tonal" class="mb-4" density="compact">{{ info }}</v-alert>

        <v-form ref="form" @submit.prevent="handleLogin">
          <v-text-field v-model.trim="email" label="Email" type="email" autocomplete="email"
            prepend-inner-icon="mdi-email" variant="outlined" :rules="[rules.required, rules.email]" />
          <v-text-field v-model="password" label="Password" :type="show ? 'text' : 'password'" autocomplete="current-password"
            prepend-inner-icon="mdi-lock" :append-inner-icon="show ? 'mdi-eye-off' : 'mdi-eye'"
            @click:append-inner="show = !show" variant="outlined" :rules="[rules.required]" />

          <v-btn type="submit" color="primary" block size="large" class="mt-2" :loading="loading">Masuk</v-btn>
          <div class="d-flex justify-space-between mt-4 text-body-medium">
            <a href="#" @click.prevent="forgot">Lupa password?</a>
            <router-link to="/register">Daftar akun baru</router-link>
          </div>
        </v-form>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script>
import { homeForRole } from '../services/format'

export default {
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
