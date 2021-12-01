<script lang="ts">
	import { onMount } from 'svelte';
	import SButton from './SButton.svelte';
	import SCard from './SCard.svelte';
	import SCardActions from './SCardActions.svelte';
	import SCardContent from './SCardContent.svelte';
	import SCardTitle from './SCardTitle.svelte';

	let canvas;
	let size = 600;

	onMount(async () => {
		updatePlayers();
	});

	async function updatePlayers() {
		const canvasCtx = canvas.getContext('2d');
		canvasCtx.clearRect(0, 0, size, size);

		const resp = await fetch('/api/players');
		const players = await resp.json();
		players.forEach((item) => {
			const positionX = size * item.x;
			const positionY = size * item.y;

			const edge = 10;
			var offsetX = 10;
			if (positionX > size - edge) {
				offsetX -= 12 * item.id.playerName.length;
			}

			var offsetY = 5;
			if (positionY < edge) {
				offsetY = 10;
			}
			if (positionY > size - edge) {
				offsetY = -5;
			}

			canvasCtx.font = '12px Arial';
			canvasCtx.fillStyle = 'black';
			canvasCtx.fillText(item.id.playerName, positionX + offsetX, positionY + offsetY);
			canvasCtx.beginPath();
			canvasCtx.fillStyle = 'red';
			canvasCtx.arc(positionX, positionY, 5, 0, 2 * Math.PI);
			canvasCtx.fill();
			canvasCtx.stroke();
		});
	}
</script>

<SCard width="600px" height="700px">
	<SCardTitle>Map</SCardTitle>
	<SCardContent>
		<img width={size + 'px'} height={size + 'px'} alt="map" src="/map.png" />
		<canvas width={size + 'px'} height={size + 'px'} bind:this={canvas} />
	</SCardContent>
	<SCardActions>
		<SButton on:click={updatePlayers}>Update</SButton>
	</SCardActions>
</SCard>

<style>
	img {
		position: absolute;
		z-index: 1;
	}
	canvas {
		position: relative;
		z-index: 20;
	}
</style>
