<script>
	import Error from './error.svelte';
</script>

<Error imgURL="/ineed/fun/err400.jpg" code="HTTP/4xx" short="странный запрос">
	<div style="font-weight: bold;">Тут может быть два варианта:</div>
	<div>1. Я накосячил и вы попали в просак (извините).</div>
	<div>2. Вы сами знаете как тут оказались.</div>
    <p />
	<div style="font-weight: bold;">Если у вас первый вариант, выберите одно из двух:</div>
    <div>1. Напишите мне о том, как вы тут оказались.</div>
    <div>2. Просто вернитесь на главную.</div>
</Error>
