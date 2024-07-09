// src/components/FAQ.js
import React from 'react';
import { Container, Typography, Box } from '@mui/material';

const FAQ = () => {
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
        <Typography component="h1" variant="h5">
          FAQ
        </Typography>
        <Typography component="p" variant="body1" sx={{ mt: 3 }}>
          This is some test text to check the rendering of the FAQ page.
        </Typography>
      </Box>
    </Container>
  );
};

export default FAQ;
