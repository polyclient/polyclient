import { SveltePlugin } from 'bun-plugin-svelte';

Bun.build({
	entrypoints: ['src/main.ts'],
	outdir: 'dist',
	target: 'browser',
	plugins: [
		SveltePlugin({
			development: true,
		}),
	],
});
