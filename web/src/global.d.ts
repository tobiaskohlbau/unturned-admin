/// <reference types="@sveltejs/kit" />

import monaco from 'monaco-editor';

declare global {
	declare interface Window {
		MonacoEnvironment: monaco.Environment;
	}
}
