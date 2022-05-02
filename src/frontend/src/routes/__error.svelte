<script context="module" lang="ts">
	import type { Load } from '@sveltejs/kit';

	export const load: Load = async (event) => {
		return {
			props: {
				statusCode: event.status
			}
		};
	};
</script>

<script lang="ts">
	export let statusCode: number | null;
</script>

<div class="error">
	{#if statusCode}
		<div class="status">
			{statusCode}
		</div>
		<div class="message">
			{#if statusCode === 404}
				Страница не найдена.
			{:else if statusCode === 500}
				Ошибка сервера. Попробуйте позже.
			{:else if statusCode > 399 && statusCode < 500}
				Странный запрос. Попробуйте вернуться на главную страницу.
			{:else}
				Неизвестная ошибка. Попробуйте вернуться на главную страницу.
			{/if}
		</div>
	{:else}
		<div class="message">Неизвестная ошибка. Попробуйте позже.</div>
	{/if}
	<a class="gohome" href="/"><b>на главную</b></a>
</div>

<style lang="scss">
	.error {
		display: flex;
		flex-direction: column;
		align-items: center;
		height: 100%;
		.status {
			height: 94px;
			font-size: 4rem;
			display: flex;
			align-items: center;
		}
		.gohome {
			border-bottom: 2px solid red;
		}
	}
</style>
