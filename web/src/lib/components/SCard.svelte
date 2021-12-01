<script lang="ts">
	import { css } from '$lib/css';
	export let width: string;
	export let height: string;
</script>

<div use:css={{ width, height }} class="card">
	<slot />
</div>

<style>
	.card {
		--width: auto;
		--height: auto;
		width: var(--width);
		height: var(--height);

		@apply shadow m-4 p-6 border border-solid border-basics-100 rounded-md bg-white text-black transition-all flex flex-col;
	}
</style>
