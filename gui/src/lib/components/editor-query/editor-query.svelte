<script lang="ts">
	import extensionMoveToLine from '$lib/codemirror/extensions/move-to-line';
	import { indentWithTab } from '@codemirror/commands';
	import { SQLite, sql } from '@codemirror/lang-sql';
	import { Compartment, EditorState } from '@codemirror/state';
	import { keymap } from '@codemirror/view';
	import { EditorView, basicSetup } from 'codemirror';

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
			extensionMoveToLine,
			basicSetup,
			keymap.of([indentWithTab]),
			language.of(
				sql({
					dialect: SQLite,
					schemas: [],
					tables: [],
				}),
			),
			tabSize.of(EditorState.tabSize.of(4)),
		],
		doc,
	});

	$effect(() => {
		const editor = new EditorView({ state, parent });
	});
</script>

<div
	bind:this={parent}
	class="h-full overflow-auto"
></div>
