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
    <div class="pagination__paginator">
      <div class="pagination__paginator--prev">
        {#if currentPageData !== 1}
          <div
            class="pagination__paginator--prev--butt pointer center"
            on:click={onPrevButton}
          ></div>
        {/if}
      </div>
      <div class="pagination__paginator--page">
        <input
          class="pagination__input"
          type="text"
          placeholder="page"
          bind:value={inputPage}
          on:input={onPageInput}
        />
      </div>
      <div class="pagination__paginator--next">
        {#if currentPageData < totalPagesData}
          <div class="pagination__paginator--next--butt pointer center" on:click={onNextButton}></div>
        {/if}
      </div>
    </div>
    <div class="pagination__total">
      <div class="pagination__total--count center">{totalPagesData}</div>
    </div>
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

  $border-radius: var(--border-radius);
  .pagination {
    width: 100%;
    height: max-content;
    display: flex;
    flex-direction: column;
    gap: 12px;
    &__paginator,
    &__total {
      width: 100%;
      background: var(--color-level-1);
      border-radius: $border-radius;
      border: var(--color-border) 1px solid;
    }
    &__paginator {
      height: 64px;
      display: grid;
      grid-template-rows: 1fr;
      grid-template-columns: 1fr 1fr 1fr;
      &--prev--butt,
      &--next--butt {
        background: var(--color-level-2);
        width: 100%;
        height: 100%;
      }
      &--prev {
        &--butt {
          border-radius: $border-radius 0 0 $border-radius;
        }
      }
      &--page {
        height: 100%;
        width: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        input {
          width: 75%;
          height: 75%;
          background: var(--color-level-2);
          border: none;
          font-size: 1.4rem;
          text-align: center;
          text-indent: 0;
        }
      }
      &--next {
        &--butt {
          border-radius: 0 $border-radius $border-radius 0;
        }
      }
    }
    &__total {
      min-height: 54px;
      font-size: 1.3rem;
      padding: 12px;
    }
  }
</style>
