<template>
  <v-app :theme="theme">
    <v-app-bar class="px-3">
      <v-spacer />

      <v-btn
        :prepend-icon="theme === 'light' ? 'mdi-weather-sunny' : 'mdi-weather-night'"
        slim
        text="Toggle Theme"
        @click="onClick"
      />
    </v-app-bar>

    <v-main>
      <v-container class="fill-height pa-3" fluid style="min-height: 0;">
        <v-responsive class="mx-auto fill-height" style="min-height: 0;">
          <router-view />
        </v-responsive>
      </v-container>
    </v-main>
  </v-app>
</template>

<script lang="ts" setup>
  import { ref } from 'vue'

  const theme = ref('dark')

  function onClick () {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
  }
</script>
