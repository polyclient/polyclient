<script lang="ts">
	import extensionMoveToLine from '$lib/codemirror/extensions/move-to-line';
	import { indentWithTab } from '@codemirror/commands';
	import { SQLite, sql } from '@codemirror/lang-sql';
	import { Compartment, EditorState } from '@codemirror/state';
	import { keymap } from '@codemirror/view';
	import { EditorView, basicSetup } from 'codemirror';
	import { onMount } from 'svelte';

	let parent: HTMLDivElement;

	const language = new Compartment();
	const tabSize = new Compartment();

	const doc = `SELECT 
    Name,
    Address,
    PhoneNumber,
    DateOfBirth,
FROM Person
WHERE DateOfBirth > '2000-01-01'
ORDER BY Name;`;

	const state = EditorState.create({
		extensions: [
			// vim(),
			extensionMoveToLine,
			basicSetup,
			keymap.of([indentWithTab]),
			tabSize.of(EditorState.tabSize.of(4)),
			language.of(
				sql({
					dialect: SQLite,
					schemas: [],
					tables: [],
				}),
			),
		],
		doc,
	});

	onMount(() => {
		new EditorView({ state, parent });
	});
</script>

<div
	bind:this={parent}
	class="h-full overflow-auto"
></div>
