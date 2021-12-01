<script lang="ts">
	import type { User } from '$lib/models';
	import { Notify } from '$lib/store';
	import { http, httpPost } from '$lib/utils';
	import { onMount } from 'svelte';
	import SCard from './SCard.svelte';
	import SCardActions from './SCardActions.svelte';
	import SCardContent from './SCardContent.svelte';
	import SCardTitle from './SCardTitle.svelte';
	import SScroller from './SScroller.svelte';
	import SSelect from './SSelect.svelte';
	import SSwitch from './SSwitch.svelte';

	const permissions = ['MODERATOR', 'ADMIN'];
	let users: User[] = [];

	onMount(async () => {
		try {
			const res = await http<User[]>('/api/users');
			users = res.parsedBody as User[];
		} catch (err) {
			Notify({
				message: 'Failed to fetch users!' + err,
				type: 'error',
				dismissible: true,
				timeout: 3000
			});
		}
	});

	async function saveUser(user) {
		try {
			await httpPost(`/api/users/${user.username}`, JSON.stringify(user));
		} catch (err) {
			Notify({
				message: `Failed to save user: ${user.username}` + err,
				type: 'error',
				dismissible: true,
				timeout: 3000
			});
		}
	}
</script>

<SCard width="600px" height="400px">
	<SCardTitle>Users</SCardTitle>
	<SCardContent>
		<SScroller>
			<div class="line">
				<div class="headline username">Username</div>
				<div class="headline">Permissions</div>
				<div class="headline activated">Activated</div>
			</div>
			{#each users as user (user.username)}
				<div class="line">
					<span class="username">{user.username}</span>
					<SSelect
						bind:value={user.permissions}
						items={permissions}
						on:change={() => saveUser(user)}
						placeholder="None"
					/>
					<SSwitch bind:value={user.activated} on:click={() => saveUser(user)} />
				</div>
			{/each}
		</SScroller>
	</SCardContent>
	<SCardActions />
</SCard>

<style>
	.line {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.headline {
		font-weight: 600;
		flex-grow: 1;
	}

	.username,
	.activated {
		width: 100px;
		flex: 0 0 auto;
	}
</style>
