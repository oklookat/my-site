<script lang="ts">
  import { link } from "svelte-spa-router";
  import active from "svelte-spa-router/active";

  const mainRoute = {
    path: /(^\/$)|(^\/settings$)/,
    className: "route-active",
  };
  const articlesRoute = { path: /^\/articles.*$/, className: "route-active" };
  const filesRoute = { path: /\/files.*/gm, className: "route-active" };
</script>

<nav class="header">
  <div class="header__items">
    <a href="/" class="header__item" use:link use:active={mainRoute}>elven</a>
    <a href="/articles" class="header__item" use:link use:active={articlesRoute}
      >articles</a
    >
    <a href="/files" class="header__item" use:link use:active={filesRoute}
      >files</a
    >
  </div>
</nav>

<style lang="scss">
  :global(a.route-active) {
    text-decoration: underline 1px;
  }

  .header {
    font-weight: bold;
    color: var(--color-header-text);
    border-bottom: var(--color-border) 1px solid;
    height: var(--header-height);
    width: 100%;
    display: flex;
    align-items: center;
    &__items {
      width: 100%;
      display: flex;
      align-items: center;
      gap: 42px;
      margin-left: 2vw;
      margin-right: 2vw;
      @media screen and (max-width: 1023px) {
        justify-content: center;
        margin-left: 0;
        margin-right: 0;
      }
    }
    &__item {
      height: 80%;
      width: max-content;
      cursor: pointer;
      display: flex;
      justify-content: center;
      align-items: center;
    }
  }
</style>
