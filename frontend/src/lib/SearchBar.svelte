<!-- this is a search bar component that I use in the library -->
<!-- it has a nice dark theme and smooth animations -->
<script lang="ts">
  // this is the search value that gets updated as the user types
  export let value = '';
</script>

<div class="group">
  <!-- this is a search icon that I drew with SVG -->
  <!-- it's positioned inside the input field on the left -->
  <svg viewBox="0 0 24 24" aria-hidden="true" class="search-icon">
    <g>
      <path
        d="M21.53 20.47l-3.66-3.66C19.195 15.24 20 13.214 20 11c0-4.97-4.03-9-9-9s-9 4.03-9 9 4.03 9 9 9c2.215 0 4.24-.804 5.808-2.13l3.66 3.66c.147.146.34.22.53.22s.385-.073.53-.22c.295-.293.295-.767.002-1.06zM3.5 11c0-4.135 3.365-7.5 7.5-7.5s7.5 3.365 7.5 7.5-3.365 7.5-7.5 7.5-7.5-3.365-7.5-7.5z"
      ></path>
    </g>
  </svg>

  <!-- the actual search input field -->
  <!-- bind:value connects this to the parent component -->
  <input
    id="query"
    class="input"
    type="search"
    placeholder="Search..."
    bind:value
    name="searchbar"
  />
</div>

<style>
  /* container for the search bar with relative positioning */
  .group {
    display: flex;
    line-height: 28px;
    align-items: center;
    position: relative;
    max-width: 190px;
  }

  /* the search input field styling */
  /* I made it dark themed to match my app's design */
  .input {
    font-family: "Montserrat", sans-serif;
    width: 100%;
    height: 45px;
    padding-left: 2.5rem; /* space for the search icon */
    box-shadow: 0 0 0 1.5px #2b2c37, 0 0 25px -17px #000;
    border: 0;
    border-radius: 12px;
    background-color: #16171d; /* dark background */
    outline: none;
    color: #bdbecb; /* light text */
    transition: all 0.25s cubic-bezier(0.19, 1, 0.22, 1); /* smooth animation */
    cursor: text;
    z-index: 0;
  }

  /* placeholder text color */
  .input::placeholder {
    color: #bdbecb;
  }

  /* hover effect - makes the border glow a bit */
  .input:hover {
    box-shadow: 0 0 0 2.5px #2f303d, 0px 0px 25px -15px #000;
  }

  /* when clicked, it scales down slightly for feedback */
  .input:active {
    transform: scale(0.95);
  }

  /* focus state - stronger glow when the input is selected */
  .input:focus {
    box-shadow: 0 0 0 2.5px #2f303d;
  }

  /* the search icon positioning and styling */
  .search-icon {
    position: absolute;
    left: 1rem;
    fill: #bdbecb; /* same color as the text */
    width: 1rem;
    height: 1rem;
    pointer-events: none; /* so clicks pass through to the input */
    z-index: 1; /* above the input */
  }
</style>
