<template>
  <v-container>
    <v-row>
      <v-col cols="12" md="7">
        <v-card elevation="2" class="pa-2">
          <v-card-title class="d-flex align-center">
            <v-icon color="primary" class="mr-2">mdi-account-hard-hat</v-icon>
            Input Kunjungan Petani
          </v-card-title>
          <v-card-text>
            <v-alert v-if="pendingCount" type="warning" variant="tonal" density="compact" class="mb-4">
              {{ pendingCount }} kunjungan tersimpan di perangkat dan belum terkirim.
              <div v-for="e in queueErrors" :key="e" class="text-body-small">Gagal: {{ e }}</div>
              <template #append>
                <v-btn size="small" variant="flat" color="warning" :loading="syncing" :disabled="!online" @click="sync">
                  Kirim sekarang
                </v-btn>
              </template>
            </v-alert>

            <v-form ref="form" @submit.prevent="submitVisit">
              <v-row dense>
                <v-col cols="12" sm="6">
                  <v-text-field v-model.trim="form.farmerName" label="Nama Petani *" variant="outlined" :rules="[req]" />
                </v-col>
                <v-col cols="12" sm="6">
                  <v-text-field v-model.trim="form.farmerPhone" label="No. HP Petani" type="tel" variant="outlined" />
                </v-col>
                <v-col cols="12" sm="6">
                  <v-text-field v-model.trim="form.village" label="Desa / Kecamatan" variant="outlined" />
                </v-col>
                <v-col cols="12" sm="6">
                  <v-combobox v-model="form.cropType" :items="crops" label="Komoditas *" variant="outlined" :rules="[req]" />
                </v-col>
                <v-col cols="12" sm="6">
                  <v-text-field v-model="form.landArea" label="Luas Lahan (Ha) *" type="number" min="0" step="0.01"
                    variant="outlined" :rules="[req, positive]" />
                </v-col>
                <v-col cols="12" sm="6">
                  <v-text-field v-model="form.fertilizerUsage" label="Pemakaian Pupuk (kg/musim)" type="number" min="0"
                    variant="outlined" :rules="[positiveOrEmpty]" />
                </v-col>
                <v-col cols="12">
                  <v-text-field v-model.trim="form.competitorProduct" label="Produk Kompetitor yang Dipakai" variant="outlined" />
                </v-col>
                <v-col cols="12">
                  <v-textarea v-model.trim="form.notes" label="Catatan Kunjungan" rows="2" auto-grow variant="outlined" />
                </v-col>
              </v-row>

              <v-btn color="secondary" block variant="tonal" :loading="locating" @click="captureLocation" class="mb-1">
                <v-icon start>mdi-crosshairs-gps</v-icon>
                {{ form.latitude ? 'Perbarui Lokasi GPS' : 'Ambil Koordinat Geotag (GPS) *' }}
              </v-btn>
              <div class="text-body-small text-center mb-4" :class="form.latitude ? 'text-success' : 'text-medium-emphasis'">
                <template v-if="form.latitude">
                  Lokasi terkunci: {{ form.latitude.toFixed(6) }}, {{ form.longitude.toFixed(6) }}
                  (akurasi ±{{ Math.round(accuracy) }} m) ·
                  <a :href="mapUrl(form.latitude, form.longitude)" target="_blank" rel="noopener">lihat peta</a>
                </template>
                <template v-else>Lokasi wajib diambil di lahan petani.</template>
              </div>

              <v-file-input v-model="photoFile" label="Foto Kunjungan Lahan *" prepend-icon="mdi-camera"
                variant="outlined" accept="image/*" capture="environment" show-size :rules="[photoRule]"
                @update:model-value="onPhoto" />
              <v-img v-if="form.photoDataUrl" :src="form.photoDataUrl" max-height="220" class="mb-4 rounded" cover />

              <v-btn type="submit" color="primary" size="large" block :loading="saving">
                <v-icon start>{{ online ? 'mdi-cloud-upload' : 'mdi-content-save' }}</v-icon>
                {{ online ? 'Kirim Kunjungan' : 'Simpan Offline' }}
              </v-btn>
            </v-form>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" md="5">
        <v-card elevation="2">
          <v-card-title class="d-flex align-center">
            <v-icon class="mr-2">mdi-history</v-icon> Kunjungan Saya
            <v-spacer />
            <v-btn icon="mdi-refresh" variant="text" size="small" :loading="loadingHistory" @click="loadHistory" />
          </v-card-title>
          <v-card-subtitle>Hari ini: {{ todayCount }} kunjungan</v-card-subtitle>
          <v-list v-if="history.length" lines="two">
            <v-list-item v-for="v in history" :key="v.id"
              :title="v.farmer_name"
              :subtitle="`${v.crop_type || '-'} · ${v.land_area ?? '-'} Ha · ${waktu(v.visit_time)}`">
              <template #prepend><v-icon color="primary">mdi-map-marker-check</v-icon></template>
            </v-list-item>
          </v-list>
          <v-card-text v-else class="text-medium-emphasis">
            {{ loadingHistory ? 'Memuat...' : online ? 'Belum ada kunjungan.' : 'Riwayat tidak tersedia saat offline.' }}
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-snackbar v-model="snack.show" :color="snack.color" timeout="4000">{{ snack.text }}</v-snackbar>
  </v-container>
</template>

<script>
import { get } from '../services/api'
import { compressImage } from '../services/image'
import { waktu } from '../services/format'

const emptyForm = () => ({
  farmerName: '', farmerPhone: '', village: '', cropType: '', landArea: '', fertilizerUsage: '',
  competitorProduct: '', notes: '', latitude: null, longitude: null, photoDataUrl: null,
})

export default {
  data() {
    return {
      form: emptyForm(),
      photoFile: null,
      accuracy: 0,
      locating: false,
      saving: false,
      history: [],
      loadingHistory: false,
      crops: ['Padi', 'Jagung', 'Kedelai', 'Cabai', 'Bawang Merah', 'Tebu', 'Kelapa Sawit', 'Kopi', 'Sayuran'],
      snack: { show: false, text: '', color: 'success' },
    }
  },
  computed: {
    online() { return this.$store.state.online },
    syncing() { return this.$store.state.syncing },
    pendingCount() { return this.$store.getters.pendingCount },
    queueErrors() { return [...new Set(this.$store.state.offlineVisits.filter((v) => v.error).map((v) => v.error))] },
    todayCount() {
      const today = new Date().toDateString()
      return this.history.filter((v) => new Date(v.visit_time).toDateString() === today).length
    },
  },
  mounted() {
    this.loadHistory()
  },
  watch: {
    online(v) { if (v) this.loadHistory() },
  },
  methods: {
    waktu,
    req: (v) => (v !== null && v !== undefined && String(v).trim() !== '') || 'Wajib diisi',
    positive: (v) => Number(v) > 0 || 'Harus lebih dari 0',
    positiveOrEmpty: (v) => v === '' || v === null || Number(v) >= 0 || 'Tidak boleh negatif',
    photoRule() { return !!this.form.photoDataUrl || 'Foto wajib diambil' },
    mapUrl: (lat, lng) => `https://www.google.com/maps?q=${lat},${lng}`,
    notify(text, color = 'success') { this.snack = { show: true, text, color } },

    captureLocation() {
      if (!navigator.geolocation) {
        this.notify('Perangkat tidak mendukung GPS.', 'error')
        return
      }
      this.locating = true
      navigator.geolocation.getCurrentPosition(
        (pos) => {
          this.form.latitude = pos.coords.latitude
          this.form.longitude = pos.coords.longitude
          this.accuracy = pos.coords.accuracy
          this.locating = false
        },
        (err) => {
          this.locating = false
          this.notify(err.code === 1 ? 'Izin lokasi ditolak. Aktifkan izin lokasi di browser.' : 'Gagal mengambil lokasi GPS. Coba lagi di area terbuka.', 'error')
        },
        { enableHighAccuracy: true, timeout: 20000, maximumAge: 0 }
      )
    },

    async onPhoto(files) {
      const file = Array.isArray(files) ? files[0] : files
      if (!file) { this.form.photoDataUrl = null; return }
      try {
        this.form.photoDataUrl = await compressImage(file)
      } catch (err) {
        this.form.photoDataUrl = null
        this.notify(err.message, 'error')
      }
    },

    async submitVisit() {
      const { valid } = await this.$refs.form.validate()
      if (!valid) return
      if (!this.form.latitude) {
        this.notify('Ambil koordinat GPS terlebih dahulu.', 'error')
        return
      }
      this.saving = true
      try {
        const visit = {
          ...this.form,
          clientRef: (crypto.randomUUID && crypto.randomUUID()) || `${Date.now()}-${Math.random().toString(36).slice(2)}`,
          createdAt: new Date().toISOString(),
        }
        const { queued } = await this.$store.dispatch('submitVisit', visit)
        this.notify(queued ? 'Tersimpan di perangkat. Akan dikirim otomatis saat online.' : 'Kunjungan berhasil dikirim.', queued ? 'warning' : 'success')
        this.form = emptyForm()
        this.photoFile = null
        await this.$nextTick()
        this.$refs.form.resetValidation()
        if (!queued) this.loadHistory()
      } catch (err) {
        this.notify(err.message, 'error')
      } finally {
        this.saving = false
      }
    },

    async sync() {
      const res = await this.$store.dispatch('syncOfflineVisits')
      if (res.synced) this.notify(`${res.synced} kunjungan berhasil dikirim.`)
      if (res.failed) this.notify(`${res.failed} kunjungan gagal dikirim, akan dicoba lagi.`, 'warning')
      this.loadHistory()
    },

    async loadHistory() {
      if (!this.online) return
      this.loadingHistory = true
      try {
        this.history = await get('/sales/visits')
      } catch (err) {
        this.notify(err.message, 'error')
      } finally {
        this.loadingHistory = false
      }
    },
  },
}
</script>
