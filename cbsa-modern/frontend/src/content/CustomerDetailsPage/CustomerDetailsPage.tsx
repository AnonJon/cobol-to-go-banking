import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import {
  Breadcrumb,
  BreadcrumbItem,
  TextInput,
  Button,
  Modal,
  Form,
  Stack,
  DataTable,
  Table,
  TableHead,
  TableRow,
  TableHeader,
  TableBody,
  TableCell,
} from '@carbon/react';
import { customerApi, accountApi } from '../../api/client';
import type { Customer, Account, UpdateCustomerRequest } from '../../api/types';

const CustomerDetailsPage: React.FC = () => {
  const [searchQuery, setSearchQuery] = useState('');
  const [customer, setCustomer] = useState<Customer | null>(null);
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [editing, setEditing] = useState(false);
  const navigate = useNavigate();
  const [editData, setEditData] = useState<UpdateCustomerRequest>({});
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const handleSearch = async () => {
    setError('');
    try {
      const c = await customerApi.get(searchQuery);
      setCustomer(c);
      const { accounts: accts } = await accountApi.listByCustomer(c.customerNumber);
      setAccounts(accts || []);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Customer not found');
      setCustomer(null);
      setAccounts([]);
    }
  };

  const handleUpdate = async () => {
    if (!customer) return;
    try {
      const updated = await customerApi.update(customer.customerNumber, editData);
      setCustomer(updated);
      setEditing(false);
      setSuccess('Customer updated successfully');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Update failed');
    }
  };

  const accountHeaders = [
    { key: 'accountNumber', header: 'Account #' },
    { key: 'accountType', header: 'Type' },
    { key: 'availableBalance', header: 'Balance' },
    { key: 'overdraftLimit', header: 'Overdraft' },
  ];

  return (
    <>
      <Breadcrumb>
        <BreadcrumbItem><Link to="/">Dashboard</Link></BreadcrumbItem>
        <BreadcrumbItem isCurrentPage>Customer Details</BreadcrumbItem>
      </Breadcrumb>

      <div className="page-header">
        <h2>Customer Details</h2>
        <p>Look up and manage customer information</p>
      </div>

      <div className="form-card" style={{ marginBottom: '1.5rem' }}>
        <Stack gap={4} orientation="horizontal">
          <TextInput id="search" labelText="Customer Number" value={searchQuery}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setSearchQuery(e.target.value)}
            onKeyDown={(e: React.KeyboardEvent) => e.key === 'Enter' && handleSearch()} />
          <Button onClick={handleSearch} style={{ alignSelf: 'flex-end' }}>Search</Button>
        </Stack>
      </div>

      {customer && !editing && (
        <div className="detail-card">
          <h3 style={{ marginTop: 0 }}>Customer {customer.customerNumber}</h3>
          <div className="detail-row"><span className="label">Name</span><span className="value">{customer.title} {customer.firstName} {customer.lastName}</span></div>
          <div className="detail-row"><span className="label">Address</span><span className="value">{customer.addressLine1}, {customer.addressLine2}, {customer.addressLine3}</span></div>
          <div className="detail-row"><span className="label">Date of Birth</span><span className="value">{customer.dateOfBirth}</span></div>
          <div className="detail-row"><span className="label">Credit Score</span><span className="value">{customer.creditScore}</span></div>
          <div style={{ marginTop: '1rem' }}>
            <Button onClick={() => {
              setEditData({ title: customer.title, firstName: customer.firstName, lastName: customer.lastName,
                addressLine1: customer.addressLine1, addressLine2: customer.addressLine2, addressLine3: customer.addressLine3 });
              setEditing(true);
            }}>Edit Customer</Button>
          </div>
        </div>
      )}

      {customer && editing && (
        <div className="form-card" style={{ marginBottom: '1.5rem' }}>
          <h3 style={{ marginTop: 0 }}>Edit Customer {customer.customerNumber}</h3>
          <Form onSubmit={(e: React.FormEvent) => { e.preventDefault(); handleUpdate(); }}>
            <Stack gap={4}>
              <TextInput id="editFirst" labelText="First Name" value={editData.firstName || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setEditData(d => ({ ...d, firstName: e.target.value }))} />
              <TextInput id="editLast" labelText="Last Name" value={editData.lastName || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setEditData(d => ({ ...d, lastName: e.target.value }))} />
              <TextInput id="editAddr1" labelText="Address Line 1" value={editData.addressLine1 || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setEditData(d => ({ ...d, addressLine1: e.target.value }))} />
              <TextInput id="editAddr2" labelText="Address Line 2" value={editData.addressLine2 || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setEditData(d => ({ ...d, addressLine2: e.target.value }))} />
              <TextInput id="editAddr3" labelText="Town / City" value={editData.addressLine3 || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setEditData(d => ({ ...d, addressLine3: e.target.value }))} />
              <Stack gap={3} orientation="horizontal">
                <Button type="submit">Save Changes</Button>
                <Button kind="secondary" onClick={() => setEditing(false)}>Cancel</Button>
              </Stack>
            </Stack>
          </Form>
        </div>
      )}

      {accounts.length > 0 && (
        <div className="detail-card">
          <h3 style={{ marginTop: 0 }}>Linked Accounts</h3>
          <DataTable rows={accounts.map(a => ({ ...a, id: a.accountNumber }))} headers={accountHeaders}>
            {({ rows, headers, getTableProps, getHeaderProps, getRowProps }: any) => (
              <Table {...getTableProps()}>
                <TableHead>
                  <TableRow>
                    {headers.map((h: any) => (<TableHeader {...getHeaderProps({ header: h })} key={h.key}>{h.header}</TableHeader>))}
                  </TableRow>
                </TableHead>
                <TableBody>
                  {rows.map((row: any) => (
                    <TableRow {...getRowProps({ row })} key={row.id}>
                      {row.cells.map((cell: any) => (<TableCell key={cell.id}>{cell.value}</TableCell>))}
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            )}
          </DataTable>
        </div>
      )}

      {error && <Modal open modalHeading="Error" passiveModal onRequestClose={() => setError('')}><p>{error}</p></Modal>}
      {success && <Modal open modalHeading="Success" passiveModal onRequestClose={() => navigate('/')}><p>{success}</p></Modal>}
    </>
  );
};

export default CustomerDetailsPage;
