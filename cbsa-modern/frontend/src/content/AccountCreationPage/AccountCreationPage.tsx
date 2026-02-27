import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import {
  Breadcrumb,
  BreadcrumbItem,
  TextInput,
  Button,
  Modal,
  Dropdown,
  Form,
  Stack,
  NumberInput,
} from '@carbon/react';
import { accountApi } from '../../api/client';
import type { CreateAccountRequest, Account } from '../../api/types';

const accountTypes = ['ISA', 'SAVING', 'CURRENT', 'LOAN', 'MORTGAGE'];

const AccountCreationPage: React.FC = () => {
  const [formData, setFormData] = useState<CreateAccountRequest>({
    customerNumber: '', accountType: 'CURRENT', interestRate: '1.50', overdraftLimit: 0,
  });
  const navigate = useNavigate();
  const [result, setResult] = useState<Account | null>(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    try { setResult(await accountApi.create(formData)); }
    catch (err: any) { setError(err.response?.data?.error || err.message); }
    finally { setLoading(false); }
  };

  return (
    <>
      <Breadcrumb>
        <BreadcrumbItem><Link to="/">Dashboard</Link></BreadcrumbItem>
        <BreadcrumbItem isCurrentPage>Create Account</BreadcrumbItem>
      </Breadcrumb>
      <div className="page-header">
        <h2>Create Account</h2>
        <p>Open a new bank account for an existing customer</p>
      </div>
      <div className="form-card">
        <Form onSubmit={handleSubmit}>
          <Stack gap={6}>
            <TextInput id="customerNumber" labelText="Customer Number" value={formData.customerNumber}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFormData(prev => ({ ...prev, customerNumber: e.target.value }))} required />
            <Dropdown id="accountType" titleText="Account Type" label="Select account type" items={accountTypes}
              selectedItem={formData.accountType}
              onChange={({ selectedItem }: { selectedItem: string }) => setFormData(prev => ({ ...prev, accountType: selectedItem }))} />
            <TextInput id="interestRate" labelText="Interest Rate (%)" value={formData.interestRate}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFormData(prev => ({ ...prev, interestRate: e.target.value }))} />
            <NumberInput id="overdraftLimit" label="Overdraft Limit" value={formData.overdraftLimit} min={0} max={99999}
              onChange={(_: any, { value }: { value: string | number }) => setFormData(prev => ({ ...prev, overdraftLimit: Number(value) || 0 }))} />
            <Button type="submit" disabled={loading}>{loading ? 'Creating...' : 'Create Account'}</Button>
          </Stack>
        </Form>
      </div>
      {result && (
        <Modal open modalHeading="Account Created Successfully" passiveModal onRequestClose={() => navigate('/')}>
          <div className="detail-card">
            <div className="detail-row"><span className="label">Account Number</span><span className="value">{result.accountNumber}</span></div>
            <div className="detail-row"><span className="label">Type</span><span className="value">{result.accountType}</span></div>
            <div className="detail-row"><span className="label">Customer</span><span className="value">{result.customerNumber}</span></div>
          </div>
        </Modal>
      )}
      {error && <Modal open modalHeading="Error" passiveModal onRequestClose={() => setError('')}><p>{error}</p></Modal>}
    </>
  );
};

export default AccountCreationPage;
