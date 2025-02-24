import { writable } from 'svelte/store';
import { SignIn } from '../../wailsjs/go/cognito/CognitoAuth';
import { CompleteNewPassword } from '../../wailsjs/go/cognito/CognitoAuth';

interface AuthState {
  isAuthenticated: boolean;
  accessToken: string | null;
  idToken: string | null;
  refreshToken: string | null;
  error: string | null;
  requiresNewPassword: boolean;
}

const { subscribe, set, update } = writable<AuthState>({
  isAuthenticated: false,
  accessToken: null,
  idToken: null,
  refreshToken: null,
  error: null,
  requiresNewPassword: false
});

interface AuthResponse {
  AccessToken: string;
  IdToken: string;
  RefreshToken: string;
  NewPasswordRequired: boolean;
}

export const auth = {
  subscribe,
  login: async (username: string, password: string) => {
    try {
      const response = await SignIn(username, password);
      console.log('SignIn response:', response);
      
      // Handle new password required case
      if (response.NewPasswordRequired) {
        console.log('New password required detected!');
        update(state => ({
          ...state,
          isAuthenticated: false,
          requiresNewPassword: true,
          error: null
        }));
        console.log('State updated in store');
        console.log('Local auth properties updated:', auth);
      } else {
        // Normal login success
        update(state => ({
          ...state,
          isAuthenticated: true,
          accessToken: response.AccessToken,
          idToken: response.IdToken,
          refreshToken: response.RefreshToken,
          requiresNewPassword: false,
          error: null
        }));
      }
    } catch (error: any) {
      update(state => ({ ...state, error: error.message }));
      throw error;
    }
  },
  completeNewPassword: async (newPassword: string) => {
    try {
      const response = await CompleteNewPassword(newPassword);
      update(state => ({
        ...state,
        isAuthenticated: true,
        accessToken: response.AccessToken,
        idToken: response.IdToken,
        refreshToken: response.RefreshToken,
        requiresNewPassword: false,
        error: null
      }));
    } catch (error: any) {
      update(state => ({ ...state, error: error.message }));
      throw error;
    }
  },
  logout: () => {
    set({
      isAuthenticated: false,
      accessToken: null,
      idToken: null,
      refreshToken: null,
      error: null,
      requiresNewPassword: false
    });
  }
}; 