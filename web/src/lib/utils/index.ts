import { http, httpDelete, httpPost, httpPut, HttpResponse } from './http';

export function clickOutside(
	node: HTMLElement,
	callback: () => void
): {
	destroy?: () => void;
} {
	const listener = (event: MouseEvent) => {
		if (node && !node.contains(event.target as Node) && !event.defaultPrevented) {
			callback();
		}
	};

	document.addEventListener('click', listener, true);

	return {
		destroy() {
			document.removeEventListener('click', listener, true);
		}
	};
}

export const snackbarKey = 'snackbar';

export { http, HttpResponse, httpPut, httpDelete, httpPost };
