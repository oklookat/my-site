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
    <div class="paginator">
      <div class="prev">
        {#if currentPageData !== 1}
          <div class="prev__butt pointer center" on:click={onPrevButton} />
        {/if}
      </div>
      <div class="page">
        <input
          class="pagination__input"
          type="text"
          placeholder="page"
          bind:value={inputPage}
          on:input={onPageInput}
        />
      </div>
      <div class="next">
        {#if currentPageData < totalPagesData}
          <div class="next__butt pointer center" on:click={onNextButton} />
        {/if}
      </div>
    </div>
    <div class="total">
      <div class="count center">{totalPagesData}</div>
    </div>
  </div>
{/if}

<style lang="scss">
  @import "./src/assets/variables";

  .center {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
    width: 100%;
  }

  .pagination {
    width: 100%;
    max-width: $max-card-width;
    height: max-content;
    display: flex;
    flex-direction: column;
    gap: 12px;
    .paginator,
    .total {
      width: 100%;
      background: var(--color-level-1);
      border-radius: var(--border-radius);
      border: var(--color-border) 1px solid;
    }
    .paginator {
      height: 64px;
      display: grid;
      grid-template-rows: 1fr;
      grid-template-columns: repeat(3, 1fr);
      .prev,
      .next {
        &__butt {
          background: var(--color-level-2);
          width: 100%;
          height: 100%;
        }
      }
      .prev {
        &__butt {
          border-radius: var(--border-radius) 0 0 var(--border-radius);
        }
      }
      .page {
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
      .next {
        &__butt {
          border-radius: 0 var(--border-radius) var(--border-radius) 0;
        }
      }
    }
    .total {
      min-height: 54px;
      font-size: 1.3rem;
      padding: 12px;
    }
  }
</style>
