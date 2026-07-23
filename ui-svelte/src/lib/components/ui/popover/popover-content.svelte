<script lang="ts">
	import { cn, type WithoutChildrenOrChild } from "$lib/utils.js";
	import PopoverPortal from "./popover-portal.svelte";
	import { Popover as PopoverPrimitive } from "bits-ui";
	import type { ComponentProps } from "svelte";

	let {
		ref = $bindable(null),
		side = "top",
		sideOffset = 8,
		align = "end",
		portalProps,
		class: className,
		...restProps
	}: PopoverPrimitive.ContentProps & {
		portalProps?: WithoutChildrenOrChild<ComponentProps<typeof PopoverPortal>>;
	} = $props();
</script>

<PopoverPortal {...portalProps}>
	<PopoverPrimitive.Content
		bind:ref
		data-slot="popover-content"
		{side}
		{sideOffset}
		{align}
		class={cn(
			"data-open:animate-in data-closed:animate-out data-closed:fade-out-0 data-open:fade-in-0 data-closed:zoom-out-95 data-open:zoom-in-95 ring-foreground/10 bg-popover text-popover-foreground data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 z-50 w-80 rounded-lg p-0 shadow-md ring-1 duration-100 outline-none",
			className
		)}
		{...restProps}
	/>
</PopoverPortal>
