<template>
  <v-container>
    <div class="text-h4 mb-4 primary--text">Edukasi & Info Petani</div>
    
    <v-row>
      <!-- Edukasi Video Content -->
      <v-col cols="12" md="8">
        <v-card elevation="3" class="mb-4">
          <v-img
            src="https://images.unsplash.com/photo-1592982537447-6f2a6a0c5c1b?auto=format&fit=crop&q=80&w=800"
            height="300px"
            cover
            class="align-center justify-center text-center"
          >
            <!-- Mock play button overlay -->
            <v-btn icon color="white" size="x-large" variant="flat">
              <v-icon color="black" size="x-large">mdi-play</v-icon>
            </v-btn>
          </v-img>
          <v-card-title>Cara Aplikasi Pupuk Organik yang Efisien</v-card-title>
          <v-card-subtitle>Streaming Edukasi Tani - 3 Menit</v-card-subtitle>
          <v-card-text>
            Panduan aplikatif menekan biaya pupuk anorganik untuk memperkaya unsur hara lahan Anda.
          </v-card-text>
        </v-card>
      </v-col>

      <!-- BMKG Weather & Market Prices -->
      <v-col cols="12" md="4">
        <v-card elevation="2" class="mb-4">
          <v-card-title class="bg-primary text-white">
            <v-icon left>mdi-weather-cloudy</v-icon>
            Prakiraan Cuaca (BMKG)
          </v-card-title>
          <v-card-text class="py-4 text-center">
            <v-icon size="64" color="info" class="mb-2">mdi-weather-partly-cloudy</v-icon>
            <div class="text-h4">28°C</div>
            <div class="text-subtitle-1">Cerah Berawan</div>
            <div class="text-caption text-grey">Jawa Timur - Hari Ini</div>
          </v-card-text>
        </v-card>

        <v-card elevation="2">
          <v-card-title class="bg-secondary text-white">
            <v-icon left>mdi-cash</v-icon>
            Harga Komoditas (Realtime)
          </v-card-title>
          <v-list lines="one">
            <v-list-item title="Gabah Kering (GKP)" subtitle="Rp 6.400 / kg" append-icon="mdi-trending-up"></v-list-item>
            <v-divider></v-divider>
            <v-list-item title="Bawang Merah Besar" subtitle="Rp 24.500 / kg" append-icon="mdi-trending-down"></v-list-item>
          </v-list>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  // Logic for pulling BMKG api goes here
}
</script>
