<script lang="ts" setup>
import { Icon } from '@iconify/vue/dist/iconify.js';
import { useStorage } from '@vueuse/core';
import { endOfDay, format, startOfDay } from 'date-fns';
import { type ComponentPublicInstance, computed, onMounted, ref, toRaw, watch, watchEffect } from 'vue';
import AutoComplete from '../components/AutoComplete.vue';
import AiKeyin from '../components/AiKeyin.vue';
import { backup, type CustomerRecord, CustomerRecordType, db, restore } from '../composables/customer-record';
import { type FormInst } from 'naive-ui';

const showQrCode = ref(false);

const customerRecordTypeOptions = Object.values(CustomerRecordType).map(value => ({ label: value, value }))

const customerRecordType = ref<CustomerRecordType>(CustomerRecordType.Private);
const customerRecords = ref<CustomerRecord[]>([]);
const yearFilter = useStorage('yearFilter', new Date().getFullYear())
const newRecordIds = ref(new Set<number>())

const filter = ref<{
  invoiceDate?: [number, number]
  invoice?: string
  customerName?: string
  invoiceAmount?: number
  chequeDate?: [number, number]
  chequeNo?: string
  chequeAmount?: number
  remark?: string
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
}

watchEffect(loadCustomerRecords)

const filteredCustomerRecords = computed(() => {
  return customerRecords.value.filter(record => {
    return (!filter.value.invoiceDate || (record.invoiceDate >= filter.value.invoiceDate[0] && record.invoiceDate <= endOfDay(filter.value.invoiceDate[1]).getTime())) &&
      (!filter.value.invoice || record.invoiceNo.startsWith(filter.value.invoice)) &&
      (!filter.value.customerName || record.customerName === filter.value.customerName) &&
      (!filter.value.invoiceAmount || record.invoiceAmount === filter.value.invoiceAmount) &&
      (!filter.value.chequeDate || (record.chequeDate && record.chequeDate >= filter.value.chequeDate[0] && record.chequeDate <= filter.value.chequeDate[1])) &&
      (!filter.value.chequeNo || record.chequeNo.startsWith(filter.value.chequeNo)) &&
      (!filter.value.chequeAmount || record.chequeAmount === filter.value.chequeAmount) &&
      (!filter.value.remark || record.remark.includes(filter.value.remark))
  })
})

const customerNames = ref<string[]>([])
const needUpdateCustomerNames = ref(true)

async function updateCustomerNames() {
  if (needUpdateCustomerNames.value) {
    needUpdateCustomerNames.value = false

    const names = new Set<string>()

    await dbCustomerRecords.value.each(record => {
      if (record.customerName)
        names.add(record.customerName)
    })

    customerNames.value = Array.from(names)
  }
}

onMounted(() => updateCustomerNames())
watch(dbCustomerRecords, () => updateCustomerNames())

const customerNameOptions = computed(() => {
  return customerNames.value.map(name => ({ label: name, value: name }))
})

function newCustomerRecord(override?: Partial<CustomerRecord>): CustomerRecord {
  const customerRecord: CustomerRecord = {
    invoiceDate: override?.invoiceDate ?? startOfDay(Date.now()).getTime(),
    chequeAmount: override?.chequeAmount ?? 0,
    chequeNo: override?.chequeNo ?? '',
    customerName: override?.customerName ?? (filter.value.customerName || ''),
    invoiceNo: override?.invoiceNo ?? '',
    invoiceAmount: override?.invoiceAmount ?? 0,
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

async function insertAfterCustomerRecord(index: number) {
  const customerRecord: CustomerRecord = newCustomerRecord({
    invoiceDate: customerRecords.value[index].invoiceDate,
  })

  await dbCustomerRecords.value.add(customerRecord)
  newRecordIds.value = new Set([customerRecord.id!])

  customerRecords.value.splice(index + 1, 0, customerRecord)
}

async function insertBeforeCustomerRecord(index: number) {
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
  needUpdateCustomerNames.value = needUpdateCustomerNames.value || original.customerName !== record.customerName

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
    needUpdateCustomerNames.value = true
    await updateCustomerNames()

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
  needUpdateCustomerNames.value = records.some(record => !customerNames.value.some(name => name === record.customerName))
  await updateCustomerNames()
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
  needUpdateCustomerNames.value = true
  await updateCustomerNames()
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

</script>

<template lang="pug">
.h-screen.w-screen.overflow-hidden.p-4.flex.flex-col
  .flex.items-center.space-x-4.mb-4
    .flex.text-2xl.font-bold Customer Records
    .w-52.font-bold
      n-select(v-model:value="customerRecordType" :options="customerRecordTypeOptions")
    n-divider(vertical)
    .flex.items-center.space-x-2
      b Year:
      n-input-number.w-32(v-model:value="yearFilter" size="small" :precision="0")
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
  .flex-grow.overflow-hidden.relative.border
    n-scrollbar
      n-table.pr-2(size="small" :single-line="false" style="overflow-y: visible;" :bordered="false")
        thead.sticky.top-0.z-10
          tr
            th.w-52 Date
            th.w-32 Invoice
            th.flex.items-center Customer Name
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
            th.w-32 Cheque
            th.w-32 Amount
            th.w-32 Remark
            th.w-8
              .flex.place-content-center
                n-button(text type="success" @click="addCustomerRecord()")
                  Icon(icon="mdi:add")
          tr
            th
              n-date-picker.text-xs(v-model:value="filter.invoiceDate" size="small" type="daterange" format="YY-MM-dd" clearable)
            th
              n-input(v-model:value="filter.invoice" size="small" clearable)
            th
              n-select(v-model:value="filter.customerName" size="small" clearable :options="customerNameOptions" filterable)
            th
              n-input-number(v-model:value="filter.invoiceAmount" size="small" :show-button="false" :precision="2" clearable)
            th
              n-date-picker(v-model:value="filter.chequeDate" size="small"  format="YY-MM-dd" type="daterange" clearable)
            th
              n-input(v-model:value="filter.chequeNo" size="small" clearable)
            th
              n-input-number(v-model:value="filter.chequeAmount" size="small" :show-button="false" :precision="2" clearable)
            th
              n-input(v-model:value="filter.remark" size="small" clearable)
            th
        tbody
          tr(
            v-for="record, i in filteredCustomerRecords" 
            :key="record.id" 
            :class="{ 'bg-green-300': newRecordIds.has(record.id || 0) }"
            :ref="(el) => scrollIntoNewRecord(record, el)"
            )
            td
              n-date-picker(v-model:value="record.invoiceDate" size="small" @update:value="saveCustomerRecord(record)")
            td
              n-input.font-mono(v-model:value="record.invoiceNo" size="small" @update:value="saveCustomerRecord(record)")
            td
              auto-complete.font-mono(v-model="record.customerName" size="small" @update:value="saveCustomerRecord(record)" :options="customerNames" @blur="updateCustomerNames()")
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
              n-date-picker(v-model:value="record.chequeDate" size="small" @update:value="saveCustomerRecord(record)")
            td 
              n-input.font-mono(v-model:value="record.chequeNo" size="small" @update:value="saveCustomerRecord(record)")
            td
              n-popover(trigger="hover" :disabled="!!record.chequeAmount || !record.chequeNo" :arrow="false")
                template(#trigger)
                  n-input-number.font-mono.text-right(
                    v-model:value="record.chequeAmount"
                    size="small"
                    :input-props="{ class: record.chequeAmount == 0 ? '!text-gray-400' : record.chequeAmount < record.invoiceAmount ? '!text-red-400' : '!text-green-600' }"
                    :show-button="false"
                    :precision="2"
                    @update:value="saveCustomerRecord(record)"
                    )
                n-button(text @click="record.chequeAmount = record.invoiceAmount" size="small")
                  .font-mono {{ record.invoiceAmount.toFixed(2) }}
            td
              n-input.font-mono(v-model:value="record.remark" size="small" @update:value="saveCustomerRecord(record)")
            td
              .flex.space-x-2
                .flex.place-content-center
                  n-tooltip
                    div Insert Before
                    template(#trigger)
                      n-button(text type="success" @click="insertBeforeCustomerRecord(i)")
                        Icon(icon="tabler:row-insert-top")
                .flex.place-content-center
                  n-tooltip
                    div Insert After
                    template(#trigger)
                      n-button(text type="success" @click="insertAfterCustomerRecord(i)")
                        Icon(icon="tabler:row-insert-bottom")
                .flex.place-content-center(v-if="record.id")
                  n-tooltip
                    div Delete
                    template(#trigger)
                      n-popconfirm(@positive-click="removeCustomerRecord(record.id)") Are you sure to delete this record?
                        template(#trigger)
                          n-button(text type="error")
                            Icon(icon="mdi:delete")

n-modal(v-model:show="showQrCode" :mask-closable="false" preset="dialog" :show-icon="false" style="width: calc(100% - 48px);")
  AiKeyin(:customer-names="customerNames" @submit="aiSubmitData")
</template>

<style lang="scss">
.bg-green-300 {
  td {
    background-color: inherit;
  }
}

tr:hover {
  td {
    background-color: rgba(0, 0, 0, 0.1);
  }
}
</style>