import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import {
  Breadcrumb,
  BreadcrumbItem,
  TextInput,
  DatePicker,
  DatePickerInput,
  Button,
  Modal,
  Dropdown,
  Form,
  Stack,
} from '@carbon/react';
import { customerApi } from '../../api/client';
import type { CreateCustomerRequest, Customer } from '../../api/types';

const titles = ['Mr', 'Mrs', 'Miss', 'Ms', 'Dr', 'Prof'];

const CustomerCreationPage: React.FC = () => {
  const [formData, setFormData] = useState<CreateCustomerRequest>({
    title: 'Mr',
    firstName: '',
    lastName: '',
    addressLine1: '',
    addressLine2: '',
    addressLine3: '',
    dateOfBirth: '',
  });
  const navigate = useNavigate();
  const [result, setResult] = useState<Customer | null>(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    try {
      const customer = await customerApi.create(formData);
      setResult(customer);
    } catch (err: any) {
      setError(err.response?.data?.error || err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <Breadcrumb>
        <BreadcrumbItem><Link to="/">Dashboard</Link></BreadcrumbItem>
        <BreadcrumbItem isCurrentPage>Create Customer</BreadcrumbItem>
      </Breadcrumb>

      <div className="page-header">
        <h2>Create Customer</h2>
        <p>Register a new bank customer with credit check</p>
      </div>

      <div className="form-card">
        <Form onSubmit={handleSubmit}>
          <Stack gap={6}>
            <Dropdown
              id="title"
              titleText="Title"
              label="Select title"
              items={titles}
              selectedItem={formData.title}
              onChange={({ selectedItem }: { selectedItem: string }) =>
                setFormData(prev => ({ ...prev, title: selectedItem }))
              }
            />
            <TextInput id="firstName" labelText="First Name" value={formData.firstName}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFormData(prev => ({ ...prev, firstName: e.target.value }))}
              required />
            <TextInput id="lastName" labelText="Last Name" value={formData.lastName}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFormData(prev => ({ ...prev, lastName: e.target.value }))}
              required />
            <TextInput id="address1" labelText="Address Line 1" value={formData.addressLine1}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFormData(prev => ({ ...prev, addressLine1: e.target.value }))} />
            <TextInput id="address2" labelText="Address Line 2" value={formData.addressLine2}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFormData(prev => ({ ...prev, addressLine2: e.target.value }))} />
            <TextInput id="address3" labelText="Town / City" value={formData.addressLine3}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFormData(prev => ({ ...prev, addressLine3: e.target.value }))} />
            <DatePicker datePickerType="single"
              onChange={([date]: Date[]) => {
                if (date) setFormData(prev => ({ ...prev, dateOfBirth: date.toISOString().split('T')[0] }));
              }}>
              <DatePickerInput id="dob" labelText="Date of Birth" placeholder="mm/dd/yyyy" />
            </DatePicker>
            <Button type="submit" disabled={loading}>
              {loading ? 'Creating...' : 'Create Customer'}
            </Button>
          </Stack>
        </Form>
      </div>

      {result && (
        <Modal open modalHeading="Customer Created Successfully" passiveModal onRequestClose={() => navigate('/')}>
          <div className="detail-card">
            <div className="detail-row"><span className="label">Customer Number</span><span className="value">{result.customerNumber}</span></div>
            <div className="detail-row"><span className="label">Name</span><span className="value">{result.title} {result.firstName} {result.lastName}</span></div>
            <div className="detail-row"><span className="label">Credit Score</span><span className="value">{result.creditScore}</span></div>
          </div>
        </Modal>
      )}
      {error && (
        <Modal open modalHeading="Error" passiveModal onRequestClose={() => setError('')}><p>{error}</p></Modal>
      )}
    </>
  );
};

export default CustomerCreationPage;
