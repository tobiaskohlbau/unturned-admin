<script lang="ts">
	let visible = false;

	export function open(): void {
		visible = true;
	}

	export function close(): void {
		visible = false;
	}
</script>

<div>
	<slot name="activator" class="popup_actor" :on="listener" />
	{#if visible}
		<div class="popup">
			<div class="popup_content" v-click-outside="closeConditional">
				<slot class="popup_content" />
			</div>
		</div>
	{/if}
</div>

<style>
	.popup {
		position: absolute;
		left: 0;
		top: 0;
		background-color: rgba(22, 22, 22, 0.5);
		width: 100%;
		height: 100%;
		z-index: 100;
	}

	.popup_content {
		display: inline-block;
		position: relative;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		-webkit-box-shadow: 0 0 10px rgba(0, 0, 0, 0.4);
		-moz-box-shadow: 0 0 10px rgba(0, 0, 0, 0.4);
		box-shadow: 0 0 10px rgba(0, 0, 0, 0.4);
	}
</style>
