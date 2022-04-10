<script lang="ts">
import ElvenLink from "$lib_elven/components/elven_link.svelte";
</script>

<svelte:head>
  <title>elven</title>
</svelte:head>

<div class="elven-hello base__container">
    <ElvenLink path="/settings">settings</ElvenLink>
</div>

<style>
    /* .settings__link {
        background-color: var(--color-level-1);
        border-radius: var(--border-radius);
        height: 64px;
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        gap: 14px;
    } */
</style>
