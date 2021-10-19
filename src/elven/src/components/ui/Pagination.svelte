<script lang="ts">
  import { createEventDispatcher } from "svelte";

  export let totalPages: number;
  export let currentPage: number;

  const dispatch = createEventDispatcher();

  function dispatchChanged(page: number) {
    dispatch("changed", page);
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
    currentPageData = value;
    inputPage = value.toString();
  }

  function onPageInput() {
    if (pageInputTimeoutID) {
      clearTimeout(pageInputTimeoutID);
    }
    pageInputTimeoutID = setTimeout(() => {
      let bad = isNaN(inputPage);
      if (bad) {
        return;
      }
      const inputPageInt = parseInt(inputPage, 10);
      bad =
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

  function onNextButton() {
    dispatchChanged(currentPageData + 1);
  }

  function onPrevButton() {
    dispatchChanged(currentPageData - 1);
  }
</script>

<main>
  {#if active}
    <div class="pagination__container">
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
</main>

<style scoped>
  .pagination__container {
    border-radius: 8px;
    background-color: var(--color-level-1);
    height: 82px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
  }

  .pagination__paginator {
    height: 36px;
    width: 100%;
    display: flex;
    flex-direction: row;
  }

  .pagination__next-page,
  .pagination__prev-page {
    width: 25%;
  }

  .pagination__next-page-butt {
    border-top-right-radius: 8px;
  }

  .pagination__prev-page-butt {
    border-top-left-radius: 8px;
  }

  .pagination__next-page-butt,
  .pagination__prev-page-butt {
    cursor: pointer;
    width: 100%;
    height: 100%;
  }

  .pagination__next-page-butt:hover,
  .pagination__prev-page-butt:hover {
    background-color: var(--color-hover);
  }

  .pagination__next-page-butt,
  .pagination__prev-page-butt,
  .pagination__pages-input {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
  }

  .pagination__pages-input {
    width: 50%;
    background-color: var(--color-level-2);
  }

  .pagination__pages-input > input {
    border: none;
    background-color: var(--color-hover);
    width: 100%;
    height: inherit;
    text-align: center;
    font-size: 1.2rem;
    box-sizing: border-box;
  }
</style>
