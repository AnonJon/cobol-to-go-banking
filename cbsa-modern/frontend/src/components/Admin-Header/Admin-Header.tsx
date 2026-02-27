import React from 'react';
import { Link } from 'react-router-dom';
import {
  HeaderNavigation,
  HeaderMenuItem,
} from '@carbon/react';

const AdminHeader: React.FC = () => (
  <HeaderNavigation aria-label="Admin Navigation">
    <HeaderMenuItem as={Link} to="/Admin/customer_creation">
      Create Customer
    </HeaderMenuItem>
    <HeaderMenuItem as={Link} to="/Admin/account_creation">
      Create Account
    </HeaderMenuItem>
    <HeaderMenuItem as={Link} to="/Admin/customer_details">
      Customer Details
    </HeaderMenuItem>
    <HeaderMenuItem as={Link} to="/Admin/account_details">
      Account Details
    </HeaderMenuItem>
    <HeaderMenuItem as={Link} to="/Admin/customer_deletion">
      Delete Customer
    </HeaderMenuItem>
    <HeaderMenuItem as={Link} to="/Admin/account_deletion">
      Delete Account
    </HeaderMenuItem>
  </HeaderNavigation>
);

export default AdminHeader;
