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

<div class="header__container">
  <nav class="header__navigation">
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
    max-width: 975px;
    height: 100%;
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: transparent;
  }
</style>
