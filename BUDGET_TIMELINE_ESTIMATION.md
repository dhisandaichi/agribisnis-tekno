# Estimasi Waktu dan Biaya (Budget & Timeline) - SFA Agritec

Berdasarkan *Business Requirement Document* (BRD) dan rancangan *Technical Requirement* yang disepakati untuk PT. Agribisnis Technology Indonesia (Agritec), berikut adalah kalkulasi estimasi waktu pengerjaan dan anggaran yang ditekan **seminimal mungkin (efisien)** untuk mencapai MVP (*Minimum Viable Product*).

---

## 1. Estimasi Waktu Pengerjaan (Timeline)
Total Estimasi Waktu: **~2 hingga 3 Bulan** (Estimasi optimal untuk menghasilkan sistem stabil, aman, dan siap pakai).

### Bulan 1: Infrastruktur, Core Backend, dan Integrasi Fase 1
* **Infrastruktur (DevOps):** Setup VPS (Ubuntu), Dockerisasi, Nginx (*Reverse Proxy*), instalasi PostgreSQL + ekstensi PostGIS, serta Redis.
* **Backend (Golang):** Pembuatan fondasi API, autentikasi JWT terpusat, API untuk input/simpan data Kunjungan Sales, Manajemen Distributor. Integrasi *library kompresi gambar* (contoh: Vips/Sharp).
* **Frontend (Vue.js):** Menghubungkan halaman frontend PWA *existing* ke API (Fungsi Login, Form Kunjungan, Geotagging API HTML5). Konfigurasi *Service Worker* agar form kunjungan bisa tersimpan sementara beroperasi di *Offline Mode*.

### Bulan 2: Manajemen Distributor, Analitik, dan Fitur Edukasi
* **Backend & Pihak Ketiga:** Setup API SMS/WhatsApp untuk peringatan hutang. Menyusun skrip agregasi data (*Python/Golang cron jobs*).
* **Frontend (Vue.js):** Melengkapi interaksi Data Stok, Hitung Hutang/COD, Dashboard *Supervisor*. Mengintegrasikan Video Pemutar Edukasi (HLS Stream) dan pemanggilan *endpoint* BMKG.
* **Analisis Data:** Integrasi Grafana untuk mengambil query data peta sebaran dan grafik dari PostGIS.

### Bulan 3: UAT (*User Acceptance Testing*), Optimasi, & Deployment
* Pengetesan menyeluruh dari skenario Sales (sistem nge-drop di pelosok, sync otomatis pas dapat sinyal).
* Evaluasi kualitas kompresi foto agar file tetap di bawah 500 KB per lembar.
* Skalabilitas & beban testing awal. Finalisasi dan **Go-Live (Fase-1 Selesai)**.

---

## 2. Estimasi Biaya (Budget / Financial Plan)
Angka di bawah ini memproyeksikan standar pengerjaan *Outsource Freelance* spesialis lokal Indonesia (tanpa via *Agency* korporat) agar sesuai prinsip anggaran minimal perusahaan.

### A. Biaya Pengembangan Tim Teknis (CAPEX / Project Base)
Karena bagian rancangan UI / visual / file *Minimal Frontend* sudah diperbarui/berdiri, porsi Frontend menjadi lebih murah (fokus di integrasi API & PWA Cache).
1. **Frontend / PWA Integration:** Rp 4.000.000 – Rp 7.000.000
2. **Backend Development (Golang Core + APIs):** Rp 10.000.000 – Rp 14.000.000
3. **DevOps & Analytics Engineer (Docker, VPS, Grafana):** Rp 4.000.000 – Rp 6.000.000
* **Total Biaya Rancang Bangun Sistem (Satu Kali): ~Rp 18.000.000 – Rp 27.000.000**

### B. Biaya Infrastruktur Server Bulanan (OPEX)
Diakali menggunakan pendekatan *self-host semi-managed* di level VPS terjangkau yang andal (*DigitalOcean / Linode / Vultr*):
1. **VPS Utama (API + Database + Cache):** 2 vCPU, 4GB RAM ➔ ~$24/bln (**~Rp 380.000**)
2. **Object Storage (S3 S3-Compatible Storage):** Untuk menyimpan file foto kunjungan / dokumen statis tanpa memberatkan VPS utama (Kapasitas ~250GB) ➔ ~$5/bln (**~Rp 80.000**)
3. **Sewa Nama Domain (.com / .id):** ~Rp 150.000 s/d Rp 250.000/Tahun ➔ **~Rp 15.000 / bln**
4. **CDN & Pengaman DDoS:** Cloudflare Basic ➔ **GRATIS**
* **Total Biaya Server / Cloud Per Bulan: ~Rp 475.000**

### C. Biaya Integrasi Sistem Pihak Ketiga (API Services Bulanan)
Tagihan mengandalkan model *Pay-As-You-Go* berdasarkan jumlah pemakaian nyata:
1. **Notifikasi Hutang Kios** (Layanan WhatsApp/SMS Non-Resmi seperti Fonnte / Wablas untuk penghematan, dibanding Twilio) ➔ **~Rp 150.000 – Rp 250.000 / bln**
2. **Layanan Streaming Video Edukasi** (Cloudflare Stream atau VOD biasa CDN) ➔ **~Rp 50.000 – Rp 150.000 / bln** (*fleksibel bergantung seberapa sering video ditonton*)
* **Total Biaya Eksternal API Per Bulan: ~Rp 200.000 – Rp 400.000**

---

## Ringkasan Eksekutif Budget

| Kategori | Model Pendanaan | Estimasi Minimal | Estimasi Maksimal |
|---|---|---|---|
| **Pengembangan (Tim Teknis)** | *Satu Kali Pembayaran (Termin)* | **Rp 18.000.000** | **Rp 27.000.000** |
| **Infrastruktur Penuh** | *Berlangganan Bulanan* | **Rp 675.000/bln** | **Rp 875.000/bln** |

---

### *Saran Tindakan Teknis Lanjutan:*
Dengan berkas frontend yang sebagian sudah siap, kita dapat langsung mulai tancap gas me-scaffold Backend (Golang). Apakah Anda ingin saya membantu **menyiapkan kerangka *boilerplate* RESTful API Golang (dengan Echo/Gin)** atau **menulis skrip migrasi database PostgreSQL** di repositori ini sekarang?
