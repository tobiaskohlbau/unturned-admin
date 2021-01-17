import { Directive, DirectiveBinding } from "vue";

interface ClickOutisdeElement extends HTMLElement {
  _clickOutside (e: Event)
}

type ClickOutsideBinding = (Event) => void;

export const ClickOutside: Directive<ClickOutisdeElement, ClickOutsideBinding> = {
  mounted(el: ClickOutisdeElement, binding: DirectiveBinding<ClickOutsideBinding>) {
    const onClick = (e: Event) => {
      if (!el.contains(e.target as Node)) {
        binding.value(e);
      }
    };

    document.addEventListener('click', onClick, true);

    el._clickOutside = onClick
  },
  unmounted(el: ClickOutisdeElement) {
    if (!el._clickOutside) return

    document.removeEventListener('click', el._clickOutside, true);
    delete el._clickOutside;
  }
}
