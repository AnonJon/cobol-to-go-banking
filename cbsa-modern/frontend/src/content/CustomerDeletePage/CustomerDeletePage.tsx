import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Breadcrumb, BreadcrumbItem, TextInput, Button, Modal, Stack } from '@carbon/react';
import { customerApi } from '../../api/client';
import type { Customer } from '../../api/types';

const CustomerDeletePage: React.FC = () => {
  const [searchQuery, setSearchQuery] = useState('');
  const [customer, setCustomer] = useState<Customer | null>(null);
  const navigate = useNavigate();
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const handleSearch = async () => {
    setError(''); setSuccess('');
    try { setCustomer(await customerApi.get(searchQuery)); }
    catch (err: any) { setError(err.response?.data?.error || 'Customer not found'); setCustomer(null); }
  };

  const handleDelete = async () => {
    if (!customer) return;
    try {
      await customerApi.delete(customer.customerNumber);
      setSuccess(`Customer ${customer.customerNumber} deleted successfully`);
      setCustomer(null); setConfirmOpen(false);
    } catch (err: any) { setError(err.response?.data?.error || 'Delete failed'); setConfirmOpen(false); }
  };

  return (
    <>
      <Breadcrumb>
        <BreadcrumbItem><Link to="/">Dashboard</Link></BreadcrumbItem>
        <BreadcrumbItem isCurrentPage>Delete Customer</BreadcrumbItem>
      </Breadcrumb>
      <div className="page-header"><h2>Delete Customer</h2><p>Permanently remove a customer and all associated accounts</p></div>
      <div className="form-card" style={{ marginBottom: '1.5rem' }}>
        <Stack gap={4} orientation="horizontal">
          <TextInput id="search" labelText="Customer Number" value={searchQuery}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setSearchQuery(e.target.value)}
            onKeyDown={(e: React.KeyboardEvent) => e.key === 'Enter' && handleSearch()} />
          <Button onClick={handleSearch} style={{ alignSelf: 'flex-end' }}>Search</Button>
        </Stack>
      </div>
      {customer && (
        <div className="detail-card">
          <h3 style={{ marginTop: 0 }}>Customer {customer.customerNumber}</h3>
          <div className="detail-row"><span className="label">Name</span><span className="value">{customer.title} {customer.firstName} {customer.lastName}</span></div>
          <div className="detail-row"><span className="label">Address</span><span className="value">{customer.addressLine1}, {customer.addressLine2}, {customer.addressLine3}</span></div>
          <div className="detail-row"><span className="label">Date of Birth</span><span className="value">{customer.dateOfBirth}</span></div>
          <div style={{ marginTop: '1rem' }}>
            <Button kind="danger" onClick={() => setConfirmOpen(true)}>Delete Customer</Button>
          </div>
        </div>
      )}
      <Modal open={confirmOpen} modalHeading="Confirm Deletion" primaryButtonText="Delete" secondaryButtonText="Cancel" danger
        onRequestSubmit={handleDelete} onRequestClose={() => setConfirmOpen(false)}>
        <p>This will permanently delete customer <strong>{customer?.customerNumber}</strong> and all associated accounts. This cannot be undone.</p>
      </Modal>
      {error && <Modal open modalHeading="Error" passiveModal onRequestClose={() => setError('')}><p>{error}</p></Modal>}
      {success && <Modal open modalHeading="Success" passiveModal onRequestClose={() => navigate('/')}><p>{success}</p></Modal>}
    </>
  );
};

export default CustomerDeletePage;
