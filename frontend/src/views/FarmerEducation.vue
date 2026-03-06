<template>
  <div>
    <div class="page-header mb-6">
      <h1 class="text-h4 font-weight-bold">Edukasi & Info Petani</h1>
      <p class="text-medium-emphasis text-body-2 mt-1">Video pembelajaran, harga komoditas, dan prakiraan cuaca terkini</p>
    </div>

    <v-row>
      <!-- Konten Edukasi -->
      <v-col cols="12" md="8">
        <!-- Filter -->
        <div class="d-flex gap-2 mb-4 flex-wrap">
          <v-chip
            v-for="cat in categories"
            :key="cat.value"
            :color="activeCategory === cat.value ? 'primary' : 'default'"
            variant="tonal"
            rounded="pill"
            class="cursor-pointer"
            @click="filterContent(cat.value)"
          >
            {{ cat.label }}
          </v-chip>
        </div>

        <!-- Loading state -->
        <v-skeleton-loader v-if="loadingContent" type="card@2"></v-skeleton-loader>

        <!-- Data dari API -->
        <template v-else-if="contents.length > 0">
          <!-- Featured video: first item -->
          <v-card rounded="xl" elevation="0" border class="mb-4">
            <div class="video-wrapper" @click="openVideo(contents[0])">
              <v-img
                :src="contents[0].thumbnail || 'https://images.unsplash.com/photo-1592982537447-6f2a6a0c5c1b?auto=format&fit=crop&q=80&w=900'"
                height="260px"
                cover
              >
                <div class="overlay d-flex align-center justify-center">
                  <v-btn icon size="x-large" color="white" variant="flat" class="elevation-6">
                    <v-icon color="primary" size="44">mdi-play</v-icon>
                  </v-btn>
                </div>
              </v-img>
            </div>
            <v-card-text class="pa-5">
              <v-chip :color="getCategoryColor(contents[0].category)" size="small" label class="mb-2">
                {{ contents[0].category || 'Edukasi' }}
              </v-chip>
              <h2 class="text-h6 font-weight-bold mb-1">{{ contents[0].title }}</h2>
              <p class="text-body-2 text-medium-emphasis mb-0">{{ contents[0].description }}</p>
            </v-card-text>
          </v-card>

          <!-- More content cards -->
          <v-row>
            <v-col cols="12" sm="6" v-for="(content, i) in contents.slice(1)" :key="i">
              <v-card rounded="xl" elevation="0" border class="cursor-pointer hover-card" @click="openVideo(content)">
                <v-img
                  :src="content.thumbnail || 'https://images.unsplash.com/photo-1464226184884-fa280b87c399?auto=format&fit=crop&q=80&w=400'"
                  height="140"
                  cover
                ></v-img>
                <v-card-text class="pa-4">
                  <v-chip :color="getCategoryColor(content.category)" size="x-small" label class="mb-2">{{ content.category || 'Edukasi' }}</v-chip>
                  <div class="text-subtitle-2 font-weight-bold">{{ content.title }}</div>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>
        </template>

        <!-- Fallback: no content from API -->
        <v-card v-else rounded="xl" elevation="0" border class="mb-4">
          <v-img
            src="https://images.unsplash.com/photo-1592982537447-6f2a6a0c5c1b?auto=format&fit=crop&q=80&w=900"
            height="260px"
            cover
          >
            <div class="overlay d-flex align-center justify-center">
              <v-btn icon size="x-large" color="white" variant="flat">
                <v-icon color="primary" size="44">mdi-play</v-icon>
              </v-btn>
            </div>
          </v-img>
          <v-card-text class="pa-5">
            <v-chip color="primary" size="small" label class="mb-2">Video Edukasi</v-chip>
            <h2 class="text-h6 font-weight-bold mb-1">Cara Aplikasi Pupuk Organik yang Efisien</h2>
            <p class="text-body-2 text-medium-emphasis mb-0">Panduan praktis menekan biaya pupuk anorganik untuk memperkaya unsur hara lahan Anda.</p>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Sidebar: Cuaca & Harga -->
      <v-col cols="12" md="4">
        <!-- Widget Cuaca (Live dari API) -->
        <v-card rounded="xl" elevation="0" border class="mb-4 overflow-hidden">
          <div class="weather-header pa-5 text-white">
            <div class="text-overline opacity-80">BMKG – Prakiraan Cuaca</div>
            <v-skeleton-loader v-if="!weather" type="text@3" theme="dark"></v-skeleton-loader>
            <template v-else>
              <v-row align="center" no-gutters class="mt-2">
                <v-col>
                  <div class="text-h2 font-weight-bold">{{ weather.temperature || 28 }}°C</div>
                  <div class="text-body-2 mt-1 opacity-90">{{ weather.condition || 'Cerah Berawan' }}</div>
                  <div class="text-caption opacity-70 mt-1">{{ today }}</div>
                </v-col>
                <v-col cols="auto">
                  <v-icon size="80" color="white">mdi-weather-partly-cloudy</v-icon>
                </v-col>
              </v-row>
              <v-divider class="my-3 opacity-30"></v-divider>
              <v-row dense>
                <v-col cols="6" class="text-center">
                  <v-icon size="18">mdi-water-percent</v-icon>
                  <div class="text-caption mt-1">{{ weather.humidity || 72 }}% Lembab</div>
                </v-col>
                <v-col cols="6" class="text-center">
                  <v-icon size="18">mdi-weather-rainy</v-icon>
                  <div class="text-caption mt-1">{{ weather.rain_chance || 20 }}% Hujan</div>
                </v-col>
              </v-row>
            </template>
          </div>
        </v-card>

        <!-- Harga Komoditas (Live dari API) -->
        <v-card rounded="xl" elevation="0" border>
          <v-card-title class="pa-5 pb-2 text-h6 font-weight-bold">
            <v-icon start color="success">mdi-cash</v-icon>
            Harga Komoditas
          </v-card-title>
          <v-skeleton-loader v-if="!commodityPrices.length" type="list-item-two-line@4" class="px-4"></v-skeleton-loader>
          <v-list v-else lines="one" class="pa-0">
            <v-list-item
              v-for="(item, i) in commodityPrices"
              :key="i"
              :title="item.name"
              :subtitle="'Rp ' + item.price_per_kg?.toLocaleString('id-ID') + ' / kg'"
            >
              <template v-slot:prepend>
                <v-icon
                  :icon="item.trend === 'up' ? 'mdi-trending-up' : 'mdi-trending-down'"
                  :color="item.trend === 'up' ? 'success' : 'error'"
                ></v-icon>
              </template>
              <template v-slot:append>
                <v-chip :color="item.trend === 'up' ? 'success' : 'error'" size="x-small" label>
                  {{ item.change_pct > 0 ? '+' : '' }}{{ item.change_pct }}%
                </v-chip>
              </template>
            </v-list-item>
          </v-list>
        </v-card>
      </v-col>
    </v-row>

    <!-- Video Dialog -->
    <v-dialog v-model="videoDialog" max-width="700">
      <v-card rounded="xl">
        <v-card-title class="d-flex align-center justify-space-between pa-4">
          {{ selectedContent?.title }}
          <v-btn icon size="small" @click="videoDialog = false"><v-icon>mdi-close</v-icon></v-btn>
        </v-card-title>
        <div v-if="selectedContent?.video_url" class="pa-4 pt-0">
          <iframe :src="selectedContent.video_url" width="100%" height="360" frameborder="0" allowfullscreen class="rounded-lg"></iframe>
        </div>
        <v-card-text v-else class="pa-5 pt-0">
          <p class="text-body-2 text-medium-emphasis">{{ selectedContent?.description }}</p>
          <v-alert type="info" variant="tonal" class="mt-3" density="compact" rounded="lg">
            Link video belum tersedia. Hubungi admin untuk upload konten.
          </v-alert>
        </v-card-text>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
export default {
  data() {
    return {
      activeCategory: '',
      loadingContent: false,
      videoDialog: false,
      selectedContent: null,
      today: new Date().toLocaleDateString('id-ID', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' }),
      categories: [
        { label: 'Semua', value: '' },
        { label: 'Pertanian', value: 'education' },
        { label: 'Harga Pasar', value: 'market_price' },
        { label: 'Teknologi', value: 'technology' },
      ]
    }
  },
  computed: {
    contents() { return this.$store.state.contents },
    weather() { return this.$store.state.weather },
    commodityPrices() { return this.$store.state.commodityPrices },
  },
  mounted() {
    this.loadData()
  },
  methods: {
    async loadData() {
      this.loadingContent = true
      await Promise.all([
        this.$store.dispatch('fetchContents', this.activeCategory),
        this.$store.dispatch('fetchWeather', 'jawa-timur'),
        this.$store.dispatch('fetchCommodityPrices'),
      ])
      this.loadingContent = false
    },
    async filterContent(category) {
      this.activeCategory = category
      this.loadingContent = true
      await this.$store.dispatch('fetchContents', category)
      this.loadingContent = false
    },
    openVideo(content) {
      this.selectedContent = content
      this.videoDialog = true
    },
    getCategoryColor(cat) {
      const map = { education: 'primary', market_price: 'success', technology: 'info' }
      return map[cat] || 'grey'
    }
  }
}
</script>

<style scoped>
.weather-header {
  background: linear-gradient(135deg, #1565C0 0%, #1E88E5 100%);
}
.hover-card {
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}
.hover-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0,0,0,0.12) !important;
}
.video-wrapper { cursor: pointer; }
.overlay {
  width: 100%;
  height: 100%;
  background: rgba(0,0,0,0.35);
}
</style>
