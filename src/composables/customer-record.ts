import Dexie, { type EntityTable } from "dexie";

export interface CustomerRecord {
  id?: number;
  invoiceDate: number;
  invoiceNo: string;
  customerName: string;
  invoiceAmount: number;
  chequeDate?: number;
  chequeNo: string;
  chequeAmount: number;
  remark: string;
}

export enum CustomerRecordType {
  Private = 'Private',
  Government = 'Government',
}

export const db = new Dexie('EasyCount') as Dexie & {
  privateCustomerRecords: EntityTable<
    CustomerRecord,
    'id' // primary key "id" (for the typings only)
  >;
  governmentCustomerRecords: EntityTable<
    CustomerRecord,
    'id' // primary key "id" (for the typings only)
  >;
};

db.version(1).stores({
  privateCustomerRecords: '++id, [invoiceDate+invoiceNo],  invoiceNo, customerName, invoiceAmount, chequeDate, chequeNo, chequeAmount, remark', // primary key "id" (for the runtime!)
  governmentCustomerRecords: '++id, [invoiceDate+invoiceNo],  invoiceNo, customerName, invoiceAmount, chequeDate, chequeNo, chequeAmount, remark', // primary key "id" (for the runtime!)
});

export async function backup() {
  const privateCustomerRecords = db.privateCustomerRecords.toArray();
  const governmentCustomerRecords = db.governmentCustomerRecords.toArray();
  const data = {
    privateCustomerRecords: await privateCustomerRecords,
    governmentCustomerRecords: await governmentCustomerRecords,
  };

  return JSON.stringify(data);
}

export async function restore(dataStr: string) {
  const data = JSON.parse(dataStr);
  await db.privateCustomerRecords.clear();
  await db.governmentCustomerRecords.clear();
  await db.privateCustomerRecords.bulkAdd(data.privateCustomerRecords);
  await db.governmentCustomerRecords.bulkAdd(data.governmentCustomerRecords);
}