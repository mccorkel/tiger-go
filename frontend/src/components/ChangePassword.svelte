<script lang="ts">
  import { auth } from '../stores/auth';

  let newPassword = '';
  let error = '';

  async function handlePasswordChange() {
    try {
      error = '';
      await auth.completeNewPassword(newPassword);
    } catch (e: any) {
      error = e.message;
    }
  }
</script>

<div class="change-password-container">
  <div class="change-password-card">
    <h2>Change Password Required</h2>
    {#if error}
      <div class="error">{error}</div>
    {/if}
    <form on:submit|preventDefault={handlePasswordChange}>
      <input
        type="password"
        bind:value={newPassword}
        placeholder="New Password"
        required
      />
      <button type="submit">Update Password</button>
    </form>
  </div>
</div>

<style>
  .change-password-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    padding: 1rem;
  }

  .change-password-card {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    width: 100%;
    max-width: 400px;
  }

  .error {
    color: red;
    margin-bottom: 1rem;
  }

  input {
    width: 100%;
    padding: 0.5rem;
    margin-bottom: 1rem;
    border: 1px solid #ccc;
    border-radius: 4px;
  }

  button {
    width: 100%;
    padding: 0.75rem;
    background: #0066cc;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  button:hover {
    background: #0052a3;
  }
</style> 