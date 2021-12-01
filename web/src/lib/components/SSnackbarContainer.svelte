<script lang="ts">
	import { dismissToast, toasts } from '$lib/store';
	import SSnackbar from './SSnackbar.svelte';
</script>

{#if $toasts.length > 0}
	<div class="notification">
		{#each $toasts as toast (toast.id)}
			<SSnackbar
				type={toast.type}
				dismissible={toast.dismissible}
				on:dismiss={() => dismissToast(toast.id)}>{toast.message}</SSnackbar
			>
		{/each}
	</div>
{/if}

<style lang="postcss">
	.notification {
		position: fixed;
		width: 400px;
		left: 50%;
		margin-left: -200px;
		bottom: 24px;
		z-index: 255;
	}
</style>
