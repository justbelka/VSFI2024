// src/components/Upload.js
import React, { useState, useContext } from 'react';
import { Container, Button, Typography, Box, TextField } from '@mui/material';
import { AuthContext } from '../contexts/AuthContext';

const Upload = () => {
  const { user, updateUserBalance } = useContext(AuthContext);
  const [file, setFile] = useState(null);
  const [message, setMessage] = useState('');

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!file) {
      setMessage('Please select a file');
      return;
    }


    const formData = new FormData();
    formData.append('file', file);
    if (localStorage.getItem('user')) {
      formData.append('username', user.username);
    }

    try {
      const response = await fetch('/api/upload', {
        method: 'POST',
        body: formData,
      });
      const data = await response.json();
      if (!response.ok) {
        setMessage(data.error);
      } else {
        setMessage(data.message);
        await updateUserBalance();
      }
    } catch (error) {
      setMessage('File upload failed');
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
      <Typography component="h1" variant="h5" gutterBottom>
        Upload Shisha
      </Typography>
      <Box component="form" onSubmit={handleSubmit} noValidate sx={{ mt: 1 }}>
          <TextField
            type="file"
            fullWidth
            onChange={handleFileChange}
            InputLabelProps={{ shrink: true }}
            sx={{ mb: 2 }}
          />
          <Button
            type="submit"
            fullWidth
            variant="contained"
            sx={{ mt: 3, mb: 2 }}
          >
            Upload
          </Button>
          </Box>
      {message && <Typography variant="body1" color="error">{message}</Typography>}
    </Box>
    </Container >
  );
};

export default Upload;
