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

<style>
  /* Improve log panel contrast */
  :global(.wails-log) {
    color: #000 !important; /* Black text for maximum contrast */
    background-color: #f5f5f5 !important;
    border: 1px solid #ddd;
    padding: 10px;
    font-family: monospace;
    font-size: 14px;
    line-height: 1.4;
  }

  :global(.wails-log-entry) {
    margin-bottom: 5px;
    color: #000 !important;
  }

  /* Add specific styles for different log levels */
  :global(.wails-log-error) {
    color: #d32f2f !important; /* Red for errors */
    font-weight: bold;
  }

  :global(.wails-log-warning) {
    color: #f57c00 !important; /* Orange for warnings */
  }

  :global(.wails-log-info) {
    color: #0277bd !important; /* Blue for info */
  }
</style>

{#if authState?.isAuthenticated}
  <Dashboard />
{:else if authState?.requiresNewPassword}
  <ChangePassword />
{:else}
  <Login />
{/if}
