<template>
  <v-app>
    <v-app-bar class="brand-appbar" density="comfortable" v-if="isAuthenticated && !isPublic">
      <template #prepend>
        <router-link to="/" class="brand-logo-badge ml-2" title="Beranda Agrina 99">
          <img src="/brand/logo.png" alt="Logo Agrina 99" width="34" height="34" />
        </router-link>
      </template>
      <v-app-bar-title>
        <div style="line-height: 1.1">
          <div class="font-weight-bold">SFA Agritec</div>
          <div class="text-body-small d-none d-sm-block" style="opacity: .85">PT. Agribisnis Technology Indonesia</div>
        </div>
      </v-app-bar-title>
      <template #append>
        <v-chip v-if="!online" color="warning" variant="flat" size="small" class="mr-2" prepend-icon="mdi-wifi-off">
          Offline
        </v-chip>
        <div class="text-right mr-2 d-none d-sm-block" style="line-height: 1.2">
          <div class="text-body-medium">{{ profile?.full_name || profile?.email }}</div>
          <div class="text-body-small">{{ roleLabel }}</div>
        </div>
        <v-btn icon="mdi-key-variant" title="Ganti password" @click="pw.show = true" />
        <v-btn icon="mdi-logout" title="Keluar" @click="logout" />
      </template>
    </v-app-bar>

    <v-main>
      <router-view />
    </v-main>

    <v-dialog v-model="pw.show" max-width="420">
      <v-card title="Ganti Password">
        <v-card-text>
          <v-alert v-if="pw.error" type="error" variant="tonal" density="compact" class="mb-3">{{ pw.error }}</v-alert>
          <v-alert v-if="pw.done" type="success" variant="tonal" density="compact" class="mb-3">Password berhasil diubah.</v-alert>
          <v-form ref="pwForm">
            <v-text-field v-model="pw.old" label="Password lama" type="password" variant="outlined" :rules="[(v) => !!v || 'Wajib diisi']" />
            <v-text-field v-model="pw.new" label="Password baru" type="password" variant="outlined"
              :rules="[(v) => (v && v.length >= 8) || 'Minimal 8 karakter']" />
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn @click="closePw">Tutup</v-btn>
          <v-btn color="primary" variant="flat" :loading="pw.loading" @click="changePassword">Simpan</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-app>
</template>

<script>
import { put } from './services/api'
import { ROLE_LABELS } from './services/format'

export default {
  data() {
    return { pw: { show: false, old: '', new: '', loading: false, error: null, done: false } }
  },
  computed: {
    isAuthenticated() { return this.$store.getters.isAuthenticated },
    isPublic() { return !!this.$route.meta.public || !!this.$route.meta.guest },
    profile() { return this.$store.state.profile },
    online() { return this.$store.state.online },
    roleLabel() { return ROLE_LABELS[this.$store.getters.userRole] || '' },
  },
  methods: {
    logout() {
      this.$store.dispatch('logout')
      this.$router.push('/login')
    },
    closePw() {
      this.pw = { show: false, old: '', new: '', loading: false, error: null, done: false }
    },
    async changePassword() {
      const { valid } = await this.$refs.pwForm.validate()
      if (!valid) return
      this.pw.loading = true
      this.pw.error = null
      try {
        await put('/me/password', { old_password: this.pw.old, new_password: this.pw.new })
        this.pw.done = true
        this.pw.old = ''
        this.pw.new = ''
        this.$refs.pwForm.resetValidation()
      } catch (err) {
        this.pw.error = err.message
      } finally {
        this.pw.loading = false
      }
    },
  },
}
</script>
