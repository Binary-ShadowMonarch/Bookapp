<!-- login +page.svelte -->
<script lang="ts">
  import LoginForm from '$lib/Login.svelte';
  import NavBar from '$lib/NavBar.svelte';

  function login(username: string, password: string) {
    // call your endpoint, handle auth, etc.
  }
  function signup() {
    // navigate to sign-up page
  }
  function forgot() {
    // show reset-password flow
  }
</script>
<div class=" relative min-h-screen flex items-center justify-center">
    <LoginForm onLogin ={login} onSignUp={signup} onForgotPassword={forgot} />
</div>
