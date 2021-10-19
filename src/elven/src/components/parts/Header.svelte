<main>
  <div class="header__container">
    <div class="header__universal">
      <nav class="header__navigation">
        <div class="header__item">elven</div>
        <div class="header__item">articles</div>
        <div class="header__item">files</div>
      </nav>
    </div>
  </div>
</main>

<script lang="ts">

</script>

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

.header__universal {
  height: 100%;
  width: 100%;
  max-width: 975px;
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
  height: 100%;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: transparent;
}
</style>