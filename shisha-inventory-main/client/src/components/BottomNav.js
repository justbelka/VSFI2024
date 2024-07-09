// src/components/BottomNav.js
import React, { useState } from 'react';
import { BottomNavigation, BottomNavigationAction, Paper } from '@mui/material';
import WorkspacePremiumIcon from '@mui/icons-material/WorkspacePremium';
import FolderSharedIcon from '@mui/icons-material/FolderShared';
import UploadIcon from '@mui/icons-material/CloudUpload';
import CoinIcon from '@mui/icons-material/AccountBalanceWallet';
import FAQIcon from '@mui/icons-material/Help';
import Collections from '@mui/icons-material/Collections';

import { Link } from 'react-router-dom';

const BottomNav = () => {
  const [value, setValue] = useState(0);

  return (
    <Paper sx={{ position: 'fixed', bottom: 0, left: 0, right: 0 }} elevation={3}>
      <BottomNavigation
        value={value}
        onChange={(event, newValue) => {
          setValue(newValue);
        }}
         sx={{
          "& .MuiBottomNavigationAction-label, svg": {
            color: "white", fontWeight: 'medium'
          }
        }}
      >
        <BottomNavigationAction label="Premium Shisha" icon={<WorkspacePremiumIcon />} component={Link} to="/" />
        <BottomNavigationAction label="Users shisha" icon={<FolderSharedIcon />} component={Link} to="/users-shisha" />
        <BottomNavigationAction label="Upload shisha" icon={<UploadIcon />} component={Link}  to="/upload" />
        <BottomNavigationAction label="Transfer Coins" icon={<CoinIcon />} component={Link} to="/transfer" />
        <BottomNavigationAction label="FAQ" icon={<FAQIcon />} component={Link} to="/faq" />
        <BottomNavigationAction label="Purchased Images" icon={<Collections />} component={Link} to="/purchased" />
      </BottomNavigation>
    </Paper>
  );
};

export default BottomNav;
