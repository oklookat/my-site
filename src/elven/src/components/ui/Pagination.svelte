<script lang="ts">
  import { createEventDispatcher } from "svelte";

  export let totalPages: number;
  export let currentPage: number;

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

  $: watchTotalPages(totalPages);
  function watchTotalPages(value: number) {
    active = value > 1;
    totalPagesData = value;
  }

  $: watchCurrentPage(currentPage);
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
    <div class="pagination__paginator">
      <div class="pagination__prev-page">
        {#if currentPageData !== 1}
          <div class="pagination__prev-page-butt" on:click={onPrevButton}>
            prev
          </div>
        {/if}
      </div>
      <div class="pagination__pages-input">
        <input
          class="pagination__pages-input-num"
          type="text"
          placeholder="page number"
          bind:value={inputPage}
          on:input={onPageInput}
        />
      </div>
      <div class="pagination__next-page">
        {#if currentPageData < totalPagesData}
          <div class="pagination__next-page-butt" on:click={onNextButton}>
            next
          </div>
        {/if}
      </div>
    </div>
    <div class="pagination__total">pages: {totalPagesData}</div>
  </div>
{/if}

<style lang="scss">
  .pagination {
    border-radius: 8px;
    background-color: var(--color-level-1);
    height: 82px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
    &__paginator {
      height: 36px;
      width: 100%;
      display: flex;
      flex-direction: row;
    }
    &__next-page,
    &__prev-page {
      width: 25%;
    }
    &__next-page-butt {
      border-top-right-radius: 8px;
    }
    &__prev-page-butt {
      border-top-left-radius: 8px;
    }
    &__next-page-butt,
    &__prev-page-butt {
      cursor: pointer;
      width: 100%;
      height: 100%;
    }
    &__next-page-butt:hover,
    &__prev-page-butt:hover {
      background-color: var(--color-hover);
    }
    &__next-page-butt,
    &__prev-page-butt,
    &__pages-input {
      display: flex;
      align-items: center;
      justify-content: center;
      height: 100%;
    }
    &__pages-input {
      width: 50%;
      background-color: var(--color-level-2);
      > input {
        border: none;
        background-color: var(--color-hover);
        width: 100%;
        height: inherit;
        text-align: center;
        font-size: 1.2rem;
        box-sizing: border-box;
      }
    }
  }
</style>
