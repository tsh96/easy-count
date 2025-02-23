import { useLocalStorage } from "@vueuse/core"
import { CustomerRecord, db } from "./customer-record"

type CustomerAccount = {
  customerName: string,
  items: CustomerAccountItem[],
}

type CustomerAccountItem = {
  id: string
  invoiceDate: number | null,
  invoiceNo: string,
  invoiceAmount: number,
  chequeDate: number | null,
  chequeNo: string,
  chequeAmount: number,
  remark: string,
}

export async function migrateOldCustomerRecord() {
  const customerRecordMigrated = useLocalStorage('customerRecordMigrated', false)
  if (customerRecordMigrated.value) return

  const privateCustomerRecords = useLocalStorage('customerRecords', {
    version: '1.0.0',
    customerAccounts: [] as CustomerAccount[],
    selectedCustomerName: null as string | null,
  }, { writeDefaults: true })

  const governmentCustomerRecords = useLocalStorage('governmentRecords', {
    version: '1.0.0',
    customerAccounts: [] as CustomerAccount[],
    selectedCustomerName: null as string | null,
  }, { writeDefaults: true })

  await db.privateCustomerRecords.bulkAdd(privateCustomerRecords.value.customerAccounts.flatMap((customerAccount) => {
    return customerAccount.items.map<CustomerRecord>((item) => {
      return {
        invoiceDate: item.invoiceDate || 0,
        invoiceNo: item.invoiceNo,
        customerName: customerAccount.customerName,
        invoiceAmount: item.invoiceAmount,
        chequeDate: item.chequeDate || undefined,
        chequeNo: item.chequeNo,
        chequeAmount: item.chequeAmount,
        remark: item.remark,
      }
    })
  }))

  await db.governmentCustomerRecords.bulkAdd(governmentCustomerRecords.value.customerAccounts.flatMap((customerAccount) => {
    return customerAccount.items.map<CustomerRecord>((item) => {
      return {
        invoiceDate: item.invoiceDate || 0,
        invoiceNo: item.invoiceNo,
        customerName: customerAccount.customerName,
        invoiceAmount: item.invoiceAmount,
        chequeDate: item.chequeDate || undefined,
        chequeNo: item.chequeNo,
        chequeAmount: item.chequeAmount,
        remark: item.remark,
      }
    })
  }))

  customerRecordMigrated.value = true
}