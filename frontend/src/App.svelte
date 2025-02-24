<script lang="ts">
  import { auth } from './stores/auth';
  import LoginComponent from './components/Login.svelte';
  import DashboardComponent from './components/Dashboard.svelte';
  import ChangePasswordComponent from './components/ChangePassword.svelte';

  const Login = LoginComponent as any;
  const Dashboard = DashboardComponent as any;
  const ChangePassword = ChangePasswordComponent as any;

  let authState;
  // Subscribe to auth store changes
  auth.subscribe(state => {
    console.log('Auth state updated:', state);
    authState = state;
  });
</script>

<div style="position: fixed; bottom: 0; right: 0; background: #eee; padding: 5px; font-size: 12px;">
  Store state: isAuthenticated: {authState?.isAuthenticated}, requiresNewPassword: {authState?.requiresNewPassword}
  <br>
  <button on:click={() => {
    authState = { ...authState, requiresNewPassword: true, isAuthenticated: false };
    auth.update(state => ({...state, requiresNewPassword: true, isAuthenticated: false}));
  }}>Force Password Change Screen</button>
</div>

{#if authState?.isAuthenticated}
  <Dashboard />
{:else if authState?.requiresNewPassword}
  <ChangePassword />
{:else}
  <Login />
{/if}
