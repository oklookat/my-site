<template>
  <elven-progress></elven-progress>
  <router-view/>
</template>

<style lang="scss">
:root {
  --header-height: 64px;
  --header-anim-height: 95%;
  --header-anim-duration: 0.4s;
}

@media (prefers-color-scheme: dark) {
  :root {
    --color-text: white;
    --color-text-inactive: #E6E6E6;
    --color-body: #101010;
    --color-level-1: #151515;
    --color-level-2: #202020;
    --color-level-3: #303030;
    --color-hover: rgba(255, 255, 255, 0.06);
    --color-border: #353535;
    --color-scrollbar: #707070;
    --color-scrollbar-hover: #909090;
    --color-header-active: #000000;
    --color-text-selection: #803E34;
  }
}

@media (prefers-color-scheme: light) {
  :root {
    --color-text: black;
    --color-text-inactive: #303030;
    --color-body: #E0E0E0;
    --color-level-1: #FFFFFF;
    --color-level-2: #FFFFFF;
    --color-level-3: #FFFFFF;
    --color-hover: rgba(255, 255, 255, 0.80);
    --color-border: #CCCCCC;
    --color-scrollbar: #505050;
    --color-scrollbar-hover: #707070;
    --color-header-active: #FFFFFF;
    --color-text-selection: #FF7C69;
  }
}
html{
  height: 100%;
}
html,
body {
  margin: 0;
  font-family: 'Helvetica', 'Arial', sans-serif;
  color: var(--color-text);
  background-color: var(--color-body);
}

::-webkit-scrollbar {
  width: 6px;
  height: 6px;
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

.container {
  word-break: break-word;
  hyphens: auto;
}

input {
  color: var(--color-text);
  background-color: var(--color-level-1);
  border: 1px solid var(--color-border);
  min-height: 24px;
  border-radius: 4px;
  text-indent: 6px;
}

input:focus::placeholder {
  opacity: 15%;
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

.content {
  height: calc(100% - var(--header-height));
  margin: auto;
  padding-bottom: 16px;
  width: 95%;
  font-size: 1.1rem;
  line-height: 1.46rem;
  letter-spacing: 0.0007rem;
}

.item {
  cursor: pointer;
  border-radius: 6px;
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
.codex-editor svg{
  color: var(--color-text);
}

.ce-inline-tool:hover,
.ce-inline-toolbar__dropdown:hover,
.ce-toolbar__settings-btn:hover,
.ce-settings__button:hover{
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
.ce-toolbar__settings-btn{
  background: var(--color-level-3);
}

.ce-inline-toolbar__actions input::placeholder,
.ce-conversion-toolbar__label{
  color: var(--color-text-inactive)
}

::selection,
.ce-block--selected .ce-block__content{
  background: var(--color-text-selection);
}
// EDITOR.JS LIGHT AND DARK THEME STYLING END //
//// EDITOR.JS STYLING END ////

@media screen and (min-width: 512px) {
  .content {
    width: 500px;
  }
}
@media screen and (min-width: 1024px) {
  .content {
    width: 812px;
  }
}
</style>

<script>
export default {
  name: 'App',
}
</script>