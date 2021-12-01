<script lang="ts">
	import { Notify } from '$lib/store';
	import { onMount } from 'svelte';
	import SButton from './SButton.svelte';
	import SCard from './SCard.svelte';
	import SCardActions from './SCardActions.svelte';
	import SCardContent from './SCardContent.svelte';
	import SCardTitle from './SCardTitle.svelte';
	import SInput from './SInput.svelte';
	import SScroller from './SScroller.svelte';
	import SSwitch from './SSwitch.svelte';

	var backendURL = 'ws://localhost:8080/api/rcon';

	let autoscroll = true;
	let command = '';
	let logs: string[] = [];
	let conn;

	function onKeyPress(event: KeyboardEvent) {
		if (event.key == 'Enter') {
			submit();
		}
	}

	function submit() {
		Notify({
			message: 'Sending command...',
			type: 'info',
			dismissible: false,
			timeout: 3000
		});
		if (!conn) {
			return;
		}
		if (command === '') {
			Notify({
				message: 'Need to specify command!',
				type: 'error',
				dismissible: true,
				timeout: 3000
			});
			return;
		}
		conn.send(command);
		command = '';
	}

	onMount(() => {
		if (import.meta.env.MODE == 'production') {
			backendURL = 'wss://' + document.location.host + '/api/rcon';
		}
		conn = new WebSocket(backendURL);
		conn.onerror = (event: ErrorEvent) => {
			Notify({
				message: `Could not establish WebSocket conneciton. Try reloading the page. ${JSON.stringify(
					event,
					['message', 'arguments', 'type', 'name']
				)}`,
				type: 'error',
				dismissible: true,
				timeout: 3000
			});
			conn = null;
		};
		conn.onclose = function () {
			if (conn == null) {
				return;
			}
			Notify({
				message: 'Connection closed.',
				type: 'error',
				dismissible: true,
				timeout: 3000
			});
		};
		conn.onmessage = (event: MessageEvent<string>) => {
			logs = [...logs, ...event.data.split('\n')];
		};
	});
</script>

<SCard width="600px" height="400px">
	<SCardTitle>RCON</SCardTitle>
	<SCardContent>
		<SScroller {autoscroll}>
			{#each logs as log}
				<div class="msg">{log}</div>
			{/each}
		</SScroller>
	</SCardContent>
	<SCardActions>
		<div class="actions">
			<div class="settings">
				<SSwitch bind:value={autoscroll}>Autoscroll:</SSwitch>
			</div>
			<div>
				<SInput type="text" placeholder="Command" bind:value={command} on:keypress={onKeyPress} />
				<SButton on:click={submit} disabled={command.length == 0}>Send</SButton>
			</div>
		</div>
	</SCardActions>
</SCard>

<style>
	.msg {
		margin-top: 3px;
		padding: 5px;
		border-radius: 3px;
		background: #f5f5f5;
	}

	.actions {
		width: 100%;
	}

	.actions > div {
		display: flex;
	}

	.actions > .settings {
		justify-content: flex-end;
		margin: 8px;
	}

	.disabled {
		background: red;
	}
</style>
