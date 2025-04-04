<script lang="ts" setup>
import { type Transaction } from '../composables/personal-record';

defineProps<{
  record: Transaction
  newRecordIds: Set<number>
  accumulated: Record<number, number>
  scrollIntoNewRecord: (record: Transaction, el: any) => void
  descriptions: string[]
  saveTransaction: (record: Transaction) => void
  updateDescriptions: () => void
  insertBeforeTransaction: (record: Transaction) => void
  insertAfterTransaction: (record: Transaction) => void
  removeTransaction: (id: number) => void
}>()

</script>

<template lang="pug">
tr(
  :class="{ 'bg-green-300': newRecordIds.has(record.id || 0) }"
  :ref="(el) => scrollIntoNewRecord(record, el)"
  )
  td
    n-date-picker(v-model:value="record.date" size="small" @update:value="saveTransaction(record)")
  td
    auto-complete.font-mono(v-model="record.description" size="small" @update:model-value="saveTransaction(record)" :options="descriptions" @blur="updateDescriptions()")
  td
    n-input-number.font-mono.text-right(
      v-model:value="record.credit"
      size="small"
      :input-props="{ class: record.credit == 0 ? '!text-gray-400' : '' }"
      :show-button="false"
      :precision="2"
      @update:value="saveTransaction(record)"
      ) 
  td
    n-input-number.font-mono.text-right(
      v-model:value="record.debit"
      size="small"
      :input-props="{ class: record.debit == 0 ? '!text-gray-400' : '' }"
      :show-button="false"
      :precision="2"
      @update:value="saveTransaction(record)"
      )
  td
    .font-mono.text-right(:class="{ 'text-red-500': accumulated[record.id || 0] < 0, 'text-green-500': accumulated[record.id || 0] > 0 }")
      | {{ accumulated[record.id || 0] }}
  td
    .flex.space-x-2
      .flex.place-content-center
        n-tooltip
          div Insert Before
          template(#trigger)
            n-button(text type="success" @click="insertBeforeTransaction(record)")
              Icon(icon="tabler:row-insert-top")
      .flex.place-content-center
        n-tooltip
          div Insert After
          template(#trigger)
            n-button(text type="success" @click="insertAfterTransaction(record)")
              Icon(icon="tabler:row-insert-bottom")
      .flex.place-content-center(v-if="record.id")
        n-tooltip
          div Delete
          template(#trigger)
            n-popconfirm(@positive-click="removeTransaction(record.id)") Are you sure to delete this record?
              template(#trigger)
                n-button(text type="error")
                  Icon(icon="mdi:delete")
</template>