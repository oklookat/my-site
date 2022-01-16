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
</nav>

<style lang="scss">
  .header {
    font-weight: bold;
    background-color: var(--color-level-1);
    color: var(--color-header-text);
    height: var(--header-height);
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 42px;
    &__item {
      display: flex;
      justify-content: center;
      align-items: center;
      height: 80%;
      width: max-content;
      cursor: pointer;
    }
    &__item-active {
      text-decoration: underline 1px;
    }
  }
</style>
