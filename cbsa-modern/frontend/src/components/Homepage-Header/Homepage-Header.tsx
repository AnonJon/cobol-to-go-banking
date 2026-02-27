import React from 'react';
import { Link } from 'react-router-dom';
import {
  Header,
  HeaderName,
  HeaderNavigation,
  HeaderMenuItem,
} from '@carbon/react';

const HomepageHeader: React.FC = () => (
  <Header aria-label="CBSA Modern">
    <HeaderName as={Link} to="/" prefix="CBSA">
      Modern Banking
    </HeaderName>
    <HeaderNavigation aria-label="Navigation">
      <HeaderMenuItem as={Link} to="/profile/Admin">
        Admin
      </HeaderMenuItem>
    </HeaderNavigation>
  </Header>
);

export default HomepageHeader;
