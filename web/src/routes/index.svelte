<script context="module" lang="ts">
	import { goto } from '$app/navigation';
	import Manager from '$lib/components/Manager.svelte';
	import RCON from '$lib/components/RCON.svelte';
	import SButton from '$lib/components/SButton.svelte';
	import UnturnedMap from '$lib/components/UnturnedMap.svelte';
	import Users from '$lib/components/Users.svelte';
	import { User } from '$lib/models';
	import { user } from '$lib/store';
	import type { LoadInput, LoadOutput } from '@sveltejs/kit';

	export async function load({ session }: LoadInput): Promise<LoadOutput> {
		if (session) {
			user.set(
				new User({
					username: session.Username,
					activated: session.Activated,
					permissions: session.Permissions
				})
			);
		}
		return {};
	}
</script>

<script lang="ts">
	function logout() {
		document.cookie = 'tid=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
		goto('/login');
	}
</script>

<div class="flex">
	<div>
		<p>Welcome {$user.username}!</p>
		<SButton on:click={logout}>Logout</SButton>
	</div>
	<div class="container flex flex-grow flex-wrap justify-around gap-1">
		{#if $user.hasPermission('MODERATOR')}
			<RCON />
			<Manager />
		{/if}
		<UnturnedMap />
		{#if $user.hasPermission('ADMIN')}
			<Users />
		{/if}
	</div>
</div>
