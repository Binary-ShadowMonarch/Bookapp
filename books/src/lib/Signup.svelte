<!-- Signup.svelte -->
<script lang="ts">
    import { Book } from 'lucide-svelte';


</script>

<div class="container">
  <div class="card">
      <div class="header text-lg justify-center flex font-bold text-gray-900 dark:text-gray-100 select-none">
      <a href="/" class="flex items-center gap-2" aria-label="Homepage">
        <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-blue-400 to-purple-500">
            <Book class="h-5 w-5 text-white" />
          </div>
          <h1>Books</h1>
      </a>
    </div>

    <div class="social-buttons gap-10">
      <!-- Google -->
      <a aria-label="GoogleSignup" href="/" class="social-button">
            <svg
        version="1.1"
        width="20"
        id="Layer_1"
        xmlns="http://www.w3.org/2000/svg"
        xmlns:xlink="http://www.w3.org/1999/xlink"
        x="0px"
        y="0px"
        viewBox="0 0 512 512"
        style="enable-background:new 0 0 512 512;"
        xml:space="preserve"
      >
        <path
          style="fill:#FBBB00;"
          d="M113.47,309.408L95.648,375.94l-65.139,1.378C11.042,341.211,0,299.9,0,256
	c0-42.451,10.324-82.483,28.624-117.732h0.014l57.992,10.632l25.404,57.644c-5.317,15.501-8.215,32.141-8.215,49.456
	C103.821,274.792,107.225,292.797,113.47,309.408z"
        ></path>
        <path
          style="fill:#518EF8;"
          d="M507.527,208.176C510.467,223.662,512,239.655,512,256c0,18.328-1.927,36.206-5.598,53.451
	c-12.462,58.683-45.025,109.925-90.134,146.187l-0.014-0.014l-73.044-3.727l-10.338-64.535
	c29.932-17.554,53.324-45.025,65.646-77.911h-136.89V208.176h138.887L507.527,208.176L507.527,208.176z"
        ></path>
        <path
          style="fill:#28B446;"
          d="M416.253,455.624l0.014,0.014C372.396,490.901,316.666,512,256,512
	c-97.491,0-182.252-54.491-225.491-134.681l82.961-67.91c21.619,57.698,77.278,98.771,142.53,98.771
	c28.047,0,54.323-7.582,76.87-20.818L416.253,455.624z"
        ></path>
        <path
          style="fill:#F14336;"
          d="M419.404,58.936l-82.933,67.896c-23.335-14.586-50.919-23.012-80.471-23.012
	c-66.729,0-123.429,42.957-143.965,102.724l-83.397-68.276h-0.014C71.23,56.123,157.06,0,256,0
	C318.115,0,375.068,22.126,419.404,58.936z"
        ></path>
      </svg>
      </a>
      <a aria-label="Github" class="social-button" href="/">
      <svg height="20" width="20" viewBox="0 0 16 16" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
        <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2 .37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0 0 16 8c0-4.42-3.58-8-8-8z"/>
      </svg>
      </a>
      <div class="or-divider"><span>or</span></div>
    </div>

    <div class="form-wrapper">
      <form>
        <div class="form-group">
          <input type="text" placeholder="Name" aria-label="Name" />
        </div>
        <div class="form-group">
          <input type="email" placeholder="Email" aria-label="Email" />
        </div>
        <div class="form-group">
          <input type="password" placeholder="Password" aria-label="Password" />
        </div>
        <div class="form-group checkbox-group">
          <input type="checkbox" id="terms" />
          <label for="terms">
            I agree to the <a href="/">Terms and Conditions</a>
            <svg viewBox="0 0 20 20" width="16" height="16">
              <path d="M6.293 9.293a1 1 0 0 1 1.414 0L10 10.586l2.293-2.293a1 1 0 1 1 
                       1.414 1.414l-3 3a1 1 0 0 1-1.414 
                       0l-3-3a1 1 0 0 1 0-1.414z" fill="currentColor"/>
            </svg>
          </label>
        </div>
        <div class="form-group">
          <button type="button">Sign up</button>
        </div>
        <p class="signin‑link">
          Already have an account? 
          <a href="/sign-in">Sign in</a>
        </p>
      </form>
    </div>
  </div>
</div>

<style>
  :global(body) {
    margin: 0;
    font-family: system-ui, sans-serif;
    color: white;
    background: #1a001f;
  }

  .container {
    display: flex;
    justify-content: center;
    padding: 2rem;
  }

  .card {
    width: 100%;
    max-width: 400px;
    background: rgba(0, 0, 0, 0.2); /* dark purple acrylic */
    backdrop-filter: blur(50px);
    border: 1px solid rgba(255,255,255,0.15);
    border-radius: 1rem;
    box-shadow: 0 8px 32px rgba(0,0,0,0.3);
    overflow: hidden;
  }

  .header {
    padding: 1.5rem;
    text-align: center;
    border-bottom: 1px solid rgba(255,255,255,0.2);
  }


  .social-buttons {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 1rem;
    position: relative;
  }

  .social-button {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 3rem;
    height: 3rem;
    border: 1px solid rgba(255,255,255,0.4);
    border-radius: 0.5rem;
    transition: transform 0.2s ease;
  }

  .social-button:hover {
    transform: scale(1.05);
    background: rgba(255,255,255,0.1);
  }

  .or-divider {
    position: absolute;
    top: 100%;
    left: 50%;
    transform: translate(-50%, -50%);
    /* background: rgba(75, 0, 130, 0.8); */
    padding: 0 0.75rem;
    font-size: 0.875rem;
    color: rgba(255,255,255,0.75);
  }

  .form-wrapper {
    padding: 1.5rem;
  }

  .form-group {
    margin-bottom: 1rem;
  }

  .form-group input[type="text"],
  .form-group input[type="email"],
  .form-group input[type="password"] {
    width: 100%;
    padding: 0.75rem;
    font-size: 0.9rem;
    border: 1px solid rgba(255,255,255,0.3);
    border-radius: 0.5rem;
    background: transparent;
    color: white;
    transition: border-color 0.2s;
  }

  .form-group input:focus {
    outline: none;
    border-color: #d0b3ff;
  }

  .checkbox-group {
    display: flex;
    align-items: center;
  }

  .checkbox-group input[type="checkbox"] {
    margin: 0;
    margin-right: 0.5rem;
    accent-color: #d0b3ff;
  }

  .checkbox-group label {
    display: inline-flex;
    align-items: center;
    font-size: 0.875rem;
    color: white;
  }

  .checkbox-group label a {
    color: #d0b3ff;
    text-decoration: underline;
  }

  .checkbox-group svg {
    margin-left: 0.25rem;
  }

  button {
    width: 100%;
    padding: 0.75rem;
    font-size: 0.9rem;
    font-weight: bold;
    border: none;
    border-radius: 0.5rem;
    background: linear-gradient(145deg, #181822, #600c85);
    color: white;
    cursor: pointer;
    transition: transform 0.2s, box-shadow 0.2s;
  }

  button:hover {
    transform: scale(1.02);
    box-shadow: 0 4px 15px rgba(0,0,0,0.4);
  }

  .signin‑link {
    margin-top: 1rem;
    font-size: 0.875rem;
    text-align: center;
  }

  .signin‑link a {
    color: #d0b3ff;
    text-decoration: underline;
  }
</style>
