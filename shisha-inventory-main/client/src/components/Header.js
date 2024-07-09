// src/components/Header.js
import React, { useContext } from 'react';
import { AppBar, Toolbar, Typography, Button, Box } from '@mui/material';
import { Link, useNavigate } from 'react-router-dom';
import { AuthContext } from '../contexts/AuthContext';
import { AccountBalanceWallet as BalanceIcon, AccountBox } from '@mui/icons-material';
import logo from '../logo.png'; // Ensure this path is correct

const Header = () => {
  const { user, logout } = useContext(AuthContext);
  const navigate = useNavigate();

  const handleHomeClick = () => {
    navigate('/');
  };

  return (
    <AppBar position="static" sx={{ backgroundColor: '#343917' }}>
      <Toolbar>
        <Box sx={{ flexGrow: 1, display: 'flex', alignItems: 'center', cursor: 'pointer' }} onClick={handleHomeClick}>
          <img src={logo} alt="Logo" style={{ height: '40px', marginRight: '10px' }} />
          <Typography variant="h6" sx={{ ml: 1 }}>
            Shisha inventory
          </Typography>
        </Box>
        {user ? (
          <>
            <Box sx={{ display: 'flex', alignItems: 'center', mr: 2 }}>
              <BalanceIcon />
              <Typography variant="h7" sx={{ ml: 0 }}>
                Balance: {user.coins}
              </Typography>
            </Box>
            <Box sx={{ display: 'flex', alignItems: 'center', mr: 2 }}>
              <AccountBox />
              <Typography variant="h7" sx={{ ml: 0 }}>
                User: {user.username}
              </Typography>
            </Box>
            <Button color="inherit" onClick={logout}>Logout</Button>
          </>
        ) : (
          <>
            <Button color="inherit" component={Link} to="/login">Login</Button>
            <Button color="inherit" component={Link} to="/register">Register</Button>
          </>
        )}
      </Toolbar>
    </AppBar>
  );
};

export default Header;
