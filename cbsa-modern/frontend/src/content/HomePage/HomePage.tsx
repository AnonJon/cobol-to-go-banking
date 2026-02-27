import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Grid, Column } from '@carbon/react';
import { infoApi, customerApi, accountApi, transactionApi } from '../../api/client';

const HomePage: React.FC = () => {
  const navigate = useNavigate();
  const [companyName, setCompanyName] = useState('');
  const [stats, setStats] = useState({ customers: 0, accounts: 0, transactions: 0 });

  useEffect(() => {
    infoApi.companyName().then(setCompanyName).catch(() => {});
    customerApi.list(1, 0).then(d => setStats(s => ({ ...s, customers: d.total }))).catch(() => {});
    accountApi.list(1, 0).then(d => setStats(s => ({ ...s, accounts: d.total }))).catch(() => {});
    transactionApi.list(1, 0).then(d => setStats(s => ({ ...s, transactions: d.total }))).catch(() => {});
  }, []);

  return (
    <>
      <div className="dashboard-hero">
        <h1>{companyName || 'CBSA Modern Banking'}</h1>
        <p>Teller dashboard — manage customers, accounts, and transactions</p>
      </div>

      <Grid narrow>
        <Column lg={4} md={4} sm={4}>
          <div className="stat-card" onClick={() => navigate('/Admin/customer_details')} style={{ cursor: 'pointer' }}>
            <h4>Total Customers</h4>
            <p className="stat-value">{stats.customers.toLocaleString()}</p>
          </div>
        </Column>
        <Column lg={4} md={4} sm={4}>
          <div className="stat-card green" onClick={() => navigate('/Admin/account_details')} style={{ cursor: 'pointer' }}>
            <h4>Total Accounts</h4>
            <p className="stat-value">{stats.accounts.toLocaleString()}</p>
          </div>
        </Column>
        <Column lg={4} md={4} sm={4}>
          <div className="stat-card purple">
            <h4>Transactions Processed</h4>
            <p className="stat-value">{stats.transactions.toLocaleString()}</p>
          </div>
        </Column>
      </Grid>

      <h3 style={{ marginTop: '2rem', marginBottom: '1rem' }}>Quick Actions</h3>
      <div className="admin-grid">
        <div className="admin-tile" onClick={() => navigate('/Admin/customer_creation')}>
          <div className="tile-icon blue">+</div>
          <h4>New Customer</h4>
          <p>Register a new bank customer</p>
        </div>
        <div className="admin-tile" onClick={() => navigate('/Admin/account_creation')}>
          <div className="tile-icon green">+</div>
          <h4>New Account</h4>
          <p>Open a bank account for a customer</p>
        </div>
        <div className="admin-tile" onClick={() => navigate('/Admin/customer_details')}>
          <div className="tile-icon teal">&#8505;</div>
          <h4>Look Up Customer</h4>
          <p>View and update customer information</p>
        </div>
        <div className="admin-tile" onClick={() => navigate('/Admin/account_details')}>
          <div className="tile-icon purple">&#8505;</div>
          <h4>Look Up Account</h4>
          <p>View and update account details</p>
        </div>
      </div>
    </>
  );
};

export default HomePage;
