<script lang="ts">
  import { onDestroy } from "svelte";
  import { location, push } from "svelte-spa-router";

  let path = "";

  const s1 = location.subscribe((v) => {
    path = v;
  });

  onDestroy(() => {
    s1();
  });

  function goIndex() {
    push("/");
  }

  function goArticles() {
    push("/articles");
  }

  function goFiles() {
    push("/files");
  }
</script>

<nav class="header">
  <div class="header__items">
    <div
      class="header__item {path === '/' ? 'header__item-active' : ''}"
      on:click={() => goIndex()}
    >
      elven
    </div>
    <div
      class="header__item {path.startsWith('/articles')
        ? 'header__item-active'
        : ''}"
      on:click={() => goArticles()}
    >
      articles
    </div>
    <div
      class="header__item {path.startsWith('/files')
        ? 'header__item-active'
        : ''}"
      on:click={() => goFiles()}
    >
      files
    </div>
  </div>
</nav>

<style lang="scss">
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
    &__item-active {
      text-decoration: underline 1px;
    }
  }
</style>
