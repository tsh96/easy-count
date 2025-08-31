<script lang="ts" setup>
import { Icon } from '@iconify/vue/dist/iconify.js';
import { type CustomerRecord } from '../composables/customer-record';


defineProps<{
  record: CustomerRecord
  newRecordIds: Set<number>
  customerNames: string[]
  scrollIntoNewRecord: (record: CustomerRecord, el: any) => void
  saveCustomerRecord: (record: CustomerRecord) => void
  insertBeforeCustomerRecord: (record: CustomerRecord) => void
  insertAfterCustomerRecord: (record: CustomerRecord) => void
  removeCustomerRecord: (id: number) => void
  triggerAutoInvoiceNo: (record: CustomerRecord) => string | undefined
  updateCustomerNames: () => void
  isDuplicated: boolean
}>()
</script>

<template lang="pug">
tr(
  :class="{ 'bg-green-300': newRecordIds.has(record.id || 0), 'bg-yellow-100': !record.chequeNo }"
  :ref="(el) => scrollIntoNewRecord(record, el)"
  )
  td
    n-date-picker(v-model:value="record.invoiceDate" size="small" @update:value="saveCustomerRecord(record)" placement="right")
  td
    n-popover(trigger="hover" :disabled="(!!record.invoiceNo || !triggerAutoInvoiceNo(record)) && !isDuplicated")
      template(#trigger)
        n-input.font-mono(
          v-model:value="record.invoiceNo"
          :input-props="{ class: isDuplicated ? '!text-red-400' : '' }"
          size="small"
          @update:value="saveCustomerRecord(record)"
        )
      template(v-if="isDuplicated")
        div Invoice No Duplicated
      template(v-else)
        n-button(text @click="record.invoiceNo = triggerAutoInvoiceNo(record) || ''; saveCustomerRecord(record)" size="small")
          .font-mono {{ triggerAutoInvoiceNo(record) }}
  td
    auto-complete.font-mono(
      v-model="record.customerName"
      :options="customerNames"
      @blur="updateCustomerNames()"
      @update:model-value="saveCustomerRecord(record)"
      size="small"
    )
  td
    n-input-number.font-mono.text-right(
      v-model:value="record.invoiceAmount"
      size="small"
      :input-props="{ class: record.invoiceAmount == 0 ? '!text-red-400' : '' }"
      :show-button="false"
      :precision="2"
      @update:value="saveCustomerRecord(record)"
      ) 
  td
    n-date-picker(v-model:value="record.chequeDate" size="small" @update:value="saveCustomerRecord(record)" placement="right")
  td 
    n-input.font-mono(v-model:value="record.chequeNo" size="small" @update:value="saveCustomerRecord(record)")
  td
    n-popover(trigger="hover" :disabled="!!record.chequeAmount || !record.chequeNo" :show-arrow="false")
      template(#trigger)
        n-input-number.font-mono.text-right(
          v-model:value="record.chequeAmount"
          size="small"
          :input-props="{ class: record.chequeAmount == 0 ? '!text-gray-400' : (record.chequeAmount || 0) < (record.invoiceAmount || 0) ? '!text-red-400' : '!text-green-600' }"
          :show-button="false"
          :precision="2"
          @update:value="saveCustomerRecord(record)"
          )
      n-button(text @click="record.chequeAmount = record.invoiceAmount; saveCustomerRecord(record)" size="small")
        .font-mono {{ record.invoiceAmount?.toFixed(2) }}
  td
    n-input.font-mono(v-model:value="record.remark" size="small" @update:value="saveCustomerRecord(record)")
  td
    .flex.space-x-2
      .flex.place-content-center
        n-tooltip
          div Insert Before
          template(#trigger)
            n-button(text type="success" @click="insertBeforeCustomerRecord(record)"  tabindex="-1")
              Icon(icon="tabler:row-insert-top")
      .flex.place-content-center
        n-tooltip
          div Insert After
          template(#trigger)
            n-button(text type="success" @click="insertAfterCustomerRecord(record)"  tabindex="-1")
              Icon(icon="tabler:row-insert-bottom")
      .flex.place-content-center(v-if="record.id")
        n-tooltip
          div Delete
          template(#trigger)
            n-popconfirm(@positive-click="removeCustomerRecord(record.id)") Are you sure to delete this record?
              template(#trigger)
                n-button(text type="error" tabindex="-1")
                  Icon(icon="mdi:delete")
</template>

<style scoped lang="scss">

</style>