import type { Extension } from '@codemirror/state';
import { type EditorView, keymap } from '@codemirror/view';

function moveToLine(view: EditorView) {
	const line = Number(prompt('Go to line:'));
	if (Number.isNaN(line) || line < 1 || line > view.state.doc.lines) {
		return false;
	}

	const position = view.state.doc.line(line).from;

	view.dispatch({
		selection: { anchor: position },
		userEvent: 'select',
		scrollIntoView: true,
		sequential: true,
	});

	return true;
}

const extension: Extension = keymap.of([
	{
		key: 'Alt-l',
		run: moveToLine,
	},
]);

export default extension;
