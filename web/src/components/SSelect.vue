<template>
  <div class="container" @click="dropdown = true">
    <div class="content">
      <div class="chip" v-for="item in modelValue" :key="item">
        <span>{{ item }}</span>
      </div>
    </div>
    <transition name="dropdown">
      <div v-if="dropdown" class="dropdown-content" v-click-outside="hideDropdown">
        <div v-for="item in items" :key="item" class="dropdown-item" @click="updateSelection(item)">
          <input type="checkbox" value="test" name="test" :checked="modelValue.includes(item)">
          <span>{{ item }}</span>
        </div>
      </div>
    </transition>
    <label for="input" :class="[{ 'active': (modelValue.length > 0) || dropdown}, 'label']">
      <span class="content">{{ placeholder }}</span>
    </label>
  </div>
</template>

<style scoped lang="postcss">
.container {
  width: 100%;
  height: 42px;
  position: relative;
  cursor: pointer;
}

.container label {
  position: absolute;
  bottom: 0px;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  border-bottom: 1px solid black;
}

.container label::after {
  content: "";
  height: 100%;
  width: 100%;
  position: absolute;
  left: 0;
  bottom: -1px;
  border-bottom: 2px solid rgb(64, 158, 255);
  transform: scaleX(0);
  transition: 0.3s cubic-bezier(0.25, 0.8, 0.5, 1);
}

.content {
  position: absolute;
  bottom: 10px;
  left: 0;
  transition: .3s cubic-bezier(.25,.8,.5,1);

  .chip {
    color: #fff;
    background-color: rgb(64, 158, 255);
    border-radius: 4px;
    display: inline;
    padding: 4px;
    text-align: center;
  }

  .chip:not(:first-of-type) {
    margin-left: 8px;
  }
}

.container label.active .content {
  transform: translateY(-150%);
  opacity: 0;
  font-size: 14px;
  color: rgb(64, 158, 255);
}

.container label.active::after {
  transform: scaleX(1);
}

.dropdown-item {
  margin: 12px;

  span {
    margin-left: 8px;
  }
}


.dropdown-content {
  background-color: white;
  position: absolute;
  min-width: 200px;
  width: 100%;
  z-index: 1;
  box-shadow: rgba(0, 0, 0, 0.1) 0px 2px 12px 0px;
}

.dropdown-enter-active,
.dropdown-leave-active {
  transition: .3s cubic-bezier(.25,.8,.5,1);
}

.dropdown-enter-from, .dropdown-leave-to {
  opacity: 0;
}
</style>

<script lang="ts">
import { computed, defineComponent, Ref, ref, watchEffect } from "vue";
export default defineComponent({
  name: "SSelect",
  props: {
    placeholder: {
      type: String,
      default: "",
    },
    modelValue: {
      type: Array,
      default: [],
    },
    items: {
      type: Array,
      default: [],
    },
  },
  setup(props, { emit }) {
    const dropdown = ref(false);

    const isSelected = (val: string) => {
      return props.modelValue.includes(val);
    };

    const updateSelection = (item: string) => {
      const index = props.modelValue.indexOf(item);
      if (index !== -1) {
        props.modelValue.splice(index, 1);
      } else {
        props.modelValue.push(item);
      }
      emit("update:modelValue", props.modelValue);
    };

    const hideDropdown = (e: Event) => {
      dropdown.value = false;
    };

    return {
      isSelected,
      updateSelection,
      dropdown,
      hideDropdown,
    };
  },
  emits: ["update:modelValue"],
});
</script>
