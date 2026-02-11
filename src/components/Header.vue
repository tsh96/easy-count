<script setup lang="ts">
import { Icon } from '@iconify/vue/dist/iconify.js';
import { defineProps, computed } from 'vue'
import { useRouter } from 'vue-router'
import { logout, currentUser } from '../composables/auth'

const router = useRouter()

const props = defineProps<{
  page: 'customer-records' | 'personal-records'
}>()

const title = computed(() => {
  switch (props.page) {
    case 'customer-records':
      return 'Customer Records'
    case 'personal-records':
      return 'Personal Records'
  }
})

function handleLogout() {
  logout()
  router.push('/auth')
}
</script>

<template lang="pug">
.flex.items-center.space-x-4.mb-4
  .flex
    .flex.space-x-2.items-center
      .text-2xl.font-bold {{ title }}
      n-popover(trigger="click")
        router-link(to="/customer-records" v-if="props.page !== 'customer-records'")
          n-button(quaternary block) Customer Records
        router-link(to="/personal-records" v-if="props.page !== 'personal-records'")
          n-button(quaternary block) Personal Records
        template(#trigger)
          n-button(size="small")
            Icon(icon="material-symbols:arrow-drop-down")
  n-divider(vertical)
  slot
  .flex-grow
  .flex.items-center.space-x-2(v-if="currentUser")
    span.text-sm.text-gray-600 {{ currentUser.email }}
    n-button(size="small" type="error" @click="handleLogout")
      .flex.items-center.space-x-1
        Icon(icon="mdi:logout")
        span Logout
</template>