import Dexie, { type EntityTable } from "dexie";

export interface Transaction {
  id?: number;
  date: number;
  description: string;
  credit: number;
  debit: number;
}

export const db = new Dexie('PersonalRecord') as Dexie & {
  transactions: EntityTable<
    Transaction,
    'id' // primary key "id" (for the typings only)
  >;
};

db.version(1).stores({
  transactions: '++id, [date+id]', // primary key "id" (for the runtime!)
});

export async function backup() {
  const transactions = db.transactions.toArray();
  const data = {
    transactions: await transactions,
  };

  return JSON.stringify(data);
}

export async function restore(dataStr: string) {
  const data = JSON.parse(dataStr);
  await db.transactions.clear();
  await db.transactions.bulkAdd(data.transactions);
}