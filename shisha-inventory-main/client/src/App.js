// src/App.js
import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Container, Box } from '@mui/material';
import { ToastContainer } from 'react-toastify';
import { ThemeProvider } from '@mui/material/styles';
import 'react-toastify/dist/ReactToastify.css';
import PremiumShishaList from './components/PremiumShishaList';
import UserShishaList from './components/UserShishaList'
import Header from './components/Header';
import BottomNav from './components/BottomNav';
import Login from './components/Login';
import Register from './components/Register';
import Upload from './components/Upload';
import Transfer from './components/Transfer';
import FAQ from './components/FAQ';
import PurchasedImages from './components/PurchasedImages';
import { AuthProvider } from './contexts/AuthContext';
import theme from './theme';
import './index.css'; // Ensure global CSS is imported

const App = () => {
  return (
    <ThemeProvider theme={theme}>
      <AuthProvider>
        <Router>
          <Header />
          <Box sx={{ position: 'relative', overflow: 'hidden', minHeight: '100vh' }}>
            <video 
              autoPlay 
              loop 
              muted 
              style={{
                position: 'fixed',
                right: '0',
                bottom: '0',
                minWidth: '100%',
                minHeight: '100%',
                zIndex: '-1',
                filter: 'brightness(0.5)', // Adjust brightness if necessary
              }}>
              <source src="/woodpecker.mp4" type="video/mp4" />
            </video>
          <Box sx={{ paddingTop: 2, minHeight: '100vh' }}>
            <Container>
              <Routes>
                <Route path="/" element={<PremiumShishaList />} />
                <Route path="/users-shisha" element={<UserShishaList />} />
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />
                <Route path="/upload" element={<Upload />} />
                <Route path="/transfer" element={<Transfer />} />
                <Route path="/faq" element={<FAQ />} />
                <Route path="/purchased" element={<PurchasedImages />} />
              </Routes>
            </Container>
          </Box>
          <BottomNav />
          </Box>
          <ToastContainer />
        </Router>
      </AuthProvider>
    </ThemeProvider>
  );
};

export default App;
