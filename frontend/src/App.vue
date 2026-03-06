<template>
  <v-app>
    <v-navigation-drawer v-if="isAuthenticated" v-model="drawer" app color="surface" elevation="3" rounded="tr-xl br-xl">
      <div class="pa-4 d-flex align-center justify-center">
        <v-avatar color="primary" size="64" class="mr-3 elevation-2">
          <v-icon dark size="36">mdi-leaf</v-icon>
        </v-avatar>
        <div>
          <div class="text-h6 font-weight-bold primary--text">SFA Agritec</div>
          <div class="text-caption text-medium-emphasis">{{ userRoleFormatted }}</div>
        </div>
      </div>
      <v-divider class="mb-2"></v-divider>
      
      <v-list nav dense class="px-2">
        <v-list-item v-for="item in menuItems" :key="item.title" :to="item.path" :prepend-icon="item.icon" rounded="lg" active-color="primary" class="mb-1" link>
          <v-list-item-title class="font-weight-medium">{{ item.title }}</v-list-item-title>
        </v-list-item>
      </v-list>
      
      <template v-slot:append>
        <div class="pa-4">
          <v-btn block color="error" variant="tonal" rounded="pill" @click="logout" prepend-icon="mdi-logout">
            Keluar
          </v-btn>
        </div>
      </template>
    </v-navigation-drawer>

    <v-app-bar app v-if="isAuthenticated" color="transparent" elevation="0" class="px-4">
      <v-app-bar-nav-icon @click="drawer = !drawer" color="primary"></v-app-bar-nav-icon>
      <v-spacer></v-spacer>
      <v-badge color="error" :content="unreadCount" :model-value="unreadCount > 0" class="mr-4">
        <v-btn icon color="primary" variant="tonal" size="small">
          <v-icon>mdi-bell-outline</v-icon>
        </v-btn>
      </v-badge>
      <v-avatar size="40" color="primary" class="elevation-2 cursor-pointer">
        <span class="text-white text-button">{{ initials }}</span>
      </v-avatar>
    </v-app-bar>

    <v-main class="bg-grey-lighten-4">
      <div class="pa-4 pa-md-6 h-100">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </div>
    </v-main>
  </v-app>
</template>

<script>
export default {
  data() {
    return {
      drawer: true
    }
  },
  computed: {
    isAuthenticated() {
      return this.$store.getters.isAuthenticated
    },
    unreadCount() {
      return this.$store.getters.unreadNotifications
    },
    userRole() {
      return this.$store.getters.userRole || ''
    },
    userRoleFormatted() {
      const role = this.userRole
      if (!role) return ''
      return role.charAt(0).toUpperCase() + role.slice(1)
    },
    initials() {
      const email = this.$store.state.user?.email || 'A'
      return email.charAt(0).toUpperCase()
    },
    menuItems() {
      const role = this.userRole
      if (role === 'sales') {
        return [
          { title: 'Kunjungan Petani', icon: 'mdi-map-marker-path', path: '/sales' },
          { title: 'Edukasi Petani', icon: 'mdi-school', path: '/farmer' }
        ]
      } else if (role === 'supervisor') {
        return [
          { title: 'Dashboard Analytics', icon: 'mdi-view-dashboard', path: '/supervisor' },
          { title: 'Data Distributor', icon: 'mdi-storefront', path: '/distributor' }
        ]
      } else if (role === 'farmer') {
        return [
          { title: 'Edukasi & Info', icon: 'mdi-book-open-variant', path: '/farmer' }
        ]
      } else if (role === 'distributor') {
        return [
          { title: 'Manajemen Stok', icon: 'mdi-package-variant-closed', path: '/distributor' }
        ]
      }
      return []
    }
  },
  watch: {
    isAuthenticated(val) {
      if (val) this.$store.dispatch('fetchNotifications')
    }
  },
  mounted() {
    if (this.isAuthenticated) this.$store.dispatch('fetchNotifications')
  },
  methods: {
    logout() {
      this.$store.dispatch('logout')
      this.$router.push('/login')
    }
  }
}
</script>

<style>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
</style>
