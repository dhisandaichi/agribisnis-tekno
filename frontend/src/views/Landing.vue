<template>
  <div class="lp">
    <!-- NAV -->
    <header class="lp-nav" :class="{ 'lp-nav--solid': scrolled }">
      <div class="lp-wrap lp-nav__inner">
        <router-link to="/" class="lp-brand" aria-label="Agrina 99 beranda">
          <span class="brand-logo-badge"><img src="/brand/logo.png" alt="" width="40" height="40" /></span>
          <span>
            <span class="lp-brand__kicker">Pupuk Organik <b>Premium</b></span>
            <span class="lp-brand__name">AGRINA 99</span>
          </span>
        </router-link>
        <nav class="lp-nav__links">
          <a href="#produk">Produk</a>
          <a href="#cara-pakai">Cara Pakai</a>
          <a href="#tanaman-buah">Tanaman Buah</a>
          <a href="#kemitraan">Kemitraan</a>
          <a href="#faq">FAQ</a>
        </nav>
        <div class="lp-nav__cta">
          <v-btn v-if="authed" :to="homePath" variant="text" color="white" class="d-none d-sm-flex">Buka Aplikasi</v-btn>
          <v-btn v-else to="/login" variant="text" color="white" class="d-none d-sm-flex">Masuk</v-btn>
          <v-btn href="#pesan" color="accent" variant="flat" class="font-weight-bold">Pesan</v-btn>
        </div>
      </div>
    </header>

    <!-- HERO -->
    <section class="lp-hero">
      <div class="lp-hero__bg" />
      <div class="lp-wrap lp-hero__grid">
        <div>
          <div class="lp-chip"><v-icon size="16">mdi-leaf</v-icon> Penyubur Tanah, Daun dan Buah</div>
          <h1>Tanah Subur, <span class="hl-green">Panen Melimpah</span>, Petani Sejahtera.</h1>
          <p class="lp-hero__lead">
            Pupuk Organik Premium <b>Agrina 99</b>, <b>Dolomite Super Premium M100</b> dan <b>Top Super Booster</b>
            diproduksi oleh PT. Agribisnis Technology Indonesia untuk memperbaiki tanah, menyeimbangkan pH,
            dan mendorong pembungaan serta pembuahan.
          </p>
          <div class="lp-hero__actions">
            <v-btn href="#pesan" color="accent" size="x-large" variant="flat" class="font-weight-bold" prepend-icon="mdi-whatsapp">
              Minta Penawaran
            </v-btn>
            <v-btn href="#produk" size="x-large" variant="outlined" color="white" append-icon="mdi-arrow-down">Lihat Produk</v-btn>
          </div>
          <ul class="lp-hero__trust">
            <li><v-icon>mdi-factory</v-icon> Produksi Sumatera Barat</li>
            <li><v-icon>mdi-recycle</v-icon> Berbahan organik</li>
            <li><v-icon>mdi-truck-delivery-outline</v-icon> Melayani kios & partai besar</li>
          </ul>
        </div>
        <div class="lp-hero__visual">
          <div class="lp-hero__stage" aria-hidden="true" />
          <SackViewer />
          <div class="lp-hero__badge">
            <div class="lp-hero__badge-num">50 kg</div>
            <div>berat bersih<br />per sak</div>
          </div>
        </div>
      </div>
    </section>

    <!-- KAPASITAS PRODUKSI -->
    <section class="lp-capacity" aria-labelledby="kapasitas-title">
      <div class="lp-wrap lp-capacity__inner">
        <div class="lp-capacity__main">
          <div class="lp-capacity__label">Kapasitas Produksi</div>
          <h2 id="kapasitas-title">Pupuk Organik Premium Agrina 99</h2>
          <div class="lp-capacity__num">10.000 <span>Ton/Bulan</span></div>
        </div>
        <ul class="lp-capacity__facts">
          <li v-for="f in capacityFacts" :key="f.label">
            <v-icon size="26">{{ f.icon }}</v-icon>
            <div><b>{{ f.value }}</b><span>{{ f.label }}</span></div>
          </li>
        </ul>
      </div>
    </section>

    <!-- MASALAH → SOLUSI -->
    <section class="lp-section">
      <div class="lp-wrap">
        <div class="lp-head">
          <div class="lp-eyebrow">Kenapa hasil panen tidak maksimal?</div>
          <h2>Masalahnya sering ada di tanah, bukan di bibit.</h2>
          <p>Banyak lahan pertanian dan perkebunan di Indonesia bersifat masam dan miskin bahan organik setelah bertahun-tahun memakai pupuk kimia saja. Akibatnya pupuk yang ditabur tidak terserap dengan baik.</p>
        </div>
        <div class="lp-problems">
          <div v-for="p in problems" :key="p.title" class="lp-problem">
            <div class="lp-problem__icon" :style="{ background: p.bg, color: p.fg }"><v-icon size="28">{{ p.icon }}</v-icon></div>
            <h3>{{ p.title }}</h3>
            <p>{{ p.text }}</p>
            <div class="lp-problem__fix"><v-icon size="18" color="primary">mdi-check-circle</v-icon> <span>Solusi: <b>{{ p.fix }}</b></span></div>
          </div>
        </div>
      </div>
    </section>

    <!-- PRODUK -->
    <section id="produk" class="lp-section lp-section--tint">
      <div class="lp-wrap">
        <div class="lp-head">
          <div class="lp-eyebrow">Produk kami</div>
          <h2>Satu rangkaian untuk tanah yang sehat</h2>
          <p>Pakai satu per satu sesuai kebutuhan lahan, atau gabungkan ketiganya sebagai program pemupukan lengkap.</p>
        </div>
        <div class="lp-products">
          <article v-for="prod in products" :key="prod.id" class="lp-product" :style="{ '--accent': prod.color }">
            <div class="lp-product__img">
              <img :src="prod.img" :alt="'Kemasan ' + prod.name" loading="lazy" />
              <span class="lp-product__weight">{{ prod.weight }}</span>
            </div>
            <div class="lp-product__body">
              <div class="lp-product__tag">{{ prod.tag }}</div>
              <h3>{{ prod.name }}</h3>
              <p>{{ prod.desc }}</p>
              <ul>
                <li v-for="b in prod.benefits" :key="b"><v-icon size="18">mdi-check</v-icon> {{ b }}</li>
              </ul>
              <div class="lp-product__crops"><v-icon size="16">mdi-sprout-outline</v-icon> {{ prod.crops }}</div>
              <v-btn :href="'#pesan'" block variant="flat" :color="prod.btn" class="mt-auto font-weight-bold" @click="pick(prod.id)">
                Pesan {{ prod.short }}
              </v-btn>
            </div>
          </article>
        </div>
      </div>
    </section>

    <!-- CARA PAKAI -->
    <section id="cara-pakai" class="lp-section">
      <div class="lp-wrap">
        <div class="lp-head">
          <div class="lp-eyebrow">Program Tanah Sehat Agrina</div>
          <h2>Tiga langkah dari olah tanah sampai panen</h2>
          <p>Urutan pemakaian yang kami sarankan. Dosis menyesuaikan jenis tanaman, kondisi tanah, dan anjuran pada kemasan atau tim lapangan kami.</p>
        </div>
        <ol class="lp-steps">
          <li v-for="(s, i) in steps" :key="s.title" class="lp-step">
            <div class="lp-step__num">{{ i + 1 }}</div>
            <div class="lp-step__when">{{ s.when }}</div>
            <h3>{{ s.title }}</h3>
            <p>{{ s.text }}</p>
          </li>
        </ol>

        <div class="lp-crops">
          <div class="lp-crops__title">Cocok untuk berbagai komoditas</div>
          <div class="lp-crops__list">
            <span v-for="c in crops" :key="c.name" class="lp-crop"><v-icon size="20">{{ c.icon }}</v-icon> {{ c.name }}</span>
          </div>
        </div>
      </div>
    </section>

    <!-- TANAMAN BUAH -->
    <section id="tanaman-buah" class="lp-section lp-section--tint">
      <div class="lp-wrap">
        <div class="lp-head">
          <div class="lp-eyebrow">Edukasi petani</div>
          <h2>Alasan pentingnya pemupukan untuk tanaman buah</h2>
          <p>Tanaman buah dipanen berkali-kali dari pohon yang sama selama bertahun-tahun. Hara yang terbawa setiap panen harus dikembalikan agar pohon tetap berbuah lebat dan berkualitas.</p>
        </div>
        <div class="lp-fruit">
          <div v-for="(r, i) in fruitReasons" :key="r.title" class="lp-fruit__item">
            <div class="lp-fruit__num">{{ String(i + 1).padStart(2, '0') }}</div>
            <v-icon size="30" class="lp-fruit__icon">{{ r.icon }}</v-icon>
            <h3>{{ r.title }}</h3>
            <p>{{ r.text }}</p>
          </div>
        </div>
        <div class="lp-fruit__program">
          <div class="lp-fruit__program-title"><v-icon>mdi-calendar-check-outline</v-icon> Waktu pemupukan yang disarankan</div>
          <ol>
            <li><b>Setelah panen.</b> Agrina 99 di piringan pohon untuk memulihkan tanah, ditambah Dolomite bila tanah masam.</li>
            <li><b>Menjelang berbunga.</b> Top Super Booster agar bunga keluar serempak dan tidak mudah rontok.</li>
            <li><b>Saat pembesaran buah.</b> Ulangi Agrina 99 sesuai anjuran supaya buah padat, manis, dan warnanya merata.</li>
          </ol>
        </div>
      </div>
    </section>

    <!-- KEMITRAAN -->
    <section id="kemitraan" class="lp-section lp-section--navy">
      <div class="lp-wrap">
        <div class="lp-head lp-head--light">
          <div class="lp-eyebrow">Beli dan bermitra</div>
          <h2>Pilih cara yang paling pas untuk Anda</h2>
        </div>
        <div class="lp-partners">
          <div v-for="pt in partners" :key="pt.title" class="lp-partner">
            <v-icon size="34" class="mb-3">{{ pt.icon }}</v-icon>
            <h3>{{ pt.title }}</h3>
            <p>{{ pt.text }}</p>
            <v-btn v-if="pt.to" :to="pt.to" color="white" variant="flat" class="text-secondary font-weight-bold mt-auto">{{ pt.cta }}</v-btn>
            <v-btn v-else href="#pesan" color="accent" variant="flat" class="font-weight-bold mt-auto" @click="buyer = pt.buyer">{{ pt.cta }}</v-btn>
          </div>
        </div>
        <div class="lp-app">
          <v-icon size="40" color="#7ee2a8">mdi-cellphone-check</v-icon>
          <div>
            <b>Kios dan distributor mitra memesan lewat aplikasi SFA Agritec.</b>
            Pesanan langsung diterima tim kami, statusnya bisa dipantau, dan sales lapangan rutin berkunjung.
          </div>
        </div>
      </div>
    </section>

    <!-- FORM PENAWARAN -->
    <section id="pesan" class="lp-section lp-section--tint">
      <div class="lp-wrap lp-order">
        <div>
          <div class="lp-eyebrow">Minta penawaran</div>
          <h2>Dapatkan harga terbaik untuk kebutuhan Anda</h2>
          <p class="text-medium-emphasis">
            Isi formulir singkat ini. Pesan akan terbuka di WhatsApp dan tim penjualan kami membalas dengan harga, stok, dan ongkos kirim ke lokasi Anda.
          </p>
          <ul class="lp-order__why">
            <li><v-icon color="primary">mdi-cash-multiple</v-icon> Harga khusus untuk pembelian partai dan kios mitra</li>
            <li><v-icon color="primary">mdi-account-tie-outline</v-icon> Pendampingan dari sales lapangan kami</li>
            <li><v-icon color="primary">mdi-map-marker-radius-outline</v-icon> Berbasis di Payakumbuh, Sumatera Barat</li>
          </ul>
        </div>
        <v-card class="lp-order__card" elevation="0">
          <v-form ref="orderForm" @submit.prevent="sendOrder">
            <v-text-field v-model.trim="order.name" label="Nama" variant="outlined" color="primary" density="comfortable"
              prepend-inner-icon="mdi-account-outline" :rules="[req]" />
            <v-text-field v-model.trim="order.location" label="Kabupaten / kota" variant="outlined" color="primary" density="comfortable"
              prepend-inner-icon="mdi-map-marker-outline" :rules="[req]" />
            <v-select v-model="buyer" :items="buyerTypes" label="Saya adalah" variant="outlined" color="primary" density="comfortable"
              prepend-inner-icon="mdi-badge-account-outline" />
            <div class="text-body-medium font-weight-medium mb-1">Produk yang diminati</div>
            <div class="lp-order__products mb-3">
              <v-checkbox v-for="prod in products" :key="prod.id" v-model="order.products" :value="prod.id" :label="prod.short"
                color="primary" density="compact" hide-details />
            </div>
            <v-text-field v-model.trim="order.qty" label="Perkiraan jumlah (karung)" variant="outlined" color="primary" density="comfortable"
              type="number" min="1" prepend-inner-icon="mdi-package-variant-closed" />
            <v-text-field v-model.trim="order.crop" label="Komoditas / luas lahan (opsional)" variant="outlined" color="primary" density="comfortable"
              prepend-inner-icon="mdi-sprout-outline" placeholder="mis. cabai 1 ha" />
            <v-alert v-if="orderError" type="warning" variant="tonal" density="compact" class="mb-3">{{ orderError }}</v-alert>
            <v-btn type="submit" color="primary" size="x-large" block variant="flat" class="font-weight-bold" prepend-icon="mdi-whatsapp">
              Kirim via WhatsApp
            </v-btn>
          </v-form>
        </v-card>
      </div>
    </section>

    <!-- FAQ -->
    <section id="faq" class="lp-section">
      <div class="lp-wrap lp-wrap--narrow">
        <div class="lp-head">
          <div class="lp-eyebrow">Pertanyaan umum</div>
          <h2>Yang sering ditanyakan pembeli</h2>
        </div>
        <v-expansion-panels variant="accordion" class="lp-faq">
          <v-expansion-panel v-for="f in faqs" :key="f.q" :title="f.q" :text="f.a" />
        </v-expansion-panels>
      </div>
    </section>

    <!-- FOOTER -->
    <footer class="lp-footer">
      <div class="lp-social">
        <div class="lp-wrap lp-social__inner">
          <span class="lp-social__title">Ikuti kami</span>
          <a v-for="sc in social" :key="sc.label" :href="sc.href" target="_blank" rel="noopener" class="lp-social__link">
            <span class="lp-social__icon"><v-icon size="18">{{ sc.icon }}</v-icon></span> {{ sc.label }}
          </a>
        </div>
      </div>
      <div class="lp-wrap lp-footer__grid">
        <div>
          <div class="lp-brand mb-3">
            <span class="brand-logo-badge"><img src="/brand/logo.png" alt="" width="44" height="44" /></span>
            <span>
              <span class="lp-brand__kicker">PT. Agribisnis Technology Indonesia</span>
              <span class="lp-brand__name">AGRINA 99</span>
            </span>
          </div>
          <p>Tanah Subur, Panen Melimpah, Petani Sejahtera.</p>
          <p class="mb-0"><v-icon size="18">mdi-map-marker-outline</v-icon> Payakumbuh, Sumatera Barat</p>
        </div>
        <div>
          <div class="lp-footer__title">Produk</div>
          <a v-for="prod in products" :key="prod.id" href="#produk">{{ prod.name }}</a>
        </div>
        <div>
          <div class="lp-footer__title">Hubungi kami</div>
          <a :href="waLink('Halo Agrina 99, saya ingin bertanya tentang produk.')" target="_blank" rel="noopener">
            <v-icon size="18">mdi-whatsapp</v-icon> {{ contact.phoneLabel }}
          </a>
          <a v-for="sc in social" :key="'f-' + sc.label" :href="sc.href" target="_blank" rel="noopener"><v-icon size="18">{{ sc.icon }}</v-icon> {{ sc.label }}</a>
        </div>
        <div>
          <div class="lp-footer__title">Untuk mitra</div>
          <router-link to="/login">Masuk SFA Agritec</router-link>
          <router-link to="/register?role=distributor">Daftar kios / distributor</router-link>
          <router-link to="/register?role=farmer">Daftar petani</router-link>
        </div>
      </div>
      <div class="lp-wrap lp-footer__bottom">© {{ year }} PT. Agribisnis Technology Indonesia. Hak cipta dilindungi.</div>
    </footer>

    <!-- Tombol WhatsApp melayang -->
    <a class="lp-wa" :href="waLink('Halo Agrina 99, saya ingin memesan pupuk.')" target="_blank" rel="noopener" aria-label="Chat WhatsApp">
      <v-icon size="30">mdi-whatsapp</v-icon>
    </a>
  </div>
</template>

<script>
import { homeForRole } from '../services/format'
import SackViewer from '../components/sack/SackViewer.vue'

// Nomor informasi dari materi resmi Agrina 99. Ganti di sini bila tim penjualan memakai nomor lain.
const CONTACT = { wa: '6282170425387', phoneLabel: '0821-7042-5387' }

// Akun resmi dari materi promosi Agrina 99.
// Alamat Facebook belum diketahui pasti, jadi sementara diarahkan ke pencarian nama halamannya.
const SOCIAL = [
  { icon: 'mdi-web', label: 'agritechindonesia.id', href: 'https://agritechindonesia.id' },
  { icon: 'mdi-facebook', label: 'Pupuk Organik Agrina', href: 'https://www.facebook.com/search/top?q=pupuk%20organik%20agrina' },
  { icon: 'mdi-instagram', label: 'pupukorganikagrina99', href: 'https://www.instagram.com/pupukorganikagrina99' },
  { icon: 'mdi-music-note', label: 'pupukorganikagrina99', href: 'https://www.tiktok.com/@pupukorganikagrina99' },
]

export default {
  components: { SackViewer },
  data() {
    return {
      scrolled: false,
      year: new Date().getFullYear(),
      contact: CONTACT,
      social: SOCIAL,
      capacityFacts: [
        { icon: 'mdi-package-variant-closed', value: '± 200.000 sak', label: 'setara kemasan 50 kg per bulan' },
        { icon: 'mdi-factory', value: 'Payakumbuh', label: 'pabrik di Sumatera Barat' },
        { icon: 'mdi-truck-fast-outline', value: 'Partai besar', label: 'siap pasok distributor dan perkebunan' },
      ],
      fruitReasons: [
        { icon: 'mdi-basket-outline', title: 'Mengganti hara yang terbawa panen',
          text: 'Setiap kilogram buah yang dipetik membawa nitrogen, kalium dan unsur lain keluar dari kebun. Tanpa pemupukan, cadangan hara tanah terus menipis.' },
        { icon: 'mdi-flower-outline', title: 'Bunga lebih banyak, tidak mudah rontok',
          text: 'Fase berbunga butuh fosfor, boron dan kalsium yang cukup. Kekurangan unsur ini membuat bunga dan bakal buah mudah gugur.' },
        { icon: 'mdi-fruit-grapes-outline', title: 'Buah lebih besar, manis dan berwarna',
          text: 'Kalium, magnesium dan bahan organik membantu pembentukan gula dan warna, sehingga buah lebih bernas dan harga jualnya lebih baik.' },
        { icon: 'mdi-sprout-outline', title: 'Akar kuat di tanah yang gembur',
          text: 'Pohon buah tumbuh bertahun-tahun di tempat yang sama. Bahan organik menjaga tanah tetap gembur, kaya mikroba baik, dan mampu menahan air.' },
        { icon: 'mdi-shield-sun-outline', title: 'Lebih tahan kemarau dan penyakit',
          text: 'Tanaman yang cukup hara dan tumbuh di tanah sehat lebih kuat menghadapi kekeringan, serangan hama, dan penyakit.' },
        { icon: 'mdi-chart-line', title: 'Panen stabil dari musim ke musim',
          text: 'Pemupukan rutin membantu pohon pulih setelah berbuah lebat, sehingga hasil tidak naik turun tajam antar musim.' },
      ],
      buyer: 'Petani',
      buyerTypes: ['Petani', 'Kelompok tani / Gapoktan', 'Kios / toko pertanian', 'Distributor', 'Perkebunan / perusahaan', 'Instansi / dinas'],
      order: { name: '', location: '', products: [], qty: '', crop: '' },
      orderError: null,
      req: (v) => !!v || 'Wajib diisi',
      problems: [
        { icon: 'mdi-flask-outline', bg: '#E8EAF6', fg: '#2E3192', title: 'Tanah masam (pH rendah)',
          text: 'Pada tanah masam, unsur hara seperti fosfor terikat dan sulit diserap akar walau pupuk sudah banyak ditabur.',
          fix: 'Dolomite Super Premium M100' },
        { icon: 'mdi-terrain', bg: '#E6F5EC', fg: '#007A3A', title: 'Tanah keras dan cepat kering',
          text: 'Bahan organik yang rendah membuat tanah padat, miskin mikroba baik, dan tidak mampu menahan air serta pupuk.',
          fix: 'Pupuk Organik Agrina 99' },
        { icon: 'mdi-flower-outline', bg: '#FDEEE4', fg: '#C2511A', title: 'Bunga rontok, buah sedikit',
          text: 'Kekurangan boron, kalsium dan magnesium di fase generatif membuat bunga mudah gugur dan buah kurang bernas.',
          fix: 'Top Super Booster' },
      ],
      products: [
        { id: 'agrina', short: 'Agrina 99', name: 'Pupuk Organik Agrina 99', tag: 'Pupuk organik premium', weight: '50 kg',
          img: '/brand/pack-agrina99.jpg', color: '#009A4A', btn: 'primary',
          desc: 'Pupuk organik padat untuk pupuk dasar. Menghidupkan kembali tanah agar pupuk lain bekerja lebih efektif.',
          benefits: ['Menggemburkan dan memperbaiki struktur tanah', 'Menambah bahan organik dan mikroba baik', 'Tanah lebih kuat menahan air dan hara', 'Mengurangi ketergantungan pada pupuk kimia'],
          crops: 'Padi, jagung, hortikultura, sawit, perkebunan' },
        { id: 'dolomite', short: 'Dolomite M100', name: 'Dolomite Super Premium M100', tag: 'Penetral pH tanah', weight: '50 kg',
          img: '/brand/pack-dolomite.jpg', color: '#2E3192', btn: 'secondary',
          desc: 'Kapur pertanian dolomit butiran halus (mesh 100) sumber kalsium dan magnesium untuk lahan masam.',
          benefits: ['Menaikkan pH tanah masam', 'Sumber kalsium (CaO) dan magnesium (MgO)', 'Butiran halus, lebih cepat bereaksi', 'Membuat fosfor lebih mudah diserap'],
          crops: 'Lahan masam, gambut, sawah, kebun sawit' },
        { id: 'booster', short: 'Top Super Booster', name: 'Top Super Booster Agrina 99', tag: 'Pembungaan & pembuahan', weight: '25 kg',
          img: '/brand/pack-booster.jpg', color: '#F26B21', btn: 'accent',
          desc: 'Pupuk pelengkap dengan nitrogen, fosfor, boron, magnesium dan kalsium untuk fase berbunga dan berbuah.',
          benefits: ['Merangsang pembungaan serempak', 'Boron membantu mengurangi bunga rontok', 'Buah lebih padat dan berisi', 'Daun hijau berkat magnesium'],
          crops: 'Cabai, bawang, tomat, jagung, buah-buahan' },
      ],
      steps: [
        { when: '1 sampai 2 minggu sebelum tanam', title: 'Netralkan tanah dengan Dolomite', text: 'Tabur merata saat pengolahan tanah agar pH naik dan tanah siap menerima pupuk.' },
        { when: 'Saat olah tanah / pupuk dasar', title: 'Suburkan dengan Agrina 99', text: 'Campurkan ke bedengan atau lubang tanam sebagai pupuk dasar. Untuk tanaman tahunan, aplikasikan di piringan.' },
        { when: 'Menjelang fase berbunga', title: 'Dorong hasil dengan Top Super Booster', text: 'Berikan saat tanaman mulai masuk fase generatif untuk bunga yang lebih banyak dan buah yang berisi.' },
      ],
      crops: [
        { icon: 'mdi-grass', name: 'Padi' }, { icon: 'mdi-corn', name: 'Jagung' }, { icon: 'mdi-chili-mild', name: 'Cabai' },
        { icon: 'mdi-seed-outline', name: 'Bawang merah' }, { icon: 'mdi-leaf', name: 'Sayuran daun' }, { icon: 'mdi-palm-tree', name: 'Kelapa sawit' },
        { icon: 'mdi-fruit-cherries', name: 'Buah-buahan' }, { icon: 'mdi-coffee-outline', name: 'Kopi & kakao' },
      ],
      partners: [
        { icon: 'mdi-sprout', title: 'Petani & kelompok tani', text: 'Beli di kios mitra terdekat atau pesan langsung dalam jumlah partai untuk kelompok.', cta: 'Tanya harga', buyer: 'Kelompok tani / Gapoktan' },
        { icon: 'mdi-storefront-outline', title: 'Kios & toko pertanian', text: 'Jadi kios mitra Agrina 99. Pesan stok lewat aplikasi, dapat kunjungan rutin dan dukungan promosi.', cta: 'Daftar jadi mitra', to: '/register?role=distributor' },
        { icon: 'mdi-truck-outline', title: 'Distributor & perkebunan', text: 'Pasokan skala besar untuk distributor wilayah, perkebunan dan pengadaan instansi.', cta: 'Ajukan penawaran', buyer: 'Distributor' },
      ],
      faqs: [
        { q: 'Apa beda Agrina 99, Dolomite dan Top Super Booster?', a: 'Dolomite menetralkan tanah masam dan menambah kalsium serta magnesium. Agrina 99 adalah pupuk organik yang memperbaiki struktur dan kesuburan tanah. Top Super Booster adalah pupuk pelengkap untuk fase berbunga dan berbuah. Ketiganya bisa dipakai bersamaan sebagai satu program.' },
        { q: 'Apakah bisa dipakai bersama pupuk kimia (Urea, NPK)?', a: 'Bisa. Agrina 99 justru membantu pupuk kimia lebih efektif karena tanah mampu menahan hara lebih lama, sehingga dalam jangka panjang pemakaian pupuk kimia dapat dikurangi.' },
        { q: 'Berapa dosis yang tepat untuk lahan saya?', a: 'Dosis tergantung jenis tanaman, umur, dan kondisi tanah. Ikuti anjuran pada kemasan, atau kirim formulir penawaran dengan menyebut komoditas dan luas lahan, tim kami akan memberi rekomendasi.' },
        { q: 'Berapa minimal pembelian dan apakah bisa dikirim ke luar Sumatera Barat?', a: 'Petani bisa membeli eceran melalui kios mitra. Untuk pembelian partai dan pengiriman antar daerah, hubungi tim penjualan melalui WhatsApp agar dihitungkan harga dan ongkos kirimnya.' },
        { q: 'Bagaimana cara menjadi kios atau distributor resmi?', a: 'Daftar sebagai Kios / Distributor di aplikasi SFA Agritec, lalu tim sales kami akan menghubungi dan menghubungkan akun Anda ke data kios. Setelah itu Anda bisa memesan stok langsung dari aplikasi.' },
      ],
    }
  },
  computed: {
    authed() { return this.$store.getters.isAuthenticated },
    homePath() { return homeForRole(this.$store.getters.userRole) },
  },
  mounted() {
    this.onScroll()
    window.addEventListener('scroll', this.onScroll, { passive: true })
  },
  beforeUnmount() {
    window.removeEventListener('scroll', this.onScroll)
  },
  methods: {
    onScroll() { this.scrolled = window.scrollY > 24 },
    pick(id) { if (!this.order.products.includes(id)) this.order.products.push(id) },
    waLink(text) { return `https://wa.me/${this.contact.wa}?text=${encodeURIComponent(text)}` },
    async sendOrder() {
      const { valid } = await this.$refs.orderForm.validate()
      if (!valid) return
      if (!this.order.products.length) {
        this.orderError = 'Pilih minimal satu produk.'
        return
      }
      this.orderError = null
      const names = this.products.filter((p) => this.order.products.includes(p.id)).map((p) => p.name)
      const lines = [
        'Halo Agrina 99, saya ingin meminta penawaran.',
        '',
        `Nama: ${this.order.name}`,
        `Lokasi: ${this.order.location}`,
        `Saya: ${this.buyer}`,
        `Produk: ${names.join(', ')}`,
      ]
      if (this.order.qty) lines.push(`Perkiraan jumlah: ${this.order.qty} karung`)
      if (this.order.crop) lines.push(`Komoditas / lahan: ${this.order.crop}`)
      window.open(this.waLink(lines.join('\n')), '_blank', 'noopener')
    },
  },
}
</script>

<style scoped>
.lp { background: #fff; color: #1d2433; overflow-x: hidden; }
.lp-wrap { width: 100%; max-width: 1180px; margin: 0 auto; padding: 0 20px; }
.lp-wrap--narrow { max-width: 820px; }

/* NAV */
.lp-nav { position: fixed; inset: 0 0 auto 0; z-index: 20; transition: background .2s, box-shadow .2s; }
.lp-nav--solid { background: rgba(31, 34, 102, .96); box-shadow: 0 4px 20px rgba(0, 0, 0, .18); backdrop-filter: blur(6px); }
.lp-nav__inner { display: flex; align-items: center; gap: 16px; height: 72px; }
.lp-brand { display: inline-flex; align-items: center; gap: 10px; color: #fff; text-decoration: none; }
.lp-brand__kicker { display: block; font-size: 10px; letter-spacing: .08em; text-transform: uppercase; opacity: .9; }
.lp-brand__kicker b { color: #ffb37f; }
.lp-brand__name { display: block; font-size: 20px; font-weight: 800; line-height: 1.1; }
.lp-nav__links { display: flex; gap: 26px; margin-left: auto; }
.lp-nav__links a { color: #fff; text-decoration: none; font-weight: 500; opacity: .9; }
.lp-nav__links a:hover { opacity: 1; text-decoration: underline; text-underline-offset: 6px; }
.lp-nav__cta { display: flex; gap: 6px; align-items: center; }

/* HERO */
.lp-hero { position: relative; color: #fff; padding: 128px 0 88px;
  background: linear-gradient(140deg, #1f2266 0%, #2e3192 42%, #007a3a 100%); overflow: hidden; }
.lp-hero__bg { position: absolute; inset: 0; background: url('/brand/field.jpg') center / cover; opacity: .2; mix-blend-mode: luminosity; }
.lp-hero__bg::after { content: ''; position: absolute; inset: 0; background: radial-gradient(circle at 90% 0%, rgba(0, 154, 74, .6), transparent 50%); }
.lp-hero__grid { position: relative; display: grid; grid-template-columns: 1.05fr 1fr; gap: 48px; align-items: center; }
.lp-chip { display: inline-flex; align-items: center; gap: 6px; background: rgba(255, 255, 255, .14); border: 1px solid rgba(255, 255, 255, .25);
  border-radius: 999px; padding: 6px 14px; font-size: 14px; font-weight: 600; margin-bottom: 18px; }
.lp-hero h1 { font-size: clamp(34px, 5vw, 56px); font-weight: 800; line-height: 1.08; margin: 0 0 18px; letter-spacing: -.01em; }
.hl-green { color: #7ee2a8; }
.lp-hero__lead { font-size: 18px; line-height: 1.65; opacity: .92; max-width: 560px; margin: 0 0 28px; }
.lp-hero__actions { display: flex; flex-wrap: wrap; gap: 12px; margin-bottom: 30px; }
.lp-hero__trust { list-style: none; padding: 0; margin: 0; display: flex; flex-wrap: wrap; gap: 10px 22px; font-size: 14px; opacity: .92; }
.lp-hero__trust li { display: flex; align-items: center; gap: 8px; }
.lp-hero__trust .v-icon { color: #7ee2a8; }
.lp-hero__visual { position: relative; }
.lp-hero__badge { position: absolute; left: -18px; bottom: -22px; display: flex; align-items: center; gap: 12px; background: #fff; color: #1f2266;
  border-radius: 16px; padding: 12px 18px; box-shadow: 0 12px 30px rgba(0, 0, 0, .25); font-size: 13px; font-weight: 600; line-height: 1.3; }
.lp-hero__badge-num { font-size: 28px; font-weight: 800; color: var(--brand-orange); }

/* SAK 3D DI HERO */
.lp-hero__stage { position: absolute; left: 50%; top: 52%; width: 82%; aspect-ratio: 1; transform: translate(-50%, -50%); border-radius: 50%;
  background: radial-gradient(circle, rgba(126, 226, 168, .38) 0%, rgba(126, 226, 168, .12) 45%, transparent 70%); pointer-events: none; }

/* KAPASITAS */
.lp-capacity { background: linear-gradient(100deg, #007a3a, #009a4a 55%, #0b8f4a); color: #fff; padding: 34px 0; position: relative; }
.lp-capacity__inner { display: grid; grid-template-columns: 1.1fr 1.4fr; gap: 32px; align-items: center; }
.lp-capacity__label { font-size: 13px; font-weight: 700; letter-spacing: .12em; text-transform: uppercase; color: #d4ffe6; }
.lp-capacity h2 { font-size: 22px; font-weight: 700; margin: 2px 0 6px; }
.lp-capacity__num { display: inline-block; background: #fff; color: #007a3a; font-size: clamp(28px, 4vw, 40px); font-weight: 800; border-radius: 14px; padding: 4px 18px; line-height: 1.3; }
.lp-capacity__num span { font-size: .6em; font-weight: 700; }
.lp-capacity__facts { list-style: none; margin: 0; padding: 0; display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 14px; }
.lp-capacity__facts li { display: flex; gap: 12px; align-items: flex-start; background: rgba(255, 255, 255, .12); border: 1px solid rgba(255, 255, 255, .22); border-radius: 16px; padding: 14px; }
.lp-capacity__facts b { display: block; font-size: 17px; }
.lp-capacity__facts span { display: block; font-size: 13px; opacity: .88; line-height: 1.4; }

/* TANAMAN BUAH */
.lp-fruit { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 20px; }
.lp-fruit__item { position: relative; background: #fff; border-radius: 20px; padding: 24px; box-shadow: 0 1px 2px rgba(16, 24, 40, .05), 0 10px 30px rgba(16, 24, 40, .06); overflow: hidden; }
.lp-fruit__num { position: absolute; top: 10px; right: 16px; font-size: 44px; font-weight: 800; color: #e6f5ec; line-height: 1; }
.lp-fruit__icon { color: var(--brand-orange); margin-bottom: 10px; }
.lp-fruit__item h3 { position: relative; font-size: 18px; color: #1f2266; margin: 0 0 8px; }
.lp-fruit__item p { color: #5b6475; line-height: 1.6; margin: 0; font-size: 15px; }
.lp-fruit__program { margin-top: 28px; background: #fff; border-left: 6px solid var(--brand-green); border-radius: 16px; padding: 20px 24px; }
.lp-fruit__program-title { display: flex; gap: 8px; align-items: center; font-weight: 700; color: #1f2266; margin-bottom: 8px; }
.lp-fruit__program-title .v-icon { color: var(--brand-green); }
.lp-fruit__program ol { margin: 0; padding-left: 20px; display: grid; gap: 6px; color: #3b4456; line-height: 1.6; }

/* MEDIA SOSIAL */
.lp-social { background: #009a4a; padding: 12px 0; margin: -56px 0 40px; }
.lp-social__inner { display: flex; flex-wrap: wrap; gap: 8px 22px; align-items: center; justify-content: center; }
.lp-social__title { font-weight: 700; color: #fff; margin-right: 4px; }
.lp-footer .lp-social__link { color: #fff; margin: 0; font-weight: 600; font-size: 14px; }
.lp-social__icon { width: 28px; height: 28px; border-radius: 50%; background: #fff; color: #009a4a; display: inline-grid; place-items: center; }
.lp-social__icon .v-icon { color: #009a4a; }

/* SECTIONS */
.lp-section { padding: 88px 0; scroll-margin-top: 64px; }
.lp-section--tint { background: #f3f8f4; }
.lp-section--navy { background: linear-gradient(160deg, #1f2266, #2e3192); color: #fff; }
.lp-head { text-align: center; max-width: 720px; margin: 0 auto 48px; }
.lp-eyebrow { color: var(--brand-green); font-weight: 700; text-transform: uppercase; letter-spacing: .12em; font-size: 13px; margin-bottom: 10px; }
.lp-head h2, .lp-order h2 { font-size: clamp(26px, 3.4vw, 38px); font-weight: 800; line-height: 1.18; color: #1f2266; margin: 0 0 14px; }
.lp-head p { color: #5b6475; font-size: 17px; line-height: 1.65; margin: 0; }
.lp-head--light h2 { color: #fff; }
.lp-head--light .lp-eyebrow { color: #7ee2a8; }

.lp-problems { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 22px; }
.lp-problem { border: 1px solid #e7ebf0; border-radius: 20px; padding: 26px; display: flex; flex-direction: column; }
.lp-problem__icon { width: 56px; height: 56px; border-radius: 16px; display: grid; place-items: center; margin-bottom: 16px; }
.lp-problem h3 { font-size: 19px; color: #1f2266; margin: 0 0 8px; }
.lp-problem p { color: #5b6475; line-height: 1.6; margin: 0 0 16px; }
.lp-problem__fix { margin-top: auto; display: flex; gap: 8px; align-items: center; font-size: 14px; background: #f3f8f4; padding: 10px 12px; border-radius: 12px; }

.lp-products { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 24px; }
.lp-product { background: #fff; border-radius: 22px; overflow: hidden; display: flex; flex-direction: column;
  box-shadow: 0 1px 2px rgba(16, 24, 40, .05), 0 16px 40px rgba(16, 24, 40, .08); border-top: 6px solid var(--accent); }
.lp-product__img { position: relative; background: linear-gradient(180deg, #eef3ef, #fff); height: 260px; display: grid; place-items: center; overflow: hidden; }
.lp-product__img img { height: 100%; width: auto; max-width: 100%; object-fit: contain; }
.lp-product__weight { position: absolute; top: 14px; right: 14px; background: var(--accent); color: #fff; font-weight: 700; font-size: 13px; padding: 4px 10px; border-radius: 999px; }
.lp-product__body { padding: 22px 22px 24px; display: flex; flex-direction: column; flex: 1; }
.lp-product__tag { color: var(--accent); font-weight: 700; font-size: 12px; text-transform: uppercase; letter-spacing: .08em; }
.lp-product h3 { font-size: 21px; color: #1f2266; margin: 4px 0 8px; }
.lp-product p { color: #5b6475; line-height: 1.6; margin: 0 0 14px; }
.lp-product ul { list-style: none; padding: 0; margin: 0 0 14px; display: grid; gap: 6px; }
.lp-product li { display: flex; gap: 8px; align-items: flex-start; font-size: 15px; }
.lp-product li .v-icon { color: var(--accent); margin-top: 2px; }
.lp-product__body > .v-btn, .lp-partner > .v-btn { flex: 0 0 auto; }
.lp-product__crops { font-size: 13px; color: #5b6475; background: #f6f7f9; border-radius: 10px; padding: 8px 10px; margin-bottom: 18px; display: flex; gap: 6px; align-items: center; }

.lp-steps { list-style: none; padding: 0; margin: 0; display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 22px; counter-reset: s; }
.lp-step { position: relative; background: #fff; border: 1px solid #e7ebf0; border-radius: 20px; padding: 28px 24px 24px; }
.lp-step__num { width: 44px; height: 44px; border-radius: 50%; background: var(--brand-green); color: #fff; font-weight: 800; font-size: 20px; display: grid; place-items: center; margin-bottom: 14px; }
.lp-step:nth-child(2) .lp-step__num { background: var(--brand-navy); }
.lp-step:nth-child(3) .lp-step__num { background: var(--brand-orange); }
.lp-step__when { font-size: 13px; font-weight: 600; color: var(--brand-green); text-transform: uppercase; letter-spacing: .05em; }
.lp-step h3 { font-size: 19px; color: #1f2266; margin: 6px 0 8px; }
.lp-step p { color: #5b6475; line-height: 1.6; margin: 0; }
.lp-crops { margin-top: 44px; text-align: center; }
.lp-crops__title { font-weight: 700; color: #1f2266; margin-bottom: 14px; }
.lp-crops__list { display: flex; flex-wrap: wrap; justify-content: center; gap: 10px; }
.lp-crop { display: inline-flex; align-items: center; gap: 8px; padding: 8px 16px; border-radius: 999px; background: #ecf8f1; color: #007a3a; font-weight: 600; font-size: 14px; }

.lp-partners { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 22px; }
.lp-partner { background: rgba(255, 255, 255, .08); border: 1px solid rgba(255, 255, 255, .16); border-radius: 20px; padding: 26px; display: flex; flex-direction: column; }
.lp-partner .v-icon { color: #7ee2a8; }
.lp-partner h3 { font-size: 20px; margin: 0 0 8px; }
.lp-partner p { opacity: .85; line-height: 1.6; margin: 0 0 20px; }
.lp-app { margin-top: 28px; display: flex; gap: 16px; align-items: center; background: rgba(0, 154, 74, .22); border: 1px solid rgba(126, 226, 168, .35); border-radius: 18px; padding: 18px 22px; line-height: 1.55; }

.lp-order { display: grid; grid-template-columns: 1fr 1fr; gap: 48px; align-items: start; }
.lp-order p { font-size: 17px; line-height: 1.65; }
.lp-order__why { list-style: none; padding: 0; margin: 22px 0 0; display: grid; gap: 12px; }
.lp-order__why li { display: flex; gap: 10px; align-items: center; font-weight: 500; }
.lp-order__card { padding: 26px; border-radius: 22px !important; box-shadow: 0 16px 40px rgba(16, 24, 40, .1) !important; border-top: 6px solid var(--brand-green); }
.lp-order__products { display: flex; flex-wrap: wrap; gap: 0 8px; }

.lp-faq :deep(.v-expansion-panel-title) { font-weight: 600; color: #1f2266; }
.lp-faq :deep(.v-expansion-panel-text) { color: #5b6475; line-height: 1.65; }

/* FOOTER */
.lp-footer { background: #161850; color: rgba(255, 255, 255, .85); padding: 56px 0 24px; }
.lp-footer__grid { display: grid; grid-template-columns: 1.4fr 1fr 1fr 1fr; gap: 32px; }
.lp-footer__title { color: #fff; font-weight: 700; margin-bottom: 12px; }
.lp-footer a { display: flex; align-items: center; gap: 6px; color: rgba(255, 255, 255, .8); text-decoration: none; margin-bottom: 8px; font-size: 15px; }
.lp-footer a:hover { color: #fff; }
.lp-footer__bottom { border-top: 1px solid rgba(255, 255, 255, .12); margin-top: 36px; padding-top: 18px; font-size: 13px; opacity: .7; }

.lp-wa { position: fixed; right: 18px; bottom: 18px; z-index: 30; width: 58px; height: 58px; border-radius: 50%; background: #25d366; color: #fff;
  display: grid; place-items: center; box-shadow: 0 10px 24px rgba(0, 0, 0, .25); text-decoration: none; }

@media (max-width: 959px) {
  .lp-nav__links { display: none; }
  .lp-nav__cta { margin-left: auto; }
  .lp-hero { padding: 104px 0 64px; }
  .lp-hero__grid, .lp-order { grid-template-columns: 1fr; gap: 36px; }
  .lp-hero__badge { left: 12px; bottom: -18px; }
  .lp-problems, .lp-products, .lp-steps, .lp-partners { grid-template-columns: 1fr; }
  .lp-fruit { grid-template-columns: 1fr 1fr; }
  .lp-capacity__inner { grid-template-columns: 1fr; gap: 20px; }
  .lp-footer__grid { grid-template-columns: 1fr 1fr; }
  .lp-section { padding: 64px 0; }
}
@media (max-width: 599px) {
  .lp-wrap { padding: 0 16px; }
  .lp-brand__name { font-size: 17px; }
  .lp-hero__lead { font-size: 16px; }
  .lp-hero__actions .v-btn { flex: 1 1 100%; }
  .lp-footer__grid { grid-template-columns: 1fr; }
  .lp-app { flex-direction: column; align-items: flex-start; }
  .lp-fruit, .lp-capacity__facts { grid-template-columns: 1fr; }
  .lp-social__inner { justify-content: flex-start; }
}
</style>
