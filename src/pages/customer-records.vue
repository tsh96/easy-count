<script lang="ts" setup>
import { Icon } from '@iconify/vue/dist/iconify.js';
import { useStorage } from '@vueuse/core';
import { endOfDay, format } from 'date-fns';
import { type ComponentPublicInstance, computed, onMounted, ref, toRaw, watch, watchEffect } from 'vue';
import AutoComplete from '../components/AutoComplete.vue';
import AiKeyin from '../components/AiKeyin.vue';
import { backup, type CustomerRecord, CustomerRecordType, db, restore } from '../composables/customer-record';
import { type FormInst } from 'naive-ui';
import Header from '../components/Header.vue';
import { migrateOldCustomerRecord } from '../composables/old-customer-record';
import { map, sumBy, uniq } from 'lodash';
import { formatNumber } from '../composables/format-number';

const showQrCode = ref(false);

const customerRecordTypeOptions = Object.values(CustomerRecordType).map(value => ({ label: value, value }))

const customerRecordType = ref<CustomerRecordType>(CustomerRecordType.Private);
const customerRecords = ref<CustomerRecord[]>([]);
const yearFilter = useStorage('yearFilter', new Date().getFullYear())
const newRecordIds = ref(new Set<number>())
watch(newRecordIds, () => setTimeout(() => newRecordIds.value = new Set(), 5000))

const filter = ref<{
  invoiceDate?: [number, number]
  invoice?: string
  customerName?: string
  invoiceAmount?: number
  chequeDate?: [number, number]
  chequeNo?: string
  chequeAmount?: number
  remark?: string
  hasNoCheque?: boolean
}>({})

const dbCustomerRecords = computed(() => {
  return customerRecordType.value === CustomerRecordType.Private ? db.privateCustomerRecords : db.governmentCustomerRecords
})

async function loadCustomerRecords() {
  const startOfYear = new Date(`${yearFilter.value}-01-01`).getTime()
  const startOfNextYear = new Date(`${yearFilter.value + 1}-01-01`).getTime()

  customerRecords.value = await dbCustomerRecords.value
    .where("invoiceDate").between(startOfYear, startOfNextYear, true, false)
    .toArray();

  sortCustomerRecords()
}

watchEffect(loadCustomerRecords)

const customerSummary = computed(() => {
  const summary: Record<string, number> = {}
  customerRecords.value.forEach(record => {
    const customerName = record.customerName || ''
    if (!summary[customerName]) {
      summary[customerName] = 0
    }
    summary[customerName] += record.invoiceAmount || 0
  })

  return map(summary, (value, key) => {
    return { customerName: key, total: +value.toFixed(2) }
  }).sort((a, b) => {
    return a.customerName.localeCompare(b.customerName, undefined, { numeric: true, sensitivity: 'base' })
  })
})

const grandTotal = computed(() => {
  return sumBy(customerSummary.value, summary => summary.total)
})

const showCustomerSummary = ref(false)

const filteredCustomerRecords = computed(() => {
  return customerRecords.value.filter(record => {
    return (!filter.value.invoiceDate || !record.invoiceDate || (record.invoiceDate >= filter.value.invoiceDate[0] && record.invoiceDate <= endOfDay(filter.value.invoiceDate[1]).getTime())) &&
      (!filter.value.invoice || record.invoiceNo.startsWith(filter.value.invoice)) &&
      (!filter.value.customerName || record.customerName === filter.value.customerName) &&
      (!filter.value.invoiceAmount || record.invoiceAmount === filter.value.invoiceAmount) &&
      (!filter.value.chequeDate || (record.chequeDate && record.chequeDate >= filter.value.chequeDate[0] && record.chequeDate <= filter.value.chequeDate[1])) &&
      (!filter.value.chequeNo || record.chequeNo.startsWith(filter.value.chequeNo)) &&
      (!filter.value.chequeAmount || record.chequeAmount === filter.value.chequeAmount) &&
      (!filter.value.remark || record.remark.includes(filter.value.remark)) &&
      (!filter.value.hasNoCheque || !record.chequeNo)
  })
})

const page = ref(1)
const pageSize = 25
const pageCount = computed(() => Math.ceil(filteredCustomerRecords.value.length / pageSize))
watchEffect(() => {
  if (page.value > pageCount.value) {
    page.value = 1
  }
})
const paginatedCustomerRecords = computed(() => {
  const start = (page.value - 1) * pageSize
  const end = start + pageSize
  console.log(start, end)
  return filteredCustomerRecords.value.slice(start, end)
})

const duplicatedInvoiceNos = computed(() => {
  const counts = new Map<string, number>()
  customerRecords.value.forEach(record => {
    if (record.invoiceNo) {
      counts.set(record.invoiceNo, (counts.get(record.invoiceNo) || 0) + 1)
    }
  })
  return new Set(Array.from(counts.entries()).filter(([_, count]) => count > 1).map(([invoiceNo, _]) => invoiceNo))
})

const totalInvoiceAmount = computed(() => sumBy(filteredCustomerRecords.value, 'invoiceAmount') || 0)

const totalChequeAmount = computed(() => sumBy(filteredCustomerRecords.value, 'chequeAmount') || 0)

const customerNames = computed(() => {
  return uniq(customerRecords.value.map(record => record.customerName))
})

onMounted(async () => {
  await migrateOldCustomerRecord()
})

const customerNameOptions = computed(() => {
  return customerNames.value.map(name => ({ label: name, value: name }))
})

function newCustomerRecord(override?: Partial<CustomerRecord>): CustomerRecord {
  const customerRecord: CustomerRecord = {
    invoiceDate: override?.invoiceDate,
    chequeAmount: override?.chequeAmount,
    chequeNo: override?.chequeNo ?? '',
    customerName: override?.customerName ?? (filter.value.customerName || ''),
    invoiceNo: override?.invoiceNo ?? '',
    invoiceAmount: override?.invoiceAmount,
    remark: override?.remark ?? ''
  }

  return customerRecord
}

async function addCustomerRecord() {
  const customerRecord: CustomerRecord = newCustomerRecord()

  await dbCustomerRecords.value.add(customerRecord)
  newRecordIds.value = new Set([customerRecord.id!])

  customerRecords.value.push(customerRecord)
}

async function removeCustomerRecord(id: number) {
  await dbCustomerRecords.value.delete(id)

  customerRecords.value = customerRecords.value.filter(record => record.id !== id)
}

async function insertAfterCustomerRecord(record: CustomerRecord) {
  const index = customerRecords.value.findIndex(r => r.id === record.id)

  const customerRecord: CustomerRecord = newCustomerRecord({
    invoiceDate: customerRecords.value[index].invoiceDate,
  })

  await dbCustomerRecords.value.add(customerRecord)
  newRecordIds.value = new Set([customerRecord.id!])

  customerRecords.value.splice(index + 1, 0, customerRecord)
}

async function insertBeforeCustomerRecord(record: CustomerRecord) {
  const index = customerRecords.value.findIndex(r => r.id === record.id)

  const customerRecord: CustomerRecord = newCustomerRecord({
    invoiceDate: customerRecords.value[index].invoiceDate,
  })

  await dbCustomerRecords.value.add(customerRecord)

  newRecordIds.value = new Set([customerRecord.id!])

  customerRecords.value.splice(index, 0, customerRecord)
}

async function saveCustomerRecord(record: CustomerRecord) {
  const original = await dbCustomerRecords.value.get(record.id)
  if (!original) return

  await dbCustomerRecords.value.put(toRaw(record))
}

const restoring = ref(false)

function restoreFromFile() {
  // open file dialog
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.json'
  restoring.value = true

  input.onchange = async () => {
    const file = input.files?.[0]
    if (!file) return

    const text = await file.text()
    await restore(text)
    await loadCustomerRecords()

    restoring.value = false
  }

  input.click()
}

const backingUp = ref(false)

async function backupToFile() {
  backingUp.value = true

  const data = await backup()

  // download the data
  const blob = new Blob([data], { type: 'application/json' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `customer_records_${format(new Date(), "yyyy-MM-dd'T'HHmmss")}.json`;
  a.click();
  backingUp.value = false
}

async function aiSubmitData(records: CustomerRecord[]) {
  await dbCustomerRecords.value.bulkAdd(toRaw(records), { allKeys: true })

  await loadCustomerRecords()
}

function scrollIntoNewRecord(record: CustomerRecord, el: Element | ComponentPublicInstance | null) {
  if (newRecordIds.value.has(record.id || 0)) {
    if (el instanceof Element)
      el?.scrollIntoView({ behavior: 'smooth', block: 'center' })
  }
}

const replacingCustomerName = ref(false)

const replaceCustomerNameForm = ref<{
  oldName: string
  newName: string
}>({
  oldName: '',
  newName: ''
})

async function replaceCustomerName(oldName: string, newName: string) {
  replacingCustomerName.value = true
  await dbCustomerRecords.value.where('customerName').equals(oldName).modify({ customerName: newName })
  customerRecords.value = []
  await loadCustomerRecords()
  replacingCustomerName.value = false
}

const replaceCustomerNameFormRef = ref<FormInst>()

function submitReplaceCustomerName() {
  replaceCustomerNameFormRef.value?.validate((errors) => {
    if (!errors) {
      replaceCustomerName(replaceCustomerNameForm.value.oldName, replaceCustomerNameForm.value.newName)
    }
  })
}

function sortCustomerRecords() {
  customerRecords.value.sort((a, b) => {
    return (a.invoiceDate || Infinity) - (b.invoiceDate || Infinity) || a.invoiceNo.localeCompare(b.invoiceNo, undefined, { numeric: true, sensitivity: 'base' })
  })
}

const autoInvoiceNo = ref('')

function triggerAutoInvoiceNo(record: CustomerRecord) {
  const index = customerRecords.value.findIndex(r => r.id === record.id)
  
  autoInvoiceNo.value = ''
  if (index === 0) return
  if (index > customerRecords.value.length - 1) return
  const lastRecord = customerRecords.value[index - 1]

  const lastInvoiceNo = lastRecord.invoiceNo
  const lastInvoiceNoNumber = parseInt(lastInvoiceNo.replace(/\D/g, ''))
  const lastInvoiceNoPrefix = lastInvoiceNo.replace(/\d/g, '')

  if (isNaN(lastInvoiceNoNumber)) return

  const nextInvoiceNo = `${lastInvoiceNoPrefix}${lastInvoiceNoNumber + 1}`
  autoInvoiceNo.value = nextInvoiceNo
  return nextInvoiceNo
}

</script>

<template lang="pug">
.h-screen.w-screen.overflow-hidden.p-4.flex.flex-col
  Header(page="customer-records")
    .w-52.font-bold
      n-select(v-model:value="customerRecordType" :options="customerRecordTypeOptions")
    n-divider(vertical)
    .flex.items-center.space-x-2
      b Year:
      n-input-number.w-32(v-model:value="yearFilter" size="small" :precision="0")
    n-button(type="info" @click="showCustomerSummary = true")
      .flex.items-center.space-x-2
        Icon(icon="mdi:calendar-month")
        div Yearly Summary
    .flex-grow
    n-button(type="success" @click="backupToFile()")
      .flex.items-center.space-x-2
        Icon(icon="mdi:download")
        div Backup
    n-button(type="info" @click="restoreFromFile()")
      .flex.items-center.space-x-2
        Icon(icon="mdi:restore")
        div Restore
    n-divider(vertical)
    //- n-button(type="success" @click="showQrCode = true")
      .flex.items-center.space-x-2
        Icon(icon="mdi:qrcode")
        div AI Keyin
  .flex.flex-col.flex-grow.overflow-hidden.relative.border
    .flex-grow.overflow-hidden
      n-scrollbar
        n-table.pr-2(size="small" :single-line="false" style="overflow-y: visible;" :bordered="false")
          thead.sticky.top-0.z-10
            tr
              th.w-52
                .items-center.flex.space-x-1
                  span Date
                  n-tooltip Reload and Sort
                    template(#trigger)
                      n-button(text @click="loadCustomerRecords()")
                        Icon(icon="mdi:refresh")
              th.w-32 Invoice
              th
                .flex.items-center Customer Name
                  n-popover(trigger="click")
                    template(#trigger)
                      n-button(type="warning" text size="small")
                        Icon.mx-2(icon="mdi:swap-horizontal")
                    n-form(:model="replaceCustomerNameForm" ref="replaceCustomerNameFormRef" size="small" style="width: 800px;")
                      n-form-item(label="Replace Customer Name" :rule="{ required: true }" path="oldName" )
                        auto-complete.font-mono(v-model="replaceCustomerNameForm.oldName" size="small" clearable :options="customerNames")
                      n-form-item(label="With" :rule="{ required: true }" path="newName")
                        auto-complete.font-mono(v-model="replaceCustomerNameForm.newName" size="small" clearable :options="customerNames")
                      n-button(type="success" size="small" block @click="submitReplaceCustomerName()" :loading="replacingCustomerName")
                        .flex.items-center.space-x-2
                          Icon(icon="mdi:swap-horizontal")
                          div Replace
              th.w-32 Amount
              th.w-52 Date
              th.w-32 
                .flex.items-center.space-x-1
                  span Cheque
                  n-tooltip Show Records without Cheque
                    template(#trigger)
                      n-switch(v-model:value="filter.hasNoCheque")
                        template(#icon)
                          Icon(icon="material-symbols:receipt-long-off-rounded")
              th.w-32 Amount
              th.w-32 Remark
              th.w-8
                .flex.place-content-center
                  n-tooltip Add New Record
                    template(#trigger)
                      n-button(text type="success" @click="addCustomerRecord()")
                        Icon(icon="mdi:add")
            tr.th-border-b
              th
                n-date-picker.text-xs(v-model:value="filter.invoiceDate" size="small" type="daterange" format="YY-MM-dd" clearable)
              th
                n-input(v-model:value="filter.invoice" size="small" clearable)
              th
                n-select(v-model:value="filter.customerName" size="small" clearable :options="customerNameOptions" filterable)
              th(class="!text-right")
                span.font-mono.font-bold {{ formatNumber(totalInvoiceAmount) }}
              th
                n-date-picker(v-model:value="filter.chequeDate" size="small"  format="YY-MM-dd" type="daterange" clearable)
              th
                n-input(v-model:value="filter.chequeNo" size="small" clearable)
              th(class="!text-right")
                span.font-mono.font-bold {{ formatNumber(totalChequeAmount) }}
              th
                n-input(v-model:value="filter.remark" size="small" clearable)
              th
          tbody
            TransitionGroup(name="list")
              CustomerRecordTr(
                v-for="(record) in paginatedCustomerRecords"
                :key="record.id"
                :record="record"
                :new-record-ids="newRecordIds"
                :accumulated="customerRecords"
                :scroll-into-new-record="scrollIntoNewRecord"
                :customer-names="customerNames"
                :save-customer-record="saveCustomerRecord"
                :insert-before-customer-record="insertBeforeCustomerRecord"
                :insert-after-customer-record="insertAfterCustomerRecord"
                :remove-customer-record="removeCustomerRecord"
                :trigger-auto-invoice-no="triggerAutoInvoiceNo"
                :is-duplicated="duplicatedInvoiceNos.has(record.invoiceNo || '')"
              )
            tr.h-96
    .flex.justify-center.my-2
      n-pagination(v-model:page="page" :page-count="pageCount")
        template(#suffix) Total {{ filteredCustomerRecords.length }} records
n-modal(v-model:show="showQrCode" :mask-closable="false" preset="dialog" :show-icon="false" style="width: calc(100% - 48px);")
  AiKeyin(:customer-names="customerNames" @submit="aiSubmitData")

n-modal(v-model:show="showCustomerSummary" :mask-closable="false" preset="dialog" :show-icon="false" style="width: calc(100% - 48px);")
  .p-4
    .text-lg.font-bold Customer Summary
    n-table(size="small" :single-line="false")
      thead.sticky.top-0.z-10
        tr
          th Customer Name
          th Total Amount 
            .font-mono.float-right.font-bold {{ formatNumber(grandTotal) }}
      tbody
        TransitionGroup(name="list")  
          tr(
            v-for="summary in customerSummary"
            :key="summary.customerName"
            class="bg-green-300 list-move"
          )
            td {{ summary.customerName }}
            td.font-mono.text-right {{ formatNumber(summary.total) }}
</template>

<style lang="scss">
.bg-green-300,
.bg-yellow-100 {
  td {
    background-color: inherit;
  }
}

tr:hover {
  td {
    background-color: rgba(0, 0, 0, 0.1);
  }
}

tr.th-border-b {
  th {
    border-bottom: 2px solid #c1c1c2;
  }
}

.list-move,
/* apply transition to moving elements */
.list-enter-active,
.list-leave-active {
  transition: all 0.5s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(-300px);
}
</style>