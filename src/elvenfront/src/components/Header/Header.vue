<template>
  <div class="header-container">

    <div class="header">
      <nav class="header-navigation">
        <RouterLink class="header-mobile-menu-item" :to="{name: 'Index'}">
          <svg id="home_butt" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 275 310.58"><path d="M373,153.25H126.93c-8,0-14.47,5.76-14.47,12.88V347.39c0,7.12,6.48,12.89,14.47,12.89h80.94V231h84.18V360.28H373c8,0,14.46-5.77,14.46-12.89V166.13C387.46,159,381,153.25,373,153.25Z" transform="translate(-112.46 -49.7)"/><path d="M255.45,51.66,382.13,155.09a8.68,8.68,0,0,1-5.49,15.41H123.29a8.69,8.69,0,0,1-5.5-15.41L244.47,51.66A8.68,8.68,0,0,1,255.45,51.66Z" transform="translate(-112.46 -49.7)"/></svg>
        </RouterLink>
        <RouterLink class="header-mobile-menu-item" :to="{name: 'Articles'}">
          <svg id="Слой_1" data-name="Слой 1" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 315.43 310"><path d="M416.83,113.5H115.49c-3.89,0-7,3.89-7,8.7V414.8c0,4.81,3.15,8.7,7,8.7H416.83c3.89,0,7.05-3.89,7.05-8.7V122.2C423.88,117.39,420.72,113.5,416.83,113.5ZM202,360.89c0,2.52-1.65,4.56-3.7,4.56H135.7c-2,0-3.7-2-3.7-4.56V339.14c0-2.52,1.65-4.56,3.7-4.56h62.6c2,0,3.7,2,3.7,4.56Zm25-63c0,2.52-1.65,4.56-3.7,4.56H135.7c-2,0-3.7-2-3.7-4.56V276.15c0-2.52,1.65-4.56,3.7-4.56h87.6c2,0,3.7,2,3.7,4.56Zm35-61.75c0,2.52-1.65,4.56-3.7,4.56H135.7c-2,0-3.7-2-3.7-4.56V214.4c0-2.52,1.65-4.57,3.7-4.57H258.3c2.05,0,3.7,2,3.7,4.57Zm60-60.52c0,2.52-1.65,4.56-3.7,4.56H135.7c-2,0-3.7-2-3.7-4.56v-34.1c0-2.52,1.65-4.56,3.7-4.56H318.3c2.05,0,3.7,2,3.7,4.56Z" transform="translate(-108.45 -113.5)"/></svg>
        </RouterLink>
        <RouterLink class="header-mobile-menu-item" :to="{name: 'Files'}">
          <svg id="files_butt" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 285.19 310"><path d="M380.37,132.14h-260a12.31,12.31,0,0,0-12.32,12.29V392.22a12.3,12.3,0,0,0,12.32,12.28h260a12.3,12.3,0,0,0,12.32-12.28V144.43A12.31,12.31,0,0,0,380.37,132.14ZM352,337.5a12.2,12.2,0,0,1-12.22,12.18H160.22A12.2,12.2,0,0,1,148,337.5V332a12.2,12.2,0,0,1,12.22-12.18H339.78A12.2,12.2,0,0,1,352,332Zm0-54.83a12.2,12.2,0,0,1-12.22,12.18H160.22A12.2,12.2,0,0,1,148,282.67v-5.54A12.2,12.2,0,0,1,160.22,265H339.78A12.2,12.2,0,0,1,352,277.13Zm0-54.82A12.2,12.2,0,0,1,339.78,240H160.22A12.2,12.2,0,0,1,148,227.85v-5.54a12.2,12.2,0,0,1,12.22-12.18H339.78A12.2,12.2,0,0,1,352,222.31Z" transform="translate(-107.5 -94.5)"/><rect width="142.31" height="85.11" rx="17.59"/></svg>
        </RouterLink>
      </nav>
    </div>


  </div>
</template>

<script>
import {defineComponent} from 'vue'

export default defineComponent({
  name: 'Header',
  data() {
    return {}
  }
})
</script>

<style lang="scss">
.header-container {
  background-color: var(--color-header);
  color: var(--color-header-text);
  height: var(--header-height);
  width: 100%;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
}

.header{
  height: 100%;
  max-width: 975px;
}

.header-navigation {
  height: 100%;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: transparent;
  gap: 48px;
}

.header-navigation svg {
  width: 20px;
  height: 20px;
}

@media screen and (prefers-color-scheme: dark) {
  .header-navigation svg {
    fill: #CCCCCC;
  }
}
@media screen and (prefers-color-scheme: light) {
  .header-navigation svg {
    fill: #444444;
  }
}


.header-mobile-menu-item {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  cursor: pointer;
}
</style>