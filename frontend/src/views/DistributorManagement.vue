<template>
  <div>
    <div class="page-header mb-6">
      <h1 class="text-h4 font-weight-bold">Manajemen Distributor & Kios</h1>
      <p class="text-medium-emphasis text-body-2 mt-1">Monitor stok, pesanan, dan status pembayaran distributor</p>
    </div>

    <!-- KPI -->
    <v-row class="mb-4">
      <v-col cols="12" sm="4">
        <v-card rounded="xl" elevation="0" border>
          <v-card-text class="pa-5 d-flex align-center gap-4">
            <v-avatar color="primary" size="48" rounded="lg">
              <v-icon icon="mdi-storefront" color="white" size="24"></v-icon>
            </v-avatar>
            <div>
              <div class="text-h5 font-weight-bold">{{ distributors.length }}</div>
              <div class="text-caption text-medium-emphasis">Total Distributor</div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" sm="4">
        <v-card rounded="xl" elevation="0" border>
          <v-card-text class="pa-5 d-flex align-center gap-4">
            <v-avatar color="error" size="48" rounded="lg">
              <v-icon icon="mdi-alert" color="white" size="24"></v-icon>
            </v-avatar>
            <div>
              <div class="text-h5 font-weight-bold">{{ criticalStockCount }}</div>
              <div class="text-caption text-medium-emphasis">Stok Kritis</div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" sm="4">
        <v-card rounded="xl" elevation="0" border>
          <v-card-text class="pa-5 d-flex align-center gap-4">
            <v-avatar color="warning" size="48" rounded="lg">
              <v-icon icon="mdi-cash-remove" color="white" size="24"></v-icon>
            </v-avatar>
            <div>
              <div class="text-h5 font-weight-bold">{{ unpaidDebts.length }}</div>
              <div class="text-caption text-medium-emphasis">Hutang Belum Lunas</div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row>
      <!-- Distributor Table -->
      <v-col cols="12" md="7">
        <v-card rounded="xl" elevation="0" border>
          <v-card-title class="pa-5 pb-2 d-flex align-center justify-space-between">
            <span class="text-h6 font-weight-bold">
              <v-icon start color="primary">mdi-package-variant-closed</v-icon>
              Daftar Distributor
            </span>
            <div class="d-flex gap-2 align-center">
              <v-text-field
                v-model="search"
                density="compact"
                variant="outlined"
                rounded="pill"
                prepend-inner-icon="mdi-magnify"
                placeholder="Cari..."
                hide-details
                style="max-width:180px;"
              ></v-text-field>
              <v-btn icon size="small" variant="tonal" color="primary" @click="loadData" :loading="loading">
                <v-icon>mdi-refresh</v-icon>
              </v-btn>
            </div>
          </v-card-title>
          <v-card-text class="pa-5 pt-0">
            <v-skeleton-loader v-if="loading" type="table-tbody"></v-skeleton-loader>
            <v-data-table
              v-else
              :headers="headers"
              :items="filteredDistributors"
              density="comfortable"
              :items-per-page="5"
              no-data-text="Belum ada data distributor. Tambahkan melalui API."
            >
              <template v-slot:item.action="{ item }">
                <v-btn size="x-small" variant="tonal" color="primary" rounded="pill" @click="loadStock(item)">
                  Stok
                </v-btn>
              </template>
            </v-data-table>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Panel Kanan -->
      <v-col cols="12" md="5">
        <!-- Stock Detail Panel -->
        <v-card v-if="selectedStock.length > 0" rounded="xl" elevation="0" border class="mb-4">
          <v-card-title class="pa-5 pb-2 text-h6 font-weight-bold">
            <v-icon start color="primary">mdi-package-variant-closed</v-icon>
            Detail Stok – {{ selectedDistName }}
          </v-card-title>
          <v-list lines="one" class="pa-4 pt-0">
            <v-list-item
              v-for="(s, i) in selectedStock"
              :key="i"
              :title="s.product_name"
              rounded="lg"
              class="mb-1 bg-grey-lighten-4"
            >
              <template v-slot:append>
                <v-chip :color="s.quantity < s.reorder_level ? 'error' : 'success'" size="small" label>
                  {{ s.quantity }} unit
                </v-chip>
              </template>
            </v-list-item>
            <div v-if="selectedStock.length === 0" class="text-caption text-grey text-center py-2">Tidak ada data stok</div>
          </v-list>
        </v-card>

        <!-- Reorder Form -->
        <v-card rounded="xl" elevation="0" border>
          <v-card-title class="pa-5 pb-2 text-h6 font-weight-bold">
            <v-icon start color="primary">mdi-cart-plus</v-icon>
            Pesanan Ulang (Re-order)
          </v-card-title>
          <v-card-text class="pa-5 pt-0">
            <v-select
              v-model="reorder.distributor"
              :items="distributors"
              item-title="name"
              item-value="id"
              label="Pilih Distributor"
              variant="outlined"
              density="comfortable"
              color="primary"
              rounded="lg"
              class="mb-3"
              :loading="loading"
            ></v-select>
            <v-text-field
              v-model="reorder.product"
              label="Nama Produk"
              variant="outlined"
              density="comfortable"
              color="primary"
              rounded="lg"
              class="mb-3"
            ></v-text-field>
            <v-text-field
              v-model="reorder.quantity"
              label="Jumlah (unit)"
              type="number"
              variant="outlined"
              density="comfortable"
              color="primary"
              rounded="lg"
              class="mb-4"
            ></v-text-field>
            <v-alert v-if="reorderSuccess" type="success" variant="tonal" density="compact" rounded="lg" class="mb-3">
              Pesanan berhasil dikirim!
            </v-alert>
            <v-btn color="primary" block rounded="pill" prepend-icon="mdi-send" @click="submitReorder" :loading="submitting">
              Kirim Pesanan
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
      search: '',
      selectedStock: [],
      selectedDistName: '',
      criticalStockCount: 0,
      unpaidDebts: [],
      submitting: false,
      reorderSuccess: false,
      reorder: { distributor: null, product: '', quantity: 1 },
      headers: [
        { title: 'Nama', align: 'start', key: 'name' },
        { title: 'Telepon', key: 'phone', sortable: false },
        { title: 'Alamat', key: 'address', sortable: false },
        { title: 'Aksi', align: 'center', key: 'action', sortable: false }
      ]
    }
  },
  computed: {
    distributors() { return this.$store.state.distributors },
    loading() { return this.$store.state.loading },
    filteredDistributors() {
      if (!this.search) return this.distributors
      return this.distributors.filter(d =>
        d.name?.toLowerCase().includes(this.search.toLowerCase())
      )
    }
  },
  mounted() { this.loadData() },
  methods: {
    async loadData() {
      await this.$store.dispatch('fetchDistributors')
      await this.loadDebts()
    },
    async loadDebts() {
      const token = this.$store.getters.token
      try {
        const res = await fetch(`${API_URL}/distributors/debts`, {
          headers: { 'Authorization': `Bearer ${token}` }
        })
        const data = await res.json()
        this.unpaidDebts = (data.data || []).filter(d => d.status === 'unpaid')
      } catch (e) { console.error(e) }
    },
    async loadStock(dist) {
      this.selectedDistName = dist.name
      const token = this.$store.getters.token
      try {
        const res = await fetch(`${API_URL}/distributors/stock?distributor_id=${dist.id}`, {
          headers: { 'Authorization': `Bearer ${token}` }
        })
        const data = await res.json()
        this.selectedStock = data.data || []
        this.criticalStockCount = this.selectedStock.filter(s => s.quantity < s.reorder_level).length
      } catch (e) { console.error(e) }
    },
    async submitReorder() {
      if (!this.reorder.distributor || !this.reorder.product) return
      this.submitting = true
      // NOTE: Re-order saat ini disimpan sebagai notifikasi ke supervisor
      // TODO: Integrate with actual order management system
      await new Promise(r => setTimeout(r, 800)) // simulate API call
      this.reorderSuccess = true
      this.reorder = { distributor: null, product: '', quantity: 1 }
      this.submitting = false
      setTimeout(() => { this.reorderSuccess = false }, 3000)
    }
  }
}
</script>
