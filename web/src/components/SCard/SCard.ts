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

export const SCardTitle = createComponent('SCardContent', 'card__title');
export const SCardContent = createComponent('SCardContent', 'card__content');
export const SCardActions = createComponent('SCardActions', 'card__actions');
