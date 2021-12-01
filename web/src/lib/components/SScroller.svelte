<script lang="ts">
	import { afterUpdate, beforeUpdate } from 'svelte';

	let container;
	let needsAutoscroll: boolean;

	export let autoscroll = false;

	beforeUpdate(() => {
		needsAutoscroll =
			container && container.offsetHeight + container.scrollTop > container.scrollHeight - 20;
	});

	afterUpdate(() => {
		if (needsAutoscroll && autoscroll) container.scrollTo(0, container.scrollHeight);
	});
</script>

<div class="container" bind:this={container}>
	<slot />
</div>

<style scoped>
	.container {
		height: 100%;
		overflow: auto;
	}
</style>
