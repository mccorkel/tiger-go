<script lang="ts">
  import { auth } from '../stores/auth';
  
  function handleLogout() {
    auth.logout();
  }
</script>

<div class="dashboard">
  <header class="top-bar">
    <h1>Dashboard</h1>
    <button on:click={handleLogout}>Logout</button>
  </header>

  <main class="content">
    <div class="card">
      <h2>Welcome to your dashboard</h2>
      <p>You are successfully logged in!</p>
    </div>
  </main>
</div>

<style>
  .dashboard {
    min-height: 100vh;
    background: #f5f5f5;
  }

  .top-bar {
    background: white;
    padding: 1rem 2rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  }

  .content {
    padding: 2rem;
  }

  .card {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    max-width: 800px;
    margin: 0 auto;
  }

  button {
    padding: 0.5rem 1rem;
    background: #0066cc;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  button:hover {
    background: #0052a3;
  }

  h1 {
    margin: 0;
  }
</style> 