<template>
  <div>
    <div class="page-header mb-6">
      <h1 class="text-h4 font-weight-bold">Input Kunjungan Petani</h1>
      <p class="text-medium-emphasis text-body-2 mt-1">Rekam data petani, foto lahan, dan koordinat GPS</p>
    </div>

    <!-- Offline Banner -->
    <v-alert v-if="!isOnline" type="warning" icon="mdi-wifi-off" variant="tonal" rounded="lg" class="mb-4" density="compact">
      Mode Offline – Data akan disinkronkan otomatis saat koneksi tersedia.
    </v-alert>

    <v-alert v-if="submitSuccess" type="success" icon="mdi-check-circle" variant="tonal" rounded="lg" class="mb-4" density="compact" closable @click:close="submitSuccess = false">
      {{ isOnline ? 'Kunjungan berhasil disimpan ke server!' : 'Kunjungan disimpan sebagai draft offline.' }}
    </v-alert>

    <v-row>
      <!-- Form -->
      <v-col cols="12" md="7">
        <v-card rounded="xl" elevation="0" border>
          <v-card-title class="pa-5 pb-2 text-h6 font-weight-bold">
            <v-icon start color="primary">mdi-account-hard-hat</v-icon>
            Data Petani & Lahan
          </v-card-title>
          <v-card-text class="pa-5">
            <v-form @submit.prevent="submitVisit" ref="visitForm">
              <v-text-field
                v-model="form.farmerName"
                label="Nama Petani"
                prepend-inner-icon="mdi-account"
                variant="outlined"
                density="comfortable"
                color="primary"
                rounded="lg"
                :rules="[v => !!v || 'Nama petani wajib diisi']"
                class="mb-3"
              ></v-text-field>

              <v-text-field
                v-model="form.cropType"
                label="Komoditas Utama"
                prepend-inner-icon="mdi-sprout"
                variant="outlined"
                density="comfortable"
                color="primary"
                rounded="lg"
                :rules="[v => !!v || 'Komoditas wajib diisi']"
                class="mb-3"
              ></v-text-field>

              <v-row>
                <v-col cols="6">
                  <v-text-field
                    v-model="form.landArea"
                    label="Luas Lahan (Ha)"
                    type="number"
                    prepend-inner-icon="mdi-ruler-square"
                    variant="outlined"
                    density="comfortable"
                    color="primary"
                    rounded="lg"
                  ></v-text-field>
                </v-col>
                <v-col cols="6">
                  <v-text-field
                    v-model="form.fertilizerUsage"
                    label="Pupuk (Kg/Ha)"
                    type="number"
                    prepend-inner-icon="mdi-flask"
                    variant="outlined"
                    density="comfortable"
                    color="primary"
                    rounded="lg"
                  ></v-text-field>
                </v-col>
              </v-row>

              <v-textarea
                v-model="form.notes"
                label="Catatan Kunjungan"
                prepend-inner-icon="mdi-note-text"
                variant="outlined"
                density="comfortable"
                color="primary"
                rounded="lg"
                rows="3"
                class="mb-3"
              ></v-textarea>

              <!-- Geolocation -->
              <v-card variant="tonal" color="primary" rounded="lg" class="mb-4 pa-3">
                <div class="d-flex align-center justify-space-between flex-wrap gap-2">
                  <div>
                    <div class="font-weight-bold text-caption text-uppercase text-primary">Koordinat GPS</div>
                    <div v-if="form.latitude" class="text-caption text-success mt-1">
                      <v-icon size="14" color="success">mdi-check-circle</v-icon>
                      Lat: {{ form.latitude.toFixed(6) }}, Lng: {{ form.longitude.toFixed(6) }}
                    </div>
                    <div v-else class="text-caption text-grey mt-1">Belum diambil</div>
                  </div>
                  <v-btn color="primary" size="small" rounded="pill" :loading="gpsLoading" @click="captureLocation" prepend-icon="mdi-map-marker">
                    Ambil GPS
                  </v-btn>
                </div>
              </v-card>

              <!-- Photo Upload -->
              <v-file-input
                v-model="form.photo"
                label="Foto Kunjungan Lahan"
                prepend-inner-icon="mdi-camera"
                prepend-icon=""
                variant="outlined"
                density="comfortable"
                color="primary"
                rounded="lg"
                accept="image/*"
                capture="environment"
                show-size
                class="mb-4"
              ></v-file-input>

              <!-- Preview foto -->
              <div v-if="photoPreview" class="mb-4">
                <v-img :src="photoPreview" height="180" cover rounded="lg" class="border"></v-img>
              </div>

              <v-btn type="submit" color="primary" block size="x-large" rounded="pill" :loading="submitting" elevation="2" prepend-icon="mdi-cloud-upload">
                {{ isOnline ? 'Kirim ke Server' : 'Simpan Draft Offline' }}
              </v-btn>
            </v-form>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Sidebar: riwayat offline -->
      <v-col cols="12" md="5">
        <v-card rounded="xl" elevation="0" border class="mb-4">
          <v-card-title class="pa-5 pb-2 text-h6 font-weight-bold">
            <v-icon start color="warning">mdi-clock-outline</v-icon>
            Draft Offline ({{ offlineVisits.length }})
          </v-card-title>
          <v-card-text class="pa-5 pt-2">
            <div v-if="offlineVisits.length === 0" class="text-center text-grey text-caption py-4">
              <v-icon size="40" color="grey-lighten-2">mdi-inbox-outline</v-icon>
              <div class="mt-2">Tidak ada data offline</div>
            </div>
            <v-list v-else lines="two" class="pa-0">
              <v-list-item
                v-for="(visit, i) in offlineVisits"
                :key="i"
                :title="visit.farmerName || 'Petani #' + (i+1)"
                :subtitle="visit.cropType"
                prepend-icon="mdi-account-outline"
                rounded="lg"
                class="mb-1 bg-grey-lighten-4"
              >
                <template v-slot:append>
                  <v-chip color="warning" size="x-small" label>Offline</v-chip>
                </template>
              </v-list-item>
            </v-list>
            <v-btn v-if="offlineVisits.length > 0 && isOnline" class="mt-3" color="success" block variant="tonal" rounded="pill" size="small" @click="syncOfflineData" prepend-icon="mdi-sync">
              Sync Sekarang
            </v-btn>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script>
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

export default {
  data() {
    return {
      isOnline: navigator.onLine,
      gpsLoading: false,
      submitting: false,
      submitSuccess: false,
      photoPreview: null,
      form: {
        farmerName: '',
        cropType: '',
        landArea: '',
        fertilizerUsage: '',
        notes: '',
        latitude: null,
        longitude: null,
        photo: null
      }
    }
  },
  computed: {
    offlineVisits() {
      return this.$store.state.offlineVisits
    }
  },
  watch: {
    'form.photo'(file) {
      if (file && file[0]) {
        const reader = new FileReader()
        reader.onload = (e) => { this.photoPreview = e.target.result }
        reader.readAsDataURL(file[0])
      } else {
        this.photoPreview = null
      }
    }
  },
  mounted() {
    window.addEventListener('online', () => { this.isOnline = true })
    window.addEventListener('offline', () => { this.isOnline = false })
  },
  beforeUnmount() {
    window.removeEventListener('online', () => {})
    window.removeEventListener('offline', () => {})
  },
  methods: {
    captureLocation() {
      if (!navigator.geolocation) {
        alert('Browser tidak mendukung Geolokasi HTML5.')
        return
      }
      this.gpsLoading = true
      navigator.geolocation.getCurrentPosition(
        (pos) => {
          this.form.latitude = pos.coords.latitude
          this.form.longitude = pos.coords.longitude
          this.gpsLoading = false
        },
        (err) => {
          console.error(err)
          alert('Gagal mengambil lokasi GPS. Pastikan izin lokasi diberikan.')
          this.gpsLoading = false
        }
      )
    },
    async submitVisit() {
      const { valid } = await this.$refs.visitForm.validate()
      if (!valid) return

      this.submitting = true
      const token = this.$store.getters.token

      if (!this.isOnline) {
        this.$store.dispatch('saveVisitOffline', { ...this.form })
        this.submitSuccess = true
        this.resetForm()
        this.submitting = false
        return
      }

      try {
        const formData = new FormData()
        formData.append('farmer_name', this.form.farmerName)
        formData.append('crop_type', this.form.cropType)
        formData.append('land_area', this.form.landArea)
        formData.append('fertilizer_usage', this.form.fertilizerUsage)
        formData.append('notes', this.form.notes)
        if (this.form.latitude) {
          formData.append('latitude', this.form.latitude)
          formData.append('longitude', this.form.longitude)
        }
        if (this.form.photo && this.form.photo[0]) {
          formData.append('photo', this.form.photo[0])
        }

        const response = await fetch(`${API_URL}/sales/visits`, {
          method: 'POST',
          headers: { 'Authorization': `Bearer ${token}` },
          body: formData
        })

        if (!response.ok) throw new Error('Gagal kirim data ke server')
        this.submitSuccess = true
        this.resetForm()
      } catch (err) {
        console.error(err)
        // Fallback to offline
        this.$store.dispatch('saveVisitOffline', { ...this.form })
        this.submitSuccess = true
        this.resetForm()
      } finally {
        this.submitting = false
      }
    },
    async syncOfflineData() {
      // Logic to sync all offline visits to the server
      alert('Fitur sync akan segera tersedia.')
    },
    resetForm() {
      this.form = { farmerName: '', cropType: '', landArea: '', fertilizerUsage: '', notes: '', latitude: null, longitude: null, photo: null }
      this.photoPreview = null
      this.$refs.visitForm.reset()
    }
  }
}
</script>
