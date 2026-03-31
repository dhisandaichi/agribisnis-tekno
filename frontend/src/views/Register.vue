<template>
  <v-container fluid class="fill-height d-flex align-center justify-center">
    <v-card width="400" elevation="4">
      <v-card-title class="text-center py-4 text-h5 primary--text">Daftar Akun Agritec</v-card-title>
      <v-card-text>
        <v-alert v-if="error" type="error" class="mb-4" dense>{{ error }}</v-alert>
        <v-alert v-if="success" type="success" class="mb-4" dense>Pendaftaran berhasil. Silakan login.</v-alert>

        <v-form @submit.prevent="handleRegister">
          <v-text-field
            v-model="email"
            label="Email"
            prepend-inner-icon="mdi-email"
            required
            variant="outlined"
          ></v-text-field>
          <v-text-field
            v-model="password"
            label="Password"
            type="password"
            prepend-inner-icon="mdi-lock"
            required
            variant="outlined"
          ></v-text-field>

          <v-select
            v-model="role"
            :items="['sales', 'supervisor', 'farmer', 'distributor']"
            label="Pilih Role (Jabatan)"
            variant="outlined"
          ></v-select>

          <v-btn type="submit" color="primary" block size="large" class="mt-4" :loading="loading">
            Daftar Sekarang
          </v-btn>
          <div class="text-center mt-3">
            <router-link to="/login">Sudah punya akun? Masuk.</router-link>
          </div>
        </v-form>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script>
export default {
  data() {
    return {
      email: '',
      password: '',
      role: 'sales',
      error: null,
      success: false,
      loading: false
    }
  },
  methods: {
    async handleRegister() {
      this.error = null
      this.success = false
      this.loading = true

      try {
        await this.$store.dispatch('register', {
          email: this.email,
          password: this.password,
          role: this.role
        })
        this.success = true
        setTimeout(() => {
          this.$router.push('/login')
        }, 1500)
      } catch (err) {
        this.error = err.message || 'Gagal mendaftar. Email mungkin sudah terpakai.'
      } finally {
        this.loading = false
      }
    }
  }
}
</script>
