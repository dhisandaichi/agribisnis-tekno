import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import { createVuetify } from 'vuetify'

// Warna identitas PT. Agribisnis Technology Indonesia / Agrina 99:
// hijau daun dan biru indigo dari logo, oranye dari aksen "PREMIUM".
export const BRAND = {
  green: '#009A4A',
  greenDark: '#007A3A',
  navy: '#2E3192',
  navyDark: '#1F2266',
  orange: '#F26B21',
}

// Komponen di-import otomatis oleh vite-plugin-vuetify (tree-shaking)
export default createVuetify({
  theme: {
    defaultTheme: 'light',
    themes: {
      light: {
        colors: {
          primary: BRAND.green,
          secondary: BRAND.navy,
          accent: BRAND.orange,
          background: '#F6F8F5',
          error: '#D32F2F',
          info: '#1E6FD9',
          success: '#1B9E4B',
          warning: '#F59E0B',
        },
      },
    },
  },
  defaults: {
    VBtn: { rounded: 'lg' },
    VCard: { rounded: 'lg' },
    VTextField: { rounded: 'lg' },
    VSelect: { rounded: 'lg' },
  },
})
