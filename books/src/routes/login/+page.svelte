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

<!-- <NavBar> -->

  <!-- <div class ="flex justify-center " slot = "name">Books</div> -->

  <!-- search -->
  <!-- <div slot="center" class="relative text-lg font-bold text-gray-900 dark:text-gray-100 select-none flex items-center"> -->
        <!-- <h1>Login</h1> -->
  <!-- </div> -->
<!-- </NavBar> -->

<div class=" relative min-h-screen flex items-center justify-center">
    <LoginForm onLogin ={login} onSignUp={signup} onForgotPassword={forgot} />
</div>
