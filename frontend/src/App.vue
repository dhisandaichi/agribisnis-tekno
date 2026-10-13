<template>
  <v-app>
    <v-app-bar color="primary" v-if="user">
      <v-toolbar-title>SFA Agritec</v-toolbar-title>
      <v-spacer></v-spacer>
      <span class="text-body-2 mr-2 d-none d-sm-inline">{{ user.full_name || user.email }}</span>
      <v-btn icon aria-label="Keluar" @click="logout">
        <v-icon>mdi-logout</v-icon>
      </v-btn>
    </v-app-bar>

    <v-main>
      <v-container>
        <router-view />
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
export default {
  computed: {
    user() {
      return this.$store.state.user
    }
  },
  methods: {
    async logout() {
      await this.$store.dispatch('logout')
      this.$router.push('/login')
    }
  }
}
</script>
