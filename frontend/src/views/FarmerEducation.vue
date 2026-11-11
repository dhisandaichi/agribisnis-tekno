<template>
  <v-container>
    <div class="text-headline-small mb-4">Edukasi & Info Petani</div>

    <v-row>
      <v-col cols="12" md="8">
        <v-alert v-if="!loading && !contents.length" type="info" variant="tonal">Belum ada konten edukasi.</v-alert>
        <v-card v-for="c in contents" :key="c.id" elevation="2" class="mb-4">
          <v-img :src="c.thumbnail_url || defaultThumb" height="220" cover class="align-center justify-center text-center">
            <v-btn v-if="c.video_url" icon size="x-large" variant="flat" color="white" :href="c.video_url" target="_blank" rel="noopener">
              <v-icon color="black" size="x-large">mdi-play</v-icon>
            </v-btn>
          </v-img>
          <v-card-title class="text-wrap">{{ c.title }}</v-card-title>
          <v-card-subtitle v-if="c.duration_minutes">Edukasi Tani · {{ c.duration_minutes }} menit</v-card-subtitle>
          <v-card-text class="text-pre-wrap">{{ c.description }}</v-card-text>
          <v-card-actions v-if="c.video_url">
            <v-btn color="primary" :href="c.video_url" target="_blank" rel="noopener" prepend-icon="mdi-play-circle">Tonton</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>

      <v-col cols="12" md="4">
        <v-card elevation="2" class="mb-4">
          <v-card-title class="bg-primary text-white d-flex align-center">
            <v-icon class="mr-2">mdi-weather-partly-cloudy</v-icon> Prakiraan Cuaca
          </v-card-title>
          <v-card-text class="py-4">
            <div v-if="weather.loading" class="text-center"><v-progress-circular indeterminate color="primary" /></div>
            <div v-else-if="weather.error" class="text-medium-emphasis">{{ weather.error }}</div>
            <template v-else-if="weather.now">
              <div class="d-flex align-center justify-center">
                <img v-if="weather.now.icon" :src="weather.now.icon" alt="" width="64" height="64" />
                <div class="text-display-small ml-2">{{ weather.now.temp }}°C</div>
              </div>
              <div class="text-center text-title-medium">{{ weather.now.desc }}</div>
              <div class="text-center text-body-small text-medium-emphasis">
                Kelembapan {{ weather.now.humidity }}% · Angin {{ weather.now.wind }} km/j
              </div>
              <div class="text-center text-body-small mt-1">{{ weather.location }}</div>
              <v-divider class="my-3" />
              <div class="d-flex justify-space-between text-center">
                <div v-for="f in weather.next" :key="f.time" class="text-body-small">
                  <div>{{ jam(f.time) }}</div>
                  <img v-if="f.icon" :src="f.icon" alt="" width="32" height="32" />
                  <div>{{ f.temp }}°</div>
                </div>
              </div>
              <div class="text-body-small text-medium-emphasis mt-2">Sumber: BMKG</div>
            </template>
          </v-card-text>
        </v-card>

        <v-card elevation="2">
          <v-card-title class="bg-secondary text-white d-flex align-center">
            <v-icon class="mr-2">mdi-cash</v-icon> Harga Komoditas
          </v-card-title>
          <v-list lines="two">
            <template v-for="(p, i) in prices" :key="p.id">
              <v-divider v-if="i" />
              <v-list-item :title="p.name" :subtitle="`${rupiah(p.price)} / ${p.unit} · ${tanggal(p.updated_at)}`">
                <template #append><v-icon :color="trendColor(p.trend)">{{ trendIcon(p.trend) }}</v-icon></template>
              </v-list-item>
            </template>
            <v-list-item v-if="!prices.length && !loading" subtitle="Belum ada data harga." />
          </v-list>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { get } from '../services/api'
import { rupiah, tanggal } from '../services/format'

export default {
  data() {
    return {
      loading: true,
      contents: [],
      prices: [],
      defaultThumb: 'https://images.unsplash.com/photo-1592982537447-6f2a6a0c5c1b?auto=format&fit=crop&q=80&w=800',
      weather: { loading: true, error: null, now: null, next: [], location: '' },
    }
  },
  async mounted() {
    this.loadWeather()
    const [c, p] = await Promise.allSettled([get('/content/education'), get('/content/prices')])
    this.contents = c.status === 'fulfilled' ? c.value : []
    this.prices = p.status === 'fulfilled' ? p.value : []
    this.loading = false
  },
  methods: {
    rupiah, tanggal,
    trendIcon: (t) => ({ up: 'mdi-trending-up', down: 'mdi-trending-down' }[t] || 'mdi-trending-neutral'),
    trendColor: (t) => ({ up: 'success', down: 'error' }[t] || 'grey'),
    jam: (t) => (t ? t.slice(11, 16) : ''),
    async loadWeather() {
      try {
        const adm4 = import.meta.env.VITE_BMKG_ADM4
        const json = await get(`/weather${adm4 ? `?adm4=${encodeURIComponent(adm4)}` : ''}`)
        const now = Date.now()
        const list = json.forecasts.filter((f) => new Date(f.time.replace(' ', 'T')).getTime() >= now - 3 * 3600e3)
        this.weather = { loading: false, error: null, now: list[0], next: list.slice(1, 5), location: json.location }
        if (!list.length) this.weather.error = 'Data prakiraan belum tersedia.'
      } catch {
        this.weather = { loading: false, error: 'Prakiraan cuaca sedang tidak tersedia.', now: null, next: [], location: '' }
      }
    },
  },
}
</script>
