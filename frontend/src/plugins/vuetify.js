import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import { createVuetify } from 'vuetify'

// Komponen di-import otomatis oleh vite-plugin-vuetify (tree-shaking)
export default createVuetify({
  theme: {
    defaultTheme: 'light',
    themes: {
      light: {
        colors: {
          primary: '#4CAF50', // Agribisnis Green Theme
          secondary: '#689F38',
          error: '#E53935',
          info: '#1E88E5',
          success: '#43A047',
          warning: '#FFA000',
        },
      },
    },
  },
})
