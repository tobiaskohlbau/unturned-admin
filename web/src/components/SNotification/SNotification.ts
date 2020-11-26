import SNotificationComponent from "./SNotification.vue"

import {
  createVNode,
  nextTick,
  render,
  VNode,
} from "vue";

let seed = 1;
const instances: Array<{vm: VNode, $el: HTMLElement}> = []

export function Notify(message: string) {
  const container = document.createElement('div');

  const id = 'message_' + seed++;

  const vm = createVNode(
    SNotificationComponent,
    {
      message: message,
      onClose: () => {
        close(id)
      },
      id: id,
    },
    null,
  )

  render(vm, container)
  instances.push({ vm, $el: container });
  document.body.appendChild(container)
};

export function close(id: string): void {
  const idx = instances.findIndex(({ vm }) => {
    return id === vm.component?.props.id
  })
  if (idx === -1) {
    return
  }

  const { vm, $el } = instances[idx]
  if (!vm) return

  render(null, $el)
  nextTick(() => {
    document.body.removeChild($el)
  })


  instances.splice(idx, 1)
}
