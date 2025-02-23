import { useStorage } from "@vueuse/core"
import { db } from "./personal-record"

export interface Transaction {
  date: number
  description: string
  credit: number
  debit: number
}

export async function migrateOldPersonalRecord() {
  const personalRecordMigrated = useStorage('personalRecordMigrated', false)
  if (personalRecordMigrated.value) return

  const oldTransactions = useStorage<Transaction[]>('transactions', [])
  await db.transactions.bulkAdd(oldTransactions.value.map((transaction) => {
    return {
      credit: transaction.credit,
      date: transaction.date,
      debit: transaction.debit,
      description: transaction.description
    }
  }))
  personalRecordMigrated.value = true
}