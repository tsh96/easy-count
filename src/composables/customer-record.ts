const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3001';

export interface CustomerRecord {
  id?: number;
  invoiceDate?: number;
  invoiceNo: string;
  customerName: string;
  invoiceAmount?: number;
  chequeDate?: number;
  chequeNo: string;
  chequeAmount?: number;
  remark: string;
}

export enum CustomerRecordType {
  Private = 'Private',
  Government = 'Government',
}

// API helper function
async function apiCall(endpoint: string, options?: RequestInit) {
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      'x-user-id': 'default', // TODO: Replace with actual user authentication
      ...options?.headers,
    },
  });

  if (!response.ok) {
    throw new Error(`API call failed: ${response.statusText}`);
  }

  return response.json();
}

// Helper function to create table API
function createTableAPI(recordType: CustomerRecordType) {
  return {
    async toArray(): Promise<CustomerRecord[]> {
      return apiCall(`/api/customer-records/${recordType}`);
    },

    where(field: string) {
      return {
        equals(value: string) {
          return {
            async modify(updates: Partial<CustomerRecord>): Promise<void> {
              // For replace customer name functionality
              if (field === 'customerName' && updates.customerName) {
                await apiCall(`/api/customer-records/${recordType}/replace-name`, {
                  method: 'POST',
                  body: JSON.stringify({
                    oldName: value,
                    newName: updates.customerName,
                  }),
                });
              }
            },
          };
        },
        between(lower: number, upper: number, _includeLower: boolean, _includeUpper: boolean) {
          return {
            async toArray(): Promise<CustomerRecord[]> {
              return apiCall(
                `/api/customer-records/${recordType}?startDate=${lower}&endDate=${upper}`
              );
            },
          };
        },
      };
    },

    async add(record: CustomerRecord): Promise<number> {
      const result = await apiCall(`/api/customer-records/${recordType}`, {
        method: 'POST',
        body: JSON.stringify(record),
      });
      record.id = result.id;
      return result.id;
    },

    async get(id?: number): Promise<CustomerRecord | undefined> {
      if (!id) return undefined;
      try {
        const records = await this.toArray();
        return records.find(r => r.id === id);
      } catch {
        return undefined;
      }
    },

    async put(record: CustomerRecord): Promise<number> {
      if (!record.id) throw new Error('Customer record ID is required');
      await apiCall(`/api/customer-records/${recordType}/${record.id}`, {
        method: 'PUT',
        body: JSON.stringify(record),
      });
      return record.id;
    },

    async delete(id: number): Promise<void> {
      await apiCall(`/api/customer-records/${recordType}/${id}`, {
        method: 'DELETE',
      });
    },

    async clear(): Promise<void> {
      const records = await this.toArray();
      for (const r of records) {
        if (r.id) await this.delete(r.id);
      }
    },

    async bulkAdd(records: CustomerRecord[], _options?: { allKeys?: boolean }): Promise<void> {
      await apiCall(`/api/customer-records/${recordType}/bulk`, {
        method: 'POST',
        body: JSON.stringify({ records }),
      });
    },
  };
}

// Database API wrapper to maintain compatibility
export const db = {
  privateCustomerRecords: createTableAPI(CustomerRecordType.Private),
  governmentCustomerRecords: createTableAPI(CustomerRecordType.Government),
};

export async function backup() {
  const data = await apiCall('/api/backup/customer-records');
  return JSON.stringify(data);
}

export async function restore(dataStr: string) {
  const data = JSON.parse(dataStr);
  await apiCall('/api/restore/customer-records', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}