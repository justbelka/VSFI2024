// src/contexts/AuthContext.js
import React, { createContext, useState, useEffect } from 'react';

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loggedInUser = JSON.parse(localStorage.getItem('user'));
    if (loggedInUser) {
      setUser(loggedInUser);
      fetchBalance(loggedInUser.username);
    }
    setLoading(false);
  }, []);

  const fetchBalance = async (username) => {
    try {
      const response = await fetch(`/api/balance?username=${username}`);
      if (!response.ok) {
        if (response.status === 401) {
          logout();
        }
        throw new Error('Failed to fetch balance');
      }
      const data = await response.json();
      setUser((prevUser) => ({ ...prevUser, coins: data.balance }));
    } catch (error) {
      console.error('Failed to fetch balance:', error);
    }
  };

  const login = async (userData) => {
    try {
      setUser(userData);
      localStorage.setItem('user', JSON.stringify(userData));
      fetchBalance(userData.username);
    } catch (error) {
      console.error('Login failed:', error);
      throw error;
    }
  };

  const logout = () => {
    localStorage.removeItem('user');
    setUser(null);
  };

  const updateUserBalance = async () => {
    await fetchBalance(user.username);
  };

  return (
    <AuthContext.Provider value={{ user, loading, login, logout, updateUserBalance }}>
      {children}
    </AuthContext.Provider>
  );
};
