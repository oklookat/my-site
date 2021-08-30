<template>
  <div class="container">
    <div class="service" v-if="$route.name !== 'Login' && $route.name !== 'Logout'">
      <elven-notify></elven-notify>
      <Header></Header>
      <elven-progress></elven-progress>
    </div>
    <div class="content">
      <div style="height: 16px;"></div>
      <router-view/>
      <div style="height: 64px;"></div>
    </div>
    <div class="service-2" v-if="$route.name !== 'Login' && $route.name !== 'Logout'">
      <elven-player></elven-player>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent} from 'vue'
import Header from '@/components/Header/Header.vue'

export default defineComponent({
  name: 'App',
  components: {Header},
  data() {
    return {
    }
  },
  mounted() {
    this.$elvenNotify.error('message')
  }
})
</script>

<style>
:root {
  --header-height: 52px;
  --border-radius: 6px;
}

:root {
  /* MAIN UI START */
  --color-scrollbar: #505050;
  --color-scrollbar-hover: #707070;
  /* HEADER START */
  --color-header: #FFFFFF;
  --color-header-active: #FFFFFF;
  --color-header-text: #000000;
  --color-header-hover: rgba(0, 0, 0, 0.03);
  /* HEADER END */
  --color-text: black;
  --color-text-inactive: #404040;
  --color-text-selection: #BAB0FF;
  --color-placeholder: #808080;
  --color-hover: rgba(0, 0, 0, 0.05);
  --color-border: #F2F2F2;
  --color-body: #E6E6E6;
  --color-level-1: #FFFFFF;
  --color-level-2: #FFFFFF;
  --color-level-3: #FFFFFF;
  /* MAIN UI END */
  --color-helper: #8CE6E6;
}

.router-link-active,
.router-link-exact-active {
  background-color: var(--color-hover);
}

html,
body {
  min-height: 100vh;
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  color: var(--color-text);
  background-color: var(--color-body);
  text-underline-offset: 0.4em;
}

::-webkit-scrollbar {
  width: 4px;
  height: 4px;
}

::-webkit-scrollbar-thumb {
  background: var(--color-scrollbar);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: var(--color-scrollbar-hover);
}


a {
  color: inherit;
  text-decoration: none;
}

#app {
  min-height: 100vh;
}

.active {
  background-color: var(--color-hover);
  border: none;
}

.container {
  min-height: 100vh;
  word-break: break-word;
  hyphens: auto;
  display: grid;
  grid-template-rows: min-content 1fr min-content;
  grid-auto-rows: 1px;
  flex-direction: column;
}

.service, .service-2 {
  position: sticky;
  z-index: 7777;
}

.service {
  top: 0;
  box-shadow: 0px 8px 8px 0px rgba(35, 60, 80, 0.1);
}

.service-2 {
  bottom: 0;
}

.content {
  height: 100%;
  margin: auto;
  width: 95%;
  font-size: 1.1rem;
  line-height: 1.46rem;
  letter-spacing: 0.0007rem;
}

textarea {
  font-family: inherit;
  color: inherit;
  resize: none;
  overflow: hidden;
}

input {
  color: var(--color-text);
  background-color: var(--color-level-1);
  border: 1px solid var(--color-border);
  min-height: 24px;
  border-radius: 4px;
  text-indent: 6px;
}

input::placeholder {
  opacity: 20%;
}

input:focus::placeholder {
  color: transparent;
}


::placeholder {
  color: var(--color-text);
  opacity: 45%;
}

:focus {
  outline: none;
}

.no-scroll {
  overflow: hidden;
}

.logo-text {
  font-size: 1.2rem;
  letter-spacing: 0.8rem;
  text-indent: 0.8rem;
  user-select: none;
}

.item {
  cursor: pointer;
  border-radius: 6px;
  height: 100%;
  display: flex;
  flex-direction: row;
  align-items: center;
}

.cursor-pointer {
  cursor: pointer;
}

.item:hover {
  background-color: var(--color-hover);
}


/* EDITOR.JS STYLING START */
/* MAIN START */
.codex-editor {
  z-index: unset;
}

.ce-block__content,
.ce-toolbar__content {
  max-height: unset;
  max-width: unset;
}

/* MAIN END */
/* EDITOR.JS LIGHT AND DARK THEME STYLING START */
.cdx-button {
  background: var(--color-level-1);
  color: var(--color-text-inactive);
  border: 1px solid var(--color-border);
}

.cdx-button:hover {
  background: var(--color-hover);
  border: 1px solid var(--color-hover);
}

.cdx-settings-button:hover {
  background-color: var(--color-hover);
}

.ce-inline-tool,
.codex-editor svg {
  color: var(--color-text);
}

.ce-inline-tool:hover,
.ce-inline-toolbar__dropdown:hover,
.ce-toolbar__settings-btn:hover,
.ce-settings__button:hover {
  background-color: var(--color-hover);
}

.ce-inline-tool-input,
.ce-inline-toolbar,
.ce-conversion-toolbar,
.ce-settings,
.ce-toolbar--opened {
  background-color: var(--color-level-2);
  border: 1px solid var(--color-border);
}

.ce-conversion-tool__icon,
.ce-toolbar__settings-btn {
  background: var(--color-level-3);
}

.ce-inline-toolbar__actions input::placeholder,
.ce-conversion-toolbar__label {
  color: var(--color-text-inactive)
}

.codex-editor--narrow .ce-toolbox {
  background-color: var(--color-level-1);
}

::selection,
.ce-block--selected .ce-block__content {
  background: var(--color-text-selection);
}

/* EDITOR.JS LIGHT AND DARK THEME STYLING END */
/* EDITOR.JS STYLING END */

img,
video {
  min-height: 240px;
  min-width: 240px;
  object-fit: fill;
  border-radius: var(--border-radius);
}

@media screen and (min-width: 1px) {
  img,
  video {
    max-width: 100%;
    max-height: 320px;
  }
}

@media screen and (min-width: 512px) {
  .content {
    width: 500px;
  }
}

@media screen and (min-width: 604px) {
  img,
  video {
    max-width: 60%;
    max-height: 60%;
  }
}
@media screen and (min-width: 1365px) {
  img,
  video {
    width: 40%;
    height: 40%;
  }
  .content {
    width: 712px;
  }

}
</style>
