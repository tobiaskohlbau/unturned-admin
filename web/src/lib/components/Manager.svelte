<script lang="ts">
	import type { File } from '$lib/models';
	import { editorData, Notify } from '$lib/store';
	import { http, httpDelete, httpPut, HttpResponse } from '$lib/utils';
	import { onMount } from 'svelte';
	import Editor from './Editor.svelte';
	import SButton from './SButton.svelte';
	import SCard from './SCard.svelte';
	import SCardActions from './SCardActions.svelte';
	import SCardContent from './SCardContent.svelte';
	import SCardTitle from './SCardTitle.svelte';
	import SPopup from './SPopup.svelte';
	import SScroller from './SScroller.svelte';

	let path = '.';
	let listing: Response = {
		folders: [],
		files: []
	} as Response;
	let filepath = '.';

	let imagePopup;
	let editorPopup;
	let editor;

	interface Response {
		files: File[];
		folders: string[];
	}

	const openFile = async (file: File) => {
		filepath = file.path + '/' + file.name;

		if (file.content_type.includes('image')) {
			imagePopup.open();
			return;
		}

		if (filepath == null) {
			return;
		}

		http(`/api/files?path=${filepath}`, 'text').then((response) => {
			editorData.set(response.parsedBody as string);
			editorPopup.open();
		});
	};

	const deleteFile = async (file: File) => {
		try {
			await httpDelete<string[]>(`/api/files?path=${path}/${file.name}`);
			const index = listing.files.indexOf(file);
			if (index > -1) {
				listing.files.splice(index, 1);
			}
		} catch (err) {
			Notify({
				message: `Failed to delete file! ${err}`,
				type: 'error',
				dismissible: true,
				timeout: 3000
			});
		}
	};

	const updatePath = async () => {
		let response: HttpResponse<Response>;
		try {
			response = await http<Response>(`/api/files?path=${path}`);
			listing = response.parsedBody as Response;
		} catch (err) {
			Notify({
				message: 'Failed to fetch files!',
				type: 'error',
				dismissible: true,
				timeout: 3000
			});
		}
	};

	onMount(() => {
		updatePath();
	});

	function openFolder(folder: string) {
		if (folder == '..') {
			const idx = path.lastIndexOf('/');
			path = path.substring(0, idx);
		} else {
			path = `${path}/${folder}`;
		}
		updatePath();
	}

	function saveFile() {
		editor.save();
		httpPut(`/api/files?path=${filepath}`, $editorData).then(() => {
			editorPopup.close();
		});
	}

	async function scheduleUpdate() {
		try {
			await http<string[]>('/api/update');
			Notify({
				message: 'Server update requested. Wait a few minutes for update to complete!',
				type: 'info',
				dismissible: true,
				timeout: 3000
			});
		} catch (err) {
			Notify({
				message: `Failed to update server! ${err}`,
				type: 'error',
				dismissible: true,
				timeout: 3000
			});
		}
	}

	async function cancelUpdate() {
		try {
			await httpDelete<string[]>('/api/update');
			Notify({
				message: 'Update canceled',
				type: 'info',
				dismissible: true,
				timeout: 3000
			});
		} catch (err) {
			Notify({
				message: `Failed to cancel update! ${err}`,
				type: 'error',
				dismissible: true,
				timeout: 3000
			});
		}
	}

	async function backup() {
		try {
			await http<string[]>('/api/backup');
			Notify({
				message: 'Backup completed!',
				type: 'info',
				dismissible: true,
				timeout: 3000
			});
		} catch (err) {
			Notify({
				message: `Failed to backup server! ${err}`,
				type: 'error',
				dismissible: true,
				timeout: 3000
			});
		}
	}
</script>

<SCard width="600px" height="400px">
	<SCardTitle>Manager</SCardTitle>
	<SCardContent>
		<SScroller>
			<i
				class="material-icons button"
				on:click={() => path != '.' && openFolder('..')}
				class:enabled={path != '.'}>arrow_back</i
			>
			{#each listing.folders as folder}
				<div class="line">
					<div on:click={() => openFolder(folder)}>
						<i class="material-icons">folder</i>
						{folder}
					</div>
				</div>
			{/each}
			{#each listing.files as file}
				<div class="line" v-for="file in listing.files">
					<div on:click={() => openFile(file)}>
						<i class="material-icons">text_snippet</i>
						{file.name}
					</div>
					<i class="material-icons button_enabled" on:click={() => deleteFile(file)}>delete</i>
				</div>
			{/each}
		</SScroller>
	</SCardContent>
	<SCardActions>
		<SButton on:click={scheduleUpdate}>Update</SButton>
		<SButton on:click={cancelUpdate}>Cancle Update</SButton>
		<SButton on:click={backup}>Backup</SButton>
	</SCardActions>
	<SPopup bind:this={editorPopup}>
		<SCard width="700px" height="600px">
			<SCardTitle>Editor</SCardTitle>
			<SCardContent>
				<Editor bind:this={editor} />
			</SCardContent>
			<SCardActions>
				<SButton on:click={saveFile}>Save</SButton>
				<SButton on:click={editorPopup.close}>Close</SButton>
			</SCardActions>
		</SCard>
	</SPopup>
	<SPopup bind:this={imagePopup}>
		<SCard width="700px" height="600px">
			<SCardTitle>Imageviewer</SCardTitle>
			<SCardContent>
				<img alt="content of the choosen file" height="600px" src={'/api/files?path=' + filepath} />
			</SCardContent>
			<SCardActions>
				<SButton on:click={imagePopup.close}>Close</SButton>
			</SCardActions>
		</SCard>
	</SPopup>
</SCard>

<style>
	.button.enabled {
		color: #000000;
		cursor: pointer;
	}

	.button {
		color: #bebebe;
		cursor: not-allowed;
	}

	.line {
		display: flex;
		align-content: center;
	}

	.line > :first-child {
		color: #000000;
		cursor: pointer;
		display: flex;
		align-items: center;
		flex-grow: 1;
	}
</style>
