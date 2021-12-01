export function css(
	node: HTMLElement,
	properties: Record<string, string>
): {
	update(properties: Record<string, string>): void;
} {
	function setProperties() {
		for (const prop of Object.keys(properties)) {
			node.style.setProperty(`--${prop}`, properties[prop]);
		}
	}

	setProperties();

	return {
		update(newProperties) {
			properties = newProperties;
			setProperties();
		}
	};
}
