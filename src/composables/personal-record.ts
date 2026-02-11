import { apiCall } from './auth';

export interface Transaction {
  id?: number;
  date: number;
  description: string;
  credit: number;
  debit: number;
}

// Database API wrapper to maintain compatibility
export const db = {
  transactions: {
    async toArray(): Promise<Transaction[]> {
      return apiCall('/api/transactions');
    },

    async add(transaction: Transaction): Promise<number> {
      const result = await apiCall('/api/transactions', {
        method: 'POST',
        body: JSON.stringify(transaction),
      });
      transaction.id = result.id;
      return result.id;
    },

    async get(id?: number): Promise<Transaction | undefined> {
      if (!id) return undefined;
      try {
        const transactions = await this.toArray();
        return transactions.find(t => t.id === id);
      } catch {
        return undefined;
      }
    },

    async put(transaction: Transaction): Promise<number> {
      if (!transaction.id) throw new Error('Transaction ID is required');
      await apiCall(`/api/transactions/${transaction.id}`, {
        method: 'PUT',
        body: JSON.stringify(transaction),
      });
      return transaction.id;
    },

    async delete(id: number): Promise<void> {
      await apiCall(`/api/transactions/${id}`, {
        method: 'DELETE',
      });
    },

    async clear(): Promise<void> {
      const transactions = await this.toArray();
      for (const t of transactions) {
        if (t.id) await this.delete(t.id);
      }
    },

    async bulkAdd(transactions: Transaction[]): Promise<void> {
      if (transactions.length === 0) return;
      
      await apiCall('/api/transactions/bulk', {
        method: 'POST',
        body: JSON.stringify({ transactions }),
      });
    },

    async each(callback: (transaction: Transaction) => void): Promise<void> {
      const transactions = await this.toArray();
      transactions.forEach(callback);
    },

    orderBy(key: string) {
      return {
        async toArray(): Promise<Transaction[]> {
          const transactions = await db.transactions.toArray();
          return transactions.sort((a, b) => {
            if (key === '[date+id]') {
              return a.date - b.date || (a.id || 0) - (b.id || 0);
            }
            return 0;
          });
        },
      };
    },
  },
};

export async function backup() {
  const data = await apiCall('/api/backup/transactions');
  return JSON.stringify(data);
}

export async function restore(dataStr: string) {
  const data = JSON.parse(dataStr);
  await apiCall('/api/restore/transactions', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}