import {
  defineComponent,
  h,
} from "vue";

import cssModules from './SCard.module.css';

function createComponent(name: string, clz: string) {
  return defineComponent({
    name: name,
    setup({}, { attrs, slots }) {
      return () => h(`div`, {
        ...attrs,
        class: cssModules[clz],
      },
      slots.default?.());
    },
  });
}

export const SCardTitle = createComponent('SCardContent', 'cardTitle');
export const SCardContent = createComponent('SCardContent', 'cardContent');
export const SCardActions = createComponent('SCardActions', 'cardActions');
