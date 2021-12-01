<script lang="ts">
	import { clickOutside } from '$lib/utils';
	import { createEventDispatcher } from 'svelte';
	import { fade } from 'svelte/transition';

	export let items: string[] = [];
	export let placeholder = '';
	export let value: string[] = [];

	let dropdown = false;

	function isSelected(val: string) {
		return value.includes(val);
	}

	function updateSelection(item: string) {
		const index = value.indexOf(item);
		if (index !== -1) {
			value.splice(index, 1);
			value = value;
		} else {
			value = [...value, item];
		}
	}

	let child;
	const dispatch = createEventDispatcher();

	function onClickOutside() {
		dropdown = false;
		dispatch('change');
	}
</script>

<div class="container" on:click={() => (dropdown = true)} bind:this={child}>
	<div class="content">
		{#each value as item}
			<div class="chip">
				<span>{item}</span>
			</div>
		{/each}
	</div>
	{#if dropdown}
		<div transition:fade class="dropdown-content" use:clickOutside={onClickOutside}>
			{#each items as item}
				<div class="dropdown-item" on:click={() => updateSelection(item)}>
					<input type="checkbox" checked={isSelected(item)} />
					<span>{item}</span>
				</div>
			{/each}
		</div>
	{/if}
	<label for="input" class="label" class:active={value.length > 0 || dropdown}>
		<span class="content">{placeholder}</span>
	</label>
</div>

<style>
	.container {
		width: 100%;
		height: 42px;
		position: relative;
		cursor: pointer;
	}

	.container label {
		position: absolute;
		bottom: 0px;
		left: 0;
		width: 100%;
		height: 100%;
		pointer-events: none;
		border-bottom: 1px solid black;
	}

	.container label::after {
		content: '';
		height: 100%;
		width: 100%;
		position: absolute;
		left: 0;
		bottom: -1px;
		border-bottom: 2px solid rgb(64, 158, 255);
		transform: scaleX(0);
		transition: 0.3s cubic-bezier(0.25, 0.8, 0.5, 1);
	}

	.content {
		position: absolute;
		bottom: 10px;
		left: 0;
		transition: 0.3s cubic-bezier(0.25, 0.8, 0.5, 1);
	}

	.content .chip {
		color: #fff;
		background-color: rgb(64, 158, 255);
		border-radius: 4px;
		display: inline;
		padding: 4px;
		text-align: center;
	}

	.content .chip:not(:first-of-type) {
		margin-left: 8px;
	}

	.container label.active .content {
		transform: translateY(-150%);
		opacity: 0;
		font-size: 14px;
		color: rgb(64, 158, 255);
	}

	.container label.active::after {
		transform: scaleX(1);
	}

	.dropdown-item {
		margin: 12px;
	}

	.dropdown-item span {
		margin-left: 8px;
	}

	.dropdown-content {
		background-color: white;
		position: absolute;
		min-width: 200px;
		width: 100%;
		z-index: 1;
		box-shadow: rgba(0, 0, 0, 0.1) 0px 2px 12px 0px;
	}
</style>
