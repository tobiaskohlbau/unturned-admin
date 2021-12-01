<script lang="ts">
	export let value = false;

	function toggle(): void {
		value = !value;
	}
</script>

<div class="wrapper">
	<label for="input" class="label">
		<slot />
	</label>
	<div class="container" on:click={toggle} on:click>
		<input type="checkbox" name="input" {value} />
		<div class="slider" class:on={value} />
		<div class="point" class:on={value} />
	</div>
</div>

<style>
	.wrapper {
		display: flex;
		align-items: center;
	}

	.container {
		width: 38px;
		height: 24px;
		position: relative;
		overflow: hidden;
		display: inline-block;
	}

	.label {
		margin-right: 8px;
	}

	.container input {
		position: absolute;
		opacity: 0;
	}

	.slider {
		position: absolute;
		width: 100%;
		height: 14px;
		top: calc(50% - 7px);
		box-sizing: border-box;
		border-radius: 12px;
		background-color: #6d7382;
	}

	.slider.on {
		background-color: #54b2ff;
	}

	.point {
		position: absolute;
		width: 20px;
		height: 20px;
		top: calc(50% - 10px);
		box-sizing: border-box;
		border-radius: 50%;
		transition: 0.3s cubic-bezier(0.25, 0.8, 0.5, 1);
		background-color: #595f6e;
	}

	.point.on {
		background-color: rgb(64, 158, 255);
		transform: translateX(18px);
	}
</style>
