import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Breadcrumb, BreadcrumbItem, TextInput, Button, Modal, Stack } from '@carbon/react';
import { accountApi } from '../../api/client';
import type { Account } from '../../api/types';

const AccountDeletePage: React.FC = () => {
  const [searchQuery, setSearchQuery] = useState('');
  const [account, setAccount] = useState<Account | null>(null);
  const navigate = useNavigate();
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const handleSearch = async () => {
    setError(''); setSuccess('');
    try { setAccount(await accountApi.get(searchQuery)); }
    catch (err: any) { setError(err.response?.data?.error || 'Account not found'); setAccount(null); }
  };

  const handleDelete = async () => {
    if (!account) return;
    try {
      await accountApi.delete(account.accountNumber);
      setSuccess(`Account ${account.accountNumber} deleted successfully`);
      setAccount(null); setConfirmOpen(false);
    } catch (err: any) { setError(err.response?.data?.error || 'Delete failed'); setConfirmOpen(false); }
  };

  return (
    <>
      <Breadcrumb>
        <BreadcrumbItem><Link to="/">Dashboard</Link></BreadcrumbItem>
        <BreadcrumbItem isCurrentPage>Delete Account</BreadcrumbItem>
      </Breadcrumb>
      <div className="page-header"><h2>Delete Account</h2><p>Permanently close a bank account</p></div>
      <div className="form-card" style={{ marginBottom: '1.5rem' }}>
        <Stack gap={4} orientation="horizontal">
          <TextInput id="search" labelText="Account Number" value={searchQuery}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setSearchQuery(e.target.value)}
            onKeyDown={(e: React.KeyboardEvent) => e.key === 'Enter' && handleSearch()} />
          <Button onClick={handleSearch} style={{ alignSelf: 'flex-end' }}>Search</Button>
        </Stack>
      </div>
      {account && (
        <div className="detail-card">
          <h3 style={{ marginTop: 0 }}>Account {account.accountNumber}</h3>
          <div className="detail-row"><span className="label">Customer</span><span className="value">{account.customerNumber}</span></div>
          <div className="detail-row"><span className="label">Type</span><span className="value">{account.accountType}</span></div>
          <div className="detail-row"><span className="label">Balance</span><span className="value">${account.availableBalance}</span></div>
          <div style={{ marginTop: '1rem' }}>
            <Button kind="danger" onClick={() => setConfirmOpen(true)}>Delete Account</Button>
          </div>
        </div>
      )}
      <Modal open={confirmOpen} modalHeading="Confirm Deletion" primaryButtonText="Delete" secondaryButtonText="Cancel" danger
        onRequestSubmit={handleDelete} onRequestClose={() => setConfirmOpen(false)}>
        <p>This will permanently delete account <strong>{account?.accountNumber}</strong>. This cannot be undone.</p>
      </Modal>
      {error && <Modal open modalHeading="Error" passiveModal onRequestClose={() => setError('')}><p>{error}</p></Modal>}
      {success && <Modal open modalHeading="Success" passiveModal onRequestClose={() => navigate('/')}><p>{success}</p></Modal>}
    </>
  );
};

export default AccountDeletePage;
