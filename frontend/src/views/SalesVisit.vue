<template>
  <v-container>
    <v-row>
      <v-col cols="12" md="8" offset-md="2">
        <v-card elevation="2" class="pa-4">
          <v-card-title class="text-h5 primary--text mb-4">
            <v-icon left color="primary">mdi-account-hard-hat</v-icon>
            Input Kunjungan Petani
          </v-card-title>
          
          <v-form @submit.prevent="submitVisit">
            <v-text-field v-model="form.farmerName" label="Nama Petani" variant="outlined" required></v-text-field>
            <v-text-field v-model="form.cropType" label="Komoditas" variant="outlined" required></v-text-field>
            <v-text-field v-model="form.landArea" label="Luas Lahan (Ha)" type="number" variant="outlined" required></v-text-field>
            
            <v-row class="mb-2">
              <v-col cols="12">
                <v-btn color="secondary" block @click="captureLocation">
                  <v-icon left>mdi-map-marker</v-icon>
                  Ambil Koordinat Geotag (GPS)
                </v-btn>
                <div v-if="form.latitude" class="text-caption mt-2 text-center text-success">
                  Lokasi Terkunci: {{ form.latitude }}, {{ form.longitude }}
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
              required
            ></v-file-input>

            <v-btn type="submit" color="primary" size="large" block class="mt-4">
              <v-icon left>mdi-cloud-upload</v-icon>
              Kirim / Simpan Draft (Offline)
            </v-btn>
          </v-form>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  data() {
    return {
      form: {
        farmerName: '',
        cropType: '',
        landArea: '',
        latitude: null,
        longitude: null,
        photo: null
      }
    }
  },
  methods: {
    captureLocation() {
      if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(pos => {
          this.form.latitude = pos.coords.latitude
          this.form.longitude = pos.coords.longitude
          // For UX feedback
          alert(`Lokasi didapatkan: Lat ${this.form.latitude}, Lng ${this.form.longitude}`)
        }, err => {
          console.error(err)
          alert("Gagal mengambil lokasi GPS. Pastikan izin diberikan.")
        })
      } else {
        alert("Browser tidak mendukung Geolokasi HTML5.")
      }
    },
    submitVisit() {
      alert("Menyimpan data...")
      console.log("Visit Data:", this.form)
      // Save logic (Vuex / API)
      this.$store.dispatch('saveVisitOffline', this.form)
    }
  }
}
</script>
