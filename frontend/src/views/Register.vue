<template>
  <v-container class="fill-height d-flex align-center justify-center" style="min-height: 100vh">
    <v-card width="460" max-width="100%" elevation="4" class="pa-2">
      <v-card-title class="text-center pt-4 text-headline-small">Daftar Akun Agritec</v-card-title>
      <v-card-text>
        <v-alert v-if="error" type="error" variant="tonal" class="mb-4" density="compact">{{ error }}</v-alert>
        <v-alert v-if="success" type="success" variant="tonal" class="mb-4" density="compact">{{ success }}</v-alert>

        <v-form ref="form" @submit.prevent="handleRegister" :disabled="!!success">
          <v-text-field v-model.trim="fullName" label="Nama Lengkap" prepend-inner-icon="mdi-account"
            variant="outlined" :rules="[rules.required]" />
          <v-text-field v-model.trim="phone" label="No. HP (opsional)" prepend-inner-icon="mdi-phone"
            variant="outlined" type="tel" />
          <v-text-field v-model.trim="email" label="Email" type="email" prepend-inner-icon="mdi-email"
            variant="outlined" :rules="[rules.required, rules.email]" />
          <v-text-field v-model="password" label="Password" type="password" prepend-inner-icon="mdi-lock"
            variant="outlined" :rules="[rules.required, rules.min]" autocomplete="new-password" />
          <v-select v-model="role" :items="roles" item-title="label" item-value="value"
            label="Daftar sebagai" variant="outlined" prepend-inner-icon="mdi-badge-account" />
          <div class="text-body-small text-medium-emphasis mb-2">
            Akun Supervisor tidak bisa didaftarkan sendiri. Hubungi admin untuk mendapatkan akses supervisor.
          </div>

          <v-btn type="submit" color="primary" block size="large" class="mt-2" :loading="loading">Daftar Sekarang</v-btn>
          <div class="text-center mt-4 text-body-medium">
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
      fullName: '',
      phone: '',
      email: '',
      password: '',
      role: 'farmer',
      roles: [
        { value: 'farmer', label: 'Petani' },
        { value: 'sales', label: 'Sales Lapangan' },
        { value: 'distributor', label: 'Distributor / Kios' },
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
