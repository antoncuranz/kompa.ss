"use client"
import React, { createContext, useContext, ReactNode } from 'react';

interface UserContextType {
  username: string;
}

const UserContext = createContext<UserContextType | undefined>(undefined);

export function UserProvider({
  children, username
}: {
  username: string;
  children: ReactNode;
}) {
  
  return (
    <UserContext.Provider value={{ username }}>
      {children}
    </UserContext.Provider>
  );
}

export function useUser(): UserContextType {
  const context = useContext(UserContext);
  if (!context) {
    throw new Error('useUser must be used within a UserProvider');
  }
  return context;
}