import React from 'react';
import './app.scss';
import { Theme, Header, HeaderName, HeaderGlobalBar, HeaderGlobalAction } from '@carbon/react';
import { BrowserRouter, Routes, Route, useNavigate, useLocation } from 'react-router-dom';
import {
  SideNav,
  SideNavItems,
  SideNavMenu,
  SideNavMenuItem,
  SideNavLink,
} from '@carbon/react';
import HomePage from './content/HomePage/HomePage';
import AdminPage from './content/AdminPage/AdminPage';
import CustomerCreationPage from './content/CustomerCreationPage/CustomerCreationPage';
import AccountCreationPage from './content/AccountCreationPage/AccountCreationPage';
import CustomerDetailsPage from './content/CustomerDetailsPage/CustomerDetailsPage';
import AccountDetailsPage from './content/AccountDetailsPage/AccountDetailsPage';
import CustomerDeletePage from './content/CustomerDeletePage/CustomerDeletePage';
import AccountDeletePage from './content/AccountDeletePage/AccountDeletePage';

const AppShell: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();

  return (
    <>
      <Theme theme="g100">
        <Header aria-label="CBSA Modern Banking">
          <HeaderName prefix="CBSA" onClick={() => navigate('/')} style={{ cursor: 'pointer' }}>
            Modern Banking
          </HeaderName>
          <HeaderGlobalBar>
            <HeaderGlobalAction aria-label="Admin" onClick={() => navigate('/profile/Admin')}>
              <svg width="20" height="20" viewBox="0 0 32 32" fill="currentColor">
                <path d="M16 4a5 5 0 1 1 0 10 5 5 0 0 1 0-10m0-2a7 7 0 1 0 0 14 7 7 0 0 0 0-14zM26 30h-2v-5a5 5 0 0 0-5-5h-6a5 5 0 0 0-5 5v5H6v-5a7 7 0 0 1 7-7h6a7 7 0 0 1 7 7z"/>
              </svg>
            </HeaderGlobalAction>
          </HeaderGlobalBar>
        </Header>
      </Theme>
      <div className="app-shell">
        <Theme theme="g10">
          <SideNav aria-label="Side navigation" className="app-sidenav" expanded isRail={false} isFixedNav>
            <SideNavItems>
              <SideNavLink onClick={() => navigate('/')} isActive={location.pathname === '/'}>
                Dashboard
              </SideNavLink>
              <SideNavMenu title="Customers" defaultExpanded>
                <SideNavMenuItem onClick={() => navigate('/Admin/customer_creation')}
                  isActive={location.pathname === '/Admin/customer_creation'}>
                  Create Customer
                </SideNavMenuItem>
                <SideNavMenuItem onClick={() => navigate('/Admin/customer_details')}
                  isActive={location.pathname === '/Admin/customer_details'}>
                  Customer Details
                </SideNavMenuItem>
                <SideNavMenuItem onClick={() => navigate('/Admin/customer_deletion')}
                  isActive={location.pathname === '/Admin/customer_deletion'}>
                  Delete Customer
                </SideNavMenuItem>
              </SideNavMenu>
              <SideNavMenu title="Accounts" defaultExpanded>
                <SideNavMenuItem onClick={() => navigate('/Admin/account_creation')}
                  isActive={location.pathname === '/Admin/account_creation'}>
                  Create Account
                </SideNavMenuItem>
                <SideNavMenuItem onClick={() => navigate('/Admin/account_details')}
                  isActive={location.pathname === '/Admin/account_details'}>
                  Account Details
                </SideNavMenuItem>
                <SideNavMenuItem onClick={() => navigate('/Admin/account_deletion')}
                  isActive={location.pathname === '/Admin/account_deletion'}>
                  Delete Account
                </SideNavMenuItem>
              </SideNavMenu>
            </SideNavItems>
          </SideNav>
        </Theme>
        <main className="app-main">
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/profile/Admin" element={<AdminPage />} />
            <Route path="/Admin/customer_creation" element={<CustomerCreationPage />} />
            <Route path="/Admin/account_creation" element={<AccountCreationPage />} />
            <Route path="/Admin/customer_details" element={<CustomerDetailsPage />} />
            <Route path="/Admin/account_details" element={<AccountDetailsPage />} />
            <Route path="/Admin/customer_deletion" element={<CustomerDeletePage />} />
            <Route path="/Admin/account_deletion" element={<AccountDeletePage />} />
          </Routes>
        </main>
      </div>
    </>
  );
};

const App: React.FC = () => (
  <BrowserRouter>
    <AppShell />
  </BrowserRouter>
);

export default App;
