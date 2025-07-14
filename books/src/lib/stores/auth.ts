// src/lib/stores/auth.ts
import { writable } from 'svelte/store';

export interface AuthState {
  isAuthenticated: boolean;
  user: { email: string } | null;
  accessToken: string | null;
}

const initialState: AuthState = {
  isAuthenticated: false,
  user: null,
  accessToken: null,
};

export const authStore = writable<AuthState>(initialState);

export const auth = {
  // Set user as authenticated
  setAuth: (accessToken: string, userEmail: string) => {
    authStore.set({
      isAuthenticated: true,
      user: { email: userEmail },
      accessToken,
    });
  },

  // Clear authentication
  clearAuth: () => {
    authStore.set(initialState);
  },

  // Update access token (for refresh)
  updateToken: (accessToken: string) => {
    authStore.update(state => ({
      ...state,
      accessToken,
    }));
  },

  // Check if user is authenticated
  isAuthenticated: () => {
    let isAuth = false;
    authStore.subscribe(state => {
      isAuth = state.isAuthenticated;
    })();
    return isAuth;
  },
};