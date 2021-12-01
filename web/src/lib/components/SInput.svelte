<script lang="ts">
	export let value = '';
	export let placeholder: string;
	export let type: string;

	function handleInput(event: Event & { currentTarget: EventTarget & HTMLInputElement }): void {
		value = event.currentTarget.value;
	}
</script>

<div class="w-full h-12 relative overflow-auto">
	<input
		{type}
		name="input"
		placeholder=" "
		on:input={handleInput}
		on:keypress
		{value}
		class="w-full h-full text-basics-300 pt-6 border-none outline-none"
	/>
	<label
		for="input"
		class="absolute left-0 bottom-0 h-full w-full pointer-events-none border-b border-solid border-black"
	>
		<span class="content absolute bottom-1 left-0 transition-all duration-300 ease-in-out"
			>{placeholder}</span
		>
	</label>
</div>

<style>
	label::after {
		content: '';
		transform: scaleX(0);
		@apply h-full w-full absolute left-0 bottom-[-1px] duration-300 ease-in-out border-b-2 border-solid border-primary;
	}

	input:focus + label .content,
	input:not(:placeholder-shown) + label .content {
		transform: translateY(-100%);
		@apply text-primary text-sm;
	}

	input:focus + label::after,
	input:not(:placeholder-shown) + label::after {
		transform: scaleX(1);
	}
</style>
