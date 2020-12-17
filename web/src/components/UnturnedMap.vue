<template>
  <s-card width="600px" height="650px">
    <s-card-title>Map</s-card-title>
    <div style="display: none">
      <img src="/map.png" ref="map" />
    </div>
    <canvas ref="canvas" :style="{width: size + 'px', height: size + 'px'}"></canvas>
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

    onMounted(async () => {
      const resp = await fetch("/api/players");
      const players = await resp.json();
      canvas.value.width = props.size;
      canvas.value.height = props.size;
      const ctx = canvas.value.getContext("2d");
      map.value.addEventListener("load", (e) => {
        ctx.drawImage(map.value, 0, 0, props.size, props.size);
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

          ctx.font = "12px Arial";
          ctx.fillStyle = "black";
          ctx.fillText(item.id.playerName, positionX + offsetX, positionY + offsetY);
          ctx.beginPath();
          ctx.fillStyle = "red";
          ctx.arc(positionX, positionY, 5, 0, 2 * Math.PI);
          ctx.fill();
          ctx.stroke();
          ctx.save();
        });
      });
    });

    return {
      canvas,
      map,
    };
  },
});
</script>
