// src/components/Transfer.js
import React, { useState, useContext } from 'react';
import { Container, TextField, Button, Typography, Box } from '@mui/material';
import { AuthContext } from '../contexts/AuthContext';

const Transfer = () => {
  const { user, updateUserBalance } = useContext(AuthContext);
  const [toUsername, setToUsername] = useState('');
  const [amount, setAmount] = useState('');
  const [message, setMessage] = useState('');

  const handleTransfer = async () => {
    try {
      const response = await fetch('/api/transfer', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          from_username: user.username,
          to_username: toUsername,
          amount: parseInt(amount, 10),
        }),
      });
      const data = await response.json();
      if (!response.ok) {
        setMessage(data.error);
      } else {
        setMessage(data.message);
        await updateUserBalance();
      }
    } catch (error) {
      setMessage('Transfer failed');
    }
  };

  return (
    <Container component="main" maxWidth="md">
      <Box
        sx={{
          marginTop: 8,
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
        }}
      >
        <Typography variant="h4" gutterBottom>
          Transfer Coins
        </Typography>
        <Box component="form" noValidate sx={{ mt: 1 }}>
          <TextField
            margin="normal"
            required
            fullWidth
            id="recipient"
            label="Recipient"
            name="recipient"
            autoComplete="recipient"
            autoFocus
            value={toUsername}
            onChange={(e) => setToUsername(e.target.value)}
          />
          <TextField
            margin="normal"
            required
            fullWidth
            name="amount"
            label="Amount"
            type="number"
            id="amount"
            autoComplete="amount"
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
          />
          <Button variant="contained" color="primary" fullWidth onClick={handleTransfer} sx={{ mt: 3, mb: 2 }}>
            Transfer
          </Button>
          {message && <Typography variant="body1" color="error">{message}</Typography>}
        </Box>
      </Box>
    </Container>
  );
};

export default Transfer;
