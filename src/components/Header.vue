<script setup lang="ts">
import { Icon } from '@iconify/vue/dist/iconify.js';
import { defineProps } from 'vue'

defineProps<{
  title: string
}>()
</script>

<template lang="pug">
.flex.items-center.space-x-4.mb-4
  .flex
    .flex.space-x-2.items-center
      .text-2xl.font-bold {{ title }}
      n-popover(trigger="click")
        router-link(to="/customer-records")
          n-button(quaternary block) Customer Records
        router-link(to="/personal-records")
          n-button(quaternary block) Personal Records
        template(#trigger)
          n-button(size="small")
            Icon(icon="material-symbols:arrow-drop-down")
  n-divider(vertical)
  slot
</template>