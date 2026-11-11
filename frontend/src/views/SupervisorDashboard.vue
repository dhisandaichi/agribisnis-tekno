<template>
  <v-container fluid class="pa-4">
    <div class="d-flex align-center flex-wrap mb-4">
      <div class="text-headline-small">Dashboard Supervisor</div>
      <v-spacer />
      <v-btn variant="tonal" prepend-icon="mdi-refresh" :loading="loading" @click="loadAll">Muat ulang</v-btn>
    </div>

    <v-row>
      <v-col v-for="s in stats" :key="s.title" cols="12" sm="6" lg="3">
        <v-card :color="s.color" variant="flat" class="text-white">
          <v-card-item>
            <template #prepend><v-icon size="36">{{ s.icon }}</v-icon></template>
            <v-card-title>{{ s.title }}</v-card-title>
            <v-card-subtitle class="text-wrap">{{ s.subtitle }}</v-card-subtitle>
          </v-card-item>
          <v-card-text class="text-headline-large font-weight-bold">{{ s.value }}</v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-card class="mt-4" elevation="2">
      <v-tabs v-model="tab" color="primary" show-arrows>
        <v-tab value="visits" prepend-icon="mdi-map-marker-multiple">Kunjungan</v-tab>
        <v-tab value="distributors" prepend-icon="mdi-store">Distributor & Order</v-tab>
        <v-tab value="users" prepend-icon="mdi-account-group">Pengguna</v-tab>
        <v-tab value="content" prepend-icon="mdi-book-open-variant">Konten Petani</v-tab>
      </v-tabs>
      <v-divider />

      <v-window v-model="tab">
        <!-- KUNJUNGAN -->
        <v-window-item value="visits" class="pa-4">
          <div class="d-flex flex-wrap ga-2 mb-3">
            <v-text-field v-model="visitSearch" label="Cari petani / sales / komoditas" prepend-inner-icon="mdi-magnify"
              density="compact" variant="outlined" hide-details style="max-width: 360px" />
            <v-select v-model="visitDays" :items="dayOptions" item-title="label" item-value="value" label="Periode"
              density="compact" variant="outlined" hide-details style="max-width: 200px" @update:model-value="loadVisits" />
            <v-spacer />
            <v-btn variant="outlined" prepend-icon="mdi-file-delimited" @click="exportVisits" :disabled="!visits.length">Ekspor CSV</v-btn>
          </div>
          <v-data-table :mobile="$vuetify.display.xs" :headers="visitHeaders" :items="visitRows" :search="visitSearch" :loading="loading"
            items-per-page="10" no-data-text="Belum ada kunjungan pada periode ini">
            <template #item.visit_time="{ item }">{{ waktu(item.visit_time) }}</template>
            <template #item.gps="{ item }">
              <a v-if="item.latitude" :href="`https://www.google.com/maps?q=${item.latitude},${item.longitude}`" target="_blank" rel="noopener">
                <v-icon size="small">mdi-map-marker</v-icon> Peta
              </a>
              <span v-else>-</span>
            </template>
            <template #item.photo="{ item }">
              <v-btn v-if="item.has_photo" size="small" variant="text" icon="mdi-image" @click="openPhoto(item)" />
              <span v-else>-</span>
            </template>
          </v-data-table>
          <div class="text-title-small mt-4 mb-2">Jumlah kunjungan per sales (periode terpilih)</div>
          <v-table density="compact">
            <thead><tr><th>Sales</th><th class="text-right">Kunjungan</th><th class="text-right">Total Lahan (Ha)</th></tr></thead>
            <tbody>
              <tr v-for="r in perSales" :key="r.name"><td>{{ r.name }}</td><td class="text-right">{{ r.count }}</td><td class="text-right">{{ r.area.toFixed(2) }}</td></tr>
              <tr v-if="!perSales.length"><td colspan="3" class="text-medium-emphasis">Tidak ada data</td></tr>
            </tbody>
          </v-table>
        </v-window-item>

        <!-- DISTRIBUTOR & ORDER -->
        <v-window-item value="distributors" class="pa-4">
          <div class="d-flex flex-wrap ga-2 mb-3">
            <div class="text-title-medium">Kios / Distributor</div>
            <v-spacer />
            <v-btn color="primary" prepend-icon="mdi-plus" @click="openDistributor()">Tambah Distributor</v-btn>
          </div>
          <v-data-table :mobile="$vuetify.display.xs" :headers="distHeaders" :items="distRows" items-per-page="10" no-data-text="Belum ada distributor">
            <template #item.credit_limit_days="{ item }">{{ item.credit_limit_days ? `Tempo ${item.credit_limit_days} hari` : 'COD' }}</template>
            <template #item.outstanding="{ item }">{{ rupiah(item.outstanding) }}</template>
            <template #item.status="{ item }">
              <v-chip v-if="item.overdue > 0" color="error" size="small">Terlambat {{ item.days_late }} hari</v-chip>
              <v-chip v-else-if="item.outstanding > 0" color="warning" size="small">Belum lunas</v-chip>
              <v-chip v-else color="success" size="small">Lunas</v-chip>
            </template>
            <template #item.actions="{ item }">
              <v-btn size="small" variant="text" icon="mdi-pencil" @click="openDistributor(item)" />
            </template>
          </v-data-table>

          <div class="d-flex flex-wrap ga-2 mt-6 mb-3">
            <div class="text-title-medium">Order</div>
            <v-spacer />
            <v-btn color="primary" prepend-icon="mdi-cart-plus" :disabled="!distributors.length" @click="openOrder">Tambah Order</v-btn>
          </div>
          <v-data-table :mobile="$vuetify.display.xs" :headers="orderHeaders" :items="orderRows" items-per-page="10" no-data-text="Belum ada order">
            <template #item.order_date="{ item }">{{ tanggal(item.order_date) }}</template>
            <template #item.due_date="{ item }">{{ tanggal(item.due_date) }}</template>
            <template #item.total_amount="{ item }">{{ rupiah(item.total_amount) }}</template>
            <template #item.payment_status="{ item }">
              <v-chip :color="item.payment_status === 'paid' ? 'success' : (isOverdue(item) ? 'error' : 'warning')" size="small">
                {{ item.payment_status === 'paid' ? 'Lunas' : (isOverdue(item) ? 'Jatuh tempo' : 'Belum lunas') }}
              </v-chip>
            </template>
            <template #item.actions="{ item }">
              <v-btn size="small" variant="tonal" :color="item.payment_status === 'paid' ? 'grey' : 'success'" @click="togglePaid(item)">
                {{ item.payment_status === 'paid' ? 'Batalkan lunas' : 'Tandai lunas' }}
              </v-btn>
              <v-btn size="small" variant="text" icon="mdi-pencil" @click="openOrder(item)" />
            </template>
          </v-data-table>
        </v-window-item>

        <!-- PENGGUNA -->
        <v-window-item value="users" class="pa-4">
          <v-alert type="info" variant="tonal" density="compact" class="mb-3">
            Ubah role pengguna (misalnya menjadikan seseorang Supervisor), nonaktifkan akun, atau buat password sementara untuk pengguna yang lupa password.
          </v-alert>
          <v-text-field v-model="userSearch" label="Cari pengguna" prepend-inner-icon="mdi-magnify" density="compact"
            variant="outlined" hide-details class="mb-3" style="max-width: 360px" />
          <v-data-table :mobile="$vuetify.display.xs" :headers="userHeaders" :items="users" :search="userSearch" items-per-page="10">
            <template #item.created_at="{ item }">{{ tanggal(item.created_at) }}</template>
            <template #item.role="{ item }">
              <v-select :model-value="item.role" :items="roleItems" item-title="label" item-value="value"
                density="compact" variant="outlined" hide-details style="min-width: 190px"
                :disabled="item.id === myId" @update:model-value="(v) => updateUser(item, { role: v })" />
            </template>
            <template #item.active="{ item }">
              <v-switch :model-value="item.active" color="success" hide-details density="compact" :disabled="item.id === myId"
                @update:model-value="(v) => updateUser(item, { active: v })" />
            </template>
            <template #item.actions="{ item }">
              <v-btn size="small" variant="text" prepend-icon="mdi-lock-reset" :disabled="item.id === myId" @click="resetPassword(item)">Reset</v-btn>
            </template>
          </v-data-table>
        </v-window-item>

        <!-- KONTEN PETANI -->
        <v-window-item value="content" class="pa-4">
          <div class="d-flex align-center mb-2">
            <div class="text-title-medium">Harga Komoditas</div><v-spacer />
            <v-btn color="primary" size="small" prepend-icon="mdi-plus" @click="openPrice()">Tambah</v-btn>
          </div>
          <v-table density="compact" class="mb-6">
            <thead><tr><th>Komoditas</th><th class="text-right">Harga</th><th>Tren</th><th>Diperbarui</th><th></th></tr></thead>
            <tbody>
              <tr v-for="p in prices" :key="p.id">
                <td>{{ p.name }}</td><td class="text-right">{{ rupiah(p.price) }} / {{ p.unit }}</td>
                <td><v-icon :color="trendColor(p.trend)">{{ trendIcon(p.trend) }}</v-icon></td>
                <td>{{ tanggal(p.updated_at) }}</td>
                <td class="text-right">
                  <v-btn size="small" variant="text" icon="mdi-pencil" @click="openPrice(p)" />
                  <v-btn size="small" variant="text" icon="mdi-delete" color="error" @click="remove(`/supervisor/prices/${p.id}`)" />
                </td>
              </tr>
            </tbody>
          </v-table>

          <div class="d-flex align-center mb-2">
            <div class="text-title-medium">Konten Edukasi</div><v-spacer />
            <v-btn color="primary" size="small" prepend-icon="mdi-plus" @click="openEdu()">Tambah</v-btn>
          </div>
          <v-table density="compact">
            <thead><tr><th>Judul</th><th>Video</th><th>Tayang</th><th></th></tr></thead>
            <tbody>
              <tr v-for="e in education" :key="e.id">
                <td>{{ e.title }}</td>
                <td><a v-if="e.video_url" :href="e.video_url" target="_blank" rel="noopener">Buka</a><span v-else>-</span></td>
                <td><v-icon :color="e.published ? 'success' : 'grey'">{{ e.published ? 'mdi-eye' : 'mdi-eye-off' }}</v-icon></td>
                <td class="text-right">
                  <v-btn size="small" variant="text" icon="mdi-pencil" @click="openEdu(e)" />
                  <v-btn size="small" variant="text" icon="mdi-delete" color="error" @click="remove(`/supervisor/education/${e.id}`)" />
                </td>
              </tr>
            </tbody>
          </v-table>
        </v-window-item>
      </v-window>
    </v-card>

    <!-- Dialog foto -->
    <v-dialog v-model="photo.show" max-width="720">
      <v-card>
        <v-img v-if="photo.url" :src="photo.url" max-height="70vh" contain />
        <v-card-text v-else class="text-center pa-8"><v-progress-circular indeterminate /></v-card-text>
        <v-card-actions><v-spacer /><v-btn @click="photo.show = false">Tutup</v-btn></v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Dialog form generik -->
    <v-dialog v-model="dialog.show" max-width="520" persistent>
      <v-card :title="dialog.title">
        <v-card-text>
          <v-form ref="dialogForm">
            <template v-for="f in dialog.fields" :key="f.key">
              <v-select v-if="f.type === 'select'" v-model="dialog.model[f.key]" :items="f.items" item-title="label"
                item-value="value" :label="f.label" variant="outlined" :rules="f.required ? [req] : []" :clearable="!f.required" />
              <v-switch v-else-if="f.type === 'switch'" v-model="dialog.model[f.key]" :label="f.label" color="primary" />
              <v-textarea v-else-if="f.type === 'textarea'" v-model="dialog.model[f.key]" :label="f.label" rows="2" variant="outlined" />
              <v-text-field v-else v-model="dialog.model[f.key]" :label="f.label" :type="f.type || 'text'" variant="outlined"
                :rules="f.required ? [req] : []" />
            </template>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn @click="dialog.show = false">Batal</v-btn>
          <v-btn color="primary" variant="flat" :loading="dialog.saving" @click="saveDialog">Simpan</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="tempPw.show" max-width="420">
      <v-card title="Password Sementara">
        <v-card-text>
          Berikan password ini kepada <b>{{ tempPw.email }}</b>. Minta pengguna segera menggantinya setelah login.
          <div class="text-headline-small text-center my-4 font-weight-bold" style="font-family: monospace">{{ tempPw.password }}</div>
        </v-card-text>
        <v-card-actions><v-spacer /><v-btn @click="tempPw.show = false">Tutup</v-btn></v-card-actions>
      </v-card>
    </v-dialog>

    <v-snackbar v-model="snack.show" :color="snack.color" timeout="4000">{{ snack.text }}</v-snackbar>
  </v-container>
</template>

<script>
import { get, post, put, patch, del, photoUrl } from '../services/api'
import { rupiah, tanggal, waktu, ROLE_LABELS } from '../services/format'

const day = (d) => (d ? String(d).slice(0, 10) : null)

export default {
  data() {
    return {
      tab: 'visits',
      loading: false,
      summary: {},
      visits: [],
      distributors: [],
      orders: [],
      users: [],
      prices: [],
      education: [],
      visitSearch: '',
      userSearch: '',
      visitDays: 30,
      dayOptions: [
        { value: 1, label: 'Hari ini' }, { value: 7, label: '7 hari' }, { value: 30, label: '30 hari' },
        { value: 90, label: '90 hari' }, { value: 365, label: '1 tahun' },
      ],
      visitHeaders: [
        { title: 'Waktu', key: 'visit_time' }, { title: 'Sales', key: 'sales_name' }, { title: 'Petani', key: 'farmer_name' },
        { title: 'Desa', key: 'village' }, { title: 'Komoditas', key: 'crop_type' }, { title: 'Lahan (Ha)', key: 'land_area', align: 'end' },
        { title: 'Kompetitor', key: 'competitor_product' }, { title: 'Catatan', key: 'notes' },
        { title: 'GPS', key: 'gps', sortable: false }, { title: 'Foto', key: 'photo', sortable: false },
      ],
      distHeaders: [
        { title: 'Nama Kios', key: 'name' }, { title: 'Pemilik (akun)', key: 'owner' }, { title: 'Tempo', key: 'credit_limit_days' },
        { title: 'Piutang', key: 'outstanding', align: 'end' }, { title: 'Status', key: 'status', sortable: false },
        { title: '', key: 'actions', sortable: false, align: 'end' },
      ],
      orderHeaders: [
        { title: 'Tanggal', key: 'order_date' }, { title: 'Distributor', key: 'distributor_name' }, { title: 'Produk', key: 'product' },
        { title: 'Qty', key: 'quantity', align: 'end' }, { title: 'Total', key: 'total_amount', align: 'end' },
        { title: 'Jatuh Tempo', key: 'due_date' }, { title: 'Status', key: 'payment_status' },
        { title: '', key: 'actions', sortable: false, align: 'end' },
      ],
      userHeaders: [
        { title: 'Nama', key: 'full_name' }, { title: 'Email', key: 'email' }, { title: 'No. HP', key: 'phone' },
        { title: 'Terdaftar', key: 'created_at' }, { title: 'Role', key: 'role', sortable: false },
        { title: 'Aktif', key: 'active', sortable: false }, { title: 'Password', key: 'actions', sortable: false },
      ],
      roleItems: Object.entries(ROLE_LABELS).map(([value, label]) => ({ value, label })),
      photo: { show: false, url: null },
      tempPw: { show: false, email: '', password: '' },
      dialog: { show: false, title: '', path: '', id: null, fields: [], model: {}, saving: false },
      snack: { show: false, text: '', color: 'success' },
    }
  },
  computed: {
    myId() { return this.$store.getters.userId },
    visitRows() {
      return this.visits.map((v) => ({ ...v, village: v.village || '-', crop_type: v.crop_type || '-', competitor_product: v.competitor_product || '-', notes: v.notes || '' }))
    },
    perSales() {
      const map = {}
      for (const v of this.visits) {
        map[v.sales_name] = map[v.sales_name] || { name: v.sales_name, count: 0, area: 0 }
        map[v.sales_name].count++
        map[v.sales_name].area += Number(v.land_area) || 0
      }
      return Object.values(map).sort((a, b) => b.count - a.count)
    },
    distRows() {
      return this.distributors.map((d) => ({ ...d, owner: d.owner_name || '-', days_late: d.days_late || 0 }))
    },
    orderRows() { return this.orders },
    stats() {
      const s = this.summary
      return [
        { title: 'Total Petani', subtitle: 'Terdata dari kunjungan', value: (s.farmers || 0).toLocaleString('id-ID'), icon: 'mdi-account-group', color: 'success' },
        { title: 'Kunjungan 7 Hari', subtitle: `${s.visits_today || 0} hari ini`, value: s.visits_7d || 0, icon: 'mdi-map-marker-check', color: 'primary' },
        { title: 'Distributor / Kios', subtitle: `${s.overdue_distributors || 0} lewat jatuh tempo · ${s.pending_orders || 0} order menunggu harga`, value: s.distributors || 0, icon: 'mdi-store', color: 'info' },
        { title: 'Piutang Jatuh Tempo', subtitle: `Total belum lunas ${rupiah(s.outstanding_amount)}`, value: rupiah(s.overdue_amount), icon: 'mdi-cash-clock', color: 'error' },
      ]
    },
  },
  mounted() {
    this.loadAll()
  },
  methods: {
    rupiah, tanggal, waktu,
    req: (v) => (v !== null && v !== undefined && String(v).trim() !== '') || 'Wajib diisi',
    notify(text, color = 'success') { this.snack = { show: true, text, color } },
    trendIcon: (t) => ({ up: 'mdi-trending-up', down: 'mdi-trending-down' }[t] || 'mdi-trending-neutral'),
    trendColor: (t) => ({ up: 'success', down: 'error' }[t] || 'grey'),
    isOverdue: (o) => o.payment_status !== 'paid' && o.due_date && day(o.due_date) < new Date().toLocaleDateString('sv-SE'),

    async safe(promise, fallback) {
      try { return await promise } catch (err) { this.notify(err.message, 'error'); return fallback }
    },

    async loadAll() {
      this.loading = true
      const [summary, dists, orders, users, prices, edu] = await Promise.all([
        this.safe(get('/supervisor/summary'), {}),
        this.safe(get('/supervisor/distributors'), []),
        this.safe(get('/supervisor/orders'), []),
        this.safe(get('/supervisor/users'), []),
        this.safe(get('/content/prices'), []),
        this.safe(get('/content/education'), []),
      ])
      Object.assign(this, { summary, distributors: dists, orders, users, prices, education: edu })
      await this.loadVisits()
      this.loading = false
    },

    async loadVisits() {
      this.visits = await this.safe(get(`/supervisor/visits?days=${this.visitDays}`), [])
    },

    async openPhoto(item) {
      if (this.photo.url) URL.revokeObjectURL(this.photo.url)
      this.photo = { show: true, url: null }
      try {
        this.photo.url = await photoUrl(item.id)
      } catch (err) {
        this.photo.show = false
        this.notify(err.message, 'error')
      }
    },

    exportVisits() {
      const cols = ['visit_time', 'sales_name', 'farmer_name', 'village', 'crop_type', 'land_area', 'competitor_product', 'notes', 'latitude', 'longitude']
      const esc = (v) => `"${String(v ?? '').replace(/"/g, '""')}"`
      const csv = [cols.join(','), ...this.visitRows.map((r) => cols.map((c) => esc(r[c])).join(','))].join('\n')
      const url = URL.createObjectURL(new Blob(['\ufeff' + csv], { type: 'text/csv;charset=utf-8' }))
      const a = document.createElement('a')
      a.href = url
      a.download = `kunjungan-${new Date().toISOString().slice(0, 10)}.csv`
      a.click()
      URL.revokeObjectURL(url)
    },

    async updateUser(user, change) {
      try {
        await patch(`/supervisor/users/${user.id}`, change)
        Object.assign(user, change)
        this.notify(`Pengguna ${user.email} diperbarui.`)
      } catch (err) {
        this.notify(err.message, 'error')
      }
    },

    async resetPassword(user) {
      if (!confirm(`Buat password sementara untuk ${user.email}? Password lama tidak berlaku lagi.`)) return
      try {
        const { temporary_password } = await post(`/supervisor/users/${user.id}/reset-password`)
        this.tempPw = { show: true, email: user.email, password: temporary_password }
      } catch (err) {
        this.notify(err.message, 'error')
      }
    },

    async togglePaid(order) {
      const status = order.payment_status === 'paid' ? 'unpaid' : 'paid'
      if (status === 'paid' && !(Number(order.total_amount) > 0) && !confirm('Total order masih Rp 0. Tetap tandai lunas?')) return
      try {
        await patch(`/supervisor/orders/${order.id}/status`, { payment_status: status })
        this.notify(status === 'paid' ? 'Order ditandai lunas.' : 'Status lunas dibatalkan.')
        this.loadAll()
      } catch (err) {
        this.notify(err.message, 'error')
      }
    },

    async remove(path) {
      if (!confirm('Hapus data ini?')) return
      try {
        await del(path)
        this.notify('Data dihapus.')
        this.loadAll()
      } catch (err) {
        this.notify(err.message, 'error')
      }
    },

    openDialog(title, path, fields, model = {}, id = null) {
      this.dialog = { show: true, title, path, id, fields, model: { ...model }, saving: false }
    },

    openDistributor(d) {
      const owners = this.users.filter((u) => u.role === 'distributor')
        .map((u) => ({ value: u.id, label: `${u.full_name || '-'} (${u.email})` }))
      this.openDialog(d ? 'Ubah Distributor' : 'Tambah Distributor', '/supervisor/distributors', [
        { key: 'name', label: 'Nama Kios *', required: true },
        { key: 'owner_id', label: 'Akun Pemilik (role distributor)', type: 'select', items: owners },
        { key: 'address', label: 'Alamat', type: 'textarea', str: true },
        { key: 'phone', label: 'No. HP', type: 'tel', str: true },
        { key: 'credit_limit_days', label: 'Tempo Pembayaran (hari, 0 = COD) *', type: 'number', required: true },
      ], d ? { name: d.name, owner_id: d.owner_id, address: d.address, phone: d.phone, credit_limit_days: d.credit_limit_days }
           : { credit_limit_days: 0 }, d?.id)
    },

    openOrder(o) {
      const isEdit = o && o.id
      const items = this.distributors.map((d) => ({ value: d.id, label: d.name }))
      this.openDialog(isEdit ? 'Ubah Order' : 'Tambah Order', '/supervisor/orders', [
        { key: 'distributor_id', label: 'Distributor *', type: 'select', items, required: true },
        { key: 'product', label: 'Produk *', required: true },
        { key: 'quantity', label: 'Jumlah *', type: 'number', required: true },
        { key: 'total_amount', label: 'Total Nilai (Rp) *', type: 'number', required: true },
        { key: 'order_date', label: 'Tanggal Order *', type: 'date', required: true },
        { key: 'due_date', label: 'Jatuh Tempo (kosongkan = otomatis dari tempo)', type: 'date', str: true },
      ], isEdit
        ? { distributor_id: o.distributor_id, product: o.product, quantity: o.quantity, total_amount: o.total_amount, order_date: day(o.order_date), due_date: day(o.due_date) }
        : { order_date: new Date().toLocaleDateString('sv-SE') }, isEdit ? o.id : null)
    },

    openPrice(p) {
      this.openDialog(p ? 'Ubah Harga' : 'Tambah Harga', '/supervisor/prices', [
        { key: 'name', label: 'Komoditas *', required: true },
        { key: 'price', label: 'Harga (Rp) *', type: 'number', required: true },
        { key: 'unit', label: 'Satuan *', required: true },
        { key: 'trend', label: 'Tren', type: 'select', required: true,
          items: [{ value: 'up', label: 'Naik' }, { value: 'flat', label: 'Stabil' }, { value: 'down', label: 'Turun' }] },
      ], p ? { name: p.name, price: p.price, unit: p.unit, trend: p.trend } : { unit: 'kg', trend: 'flat' }, p?.id)
    },

    openEdu(e) {
      this.openDialog(e ? 'Ubah Konten' : 'Tambah Konten', '/supervisor/education', [
        { key: 'title', label: 'Judul *', required: true },
        { key: 'description', label: 'Deskripsi', type: 'textarea' },
        { key: 'video_url', label: 'Link Video (YouTube, dll)', type: 'url' },
        { key: 'thumbnail_url', label: 'Link Gambar Thumbnail', type: 'url' },
        { key: 'duration_minutes', label: 'Durasi (menit)', type: 'number' },
        { key: 'published', label: 'Tayangkan ke petani', type: 'switch' },
      ], e ? { ...e } : { published: true }, e?.id)
    },

    async saveDialog() {
      const { valid } = await this.$refs.dialogForm.validate()
      if (!valid) return
      const d = this.dialog
      const payload = {}
      for (const f of d.fields) {
        let v = d.model[f.key]
        if (v === '' || v === undefined) v = f.str ? '' : null
        if (f.type === 'number' && v !== null && v !== '') v = Number(v)
        if (f.type === 'switch') v = !!v
        payload[f.key] = v
      }
      d.saving = true
      try {
        if (d.id) await put(`${d.path}/${d.id}`, payload)
        else await post(d.path, payload)
        d.show = false
        this.notify('Data tersimpan.')
        this.loadAll()
      } catch (err) {
        this.notify(err.message, 'error')
      } finally {
        d.saving = false
      }
    },
  },
}
</script>
