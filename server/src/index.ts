import express from 'express';
import cors from 'cors';
import dotenv from 'dotenv';
import { initDatabase, sql } from './db.js';

dotenv.config();

const app = express();
const PORT = process.env.PORT || 3001;

app.use(cors());
app.use(express.json());

// Initialize database
await initDatabase();

// TODO: SECURITY WARNING - Implement proper authentication before production deployment
// The current x-user-id header can be easily spoofed. Consider implementing:
// - JWT-based authentication
// - OAuth 2.0 (Google, GitHub, etc.)
// - Session-based authentication
// - API key authentication for single-user deployments

// Personal Records (Transactions) endpoints
app.get('/api/transactions', async (req, res) => {
  try {
    const userId = req.headers['x-user-id'] || 'default';
    const transactions = await sql`
      SELECT id, date, description, credit, debit 
      FROM transactions 
      WHERE user_id = ${userId}
      ORDER BY date, id
    `;
    
    // Convert numeric strings back to numbers
    const result = transactions.map(t => ({
      ...t,
      credit: parseFloat(t.credit),
      debit: parseFloat(t.debit),
      date: Number(t.date)
    }));
    
    res.json(result);
  } catch (error) {
    console.error('Error fetching transactions:', error);
    res.status(500).json({ error: 'Failed to fetch transactions' });
  }
});

app.post('/api/transactions', async (req, res) => {
  try {
    const userId = req.headers['x-user-id'] || 'default';
    const { date, description, credit, debit } = req.body;
    
    const result = await sql`
      INSERT INTO transactions (date, description, credit, debit, user_id)
      VALUES (${date}, ${description}, ${credit}, ${debit}, ${userId})
      RETURNING id, date, description, credit, debit
    `;
    
    const transaction = result[0];
    res.json({
      ...transaction,
      credit: parseFloat(transaction.credit),
      debit: parseFloat(transaction.debit),
      date: Number(transaction.date)
    });
  } catch (error) {
    console.error('Error creating transaction:', error);
    res.status(500).json({ error: 'Failed to create transaction' });
  }
});

app.put('/api/transactions/:id', async (req, res) => {
  try {
    const userId = req.headers['x-user-id'] || 'default';
    const { id } = req.params;
    const { date, description, credit, debit } = req.body;
    
    const result = await sql`
      UPDATE transactions 
      SET date = ${date}, description = ${description}, credit = ${credit}, debit = ${debit}, updated_at = CURRENT_TIMESTAMP
      WHERE id = ${id} AND user_id = ${userId}
      RETURNING id, date, description, credit, debit
    `;
    
    if (result.length === 0) {
      return res.status(404).json({ error: 'Transaction not found' });
    }
    
    const transaction = result[0];
    res.json({
      ...transaction,
      credit: parseFloat(transaction.credit),
      debit: parseFloat(transaction.debit),
      date: Number(transaction.date)
    });
  } catch (error) {
    console.error('Error updating transaction:', error);
    res.status(500).json({ error: 'Failed to update transaction' });
  }
});

app.delete('/api/transactions/:id', async (req, res) => {
  try {
    const userId = req.headers['x-user-id'] || 'default';
    const { id } = req.params;
    
    await sql`
      DELETE FROM transactions 
      WHERE id = ${id} AND user_id = ${userId}
    `;
    
    res.json({ success: true });
  } catch (error) {
    console.error('Error deleting transaction:', error);
    res.status(500).json({ error: 'Failed to delete transaction' });
  }
});

// Bulk transactions endpoint
app.post('/api/transactions/bulk', async (req, res) => {
  try {
    const userId = req.headers['x-user-id'] || 'default';
    const { transactions } = req.body;
    
    if (!transactions || transactions.length === 0) {
      return res.json([]);
    }
    
    const results = [];
    for (const t of transactions) {
      const result = await sql`
        INSERT INTO transactions (date, description, credit, debit, user_id)
        VALUES (${t.date}, ${t.description}, ${t.credit}, ${t.debit}, ${userId})
        RETURNING id
      `;
      results.push(result[0]);
    }
    
    res.json(results);
  } catch (error) {
    console.error('Error bulk creating transactions:', error);
    res.status(500).json({ error: 'Failed to bulk create transactions' });
  }
});

// Customer Records endpoints
app.get('/api/customer-records/:type', async (req, res) => {
  try {
    const userId = req.headers['x-user-id'] || 'default';
    const { type } = req.params;
    const { startDate, endDate } = req.query;
    
    let records;
    if (startDate && endDate) {
      records = await sql`
        SELECT id, invoice_date, invoice_no, customer_name, invoice_amount, 
               cheque_date, cheque_no, cheque_amount, remark
        FROM customer_records 
        WHERE record_type = ${type} AND user_id = ${userId}
          AND invoice_date >= ${startDate} AND invoice_date < ${endDate}
        ORDER BY invoice_date, invoice_no
      `;
    } else {
      records = await sql`
        SELECT id, invoice_date, invoice_no, customer_name, invoice_amount, 
               cheque_date, cheque_no, cheque_amount, remark
        FROM customer_records 
        WHERE record_type = ${type} AND user_id = ${userId}
        ORDER BY invoice_date, invoice_no
      `;
    }
    
    const result = records.map(r => ({
      ...r,
      invoiceDate: r.invoice_date ? Number(r.invoice_date) : undefined,
      invoiceNo: r.invoice_no,
      customerName: r.customer_name,
      invoiceAmount: r.invoice_amount ? parseFloat(r.invoice_amount) : undefined,
      chequeDate: r.cheque_date ? Number(r.cheque_date) : undefined,
      chequeNo: r.cheque_no,
      chequeAmount: r.cheque_amount ? parseFloat(r.cheque_amount) : undefined,
      remark: r.remark
    }));
    
    res.json(result);
  } catch (error) {
    console.error('Error fetching customer records:', error);
    res.status(500).json({ error: 'Failed to fetch customer records' });
  }
});

app.post('/api/customer-records/:type', async (req, res) => {
  try {
    const userId = req.headers['x-user-id'] || 'default';
    const { type } = req.params;
    const { invoiceDate, invoiceNo, customerName, invoiceAmount, chequeDate, chequeNo, chequeAmount, remark } = req.body;
    
    const result = await sql`
      INSERT INTO customer_records 
        (record_type, invoice_date, invoice_no, customer_name, invoice_amount, 
         cheque_date, cheque_no, cheque_amount, remark, user_id)
      VALUES 
        (${type}, ${invoiceDate || null}, ${invoiceNo}, ${customerName}, ${invoiceAmount || null},
         ${chequeDate || null}, ${chequeNo}, ${chequeAmount || null}, ${remark}, ${userId})
      RETURNING id, invoice_date, invoice_no, customer_name, invoice_amount, 
                cheque_date, cheque_no, cheque_amount, remark
    `;
    
    const record = result[0];
    res.json({
      ...record,
      invoiceDate: record.invoice_date ? Number(record.invoice_date) : undefined,
      invoiceNo: record.invoice_no,
      customerName: record.customer_name,
      invoiceAmount: record.invoice_amount ? parseFloat(record.invoice_amount) : undefined,
      chequeDate: record.cheque_date ? Number(record.cheque_date) : undefined,
      chequeNo: record.cheque_no,
      chequeAmount: record.cheque_amount ? parseFloat(record.cheque_amount) : undefined,
      remark: record.remark
    });
  } catch (error) {
    console.error('Error creating customer record:', error);
    res.status(500).json({ error: 'Failed to create customer record' });
  }
});

app.post('/api/customer-records/:type/bulk', async (req, res) => {
  try {
    const userId = req.headers['x-user-id'] || 'default';
    const { type } = req.params;
    const { records } = req.body;
    
    if (!records || records.length === 0) {
      return res.json([]);
    }
    
    // Use a transaction for bulk insert
    const results = [];
    for (const record of records) {
      const result = await sql`
        INSERT INTO customer_records 
          (record_type, invoice_date, invoice_no, customer_name, invoice_amount, 
           cheque_date, cheque_no, cheque_amount, remark, user_id)
        VALUES 
          (${type}, ${record.invoiceDate || null}, ${record.invoiceNo}, ${record.customerName}, 
           ${record.invoiceAmount || null}, ${record.chequeDate || null}, ${record.chequeNo}, 
           ${record.chequeAmount || null}, ${record.remark}, ${userId})
        RETURNING id
      `;
      results.push(result[0]);
    }
    
    res.json(results);
  } catch (error) {
    console.error('Error bulk creating customer records:', error);
    res.status(500).json({ error: 'Failed to bulk create customer records' });
  }
});

app.put('/api/customer-records/:type/:id', async (req, res) => {
  try {
    const userId = req.headers['x-user-id'] || 'default';
    const { type, id } = req.params;
    const { invoiceDate, invoiceNo, customerName, invoiceAmount, chequeDate, chequeNo, chequeAmount, remark } = req.body;
    
    const result = await sql`
      UPDATE customer_records 
      SET invoice_date = ${invoiceDate || null}, 
          invoice_no = ${invoiceNo}, 
          customer_name = ${customerName}, 
          invoice_amount = ${invoiceAmount || null},
          cheque_date = ${chequeDate || null}, 
          cheque_no = ${chequeNo}, 
          cheque_amount = ${chequeAmount || null}, 
          remark = ${remark},
          updated_at = CURRENT_TIMESTAMP
      WHERE id = ${id} AND record_type = ${type} AND user_id = ${userId}
      RETURNING id, invoice_date, invoice_no, customer_name, invoice_amount, 
                cheque_date, cheque_no, cheque_amount, remark
    `;
    
    if (result.length === 0) {
      return res.status(404).json({ error: 'Customer record not found' });
    }
    
    const record = result[0];
    res.json({
      ...record,
      invoiceDate: record.invoice_date ? Number(record.invoice_date) : undefined,
      invoiceNo: record.invoice_no,
      customerName: record.customer_name,
      invoiceAmount: record.invoice_amount ? parseFloat(record.invoice_amount) : undefined,
      chequeDate: record.cheque_date ? Number(record.cheque_date) : undefined,
      chequeNo: record.cheque_no,
      chequeAmount: record.cheque_amount ? parseFloat(record.cheque_amount) : undefined,
      remark: record.remark
    });
  } catch (error) {
    console.error('Error updating customer record:', error);
    res.status(500).json({ error: 'Failed to update customer record' });
  }
});

app.delete('/api/customer-records/:type/:id', async (req, res) => {
  try {
    const userId = req.headers['x-user-id'] || 'default';
    const { type, id } = req.params;
    
    await sql`
      DELETE FROM customer_records 
      WHERE id = ${id} AND record_type = ${type} AND user_id = ${userId}
    `;
    
    res.json({ success: true });
  } catch (error) {
    console.error('Error deleting customer record:', error);
    res.status(500).json({ error: 'Failed to delete customer record' });
  }
});

// Update customer names
app.post('/api/customer-records/:type/replace-name', async (req, res) => {
  try {
    const userId = req.headers['x-user-id'] || 'default';
    const { type } = req.params;
    const { oldName, newName } = req.body;
    
    await sql`
      UPDATE customer_records 
      SET customer_name = ${newName}, updated_at = CURRENT_TIMESTAMP
      WHERE customer_name = ${oldName} AND record_type = ${type} AND user_id = ${userId}
    `;
    
    res.json({ success: true });
  } catch (error) {
    console.error('Error replacing customer name:', error);
    res.status(500).json({ error: 'Failed to replace customer name' });
  }
});

// Backup endpoints
app.get('/api/backup/transactions', async (req, res) => {
  try {
    const userId = req.headers['x-user-id'] || 'default';
    const transactions = await sql`
      SELECT date, description, credit, debit 
      FROM transactions 
      WHERE user_id = ${userId}
      ORDER BY date, id
    `;
    
    const result = transactions.map(t => ({
      date: Number(t.date),
      description: t.description,
      credit: parseFloat(t.credit),
      debit: parseFloat(t.debit)
    }));
    
    res.json({ transactions: result });
  } catch (error) {
    console.error('Error backing up transactions:', error);
    res.status(500).json({ error: 'Failed to backup transactions' });
  }
});

app.post('/api/restore/transactions', async (req, res) => {
  try {
    const userId = req.headers['x-user-id'] || 'default';
    const { transactions } = req.body;
    
    // Clear existing transactions
    await sql`DELETE FROM transactions WHERE user_id = ${userId}`;
    
    // Insert new transactions using bulk endpoint logic
    if (transactions && transactions.length > 0) {
      for (const t of transactions) {
        await sql`
          INSERT INTO transactions (date, description, credit, debit, user_id)
          VALUES (${t.date}, ${t.description}, ${t.credit}, ${t.debit}, ${userId})
        `;
      }
    }
    
    res.json({ success: true });
  } catch (error) {
    console.error('Error restoring transactions:', error);
    res.status(500).json({ error: 'Failed to restore transactions' });
  }
});

app.get('/api/backup/customer-records', async (req, res) => {
  try {
    const userId = req.headers['x-user-id'] || 'default';
    
    const privateRecords = await sql`
      SELECT invoice_date, invoice_no, customer_name, invoice_amount, 
             cheque_date, cheque_no, cheque_amount, remark
      FROM customer_records 
      WHERE record_type = 'Private' AND user_id = ${userId}
      ORDER BY invoice_date, invoice_no
    `;
    
    const governmentRecords = await sql`
      SELECT invoice_date, invoice_no, customer_name, invoice_amount, 
             cheque_date, cheque_no, cheque_amount, remark
      FROM customer_records 
      WHERE record_type = 'Government' AND user_id = ${userId}
      ORDER BY invoice_date, invoice_no
    `;
    
    const mapRecord = (r: any) => ({
      invoiceDate: r.invoice_date ? Number(r.invoice_date) : undefined,
      invoiceNo: r.invoice_no,
      customerName: r.customer_name,
      invoiceAmount: r.invoice_amount ? parseFloat(r.invoice_amount) : undefined,
      chequeDate: r.cheque_date ? Number(r.cheque_date) : undefined,
      chequeNo: r.cheque_no,
      chequeAmount: r.cheque_amount ? parseFloat(r.cheque_amount) : undefined,
      remark: r.remark
    });
    
    res.json({
      privateCustomerRecords: privateRecords.map(mapRecord),
      governmentCustomerRecords: governmentRecords.map(mapRecord)
    });
  } catch (error) {
    console.error('Error backing up customer records:', error);
    res.status(500).json({ error: 'Failed to backup customer records' });
  }
});

app.post('/api/restore/customer-records', async (req, res) => {
  try {
    const userId = req.headers['x-user-id'] || 'default';
    const { privateCustomerRecords, governmentCustomerRecords } = req.body;
    
    // Clear existing records
    await sql`DELETE FROM customer_records WHERE user_id = ${userId}`;
    
    // Insert private records
    for (const r of privateCustomerRecords) {
      await sql`
        INSERT INTO customer_records 
          (record_type, invoice_date, invoice_no, customer_name, invoice_amount, 
           cheque_date, cheque_no, cheque_amount, remark, user_id)
        VALUES 
          ('Private', ${r.invoiceDate || null}, ${r.invoiceNo}, ${r.customerName}, 
           ${r.invoiceAmount || null}, ${r.chequeDate || null}, ${r.chequeNo}, 
           ${r.chequeAmount || null}, ${r.remark}, ${userId})
      `;
    }
    
    // Insert government records
    for (const r of governmentCustomerRecords) {
      await sql`
        INSERT INTO customer_records 
          (record_type, invoice_date, invoice_no, customer_name, invoice_amount, 
           cheque_date, cheque_no, cheque_amount, remark, user_id)
        VALUES 
          ('Government', ${r.invoiceDate || null}, ${r.invoiceNo}, ${r.customerName}, 
           ${r.invoiceAmount || null}, ${r.chequeDate || null}, ${r.chequeNo}, 
           ${r.chequeAmount || null}, ${r.remark}, ${userId})
      `;
    }
    
    res.json({ success: true });
  } catch (error) {
    console.error('Error restoring customer records:', error);
    res.status(500).json({ error: 'Failed to restore customer records' });
  }
});

app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});
