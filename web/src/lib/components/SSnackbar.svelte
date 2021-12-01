<script>
	import { createEventDispatcher } from 'svelte';
	import { fade } from 'svelte/transition';
	// import CloseIcon from './CloseIcon.svelte';
	// import ErrorIcon from './ErrorIcon.svelte';
	// import InfoIcon from './InfoIcon.svelte';
	// import SuccessIcon from './SuccessIcon.svelte';

	const dispatch = createEventDispatcher();

	export let type = 'error';
	export let dismissible = true;
</script>

<article transition:fade class={type}>
	{#if type === 'success'}
		<i class="material-icons">check</i>
	{:else if type === 'error'}
		<i class="material-icons">error</i>
	{:else}
		<i class="material-icons">info</i>
	{/if}

	<div class="text">
		<slot />
	</div>

	{#if dismissible}
		<button class="close" on:click={() => dispatch('dismiss')}>
			<i class="material-icons">close</i>
		</button>
	{/if}
</article>

<style lang="postcss">
	article {
		color: white;
		padding: 0.75rem 1.5rem;
		border-radius: 0.2rem;
		display: flex;
		align-items: center;
		margin: 0 auto 0.5rem auto;
		width: 20rem;
	}

	.text {
		margin-left: 1rem;
	}

	button {
		color: white;
		background: transparent;
		border: 0 none;
		padding: 0;
		margin: 0 0 0 auto;
		line-height: 1;
		font-size: 1rem;
		cursor: pointer;
	}

	.error {
		background: IndianRed;
	}

	.success {
		background: MediumSeaGreen;
	}

	.info {
		background: SkyBlue;
	}
</style>
