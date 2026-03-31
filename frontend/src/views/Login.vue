<template>
  <v-container fluid class="fill-height d-flex align-center justify-center">
    <v-card width="400" elevation="4">
      <v-card-title class="text-center py-4 text-h5 primary--text">Login Agritec</v-card-title>
      <v-card-text>
        <v-alert v-if="error" type="error" class="mb-4" dense>{{ error }}</v-alert>

        <v-form @submit.prevent="handleLogin">
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
export default {
  data() {
    return {
      email: '',
      password: '',
      error: null,
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

        // Redirect based on role fetched from Database JWT
        if (user.role === 'sales') this.$router.push('/sales')
        else if (user.role === 'supervisor') this.$router.push('/supervisor')
        else if (user.role === 'farmer') this.$router.push('/farmer')
        else this.$router.push('/login') // fallback

      } catch (err) {
        this.error = err.message || 'Login Gagal'
      } finally {
        this.loading = false
      }
    }
  }
}
</script>
