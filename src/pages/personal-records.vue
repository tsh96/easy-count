<script lang="ts" setup>
import { Icon } from '@iconify/vue/dist/iconify.js';
import { useDebounceFn, useStorage } from '@vueuse/core';
import { endOfDay, format, startOfDay } from 'date-fns';
import { type ComponentPublicInstance, computed, onMounted, ref, toRaw, watch, watchEffect } from 'vue';
import AutoComplete from '../components/AutoComplete.vue';
import { backup, type Transaction, db, restore } from '../composables/personal-record';
import { migrateOldPersonalRecord } from '../composables/old-personal-record';


const transactions = ref<Transaction[]>([]);
const yearFilter = useStorage('yearFilter', new Date().getFullYear())
const newRecordIds = ref(new Set<number>())

watch(newRecordIds, useDebounceFn(() => {
  newRecordIds.value = new Set()
}, 5000))

const filter = ref<{
  date?: [number, number]
  description?: string
}>({})

const accumulated = computed(() => {
  let obj: Record<number, number> = {}
  let sum = 0
  transactions.value.forEach(record => {
    sum += record.credit - record.debit
    obj[record.id!] = sum
  })
  return obj;
})

async function loadTransactions() {
  transactions.value = await db.transactions
    .orderBy('[date+id]')
    .toArray();
}

watchEffect(loadTransactions)

const filteredTransactions = computed(() => {
  return transactions.value.filter(record => {
    const startOfYear = new Date(`${yearFilter.value}-01-01`).getTime()
    const startOfNextYear = new Date(`${yearFilter.value + 1}-01-01`).getTime()
    if (record.date < startOfYear || record.date >= startOfNextYear) return false

    return (!filter.value.date || (record.date >= filter.value.date[0] && record.date <= endOfDay(filter.value.date[1]).getTime())) &&
      (!filter.value.description || record.description.startsWith(filter.value.description))
  })
})

const descriptions = ref<string[]>([])
const needUpdateDescriptions = ref(true)


const descriptionOptions = computed(() => {
  return descriptions.value.map(name => ({ label: name, value: name }))
})

async function updateDescriptions() {
  if (needUpdateDescriptions.value) {
    needUpdateDescriptions.value = false

    const names = new Set<string>()

    await db.transactions.each(record => {
      if (record.description)
        names.add(record.description)
    })

    descriptions.value = Array.from(names)
  }
}

onMounted(async () => {
  await migrateOldPersonalRecord()
  updateDescriptions()
})

function newTransaction(override?: Partial<Transaction>): Transaction {
  const transaction: Transaction = {
    date: override?.date ?? startOfDay(Date.now()).getTime(),
    description: override?.description ?? '',
    credit: override?.credit ?? 0,
    debit: override?.debit ?? 0,
  }

  return transaction
}


async function addTransaction() {
  const transaction: Transaction = newTransaction({
    description: filter.value.description,
    date: filter.value.date ? filter.value.date[1] : startOfDay(Date.now()).getTime(),
  })

  await db.transactions.add(transaction)
  newRecordIds.value = new Set([transaction.id!])

  transactions.value.push(transaction)
}

async function removeTransaction(id: number) {
  await db.transactions.delete(id)

  transactions.value = transactions.value.filter(record => record.id !== id)
}

async function insertAfterTransaction(index: number) {
  const customerRecord: Transaction = newTransaction({
    date: transactions.value[index].date,
  })

  await db.transactions.add(customerRecord)
  newRecordIds.value = new Set([customerRecord.id!])

  transactions.value.splice(index + 1, 0, customerRecord)
}

async function insertBeforeTransaction(index: number) {
  const customerRecord: Transaction = newTransaction({
    date: transactions.value[index].date,
  })

  await db.transactions.add(customerRecord)

  newRecordIds.value = new Set([customerRecord.id!])

  transactions.value.splice(index, 0, customerRecord)
}

async function saveTransaction(record: Transaction) {
  const original = await db.transactions.get(record.id)
  if (!original) return
  needUpdateDescriptions.value = needUpdateDescriptions.value || original.description !== record.description

  await db.transactions.put(toRaw(record))
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
    await loadTransactions()
    needUpdateDescriptions.value = true
    await updateDescriptions()

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


function scrollIntoNewRecord(record: Transaction, el: Element | ComponentPublicInstance | null) {
  if (newRecordIds.value.has(record.id || 0)) {
    if (el instanceof Element)
      el?.scrollIntoView({ behavior: 'smooth', block: 'center' })
  }
}

</script>

<template lang="pug">
.h-screen.w-screen.overflow-hidden.p-4.flex.flex-col
  Header(page="personal-records")
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
            th.w-52
              .items-center.flex.space-x-1
                span Date
                n-tooltip Reload and Sort
                  template(#trigger)
                    n-button(text @click="loadTransactions()")
                      Icon(icon="mdi:refresh")
            th Description
            th.w-32 Credit
            th.w-32 Debit
            th.w-32 Accumulated
            th.w-8
              .flex.place-content-center
                n-tooltip Add New record
                  template(#trigger)
                    n-button(text type="success" @click="addTransaction()")
                      Icon(icon="mdi:add")
          tr
            th
              n-date-picker.text-xs(v-model:value="filter.date" size="small" type="daterange" format="YY-MM-dd" clearable)
            th
              n-select(v-model:value="filter.description" size="small" clearable :options="descriptionOptions" filterable)
            th
            th
            th
            th
        tbody
          TransitionGroup(name="list")
            tr(
              v-for="record, i in filteredTransactions" 
              :key="record.id" 
              :class="{ 'bg-green-300': newRecordIds.has(record.id || 0) }"
              :ref="(el) => scrollIntoNewRecord(record, el)"
              )
              td
                n-date-picker(v-model:value="record.date" size="small" @update:value="saveTransaction(record)")
              td
                auto-complete.font-mono(v-model="record.description" size="small" @update:value="saveTransaction(record)" :options="descriptions" @blur="updateDescriptions()")
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
                        n-button(text type="success" @click="insertBeforeTransaction(i)")
                          Icon(icon="tabler:row-insert-top")
                  .flex.place-content-center
                    n-tooltip
                      div Insert After
                      template(#trigger)
                        n-button(text type="success" @click="insertAfterTransaction(i)")
                          Icon(icon="tabler:row-insert-bottom")
                  .flex.place-content-center(v-if="record.id")
                    n-tooltip
                      div Delete
                      template(#trigger)
                        n-popconfirm(@positive-click="removeTransaction(record.id)") Are you sure to delete this record?
                          template(#trigger)
                            n-button(text type="error")
                              Icon(icon="mdi:delete")
          tr.h-96
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