<template>
  <div class="header">

    <div class="header-mobile"
         @click.self="toggleMobileMenu"
         :class="{'mobile-menu-active': isMobileMenuActive}"
    >
      <div class="logo-text" @click.self="toggleMobileMenu">{{ title }}</div>
      <transition
          name="show-mobile"
          enter-active-class="show-mobile enter-active"
          leave-active-class="show-mobile leave-active"
      >
      <div class="header-mobile-menu" v-if="isMobileMenuActive">
        <RouterLink class="header-mobile-menu-item" :to="{name: 'Index'}">Главная</RouterLink>
        <RouterLink class="header-mobile-menu-item" :to="{name: 'Articles'}">Записи</RouterLink>
        <RouterLink class="header-mobile-menu-item" :to="{name: 'Files'}">Файлы</RouterLink>
        <RouterLink class="header-mobile-menu-item" :to="{name: 'Settings'}">Настройки</RouterLink>
        <RouterLink class="header-mobile-menu-item" :to="{name: 'Logout'}">Выход</RouterLink>
      </div>
      </transition>
    </div>

  </div>
</template>

<script>
import {defineComponent} from 'vue'
export default defineComponent({
  name: 'Header',
  data(){
    return{
      isMobileMenuActive: false,
      scrollPos: 0,
    }
  },
  props: ['title'],
  unmounted() {
    document.body.classList.remove('no-scroll')
  },
  methods: {
    toggleMobileMenu(){
      this.isMobileMenuActive = !this.isMobileMenuActive
      this.toggleBodyScroll()
    },

    toggleBodyScroll(){
      if(document.body.classList.contains('no-scroll')){
        document.body.classList.remove('no-scroll')
      } else{
        document.body.classList.add('no-scroll')
      }
    },

  }
})
</script>

<style scoped lang="scss">
.logo-text{
  color: var(--color-text-inactive);
}
.logo-text:hover{
  color: var(--color-text);
}
.header{
  height: var(--header-height);
  width: 100%;
}
.header-mobile{
  height: 100%;
  width: 100%;
  transition: background-color 0.6s;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}
.header-mobile.mobile-menu-active{
  background-color: var(--color-header-active);
}
.header-mobile-menu{
  z-index: 9999;
  width: 100%;
  background-color: var(--color-header-menu);
  position: absolute;
  top: var(--header-height); left: 0; bottom: 0; right: 0;
  display: block;
  font-size: 1.4rem;
  cursor: pointer;
  overflow: auto;
}
.header-mobile-menu-item{
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 110px;
  cursor: pointer;
  color: var(--color-text-inactive);
}
.header-mobile-menu-item:hover{
  background-color: var(--color-hover);
  color: var(--color-text);
}





.header-mobile-menu.show-mobile.enter-active{
  overflow: hidden;
  animation: show-mobile 0.6s;
}
.header-mobile-menu.show-mobile.leave-active {
  overflow: hidden;
  animation: show-mobile 0.6s reverse;
}
@keyframes show-mobile {
  from{
    height: 0;
  }
  to{
    height: 95%;
  }
}


.header-mobile-menu.show-mobile.enter-active > a:nth-child(1){
  animation: item-1  0.6s;
}
.header-mobile-menu.show-mobile.leave-active > a:nth-child(1){
  animation: item-1  0.6s reverse;
}

.header-mobile-menu.show-mobile.enter-active > a:nth-child(2){
  animation: item-2  0.6s;
}
.header-mobile-menu.show-mobile.leave-active > a:nth-child(2){
  animation: item-2  0.6s reverse;
}

.header-mobile-menu.show-mobile.enter-active > a:nth-child(3){
  animation: item-3  0.6s;
}
.header-mobile-menu.show-mobile.leave-active > a:nth-child(3){
  animation: item-3  0.6s reverse;
}

.header-mobile-menu.show-mobile.enter-active > a:nth-child(4){
  animation: item-4  0.6s;
}
.header-mobile-menu.show-mobile.leave-active > a:nth-child(4){
  animation: item-4  0.6s reverse;
}
.header-mobile-menu.show-mobile.enter-active > a:nth-child(5){
  animation: item-4  0.6s;
}
.header-mobile-menu.show-mobile.leave-active > a:nth-child(5){
  animation: item-4  0.6s reverse;
}



@keyframes item-1 {
  0%, 10% {
    opacity: 0;
  }
  100%{
    opacity: 100;
  }
}
@keyframes item-2 {
  0%, 35% {
    opacity: 0;
  }
  100%{
    opacity: 100;
  }
}
@keyframes item-3 {
  0%, 45% {
    opacity: 0;
  }
  100%{
    opacity: 100;
  }
}
@keyframes item-4 {
  0%, 55% {
    opacity: 0;
  }
  100%{
    opacity: 100;
  }
}

@keyframes item-5 {
  0%, 65% {
    opacity: 0;
  }
  100%{
    opacity: 100;
  }
}
</style>