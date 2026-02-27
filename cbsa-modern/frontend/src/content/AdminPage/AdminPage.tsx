import React from 'react';
import { useNavigate } from 'react-router-dom';

const adminLinks = [
  { to: '/Admin/customer_creation', label: 'Create Customer', desc: 'Register a new bank customer', icon: '+', color: 'blue' },
  { to: '/Admin/account_creation', label: 'Create Account', desc: 'Open a new bank account', icon: '+', color: 'green' },
  { to: '/Admin/customer_details', label: 'Customer Details', desc: 'View and update customer information', icon: '\u2139', color: 'teal' },
  { to: '/Admin/account_details', label: 'Account Details', desc: 'View and update account information', icon: '\u2139', color: 'purple' },
  { to: '/Admin/customer_deletion', label: 'Delete Customer', desc: 'Remove a customer and their accounts', icon: '\u2715', color: 'red' },
  { to: '/Admin/account_deletion', label: 'Delete Account', desc: 'Close a bank account', icon: '\u2715', color: 'red' },
];

const AdminPage: React.FC = () => {
  const navigate = useNavigate();

  return (
    <>
      <div className="page-header">
        <h2>Administration</h2>
        <p>Manage customers and accounts</p>
      </div>
      <div className="admin-grid">
        {adminLinks.map(({ to, label, desc, icon, color }) => (
          <div key={to} className="admin-tile" onClick={() => navigate(to)}>
            <div className={`tile-icon ${color}`}>{icon}</div>
            <h4>{label}</h4>
            <p>{desc}</p>
          </div>
        ))}
      </div>
    </>
  );
};

export default AdminPage;
