<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <div class="text-h4 mb-4">Analitik Supervisor Dashboard</div>
      </v-col>
    </v-row>

    <v-row>
      <!-- Cards Info -->
      <v-col cols="12" md="4">
        <v-card color="success" theme="dark">
          <v-card-item title="Total Petani">
            <template v-slot:subtitle>Terdaftar aktif </template>
          </v-card-item>
          <v-card-text class="text-h2 py-4">
            1,204
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="4">
        <v-card color="info" theme="dark">
          <v-card-item title="Distributor / Kios">
            <template v-slot:subtitle>Membutuhkan Restock</template>
          </v-card-item>
          <v-card-text class="text-h2 py-4">
            32
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="4">
        <v-card color="error" theme="dark">
          <v-card-item title="Hutang Distributor">
            <template v-slot:subtitle>Lewat Jatuh Tempo (30 Hari)</template>
          </v-card-item>
          <v-card-text class="text-h2 py-4">
            Rp 42.5M
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Map & Stok section placeholder -->
      <v-col cols="12" class="mt-4">
        <v-card elevation="2">
          <v-card-title>
            <v-icon left>mdi-store</v-icon>
            Manajemen Kios & Stok (Data Tabel)
          </v-card-title>
          <v-card-text>
            <v-data-table
              :headers="headers"
              :items="kios"
              class="elevation-0"
            ></v-data-table>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  data() {
    return {
      headers: [
        { title: 'Nama Kios', align: 'start', key: 'name' },
        { title: 'Status Pembayaran', align: 'start', key: 'status' },
        { title: 'Total Tagihan', align: 'end', key: 'debt' },
        { title: 'Tempo', align: 'end', key: 'due' },
      ],
      kios: [
        { name: 'Kios Tani Makmur', status: 'COD', debt: 'Lunas', due: '-' },
        { name: 'UD Subur Jaya', status: 'Hutang (>30hr)', debt: 'Rp 14.000.000', due: 'Terlambat 5 hari' },
      ]
    }
  }
}
</script>
