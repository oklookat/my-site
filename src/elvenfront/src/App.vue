<template>
  <div class="container">
    <elven-progress></elven-progress>
    <div class="content">
      <div style="height: 16px;"></div>
      <router-view/>
      <div style="height: 64px;"></div>
    </div>
    <elven-player></elven-player>
    <Header v-if="$route.name !== 'Login' && $route.name !== 'Logout'"></Header>
  </div>
</template>

<style lang="scss">
:root {
  --header-height: 52px;
  --header-anim-height: 95%;
  --header-anim-duration: 0.4s;
  --border-radius: 6px;
}

@media (prefers-color-scheme: dark) {
  :root {
    /* MAIN UI START */
    --color-scrollbar: #707070;
    --color-scrollbar-hover: #909090;
    --color-header-active: #050505;
    --color-text: white;
    --color-text-inactive: #E6E6E6;
    --color-text-selection: #803E34;
    --color-placeholder: #666666;
    --color-hover: rgba(255, 255, 255, 0.06);
    --color-border: #202020;
    --color-body: #000000;
    --color-level-1: #151515;
    --color-level-2: #252525;
    --color-level-3: #303030;
    /* MAIN UI END */
    --color-helper: #8CE6E6;
  }
}

@media (prefers-color-scheme: light) {
  :root {
    /* MAIN UI START */
    --color-scrollbar: #505050;
    --color-scrollbar-hover: #707070;
    --color-header-active: #FFFFFF;
    --color-text: black;
    --color-text-inactive: #404040;
    --color-text-selection: #FF7C69;
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
}

html {
  height: 100%;
}

html,
body {
  margin: 0;
  font-family: 'Helvetica', 'Arial', sans-serif;
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

}

.active {
  background-color: var(--color-hover);
  border: none;
}

.the-underline {
  height: 1px;
  border-bottom: 1px solid var(--color-text);
}

.container {
  word-break: break-word;
  hyphens: auto;
  height: 100vh;
  display: grid;
  grid-template-columns: 1fr;
  grid-template-rows: 1fr;
  grid-auto-rows: auto;
}

.content {
  overflow: auto;
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
  opacity: 10%;
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

.router-link-active,
.router-link-exact-active {
  background-color: var(--color-hover);
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


//// EDITOR.JS STYLING START ////
// WIDTH AND HEIGHT START //
.ce-block__content,
.ce-toolbar__content {
  max-height: unset;
  max-width: unset;
}

// WIDTH AND HEIGHT END //
// EDITOR.JS LIGHT AND DARK THEME STYLING START //
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
.ce-toolbar--opened,
{
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

img,
video {
  min-height: 240px;
  min-width: 240px;
  object-fit: fill;
  border-radius: var(--border-radius);
}

// EDITOR.JS LIGHT AND DARK THEME STYLING END //
//// EDITOR.JS STYLING END ////

@media screen and (min-width: 1px) {
  img,
  video {
    width: 100%;
    height: 100%;
  }
}

@media screen and (min-width: 604px) {
  img,
  video {
    width: 80%;
    height: 80%;
  }
}

@media screen and (min-width: 512px) {
  .content {
    width: 500px;
  }
}

@media screen and (min-width: 1000px) {
  .content {
    width: 612px;
  }
}

@media screen and (min-width: 1366px) {
  img,
  video {
    width: 40%;
    height: 40%;
  }
  .content {
    width: 812px;
  }
}
</style>

<script>
import Header from '@/components/Header/Header'

export default {
  name: 'App',
  components: {Header},
  data() {
    return {
      isPageLoaded: false,
      audioPlayerSources: undefined,
    }
  }
}
</script>

