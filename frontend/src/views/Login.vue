<template>
  <v-container fluid class="fill-height login-bg">
    <v-row justify="center" align="center" class="h-100">
      <v-col cols="12" sm="8" md="6" lg="4">
        <v-card elevation="12" rounded="xl" class="overflow-hidden glass-card">
          <div class="primary-gradient pa-6 text-center text-white">
            <v-icon size="64" class="mb-2">mdi-leaf-circle</v-icon>
            <h2 class="text-h4 font-weight-bold">SFA Agritec</h2>
            <p class="text-subtitle-1 opacity-80 mt-1">Sistem Otomasi & Manajemen Agribisnis</p>
          </div>
          
          <v-card-text class="pa-8">
            <v-alert v-if="error" type="error" variant="tonal" class="mb-4 text-caption rounded-lg" density="compact" icon="mdi-alert-circle">
              {{ error }}
            </v-alert>

            <v-form @submit.prevent="handleLogin">
              <div class="text-subtitle-2 font-weight-bold mb-1 text-grey-darken-1">Email / Username</div>
              <v-text-field
                v-model="email"
                prepend-inner-icon="mdi-email-outline"
                required
                variant="outlined"
                density="comfortable"
                color="primary"
                rounded="lg"
                bg-color="grey-lighten-4"
                placeholder="Masukkan email"
                class="mb-3"
              ></v-text-field>

              <div class="text-subtitle-2 font-weight-bold mb-1 text-grey-darken-1 mt-2">Password</div>
              <v-text-field
                v-model="password"
                :type="showPassword ? 'text' : 'password'"
                prepend-inner-icon="mdi-lock-outline"
                :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                @click:append-inner="showPassword = !showPassword"
                required
                variant="outlined"
                density="comfortable"
                color="primary"
                rounded="lg"
                bg-color="grey-lighten-4"
                placeholder="Masukkan password"
              ></v-text-field>

              <div class="d-flex justify-space-between align-center mb-6 mt-1">
                <v-checkbox label="Ingat Saya" color="primary" density="compact" hide-details></v-checkbox>
                <a href="#" class="text-primary text-caption font-weight-bold text-decoration-none">Lupa Password?</a>
              </div>

              <v-btn type="submit" color="primary" block size="x-large" class="mt-2 font-weight-bold" rounded="pill" :loading="loading" elevation="4">
                MASUK
              </v-btn>
              
              <div class="text-center mt-6 text-caption text-medium-emphasis">
                 Belum Punya Akun? 
                 <router-link to="/register" class="text-primary font-weight-bold text-decoration-none ml-1">Daftar Ekosistem</router-link>
              </div>
            </v-form>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  data() {
    return {
      email: '',
      password: '',
      showPassword: false,
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

        if (user.role === 'sales') this.$router.push('/sales')
        else if (user.role === 'supervisor') this.$router.push('/supervisor')
        else if (user.role === 'farmer') this.$router.push('/farmer')
        else if (user.role === 'distributor') this.$router.push('/distributor')
        else this.$router.push('/login') 

      } catch (err) {
        this.error = err.message || 'Login Gagal. Silakan periksa kredensial Anda.'
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.login-bg {
  background: linear-gradient(135deg, #e8f5e9 0%, #c8e6c9 100%);
  position: relative;
  overflow: hidden;
}

.login-bg::before {
  content: '';
  position: absolute;
  top: -10%;
  left: -10%;
  width: 50vw;
  height: 50vw;
  background: radial-gradient(circle, rgba(76,175,80,0.15) 0%, rgba(255,255,255,0) 70%);
  border-radius: 50%;
}

.login-bg::after {
  content: '';
  position: absolute;
  bottom: -15%;
  right: -5%;
  width: 60vw;
  height: 60vw;
  background: radial-gradient(circle, rgba(139,195,74,0.15) 0%, rgba(255,255,255,0) 70%);
  border-radius: 50%;
}

.glass-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255,255,255,0.5);
  z-index: 10;
}

.primary-gradient {
  background: linear-gradient(135deg, #2E7D32 0%, #4CAF50 100%);
}

.opacity-80 {
  opacity: 0.8;
}
</style>
