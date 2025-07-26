<!-- src/routes/signup/success/+page.svelte -->
<script lang="ts">
  import { Book, CheckCircle } from 'lucide-svelte';
</script>

<div class="flex container min-h-screen items-center justify-center">  <div class="card">
    <div class="header">
      <a href="/" class="logo" aria-label="Homepage">
        <div class="logo-icon">
          <Book class="h-5 w-5 text-white" />
        </div>
        <h1>Books</h1>
      </a>
    </div>

    <div class="success-content">
      <div class="success-icon">
        <CheckCircle class="h-16 w-16 text-green-500" />
      </div>
      
      <h2>Account Created Successfully!</h2>
      
      <p>
        Your account has been registered and your personal library space has been created. 
        You can now sign in to start building your digital book collection.
      </p>

      <div class="actions">
        <a href="/login" class="login-button">
          Sign In to Your Account
        </a>
        
        <a href="/" class="home-link">
          Back to Home
        </a>
      </div>
    </div>
  </div>
</div>

<style>
  .container {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 1rem;
  }

  .card {
    background: rgba(0, 0, 0, 0.2);    
    border-radius: 16px;
    padding: 2rem;
    width: 100%;
    max-width: 400px;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
    text-align: center;
  }

  .header {
    margin-bottom: 2rem;
  }

  .logo {
    display: inline-flex;
    gap: 0.5rem;
    text-decoration: none;
    color: #1f2937;
    font-size: 1.125rem;
    font-weight: 700;
  }

  .logo-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 2rem;
    height: 2rem;
    border-radius: 0.5rem;
    background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
  }

  .success-content {
    text-align: center;
  }

  .success-icon {
    margin-bottom: 1.5rem;
    display: flex;
    justify-content: center;
  }

  h2 {
    color: #1f2937;
    font-size: 1.5rem;
    font-weight: 700;
    margin-bottom: 1rem;
  }

  p {
    color: #6b7280;
    font-size: 1rem;
    line-height: 1.6;
    margin-bottom: 2rem;
  }

  .actions {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .login-button {
    display: inline-block;
    padding: 0.875rem 1.5rem;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    text-decoration: none;
    border-radius: 12px;
    font-weight: 600;
    transition: all 0.2s;
  }

  .login-button:hover {
    transform: translateY(-2px);
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
  }

  .home-link {
    color: #6b7280;
    text-decoration: none;
    font-size: 0.875rem;
    padding: 0.5rem;
    transition: color 0.2s;
  }

  .home-link:hover {
    color: #667eea;
  }

  @media (prefers-color-scheme: dark) {
    .card {
      background: #1f2937;
    }

    .logo {
      color: #f9fafb;
    }

    h2 {
      color: #f9fafb;
    }

    p {
      color: #d1d5db;
    }

    .home-link {
      color: #9ca3af;
    }

    .home-link:hover {
      color: #667eea;
    }
  }
</style>