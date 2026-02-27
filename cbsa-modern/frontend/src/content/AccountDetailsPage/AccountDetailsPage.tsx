import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import {
  Breadcrumb, BreadcrumbItem, TextInput, Button, Modal, Form, Stack, NumberInput, Dropdown,
} from '@carbon/react';
import { accountApi } from '../../api/client';
import type { Account, UpdateAccountRequest } from '../../api/types';

const accountTypes = ['ISA', 'SAVING', 'CURRENT', 'LOAN', 'MORTGAGE'];

const AccountDetailsPage: React.FC = () => {
  const [searchQuery, setSearchQuery] = useState('');
  const [account, setAccount] = useState<Account | null>(null);
  const [editing, setEditing] = useState(false);
  const navigate = useNavigate();
  const [editData, setEditData] = useState<UpdateAccountRequest>({});
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const handleSearch = async () => {
    setError('');
    try { setAccount(await accountApi.get(searchQuery)); }
    catch (err: any) { setError(err.response?.data?.error || 'Account not found'); setAccount(null); }
  };

  const handleUpdate = async () => {
    if (!account) return;
    try {
      setAccount(await accountApi.update(account.accountNumber, editData));
      setEditing(false);
      setSuccess('Account updated successfully');
    } catch (err: any) { setError(err.response?.data?.error || 'Update failed'); }
  };

  return (
    <>
      <Breadcrumb>
        <BreadcrumbItem><Link to="/">Dashboard</Link></BreadcrumbItem>
        <BreadcrumbItem isCurrentPage>Account Details</BreadcrumbItem>
      </Breadcrumb>
      <div className="page-header"><h2>Account Details</h2><p>Look up and manage account information</p></div>
      <div className="form-card" style={{ marginBottom: '1.5rem' }}>
        <Stack gap={4} orientation="horizontal">
          <TextInput id="search" labelText="Account Number" value={searchQuery}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setSearchQuery(e.target.value)}
            onKeyDown={(e: React.KeyboardEvent) => e.key === 'Enter' && handleSearch()} />
          <Button onClick={handleSearch} style={{ alignSelf: 'flex-end' }}>Search</Button>
        </Stack>
      </div>
      {account && !editing && (
        <div className="detail-card">
          <h3 style={{ marginTop: 0 }}>Account {account.accountNumber}</h3>
          <div className="detail-row"><span className="label">Customer</span><span className="value">{account.customerNumber}</span></div>
          <div className="detail-row"><span className="label">Type</span><span className="value">{account.accountType}</span></div>
          <div className="detail-row"><span className="label">Interest Rate</span><span className="value">{account.interestRate}%</span></div>
          <div className="detail-row"><span className="label">Overdraft Limit</span><span className="value">{Number(account.overdraftLimit).toLocaleString()}</span></div>
          <div className="detail-row"><span className="label">Available Balance</span><span className="value" style={{ color: '#24a148', fontWeight: 700 }}>${account.availableBalance}</span></div>
          <div className="detail-row"><span className="label">Actual Balance</span><span className="value">${account.actualBalance}</span></div>
          <div className="detail-row"><span className="label">Opened</span><span className="value">{account.openedDate}</span></div>
          <div style={{ marginTop: '1rem' }}>
            <Button onClick={() => { setEditData({ accountType: account.accountType, interestRate: account.interestRate, overdraftLimit: account.overdraftLimit }); setEditing(true); }}>Edit Account</Button>
          </div>
        </div>
      )}
      {account && editing && (
        <div className="form-card">
          <h3 style={{ marginTop: 0 }}>Edit Account {account.accountNumber}</h3>
          <Form onSubmit={(e: React.FormEvent) => { e.preventDefault(); handleUpdate(); }}>
            <Stack gap={4}>
              <Dropdown id="editType" titleText="Account Type" label="Select account type" items={accountTypes}
                selectedItem={editData.accountType}
                onChange={({ selectedItem }: { selectedItem: string }) => setEditData(d => ({ ...d, accountType: selectedItem }))} />
              <TextInput id="editRate" labelText="Interest Rate (%)" value={editData.interestRate || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setEditData(d => ({ ...d, interestRate: e.target.value }))} />
              <NumberInput id="editOD" label="Overdraft Limit" value={editData.overdraftLimit || 0} min={0} max={99999}
                onChange={(_: any, { value }: { value: string | number }) => setEditData(d => ({ ...d, overdraftLimit: Number(value) || 0 }))} />
              <Stack gap={3} orientation="horizontal">
                <Button type="submit">Save Changes</Button>
                <Button kind="secondary" onClick={() => setEditing(false)}>Cancel</Button>
              </Stack>
            </Stack>
          </Form>
        </div>
      )}
      {error && <Modal open modalHeading="Error" passiveModal onRequestClose={() => setError('')}><p>{error}</p></Modal>}
      {success && <Modal open modalHeading="Success" passiveModal onRequestClose={() => navigate('/')}><p>{success}</p></Modal>}
    </>
  );
};

export default AccountDetailsPage;
