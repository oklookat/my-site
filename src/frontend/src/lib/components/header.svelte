<script lang="ts">
import ElvenLink from "$lib/components/elven_link.svelte";

</script>

<nav class="header">
  <div class="header__items">
    <ElvenLink path="">elven</ElvenLink>
    <ElvenLink path="/articles">articles</ElvenLink>
    <ElvenLink path="/files">files</ElvenLink>
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
