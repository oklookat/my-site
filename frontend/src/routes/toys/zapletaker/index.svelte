<script lang="ts" context="module">
	import type { LoadEvent } from "@sveltejs/kit";


	export const load = async (e: LoadEvent) => {
		const stuff = e.stuff
		stuff.title = "заплетакер"
		return {
			status: 200,
			stuff: stuff
		}
	}
</script>

<script lang="ts">
	class ThePartner {
		public exists = false;
		public index = 0;
		// укроп
		public full = '';
		// укр
		public first = '';
		// оп
		public second = '';

		public handle(index: number, word: string) {
			this.exists = true;
			this.index = index;

			// петр
			const first = word.slice(0, cutLen);
			this.first = first;

			// ушка
			const second = word.slice(cutLen);
			this.second = second;

			finalWords.push(word);
		}
	}

	const partners: { w1?: ThePartner; w2?: ThePartner } = {
		w1: new ThePartner(),
		w2: new ThePartner()
	};

	let finalVal = '';
	let finalWords: string[] = [];
	function finalWordsMakeCursed() {
		if (partners.w1?.exists && partners.w2?.exists) {
			// swap
			finalWords[partners.w1.index] = partners.w2.first + partners.w1.second;
			finalWords[partners.w2.index] = partners.w1.first + partners.w2.second;

			// reset
			delete partners.w1;
			delete partners.w2;
		}
	}

	let inputVal = '';

	let maxCutLen = 1;

	let cutLen = 1;

	$: onCutLenChanged(cutLen);
	function onCutLenChanged(newVal: number) {
		onInputValChanged(inputVal);
	}

	$: onInputValChanged(inputVal);
	function onInputValChanged(newVal: string) {
		if (!newVal) {
			maxCutLen = 1;
			finalVal = '';
			return;
		}

		// remove duplicate whitespace
		newVal = newVal.replace(/\s+/g, ' ').trim();
		const splitted = newVal.split(' ');
		if (splitted.length < 2) {
			return;
		}

		const isPartnersExists = 'w1' in partners && 'w2' in partners;
		if (!isPartnersExists) {
			partners.w1 = new ThePartner();
			partners.w2 = new ThePartner();
		}

		for (let i = 0; i < 2; i++) {
			const word = splitted[i];

			// skip small words
			if (word.length < 3) {
				finalWords.push(word);
				continue;
			}

			const wordCopy = word;

			if (partners.w1 && !partners.w1.exists) {
				partners.w1.full = wordCopy;
				maxCutLen = wordCopy.length;
				if (wordCopy.length > 6) {
					maxCutLen--;
				}

				partners.w1.handle(i, wordCopy);

				// handle next word for partner 2
				continue;
			}

			if (partners.w1 && partners.w2 && !partners.w2.exists) {
				partners.w2.full = wordCopy;
				partners.w2.handle(i, wordCopy);
			}

			finalWordsMakeCursed();
		}

		finalVal = finalWords.join(' ');
		finalWords = [];
	}
</script>

<div class="zapletaker">
	<div class="desc">
		<h1>Заплетакер</h1>
		<div>введи два слова — получишь заплетак</div>
	</div>
	<div class="main">
		<div class="text">
			<b>Два слова:</b>
			<input type="text" bind:value={inputVal} />
		</div>
		<div class="madness">
			<b>Упоротость:</b>
			<input type="range" min="1" max={maxCutLen} bind:value={cutLen} />
		</div>
	</div>
	<div class="result">
		<h1>{finalVal ? finalVal : "Тут будет результат"}</h1>
	</div>
</div>

<style lang="scss">
	.zapletaker {
		width: 100%;
		display: flex;
		flex-direction: column;
		gap: 14px;
		align-items: center;
		justify-content: center;
		gap: 64px;
		.desc {
			display: flex;
			flex-direction: column;
			align-items: center;
		}
		.main {
			height: 100%;
			width: 100%;
			max-width: 324px;
			margin: auto;
			display: flex;
			flex-direction: column;
			gap: 24px;
			b {
				font-size: 1.4rem;
			}
			input {
				height: 44px;
				width: 100%;
				font-size: inherit;
			}
			.madness,
			.text {
				display: flex;
				flex-direction: column;
				justify-content: center;
				align-items: center;
				gap: 12px;
			}
			.result {
				align-self: center;
			}
		}
	}
</style>
