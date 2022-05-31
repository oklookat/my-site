<script lang="ts" context="module">
	import type { LoadEvent } from "@sveltejs/kit";


	export const load = async (e: LoadEvent) => {
		const stuff = e.stuff
		stuff.title = "прикольные штуки"
		return {
			status: 200,
			stuff: stuff
		}
	}
</script>

<script>
import ToyCard from "$oklookat/components/toy_card.svelte";

</script>

<div class="cooltainer">
    <div>
        <h2>Потешные вещи (😂)</h2>
    </div>
    <div class="items">
        <ToyCard name="заплетакер" path="/toys/zapletaker">
            Миска каши? Булка хлеба? Пакет сухариков?
        </ToyCard>
    </div>
</div> 

<style lang="scss">
    .cooltainer {
        width: 100%;
        max-width: 624px;
        margin: auto;
        display: flex;
        flex-direction: column;
        gap: 12px;
    }
</style>