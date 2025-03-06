<script>
	import * as Avatar from "$lib/components/ui/avatar";
	import { createEventDispatcher } from "svelte";

	const dispatch = createEventDispatcher();
	function clicked(logo) {
		dispatch("click", logo);
	}
	export let logos = [];
</script>

<div class="flex flex-row items-center justify-start pl-2">
	{#each logos as logo, i}
		<div
			on:click={(e) => {
				clicked(logo);
			}}
		>
			<Avatar.Root
				title={logo.name}
				class="z-1 -ml-2  h-8 w-8 cursor-pointer rounded-full border-2 bg-secondary p-[2px]    {!!logo.selected
					? 'border border-accent-foreground'
					: ''} transition-all hover:z-10 hover:border"
			>
				<Avatar.Image src={logo.image} alt="name" />
				<Avatar.Fallback>{logo.name[0]}</Avatar.Fallback>
			</Avatar.Root>
		</div>
	{/each}
</div>
