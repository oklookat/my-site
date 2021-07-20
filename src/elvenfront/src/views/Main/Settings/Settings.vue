<template>
  <div class="container">
    <Header></Header>
    <div class="content">
     settings
    </div>
  </div>
</template>

<script>
import {defineComponent} from "vue";
import Header from "@/components/Header/Header";

export default defineComponent({
  name: 'Settings',
  components: {Header}
})
</script>

<style scoped>

</style>