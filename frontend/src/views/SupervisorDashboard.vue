<template>
  <v-container>
    <v-row>
      <v-col cols="12" class="d-flex align-center">
        <div class="text-h4">Analitik Supervisor Dashboard</div>
        <v-spacer></v-spacer>
        <v-btn icon variant="text" aria-label="Muat ulang" :loading="loading" @click="load">
          <v-icon>mdi-refresh</v-icon>
        </v-btn>
      </v-col>
    </v-row>

    <v-alert v-if="error" type="error" class="mb-4" density="compact">{{ error }}</v-alert>

    <v-row>
      <v-col cols="12" md="4">
        <v-card color="success" theme="dark">
          <v-card-item title="Total Petani">
            <template v-slot:subtitle>Kunjungan 7 hari: {{ summary.visits_7d }}</template>
          </v-card-item>
          <v-card-text class="text-h2 py-4">
            {{ formatNumber(summary.total_farmers) }}
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="4">
        <v-card color="info" theme="dark">
          <v-card-item title="Distributor / Kios">
            <template v-slot:subtitle>Hutang lewat jatuh tempo</template>
          </v-card-item>
          <v-card-text class="text-h2 py-4">
            {{ formatNumber(summary.overdue_distributors) }}
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="4">
        <v-card color="error" theme="dark">
          <v-card-item title="Hutang Distributor">
            <template v-slot:subtitle>Total lewat jatuh tempo</template>
          </v-card-item>
          <v-card-text class="text-h4 py-6">
            {{ formatCurrency(summary.overdue_amount) }}
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" class="mt-4">
        <v-card elevation="2">
          <v-card-title>
            <v-icon start>mdi-map-marker-check</v-icon>
            Kunjungan Sales Terbaru
          </v-card-title>
          <v-card-text>
            <v-data-table
              :headers="visitHeaders"
              :items="visits"
              :loading="loading"
              items-per-page="10"
              no-data-text="Belum ada kunjungan."
              class="elevation-0"
            >
              <template #item.visit_time="{ item }">{{ formatDate(item.visit_time) }}</template>
              <template #item.location="{ item }">
                <a :href="mapsUrl(item)" target="_blank" rel="noopener">
                  {{ item.latitude.toFixed(5) }}, {{ item.longitude.toFixed(5) }}
                </a>
              </template>
              <template #item.photo="{ item }">
                <v-btn
                  v-if="item.photo_path"
                  size="small"
                  variant="text"
                  prepend-icon="mdi-image"
                  @click="openPhoto(item.photo_path)"
                >Lihat</v-btn>
              </template>
            </v-data-table>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" class="mt-4">
        <v-card elevation="2">
          <v-card-title>
            <v-icon start>mdi-store</v-icon>
            Manajemen Kios & Hutang
          </v-card-title>
          <v-card-text>
            <v-data-table
              :headers="headers"
              :items="kios"
              :loading="loading"
              no-data-text="Belum ada data distributor."
              class="elevation-0"
            >
              <template #item.credit_term_days="{ item }">
                {{ item.credit_term_days === 0 ? 'COD' : `Tempo ${item.credit_term_days} hari` }}
              </template>
              <template #item.unpaid_amount="{ item }">
                {{ Number(item.unpaid_amount) > 0 ? formatCurrency(item.unpaid_amount) : 'Lunas' }}
              </template>
              <template #item.days_overdue="{ item }">
                <v-chip v-if="item.days_overdue > 0" color="error" size="small">
                  Terlambat {{ item.days_overdue }} hari
                </v-chip>
                <span v-else-if="item.oldest_due_date">Jatuh tempo {{ formatDate(item.oldest_due_date) }}</span>
                <span v-else>-</span>
              </template>
            </v-data-table>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { supabase, VISIT_PHOTO_BUCKET } from '../services/supabase'

export default {
  data() {
    return {
      loading: false,
      error: null,
      summary: { total_farmers: 0, visits_7d: 0, overdue_distributors: 0, overdue_amount: 0 },
      visits: [],
      kios: [],
      visitHeaders: [
        { title: 'Waktu', key: 'visit_time' },
        { title: 'Sales', key: 'sales' },
        { title: 'Petani', key: 'farmer' },
        { title: 'Komoditas', key: 'crop_type' },
        { title: 'Luas (Ha)', key: 'land_area', align: 'end' },
        { title: 'Lokasi', key: 'location', sortable: false },
        { title: 'Foto', key: 'photo', sortable: false }
      ],
      headers: [
        { title: 'Nama Kios', align: 'start', key: 'name' },
        { title: 'Pembayaran', align: 'start', key: 'credit_term_days' },
        { title: 'Total Tagihan', align: 'end', key: 'unpaid_amount' },
        { title: 'Tempo', align: 'end', key: 'days_overdue' }
      ]
    }
  },
  mounted() {
    this.load()
  },
  methods: {
    async load() {
      this.loading = true
      this.error = null
      try {
        const [summary, visits, kios] = await Promise.all([
          supabase.rpc('supervisor_summary'),
          supabase
            .from('visits')
            .select('id, visit_time, latitude, longitude, photo_path, farmers(name, crop_type, land_area), profiles(email, full_name)')
            .order('visit_time', { ascending: false })
            .limit(50),
          supabase.from('distributor_balances').select('*').order('days_overdue', { ascending: false })
        ])
        for (const res of [summary, visits, kios]) {
          if (res.error) throw res.error
        }
        this.summary = summary.data
        this.visits = visits.data.map(v => ({
          ...v,
          sales: v.profiles?.full_name || v.profiles?.email,
          farmer: v.farmers?.name,
          crop_type: v.farmers?.crop_type,
          land_area: v.farmers?.land_area
        }))
        this.kios = kios.data
      } catch (err) {
        console.error(err)
        this.error = err.message || 'Gagal memuat data dashboard.'
      } finally {
        this.loading = false
      }
    },
    async openPhoto(path) {
      // Buka tab dulu agar tidak diblokir popup blocker, lalu isi dengan signed URL.
      const win = window.open('', '_blank')
      const { data, error } = await supabase.storage.from(VISIT_PHOTO_BUCKET).createSignedUrl(path, 300)
      if (error) {
        win?.close()
        this.error = error.message
        return
      }
      if (win) win.location.href = data.signedUrl
    },
    mapsUrl(item) {
      return `https://www.google.com/maps?q=${item.latitude},${item.longitude}`
    },
    formatNumber(n) {
      return new Intl.NumberFormat('id-ID').format(Number(n) || 0)
    },
    formatCurrency(n) {
      return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(Number(n) || 0)
    },
    formatDate(d) {
      return new Date(d).toLocaleString('id-ID', { dateStyle: 'medium', ...(d.length > 10 ? { timeStyle: 'short' } : {}) })
    }
  }
}
</script>
