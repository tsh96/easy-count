<script setup lang="ts">
import Peer from 'peerjs';
import { ref, toRaw } from 'vue';
import { type CustomerRecord } from '../composables/customer-record';
import { Icon } from '@iconify/vue/dist/iconify.js';

defineProps<{
  customerNames: string[];
}>();

const emit = defineEmits<{
  (e: 'submit', records: CustomerRecord[]): void;
}>();

const peer = new Peer()
const peerId = ref<string>();

peer.on('open', (id) => {
  peerId.value = id;
});

const customerRecords = ref<CustomerRecord[]>([]);

let recordId = 0;

peer.on('connection', (connection) => {
  connection.on('data', (data) => {
    if (currentStep.value === 1) {
      currentStep.value += 1;
    } else if (currentStep.value === 2) {
      const tsv = (data as string)

      let prevRecord = undefined as CustomerRecord | undefined;

      console.log(tsv);

      customerRecords.value = tsv.split('\n').flatMap((line) => {
        const columns = line.split('\t')
        if (columns[0].toLocaleLowerCase().includes("date")) return [];
        if (columns.length < 2) return [];
        prevRecord = {
          id: recordId++,
          invoiceDate: new Date(columns[0]).getTime() || prevRecord?.invoiceDate || new Date().getTime(),
          invoiceNo: columns[1] || '',
          customerName: columns[2] == "''" || !columns[2] ? prevRecord?.customerName || '' : columns[2] || '',
          invoiceAmount: parseFloat(columns[3]?.replace(/[, ]/g, '')) || 0,
          chequeDate: new Date(columns[4]).getTime() || undefined,
          chequeNo: columns[5] == "''" ? prevRecord?.chequeNo || '' : columns[5],
          chequeAmount: parseFloat(columns[6]?.replace(/[, ]/g, '')) || 0,
          remark: columns[7] || '',
        };
        return prevRecord;
      })

      currentStep.value = 3;
    }
  });
});

const currentStep = ref(1);

function addCustomerRecord() {
  customerRecords.value.push({
    id: recordId++,
    invoiceDate: new Date().getTime(),
    invoiceNo: '',
    customerName: '',
    invoiceAmount: 0,
    chequeNo: '',
    chequeAmount: 0,
    remark: '',
  })
}

function removeCustomerRecord(index: number) {
  customerRecords.value.splice(index, 1)
}

function insertAfterCustomerRecord(index: number) {
  customerRecords.value.splice(index + 1, 0, {
    id: recordId++,
    invoiceDate: customerRecords.value[index].invoiceDate,
    invoiceNo: '',
    customerName: '',
    invoiceAmount: 0,
    chequeNo: '',
    chequeAmount: 0,
    remark: '',
  })
}

function insertBeforeCustomerRecord(index: number) {
  customerRecords.value.splice(index, 0, {
    id: recordId++,
    invoiceDate: customerRecords.value[index].invoiceDate,
    invoiceNo: '',
    customerName: '',
    invoiceAmount: 0,
    chequeNo: '',
    chequeAmount: 0,
    remark: '',
  })
}

function submitCustomerRecords() {
  customerRecords.value.forEach((record) => {
    record.id = undefined;
  });
  emit('submit', toRaw(customerRecords.value));
  currentStep.value++;
}
</script>

<template lang="pug">
n-space(vertical)
  n-steps(:current="currentStep" )
    n-step(title="Scan QR")
    n-step(title="Take Photo")
    n-step(title="Check and Edit Data")
    n-step(title="Done")

  template(v-if="currentStep === 1")
    n-spin.flex.place-content-center(:show="!peerId")
      .flex.flex-col.items-center
        n-qr-code.my-10(:value="peerId" :size="256")
        .text-sm {{ peerId }}
        .text-lg Please scan the QR code from your mobile phone.

  template(v-else-if="currentStep === 2")
    n-spin.flex.place-content-center(:show="!customerRecords.length")
      .flex.flex-col.items-center.justify-center.h-56
        .text-2xl Please take a photo of the document from your mobile phone.
        .text-lg It will be sent to the computer automatically. Please wait.

  template(v-else-if="currentStep === 3")
    .min-h-96
      n-table(size="small" :single-line="false")
        thead
          tr
            th.w-52 Date
            th.w-32 Invoice
            th Customer Name
            th.w-32 Amount
            th.w-52 Date
            th.w-32 Cheque
            th.w-32 Amount
            th.w-32 Remark
            th.w-8
              .flex.place-content-center
                n-button(text type="success" @click="addCustomerRecord()")
                  Icon(icon="mdi:add")
        tbody
          tr(v-for="record, i in customerRecords" :key="record.id")
            td
              n-date-picker(v-model:value="record.invoiceDate" size="small")
            td
              n-input.font-mono(v-model:value="record.invoiceNo" size="small")
            td
              auto-complete.font-mono(v-model="record.customerName" size="small" :options="customerNames")
            td
              n-input-number.font-mono.text-right(v-model:value="record.invoiceAmount" size="small" :show-button="false" :precision="2") 
            td
              n-date-picker(v-model:value="record.chequeDate" size="small")
            td 
              n-input.font-mono(v-model:value="record.chequeNo" size="small")
            td
              n-input-number.font-mono.text-right(v-model:value="record.chequeAmount" size="small" :show-button="false" :precision="2") 
            td
              n-input.font-mono(v-model:value="record.remark" size="small")
            td
              .flex.space-x-2
                .flex.place-content-center
                  n-button(text type="success" @click="insertBeforeCustomerRecord(i)")
                    Icon(icon="tabler:row-insert-top")
                .flex.place-content-center
                  n-button(text type="success" @click="insertAfterCustomerRecord(i)")
                    Icon(icon="tabler:row-insert-bottom")
                .flex.place-content-center
                  n-button(text type="error" @click="removeCustomerRecord(i)")
                    Icon(icon="mdi:delete")
    .flex.justify-end.mt-4
      n-button(type="primary" @click="submitCustomerRecords")
        | Submit

  template(v-else-if="currentStep === 4")
    .flex.flex-col.items-center.justify-center.h-96
      .text-2xl Your data has been submitted successfully.
      .text-lg You can take another photo or close this dialog.
</template>