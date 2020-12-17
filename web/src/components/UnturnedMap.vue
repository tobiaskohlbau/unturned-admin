<template>
  <s-card width="600px" height="700px">
    <s-card-title>Map</s-card-title>
    <s-card-content>
      <div style="display: none">
        <img src="/map.png" ref="map" />
      </div>
      <canvas ref="canvas" :style="{width: size + 'px', height: size + 'px'}"></canvas>
    </s-card-content>
    <s-card-actions>
      <s-button @click.stop="updatePlayers">Update</s-button>
    </s-card-actions>
  </s-card>
</template>

<style scoped lang="postcss">
</style>

<script lang="ts">
import {
  computed,
  defineComponent,
  onMounted,
  Ref,
  ref,
  watchEffect,
} from "vue";
export default defineComponent({
  name: "UnturnedMap",
  props: {
    size: {
      type: Number,
      default: 600,
    },
  },
  setup(props, ctx) {
    const canvas: Ref<HTMLCanvasElement> = ref(null);
    const map: Ref<HTMLImageElement> = ref(null);

    const updatePlayers = async () => {
      const canvasCtx = canvas.value.getContext("2d");
      canvasCtx.clearRect(0, 0, props.size, props.size);
      canvasCtx.drawImage(map.value, 0, 0, props.size, props.size);

      const resp = await fetch("/api/players");
      const players = await resp.json();
      players.forEach((item) => {
        const positionX = props.size * item.x;
        const positionY = props.size * item.y;

        const edge = 10;
        var offsetX = 10;
        if (positionX > (props.size - edge)) {
          offsetX -= 12 * item.id.playerName.length;
        }

        var offsetY = 5;
        if (positionY < edge) {
          offsetY = 10;
        }
        if (positionY > (props.size - edge)) {
          offsetY = -5;
        }

        canvasCtx.font = "12px Arial";
        canvasCtx.fillStyle = "black";
        canvasCtx.fillText(item.id.playerName, positionX + offsetX, positionY + offsetY);
        canvasCtx.beginPath();
        canvasCtx.fillStyle = "red";
        canvasCtx.arc(positionX, positionY, 5, 0, 2 * Math.PI);
        canvasCtx.fill();
        canvasCtx.stroke();
      });
    };

    onMounted(async () => {
      canvas.value.width = props.size;
      canvas.value.height = props.size;
      await new Promise(r => map.value.onload=r);
      updatePlayers();
    });

    return {
      canvas,
      map,
      updatePlayers,
    };
  },
});
</script>
