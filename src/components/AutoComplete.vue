<script setup lang="ts">
import { ref, toRefs, computed, watch } from 'vue';
import { NAutoComplete } from 'naive-ui';
import Fuse from 'fuse.js';

interface Props {
  modelValue: string;
  options: string[];
}

const props = defineProps<Props>();
const { modelValue, options } = toRefs(props);
const emit = defineEmits(['update:modelValue']);

const searchQuery = ref(modelValue.value);

const fuse = new Fuse<string>(options.value, {
  keys: ['label'],
  threshold: 0.6,
  shouldSort: true,
  isCaseSensitive: false,
  useExtendedSearch: true,
})

watch(options, (newValue) => {
  fuse.setCollection(newValue);
});

const filteredOptions = computed(() => {
  return fuse.search(searchQuery.value).map(({ item }) => item);
});

watch(searchQuery, (newValue) => {
  emit('update:modelValue', newValue);
});
</script>

<template>
  <NAutoComplete v-model:value="searchQuery" :options="filteredOptions" />
</template>
