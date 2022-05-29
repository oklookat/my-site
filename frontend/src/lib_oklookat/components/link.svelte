<script lang="ts">
    import { page } from '$app/stores';
    export let path: string

	let highlightPath = false;
	$: onPathChanged($page.url.pathname)

    function onPathChanged(pathname: string) {
        const isPropPathRoot = !path || path === "/"
        if(isPropPathRoot) {
            highlightPath = pathname === "/"
            return
        }
        highlightPath = !!(path && path !== "/") && !!($page.url.pathname && $page.url.pathname.includes(path));
    }
 
</script>

<a class:active={highlightPath} href={path}><slot/></a>