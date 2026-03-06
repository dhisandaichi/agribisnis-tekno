<template>
  <v-container fluid class="fill-height register-bg">
    <v-row justify="center" align="center" class="h-100">
      <v-col cols="12" sm="9" md="7" lg="5">
        <v-card elevation="12" rounded="xl" class="overflow-hidden glass-card">
          <div class="primary-gradient pa-6 text-center text-white">
            <v-icon size="56" class="mb-2">mdi-account-plus</v-icon>
            <h2 class="text-h5 font-weight-bold">Bergabung dengan Ekosistem Agritec</h2>
            <p class="text-caption opacity-80 mt-1">Daftarkan akun untuk akses penuh sistem SFA</p>
          </div>

          <v-card-text class="pa-8">
            <v-alert v-if="error" type="error" variant="tonal" class="mb-4 text-caption rounded-lg" density="compact" icon="mdi-alert-circle">
              {{ error }}
            </v-alert>
            <v-alert v-if="success" type="success" variant="tonal" class="mb-4 text-caption rounded-lg" density="compact" icon="mdi-check-circle">
              Pendaftaran berhasil! Mengarahkan ke halaman login...
            </v-alert>

            <v-form @submit.prevent="handleRegister">
              <div class="text-subtitle-2 font-weight-bold mb-1 text-grey-darken-1">Email</div>
              <v-text-field
                v-model="email"
                prepend-inner-icon="mdi-email-outline"
                required
                variant="outlined"
                density="comfortable"
                color="primary"
                rounded="lg"
                bg-color="grey-lighten-4"
                placeholder="contoh@agritec.com"
                class="mb-3"
              ></v-text-field>

              <div class="text-subtitle-2 font-weight-bold mb-1 text-grey-darken-1">Password</div>
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
                placeholder="Minimal 6 karakter"
                class="mb-3"
              ></v-text-field>

              <div class="text-subtitle-2 font-weight-bold mb-1 text-grey-darken-1">Jabatan / Role</div>
              <v-select
                v-model="role"
                :items="roleOptions"
                item-title="label"
                item-value="value"
                prepend-inner-icon="mdi-briefcase-outline"
                variant="outlined"
                density="comfortable"
                color="primary"
                rounded="lg"
                bg-color="grey-lighten-4"
                class="mb-4"
              ></v-select>

              <v-btn type="submit" color="primary" block size="x-large" class="font-weight-bold" rounded="pill" :loading="loading" elevation="4">
                DAFTAR SEKARANG
              </v-btn>

              <div class="text-center mt-6 text-caption text-medium-emphasis">
                Sudah punya akun?
                <router-link to="/login" class="text-primary font-weight-bold text-decoration-none ml-1">Masuk</router-link>
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
      role: 'sales',
      roleOptions: [
        { label: 'Sales (Tenaga Penjual)', value: 'sales' },
        { label: 'Supervisor', value: 'supervisor' },
        { label: 'Petani (Farmer)', value: 'farmer' },
        { label: 'Distributor / Kios', value: 'distributor' }
      ],
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

<style scoped>
.register-bg {
  background: linear-gradient(135deg, #e8f5e9 0%, #c8e6c9 100%);
}
.glass-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}
.primary-gradient {
  background: linear-gradient(135deg, #1B5E20 0%, #43A047 100%);
}
.opacity-80 {
  opacity: 0.8;
}
</style>
