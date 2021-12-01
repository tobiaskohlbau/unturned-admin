<script lang="ts" context="module">
	let needsActivation = false;
	export async function load({ session }: LoadInput): Promise<LoadOutput> {
		if (session) {
			needsActivation = session.Activated == false;
		}
		return {};
	}
</script>

<script lang="ts">
	import { goto } from '$app/navigation';
	import SButton from '$lib/components/SButton.svelte';
	import SCard from '$lib/components/SCard.svelte';
	import SCardActions from '$lib/components/SCardActions.svelte';
	import SCardContent from '$lib/components/SCardContent.svelte';
	import SCardTitle from '$lib/components/SCardTitle.svelte';
	import SInput from '$lib/components/SInput.svelte';
	import { Notify } from '$lib/store';
	import { httpPost } from '$lib/utils';
	import type { LoadInput, LoadOutput } from '@sveltejs/kit';

	let username = '';
	let password = '';

	async function login() {
		try {
			console.log('login');
			await httpPost('/api/login', JSON.stringify({ username, password }));
			goto('/');
		} catch (e) {
			Notify({
				message: 'failed to login',
				type: 'error',
				dismissible: true,
				timeout: 3000
			});
		}
	}
</script>

<div class="w-full flex flex-wrap justify-around gap-y-1">
	<SCard width="400px" height="400px">
		<SCardTitle>Login</SCardTitle>
		<SCardContent>
			{#if needsActivation}
				An administrator needs to activate your account. <br />
				Please try again later!
			{:else}
				<SCardContent>
					<SInput type="text" placeholder="Username" bind:value={username} />
					<SInput type="password" placeholder="Password" bind:value={password} />
				</SCardContent>
				<SCardActions>
					<SButton on:click={login}>Login</SButton>
					<a href="/api/login/steam" rel="external"
						><img
							alt="Login with Steam"
							src="https://steamcdn-a.akamaihd.net/steamcommunity/public/images/steamworks_docs/english/sits_large_noborder.png"
						/></a
					>
				</SCardActions>
			{/if}
		</SCardContent>
	</SCard>
</div>
