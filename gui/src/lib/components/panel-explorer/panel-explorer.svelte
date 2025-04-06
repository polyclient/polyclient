<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import Icon from '$lib/icons/icon.svelte';
	import { iconsNames } from '$lib/icons/types';

	// Menu items.
	const items = [
		{
			title: 'Home',
			url: '#',
			icon: iconsNames.dbDatabase,
		},
		{
			title: 'Inbox',
			url: '#',
			icon: iconsNames.dbDatabase,
		},
		{
			title: 'Calendar',
			url: '#',
			icon: iconsNames.dbDatabase,
		},
		{
			title: 'Search',
			url: '#',
			icon: iconsNames.dbDatabase,
		},
		{
			title: 'Settings',
			url: '#',
			icon: iconsNames.dbDatabase,
		},
	];
</script>

<Sidebar.Provider>
	<Sidebar.Root
		collapsible="none"
		class="w-full"
	>
		<Sidebar.Content>
			<Sidebar.Group>
				<Sidebar.GroupLabel>Application</Sidebar.GroupLabel>
				<Sidebar.GroupContent>
					<Sidebar.Menu>
						{#each items as item (item.title)}
							<Sidebar.MenuItem>
								<Sidebar.MenuButton>
									{#snippet child({ props })}
										<a
											href={item.url}
											{...props}
										>
											<Icon name={item.icon} />
											<span>{item.title}</span>
										</a>
									{/snippet}
								</Sidebar.MenuButton>
							</Sidebar.MenuItem>
						{/each}
					</Sidebar.Menu>
				</Sidebar.GroupContent>
			</Sidebar.Group>
		</Sidebar.Content>
	</Sidebar.Root>
</Sidebar.Provider>
