<script lang="ts">
	import { AppHeader } from '$lib/components/app-header';
	import { PanelExplorer } from '$lib/components/panel-explorer';
	import * as Resizable from '$lib/components/ui/resizable';
	import * as Sidebar from '$lib/components/ui/sidebar';

	const { children } = $props();
</script>

<div class="h-screen grid grid-rows-[auto_1fr] overflow-hidden">
	<AppHeader />

	<Resizable.PaneGroup direction="horizontal">
		<Resizable.Pane
			collapsible
			defaultSize={20}
			minSize={0}
			maxSize={50}
		>
			<PanelExplorer />
		</Resizable.Pane>

		<Resizable.Handle withHandle />

		<Resizable.Pane
			defaultSize={100}
			minSize={0}
			maxSize={100}
		>
			<Resizable.PaneGroup direction="vertical">
				<Resizable.Pane
					defaultSize={70}
					minSize={0}
					maxSize={100}
				>
					<main class="h-full">
						{@render children()}
					</main>
				</Resizable.Pane>

				<Resizable.Handle withHandle />

				<Resizable.Pane
					collapsible
					defaultSize={30}
					minSize={0}
					maxSize={100}
				>
					<div class="flex h-full items-center justify-center p-6">
						<span class="font-semibold">Results</span>
					</div>
				</Resizable.Pane>
			</Resizable.PaneGroup>
		</Resizable.Pane>

		<Resizable.Handle withHandle />

		<Resizable.Pane
			collapsible
			defaultSize={true ? 0 : 30}
			minSize={0}
			maxSize={100}
		>
			<div class="flex h-full items-center justify-center p-6">
				<span class="font-semibold">Integrations</span>
			</div>
		</Resizable.Pane>
	</Resizable.PaneGroup>
</div>
