<template>
  <div class="container">
    <Header title="файлы"></Header>
    <div class="content">
      files
    </div>
  </div>
</template>

<script>
import {defineComponent} from "vue";
import Header from "@/components/Header/Header";

export default defineComponent({
  name: 'Files',
  components: {Header}
})
</script>

<style scoped>

</style>