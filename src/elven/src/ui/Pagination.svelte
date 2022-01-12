<script lang="ts">
  import { createEventDispatcher } from "svelte";

  /** total pages */
  export let total: number;
  /** current page */
  export let current: number;

  const dispatch = createEventDispatcher<{ changed: number }>();

  function dispatchChanged(page: number) {
    dispatch("changed", page);
  }

  function onNextButton() {
    dispatchChanged(currentPageData + 1);
  }

  function onPrevButton() {
    dispatchChanged(currentPageData - 1);
  }

  let active = false;
  let totalPagesData = 1;
  let currentPageData = 1;
  let pageInputTimeoutID: ReturnType<typeof setTimeout> | null = null;
  let inputPage = "1";

  $: watchTotalPages(total);
  function watchTotalPages(value: number) {
    active = value > 1;
    totalPagesData = value;
  }

  $: watchCurrentPage(current);
  function watchCurrentPage(value: number) {
    if (!value) {
      return;
    }
    currentPageData = value;
    inputPage = value.toString();
  }

  function onPageInput() {
    if (pageInputTimeoutID) {
      clearTimeout(pageInputTimeoutID);
    }
    pageInputTimeoutID = setTimeout(() => {
      const inputPageInt = parseInt(inputPage, 10);
      if (isNaN(inputPageInt)) {
        return;
      }
      let bad =
        inputPageInt > totalPagesData ||
        inputPageInt < 1 ||
        inputPageInt === currentPageData;
      if (bad) {
        return;
      }
      currentPageData = inputPageInt;
      dispatchChanged(currentPageData);
    }, 1000);
  }
</script>

{#if active}
  <div class="pagination">
    <div class="pagination__prev">
      {#if currentPageData !== 1}
        <div class="pagination__prev_ pointer center" on:click={onPrevButton}>
          previous
        </div>
      {/if}
    </div>
    <div class="pagination__go">
      <input
        class="pagination__input"
        type="text"
        placeholder="page"
        bind:value={inputPage}
        on:input={onPageInput}
      />
    </div>
    <div class="pagination__next">
      {#if currentPageData < totalPagesData}
        <div class="pagination__next_ pointer center" on:click={onNextButton}>next</div>
      {/if}
    </div>
    <div class="pagination__total center">{totalPagesData} pages</div>
  </div>
{/if}

<style lang="scss">
  .pointer {
    cursor: pointer;
  }

  .center {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
    width: 100%;
  }

  .pagination {
    border-radius: var(--border-radius);
    background-color: var(--color-level-1);
    height: 94px;
    display: grid;
    grid-template-columns: 2fr 1fr 2fr;
    grid-template-rows: 54px 1fr;
    align-items: center;
    &__go {
      display: flex;
      justify-content: center;
      align-items: center;
      height: 100%;
      width: 100%;
      input {
        font-size: 1.1rem;
        height: 75%;
        width: 100%;
      }
    }
    &__prev, &__next, &__total {
      width: 100%;
      height: 100%;
    }
    &__prev > div:hover,
    &__next > div:hover {
      border-radius: var(--border-radius);
      background-color: var(--color-hover);
    }
    &__total {
      grid-column: 1 / 4;
    }
  }
</style>
