<template>
  <v-container fluid class="fill-height d-flex align-center justify-center">
    <v-card width="400" elevation="4">
      <v-card-title class="text-center py-4 text-h5 text-primary">Daftar Akun Agritec</v-card-title>
      <v-card-text>
        <v-alert v-if="error" type="error" class="mb-4" density="compact">{{ error }}</v-alert>
        <v-alert v-if="needsConfirmation" type="success" class="mb-4" density="compact">
          Pendaftaran berhasil. Cek email Anda untuk konfirmasi, lalu login.
        </v-alert>

        <v-form v-if="!needsConfirmation" @submit.prevent="handleRegister">
          <v-text-field
            v-model="fullName"
            label="Nama Lengkap"
            autocomplete="name"
            prepend-inner-icon="mdi-account"
            required
            variant="outlined"
          ></v-text-field>
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
            label="Password (min. 8 karakter)"
            type="password"
            autocomplete="new-password"
            prepend-inner-icon="mdi-lock"
            required
            variant="outlined"
          ></v-text-field>

          <!-- Supervisor & distributor ditetapkan oleh admin, bukan lewat pendaftaran. -->
          <v-select
            v-model="role"
            :items="roles"
            label="Daftar sebagai"
            variant="outlined"
          ></v-select>

          <v-btn type="submit" color="primary" block size="large" class="mt-4" :loading="loading">
            Daftar Sekarang
          </v-btn>
        </v-form>
        <div class="text-center mt-3">
          <router-link to="/login">Sudah punya akun? Masuk.</router-link>
        </div>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script>
import { homeFor } from '../store'

export default {
  data() {
    return {
      fullName: '',
      email: '',
      password: '',
      role: 'sales',
      roles: [
        { title: 'Sales Lapangan', value: 'sales' },
        { title: 'Petani', value: 'farmer' }
      ],
      error: null,
      needsConfirmation: false,
      loading: false
    }
  },
  methods: {
    async handleRegister() {
      this.error = null
      if (this.password.length < 8) {
        this.error = 'Password minimal 8 karakter.'
        return
      }
      this.loading = true

      try {
        const user = await this.$store.dispatch('register', {
          email: this.email,
          password: this.password,
          role: this.role,
          fullName: this.fullName.trim()
        })
        if (user) this.$router.push(homeFor(user.role))
        else this.needsConfirmation = true
      } catch (err) {
        this.error = err.message || 'Gagal mendaftar. Email mungkin sudah terpakai.'
      } finally {
        this.loading = false
      }
    }
  }
}
</script>
