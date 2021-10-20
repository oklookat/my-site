<script lang="ts">
  import { link } from "svelte-spa-router";
  import active from "svelte-spa-router/active";
</script>

<div class="header__container">
  <div class="header__universal">
    <nav class="header__navigation">
      <div class="header__item">
        <a
          href="/"
          use:link
          use:active={{ path: "/", className: "header__item-active" }}>elven</a
        >
      </div>
      <div class="header__item">
        <a
          href="/articles"
          use:link
          use:active={{ path: "/articles/*", className: "header__item-active" }}
          >articles</a
        >
      </div>
      <div class="header__item">
        <a
          href="/files"
          use:link
          use:active={{ path: "/files/*", className: "header__item-active" }}
          >files</a
        >
      </div>
    </nav>
  </div>
</div>

<style>
  .header__container {
    /* background-color: var(--color-header); */
    font-weight: bold;
    color: var(--color-header-text);
    height: var(--header-height);
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .header__universal {
    height: 100%;
    width: 100%;
    max-width: 975px;
  }

  .header__item {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 80%;
    width: 84px;
    cursor: pointer;
  }

  .header__item-active {
    border-radius: 4px;
    background-color: var(--color-header-active);
  }

  .header__navigation {
    height: 100%;
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: transparent;
  }
</style>
