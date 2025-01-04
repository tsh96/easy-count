<script setup lang="ts">
import { NConfigProvider } from 'naive-ui';
import { lightTheme, type GlobalThemeOverrides } from 'naive-ui'

const themeOverrides: GlobalThemeOverrides = {
  common: {
    borderRadius: '8px',
  },
  Table: {
    thColor: '#F3F3F3',
  }
}


</script>

<template>
  <NConfigProvider :theme="lightTheme" :theme-overrides="themeOverrides">
    <RouterView></RouterView>
  </NConfigProvider>
</template>

<style scoped>
</style>
