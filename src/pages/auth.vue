<script lang="ts" setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { login, register } from '../composables/auth';

const router = useRouter();

const isLogin = ref(true);
const email = ref('');
const password = ref('');
const confirmPassword = ref('');
const loading = ref(false);
const error = ref('');

async function handleSubmit() {
  error.value = '';
  
  // Validation
  if (!email.value || !password.value) {
    error.value = 'Please fill in all fields';
    return;
  }

  if (!isLogin.value && password.value !== confirmPassword.value) {
    error.value = 'Passwords do not match';
    return;
  }

  if (!isLogin.value && password.value.length < 8) {
    error.value = 'Password must be at least 8 characters long';
    return;
  }

  loading.value = true;

  try {
    if (isLogin.value) {
      await login(email.value, password.value);
    } else {
      await register(email.value, password.value);
    }
    // Redirect to home page after successful login/register
    router.push('/');
  } catch (e: any) {
    error.value = e.message || 'Authentication failed';
  } finally {
    loading.value = false;
  }
}

function toggleMode() {
  isLogin.value = !isLogin.value;
  error.value = '';
  confirmPassword.value = '';
}
</script>

<template lang="pug">
.min-h-screen.flex.items-center.justify-center.bg-gray-50.p-4
  .w-full.max-w-md
    .bg-white.rounded-lg.shadow-lg.p-8
      .text-center.mb-8
        h1.text-3xl.font-bold.text-gray-900 Easy Count
        p.text-gray-600.mt-2 {{ isLogin ? 'Sign in to your account' : 'Create a new account' }}
      
      form(@submit.prevent="handleSubmit")
        .space-y-4
          div
            label.block.text-sm.font-medium.text-gray-700.mb-1 Email
            n-input(
              v-model:value="email"
              placeholder="you@example.com"
              size="large"
              :disabled="loading"
            )
          
          div
            label.block.text-sm.font-medium.text-gray-700.mb-1 Password
            n-input(
              v-model:value="password"
              type="password"
              placeholder="••••••••"
              size="large"
              :disabled="loading"
            )
          
          div(v-if="!isLogin")
            label.block.text-sm.font-medium.text-gray-700.mb-1 Confirm Password
            n-input(
              v-model:value="confirmPassword"
              type="password"
              placeholder="••••••••"
              size="large"
              :disabled="loading"
            )
          
          n-alert(v-if="error" type="error" :show-icon="false") {{ error }}
          
          n-button(
            type="primary"
            block
            size="large"
            :loading="loading"
            attr-type="submit"
          ) {{ isLogin ? 'Sign In' : 'Create Account' }}
        
        .mt-6.text-center
          n-button(
            text
            type="primary"
            @click="toggleMode"
            :disabled="loading"
          ) {{ isLogin ? 'Need an account? Sign up' : 'Already have an account? Sign in' }}
</template>

<style scoped>
/* Any additional styles if needed */
</style>
