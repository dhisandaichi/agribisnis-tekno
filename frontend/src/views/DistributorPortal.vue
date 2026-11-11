<template>
  <v-container>
    <div class="text-headline-small mb-4">Portal Distributor / Kios</div>

    <v-alert v-if="!loading && !distributor" type="info" variant="tonal">
      Akun Anda belum terhubung ke data kios. Hubungi supervisor Agritec agar akun Anda dihubungkan.
    </v-alert>

    <template v-if="distributor">
      <v-row>
        <v-col cols="12" md="4">
          <v-card elevation="2">
            <v-card-item>
              <template #prepend><v-icon size="36" color="primary">mdi-store</v-icon></template>
              <v-card-title>{{ distributor.name }}</v-card-title>
              <v-card-subtitle>{{ distributor.credit_limit_days ? `Tempo ${distributor.credit_limit_days} hari` : 'Pembayaran COD' }}</v-card-subtitle>
            </v-card-item>
            <v-card-text>
              <div v-if="distributor.address">{{ distributor.address }}</div>
              <div v-if="distributor.phone">{{ distributor.phone }}</div>
            </v-card-text>
          </v-card>
        </v-col>
        <v-col cols="12" sm="6" md="4">
          <v-card color="warning" variant="flat" class="text-white">
            <v-card-item title="Total Tagihan Belum Lunas" />
            <v-card-text class="text-headline-large font-weight-bold">{{ rupiah(outstanding) }}</v-card-text>
          </v-card>
        </v-col>
        <v-col cols="12" sm="6" md="4">
          <v-card color="error" variant="flat" class="text-white">
            <v-card-item title="Lewat Jatuh Tempo" />
            <v-card-text class="text-headline-large font-weight-bold">{{ rupiah(overdue) }}</v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <v-card class="mt-4" elevation="2">
        <v-card-title class="d-flex align-center flex-wrap">
          Riwayat Order
          <v-spacer />
          <v-btn color="primary" prepend-icon="mdi-cart-plus" @click="dialog = true">Ajukan Order</v-btn>
        </v-card-title>
        <v-data-table :headers="headers" :items="orders" :loading="loading" items-per-page="10" :mobile="$vuetify.display.xs" no-data-text="Belum ada order">
          <template #item.order_date="{ item }">{{ tanggal(item.order_date) }}</template>
          <template #item.due_date="{ item }">{{ tanggal(item.due_date) }}</template>
          <template #item.total_amount="{ item }">{{ item.total_amount > 0 ? rupiah(item.total_amount) : 'Menunggu konfirmasi' }}</template>
          <template #item.payment_status="{ item }">
            <v-chip :color="item.payment_status === 'paid' ? 'success' : (isOverdue(item) ? 'error' : 'warning')" size="small">
              {{ item.payment_status === 'paid' ? 'Lunas' : (isOverdue(item) ? 'Jatuh tempo' : 'Belum lunas') }}
            </v-chip>
          </template>
        </v-data-table>
      </v-card>
    </template>

    <v-dialog v-model="dialog" max-width="480">
      <v-card title="Ajukan Order Baru">
        <v-card-text>
          <v-form ref="form">
            <v-combobox v-model="form.product" :items="products" label="Produk *" variant="outlined"
              :rules="[(v) => !!v || 'Wajib diisi']" />
            <v-text-field v-model="form.quantity" label="Jumlah (sak/unit) *" type="number" min="1" variant="outlined"
              :rules="[(v) => Number(v) > 0 || 'Minimal 1']" />
          </v-form>
          <div class="text-body-small text-medium-emphasis">Nilai order akan dikonfirmasi oleh supervisor Agritec.</div>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn @click="dialog = false">Batal</v-btn>
          <v-btn color="primary" variant="flat" :loading="saving" @click="submit">Kirim</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-snackbar v-model="snack.show" :color="snack.color" timeout="4000">{{ snack.text }}</v-snackbar>
  </v-container>
</template>

<script>
import { get, post } from '../services/api'
import { rupiah, tanggal } from '../services/format'

export default {
  data() {
    return {
      loading: true,
      saving: false,
      dialog: false,
      distributor: null,
      orders: [],
      form: { product: '', quantity: '' },
      products: ['Pupuk Organik Granul', 'Pupuk NPK', 'Benih Padi', 'Benih Jagung', 'Pestisida'],
      headers: [
        { title: 'Tanggal', key: 'order_date' }, { title: 'Produk', key: 'product' },
        { title: 'Qty', key: 'quantity', align: 'end' }, { title: 'Total', key: 'total_amount', align: 'end' },
        { title: 'Jatuh Tempo', key: 'due_date' }, { title: 'Status', key: 'payment_status' },
      ],
      snack: { show: false, text: '', color: 'success' },
    }
  },
  computed: {
    outstanding() { return this.orders.filter((o) => o.payment_status !== 'paid').reduce((s, o) => s + Number(o.total_amount), 0) },
    overdue() { return this.orders.filter((o) => this.isOverdue(o)).reduce((s, o) => s + Number(o.total_amount), 0) },
  },
  mounted() { this.load() },
  methods: {
    rupiah, tanggal,
    isOverdue: (o) => o.payment_status !== 'paid' && o.due_date && String(o.due_date).slice(0, 10) < new Date().toLocaleDateString('sv-SE'),
    async load() {
      this.loading = true
      try {
        const data = await get('/distributor')
        this.distributor = data.distributor
        this.orders = data.orders
      } catch (err) {
        this.snack = { show: true, text: err.message, color: 'error' }
      } finally {
        this.loading = false
      }
    },
    async submit() {
      const { valid } = await this.$refs.form.validate()
      if (!valid) return
      this.saving = true
      try {
        await post('/distributor/orders', { product: this.form.product, quantity: Number(this.form.quantity) })
      } catch (err) {
        this.saving = false
        this.snack = { show: true, text: err.message, color: 'error' }
        return
      }
      this.saving = false
      this.dialog = false
      this.form = { product: '', quantity: '' }
      this.snack = { show: true, text: 'Order terkirim ke supervisor.', color: 'success' }
      this.load()
    },
  },
}
</script>
