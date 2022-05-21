# oklookat / frontend
main site + admin panel (elven)


## BUG
sveltekit has bug - layout may be mounted twice

https://github.com/sveltejs/kit/issues/2130

this bug affects [elven layout](./src/routes/elven/__layout@root.svelte) where placed elvenPlayer plugin:
- refresh files page
- play music
- go to another route 
- player will be destroyed

[temp fix](https://github.com/sveltejs/kit/issues/2130#issuecomment-1104692771): 
```bash
npx browserslist@latest --update-db
```