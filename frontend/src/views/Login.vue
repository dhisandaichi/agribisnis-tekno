<template>
  <v-container fluid class="fill-height d-flex align-center justify-center">
    <v-card width="400" elevation="4">
      <v-card-title class="text-center py-4 text-h5 text-primary">Login Agritec</v-card-title>
      <v-card-text>
        <v-alert v-if="error" type="error" class="mb-4" density="compact">{{ error }}</v-alert>

        <v-form @submit.prevent="handleLogin">
          <v-text-field
            v-model="email"
            label="Email"
            type="email"
            autocomplete="email"
            prepend-inner-icon="mdi-email"
            required
            variant="outlined"
          ></v-text-field>
          <v-text-field
            v-model="password"
            label="Password"
            type="password"
            autocomplete="current-password"
            prepend-inner-icon="mdi-lock"
            required
            variant="outlined"
          ></v-text-field>

          <v-btn type="submit" color="primary" block size="large" class="mt-4" :loading="loading">
            Masuk
          </v-btn>
          <div class="text-center mt-3">
             <router-link to="/register">Belum punya akun? Daftar.</router-link>
          </div>
        </v-form>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script>
import { homeFor } from '../store'

export default {
  data() {
    return {
      email: '',
      password: '',
      error: this.$route.query.reason === 'role'
        ? 'Role akun Anda belum memiliki akses ke aplikasi ini.'
        : null,
      loading: false
    }
  },
  methods: {
    async handleLogin() {
      this.error = null
      this.loading = true

      try {
        const user = await this.$store.dispatch('login', {
          email: this.email,
          password: this.password
        })

        const home = homeFor(user.role)
        if (!home) {
          await this.$store.dispatch('logout')
          this.error = 'Role akun Anda belum memiliki akses ke aplikasi ini.'
          return
        }
        const redirect = this.$route.query.redirect
        this.$router.push(typeof redirect === 'string' && redirect.startsWith(home) ? redirect : home)
      } catch (err) {
        this.error = err.message || 'Login Gagal'
      } finally {
        this.loading = false
      }
    }
  }
}
</script>
