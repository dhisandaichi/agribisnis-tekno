<template>
  <div>
    <div class="page-header mb-6">
      <h1 class="text-h4 font-weight-bold">Dashboard Supervisor</h1>
      <p class="text-medium-emphasis text-body-2 mt-1">Monitor aktivitas sales, distributor, dan kinerja lapangan</p>
    </div>

    <!-- KPI Cards -->
    <v-row class="mb-4">
      <v-col cols="12" sm="6" md="3">
        <v-card rounded="xl" elevation="0" border>
          <v-card-text class="pa-5">
            <div class="d-flex align-center justify-space-between mb-3">
              <v-avatar color="primary" size="44" rounded="lg">
                <v-icon icon="mdi-account-group" color="white" size="22"></v-icon>
              </v-avatar>
              <v-skeleton-loader v-if="loading" type="chip" width="60"></v-skeleton-loader>
            </div>
            <div class="text-h4 font-weight-bold">{{ totalFarmers }}</div>
            <div class="text-caption text-medium-emphasis mt-1">Total Petani Aktif</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <v-card rounded="xl" elevation="0" border>
          <v-card-text class="pa-5">
            <div class="d-flex align-center justify-space-between mb-3">
              <v-avatar color="info" size="44" rounded="lg">
                <v-icon icon="mdi-map-marker-path" color="white" size="22"></v-icon>
              </v-avatar>
            </div>
            <div class="text-h4 font-weight-bold">{{ visits.length }}</div>
            <div class="text-caption text-medium-emphasis mt-1">Total Kunjungan</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <v-card rounded="xl" elevation="0" border>
          <v-card-text class="pa-5">
            <div class="d-flex align-center justify-space-between mb-3">
              <v-avatar color="success" size="44" rounded="lg">
                <v-icon icon="mdi-storefront" color="white" size="22"></v-icon>
              </v-avatar>
            </div>
            <div class="text-h4 font-weight-bold">{{ distributors.length }}</div>
            <div class="text-caption text-medium-emphasis mt-1">Distributor / Kios</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <v-card rounded="xl" elevation="0" border>
          <v-card-text class="pa-5">
            <div class="d-flex align-center justify-space-between mb-3">
              <v-avatar color="error" size="44" rounded="lg">
                <v-icon icon="mdi-alert-circle" color="white" size="22"></v-icon>
              </v-avatar>
            </div>
            <div class="text-h4 font-weight-bold">{{ overdueCount }}</div>
            <div class="text-caption text-medium-emphasis mt-1">Hutang Jatuh Tempo</div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row>
      <!-- Tabel Kunjungan -->
      <v-col cols="12" md="8">
        <v-card rounded="xl" elevation="0" border>
          <v-card-title class="pa-5 pb-2 d-flex align-center justify-space-between">
            <span class="text-h6 font-weight-bold">
              <v-icon start color="primary">mdi-map-marker-path</v-icon>
              Kunjungan Sales (Live)
            </span>
            <v-btn size="small" color="primary" variant="tonal" rounded="pill" prepend-icon="mdi-refresh" @click="loadData" :loading="loading">
              Refresh
            </v-btn>
          </v-card-title>
          <v-card-text class="pa-5 pt-0">
            <v-skeleton-loader v-if="loading" type="table-tbody"></v-skeleton-loader>
            <v-data-table
              v-else
              :headers="visitHeaders"
              :items="visits"
              density="comfortable"
              :items-per-page="5"
              no-data-text="Belum ada data kunjungan"
            >
              <template v-slot:item.visit_time="{ item }">
                {{ formatDate(item.visit_time) }}
              </template>
              <template v-slot:item.photo_path="{ item }">
                <v-chip v-if="item.photo_path" color="success" size="x-small" label>Ada Foto</v-chip>
                <v-chip v-else color="grey" size="x-small" label>Tanpa Foto</v-chip>
              </template>
            </v-data-table>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Sidebar: Distributor -->
      <v-col cols="12" md="4">
        <v-card rounded="xl" elevation="0" border>
          <v-card-title class="pa-5 pb-2 text-h6 font-weight-bold">
            <v-icon start color="primary">mdi-storefront</v-icon>
            Distributor / Kios
          </v-card-title>
          <v-skeleton-loader v-if="loading" type="list-item-two-line@3" class="px-4"></v-skeleton-loader>
          <v-list v-else lines="two" class="pa-4 pt-0">
            <v-list-item
              v-for="(d, i) in distributors.slice(0, 5)"
              :key="i"
              :title="d.name"
              :subtitle="d.phone || d.address || '-'"
              prepend-icon="mdi-store-outline"
              rounded="lg"
              class="mb-1 bg-grey-lighten-4"
            ></v-list-item>
            <div v-if="distributors.length === 0" class="text-center text-caption text-grey py-3">
              Belum ada data distributor
            </div>
          </v-list>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script>
export default {
  data() {
    return {
      overdueCount: 0,
      totalFarmers: 0,
      visitHeaders: [
        { title: 'Petani ID', key: 'farmer_id' },
        { title: 'Sales ID', key: 'sales_id' },
        { title: 'Catatan', key: 'notes' },
        { title: 'Foto', key: 'photo_path', align: 'center' },
        { title: 'Waktu', key: 'visit_time' },
      ]
    }
  },
  computed: {
    visits() { return this.$store.state.visits },
    distributors() { return this.$store.state.distributors },
    loading() { return this.$store.state.loading },
  },
  mounted() {
    this.loadData()
  },
  methods: {
    async loadData() {
      await Promise.all([
        this.$store.dispatch('fetchVisits'),
        this.$store.dispatch('fetchDistributors'),
      ])
      // Hitung total petani unik dari kunjungan
      const uniqueFarmers = new Set(this.visits.map(v => v.farmer_id))
      this.totalFarmers = uniqueFarmers.size
    },
    formatDate(dateStr) {
      if (!dateStr) return '-'
      return new Date(dateStr).toLocaleString('id-ID', {
        day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit'
      })
    }
  }
}
</script>
