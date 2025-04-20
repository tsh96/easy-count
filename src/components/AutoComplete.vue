<script setup lang="ts">
import { ref, toRefs, computed, watch } from 'vue';
import { NAutoComplete } from 'naive-ui';

interface Props {
  modelValue: string;
  options: string[];
}

const props = defineProps<Props>();
const { modelValue, options } = toRefs(props);
const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void;
}>();

const searchQuery = ref(modelValue.value);

const filteredOptions = computed(() => {
  if (searchQuery.value === '') {
    return options.value;
  }
  return options.value.filter((option) => {
    return option.toLowerCase().includes(searchQuery.value.toLowerCase());
  });
});

watch(searchQuery, (newValue) => {
  emit('update:modelValue', newValue);
});
</script>

<template>
  <NAutoComplete v-model:value="searchQuery" :options="filteredOptions" />
</template>
