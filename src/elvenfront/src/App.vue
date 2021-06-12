<template>
  <router-view/>
</template>

<style lang="scss">
:root {
  --header-height: 46px;
  --header-anim-height: 95%;
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
    --color-header-active: #141414;
    --color-header-menu: black;
  }
}

@media (prefers-color-scheme: light) {
  :root {
    --color-text: black;
    --color-text-inactive: #303030;
    --color-body: #F5F5F5;
    --color-level-1: #E6E6E6;
    --color-level-2: #D9D9D9;
    --color-level-3: #B3B3B3;
    --color-hover: rgba(0, 0, 0, 0.06);
    --color-border: #CCCCCC;
    --color-header-active: white;
    --color-header-menu: #E6E6E6;
  }
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
  font-weight: bold;
  font-size: 1.1rem;
  letter-spacing: 0.6rem;
  user-select: none;
}

.content {
  height: var(--header-height);
  display: flex;
  margin-left: auto;
  margin-right: auto;
  width: 95%;
  font-size: 1.1rem;
  line-height: 1.46rem;
  letter-spacing: 0.0007rem;
}

.item {
  cursor: pointer;
  color: var(--color-text-inactive);
}

.item:hover {
  color: var(--color-text);
  background-color: var(--color-hover);
}

@media screen and (min-width: 640px) {
  .content {
    width: 85%;
  }
}

@media screen and (min-width: 960px) {
  .content {
    width: 75%;
  }
}
</style>
