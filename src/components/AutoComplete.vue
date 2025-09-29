<script setup lang="ts">
import { ref, toRefs, computed, watch, h } from 'vue';
import type { VNodeChild } from 'vue';
import { NAutoComplete, NEllipsis, SelectGroupOption, SelectOption } from 'naive-ui';
import { RenderLabel } from 'naive-ui/es/_internal/select-menu/src/interface';

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

  const lowerCaseQuery = searchQuery.value.toLowerCase();

  return options.value.filter((option) =>
    option.toLowerCase().includes(lowerCaseQuery),
  );
});

const renderLabel: RenderLabel = (option: SelectOption | SelectGroupOption): VNodeChild => {
  const rawLabel = option.label ?? option.value ?? '';
  const label = typeof rawLabel === 'string' ? rawLabel : String(option.value ?? '');

  return h(
    NEllipsis,
    {
      tooltip: true,
    },
    {
      default: () => label,
    },
  );
};

watch(searchQuery, (newValue) => {
  emit('update:modelValue', newValue);
});
</script>

<template>
  <NAutoComplete
    v-model:value="searchQuery"
    :options="filteredOptions"
    :render-label="renderLabel"
  />
</template> 
