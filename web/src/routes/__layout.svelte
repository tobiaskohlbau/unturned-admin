<script lang="ts">
	import SSnackbarContainer from '$lib/components/SSnackbarContainer.svelte';
</script>

<div class="flex content-center min-h-screen font-sans">
	<slot />
</div>
<SSnackbarContainer />

<style>
	@import '$lib/styles/style.css';
</style>
