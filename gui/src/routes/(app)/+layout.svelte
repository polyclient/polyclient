<script lang="ts">
	import { AppHeader } from '$lib/components/app-header';
	import { AppRail } from '$lib/components/app-rail';
	import { PanelExplorer } from '$lib/components/panel-explorer';
	import * as Resizable from '$lib/components/ui/resizable';

	const { children } = $props();
</script>

<div class="h-screen grid grid-rows-[36px_1fr_36px_22px]">
	<AppHeader />

	<Resizable.PaneGroup direction="horizontal">
		<Resizable.Pane
			collapsible
			defaultSize={16}
			minSize={5}
			maxSize={50}
		>
			<PanelExplorer />
		</Resizable.Pane>

		<Resizable.Handle withHandle />

		<Resizable.Pane
			defaultSize={100}
			minSize={10}
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
					defaultSize={40}
					minSize={10}
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
			defaultSize={true ? 0 : 34}
			minSize={0}
			maxSize={100}
		>
			<div class="flex h-full items-center justify-center p-6">
				<span class="font-semibold">Integrations</span>
			</div>
		</Resizable.Pane>
	</Resizable.PaneGroup>

	<AppRail />

	<div class="bg-blue-500"></div>
</div>

<style lang="postcss">
	:global([data-active='pointer']),
	:global([data-active='keyboard']) {
		@apply border-2;
	}
</style>
