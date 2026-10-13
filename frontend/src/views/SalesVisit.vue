<template>
  <v-container>
    <v-row>
      <v-col cols="12" md="8" offset-md="2">
        <v-alert
          v-if="pendingVisits > 0"
          type="warning"
          variant="tonal"
          class="mb-4"
          density="compact"
        >
          {{ pendingVisits }} kunjungan menunggu sinkronisasi.
          <template #append>
            <v-btn size="small" variant="text" :loading="syncing" @click="sync">Kirim sekarang</v-btn>
          </template>
        </v-alert>

        <v-card elevation="2" class="pa-4">
          <v-card-title class="text-h5 text-primary mb-4">
            <v-icon start color="primary">mdi-account-hard-hat</v-icon>
            Input Kunjungan Petani
          </v-card-title>

          <v-form ref="form" @submit.prevent="submitVisit">
            <v-text-field v-model="form.farmerName" label="Nama Petani" variant="outlined" :rules="[required]"></v-text-field>
            <v-text-field v-model="form.cropType" label="Komoditas" variant="outlined" :rules="[required]"></v-text-field>
            <v-text-field
              v-model="form.landArea"
              label="Luas Lahan (Ha)"
              type="number"
              min="0"
              step="0.01"
              variant="outlined"
              :rules="[required, nonNegative]"
            ></v-text-field>
            <v-textarea v-model="form.notes" label="Catatan (opsional)" variant="outlined" rows="2" auto-grow></v-textarea>

            <v-row class="mb-2">
              <v-col cols="12">
                <v-btn color="secondary" block :loading="locating" @click="captureLocation">
                  <v-icon start>mdi-map-marker</v-icon>
                  Ambil Koordinat Geotag (GPS)
                </v-btn>
                <div v-if="form.latitude !== null" class="text-caption mt-2 text-center text-success">
                  Lokasi Terkunci: {{ form.latitude.toFixed(6) }}, {{ form.longitude.toFixed(6) }}
                </div>
              </v-col>
            </v-row>

            <v-file-input
              v-model="form.photo"
              label="Foto Kunjungan Lahan"
              prepend-icon="mdi-camera"
              variant="outlined"
              accept="image/*"
              capture="environment"
              show-size
              :rules="[required]"
            ></v-file-input>

            <v-btn type="submit" color="primary" size="large" block class="mt-4" :loading="submitting">
              <v-icon start>mdi-cloud-upload</v-icon>
              Kirim Kunjungan
            </v-btn>
          </v-form>
        </v-card>
      </v-col>
    </v-row>

    <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="4000">
      {{ snackbar.text }}
    </v-snackbar>
  </v-container>
</template>

<script>
function emptyForm() {
  return {
    farmerName: '',
    cropType: '',
    landArea: '',
    notes: '',
    latitude: null,
    longitude: null,
    photo: null
  }
}

export default {
  data() {
    return {
      form: emptyForm(),
      locating: false,
      submitting: false,
      syncing: false,
      snackbar: { show: false, text: '', color: 'success' }
    }
  },
  computed: {
    pendingVisits() {
      return this.$store.state.pendingVisits
    }
  },
  mounted() {
    this.$store.dispatch('refreshPending')
  },
  methods: {
    required(v) {
      const empty = v === null || v === undefined || v === '' || (Array.isArray(v) && v.length === 0)
      return !empty || 'Wajib diisi'
    },
    nonNegative(v) {
      return Number(v) >= 0 || 'Tidak boleh negatif'
    },
    notify(text, color = 'success') {
      this.snackbar = { show: true, text, color }
    },
    captureLocation() {
      if (!navigator.geolocation) {
        this.notify('Browser tidak mendukung Geolokasi HTML5.', 'error')
        return
      }
      this.locating = true
      navigator.geolocation.getCurrentPosition(
        pos => {
          this.form.latitude = pos.coords.latitude
          this.form.longitude = pos.coords.longitude
          this.locating = false
        },
        err => {
          console.error(err)
          this.locating = false
          this.notify('Gagal mengambil lokasi GPS. Pastikan izin lokasi diberikan.', 'error')
        },
        { enableHighAccuracy: true, timeout: 15000 }
      )
    },
    async submitVisit() {
      const { valid } = await this.$refs.form.validate()
      if (!valid) return
      if (this.form.latitude === null) {
        this.notify('Ambil koordinat GPS terlebih dahulu.', 'error')
        return
      }

      this.submitting = true
      try {
        // v-file-input bisa mengembalikan File atau array File tergantung versi Vuetify
        const photo = Array.isArray(this.form.photo) ? this.form.photo[0] : this.form.photo
        const result = await this.$store.dispatch('submitVisit', { ...this.form, photo })
        this.notify(result === 'sent'
          ? 'Kunjungan berhasil dikirim.'
          : 'Sedang offline. Kunjungan disimpan dan akan dikirim otomatis saat online.',
          result === 'sent' ? 'success' : 'warning')
        this.form = emptyForm()
        this.$refs.form.resetValidation()
      } catch (err) {
        console.error(err)
        this.notify(err.message || 'Gagal menyimpan kunjungan.', 'error')
      } finally {
        this.submitting = false
      }
    },
    async sync() {
      this.syncing = true
      try {
        await this.$store.dispatch('syncVisits')
        if (this.pendingVisits === 0) this.notify('Semua kunjungan sudah terkirim.')
      } finally {
        this.syncing = false
      }
    }
  }
}
</script>
