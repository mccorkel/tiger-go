<script lang="ts">
  import { auth } from '../stores/auth';

  let username = '';
  let password = '';
  let error = '';

  async function handleLogin() {
    try {
      error = '';
      await auth.login(username, password);
    } catch (e: any) {
      error = e.message;
    }
  }
</script>

<div class="login-container">
  <div class="login-card">
    <h2>Login</h2>
    {#if error}
      <div class="error">{error}</div>
    {/if}
    <form on:submit|preventDefault={handleLogin}>
      <input
        type="text"
        bind:value={username}
        placeholder="Username"
        required
      />
      <input
        type="password"
        bind:value={password}
        placeholder="Password"
        required
      />
      <button type="submit">Login</button>
    </form>
  </div>
</div>

<style>
  .login-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    padding: 1rem;
  }

  .login-card {
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