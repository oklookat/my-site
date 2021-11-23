<script lang="ts">
  export let percents: number;
</script>

<div class="progress" style="width: {percents}%" />

<style>
  :root {
    --color: #404040;
    --border-radius: inherit;
  }
  .progress {
    border-radius: var(--border-radius);
    height: 100%;
    width: 0;
    background-color: var(--color);
  }
</style>
