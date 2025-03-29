<script lang="ts">
	import { PanelExplorer } from '$lib/components/panel-explorer';
	import * as Resizable from '$lib/components/ui/resizable';
	import * as Sidebar from '$lib/components/ui/sidebar';

	const { children } = $props();
</script>

<div class="h-screen grid grid-rows-[auto_1fr]">
	<header class="flex justify-between">
		<div class="flex items-center gap-4">
			<nav>
				<a href="/">Namespace</a> /
				<a href="/about">Database</a>
			</nav>
		</div>
	</header>

	<Resizable.PaneGroup direction="horizontal">
		<Resizable.Pane
			collapsible
			defaultSize={20}
			minSize={0}
			maxSize={50}
		>
			<div class="flex h-full items-center justify-center p-6">
				<span class="font-semibold">Sidebar</span>
			</div>
		</Resizable.Pane>

		<Resizable.Handle withHandle />

		<Resizable.Pane
			defaultSize={100}
			minSize={0}
			maxSize={100}
		>
			<Resizable.PaneGroup direction="vertical">
				<Resizable.Pane
					defaultSize={80}
					minSize={0}
					maxSize={100}
				>
					<div class="flex h-full items-center justify-center p-6">
						<span class="font-semibold">Main</span>
					</div>
				</Resizable.Pane>

				<Resizable.Handle withHandle />

				<Resizable.Pane
					collapsible
					defaultSize={20}
					minSize={0}
					maxSize={100}
				>
					<div class="flex h-full items-center justify-center p-6">
						<span class="font-semibold">Terminal</span>
					</div>
				</Resizable.Pane>
			</Resizable.PaneGroup>
		</Resizable.Pane>

		<Resizable.Handle withHandle />

		<Resizable.Pane
			collapsible
			defaultSize={30}
			minSize={0}
			maxSize={100}
		>
			<div class="flex h-full items-center justify-center p-6">
				<span class="font-semibold">Integrations</span>
			</div>
		</Resizable.Pane>
	</Resizable.PaneGroup>
</div>
