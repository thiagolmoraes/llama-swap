<script lang="ts">
  import { link } from "svelte-spa-router";
  import * as Popover from "$lib/components/ui/popover/index.js";
  import { Puzzle, Palette, Info } from "@lucide/svelte";

  let { children }: { children: import("svelte").Snippet<[{ props: Record<string, unknown> }]> } = $props();

  let open = $state(false);

  const items = [
    { tab: "integrations", label: "Integrations", icon: Puzzle },
    { tab: "appearance", label: "Appearance", icon: Palette },
    { tab: "about", label: "About", icon: Info },
  ];
</script>

<Popover.Root bind:open>
  <Popover.Trigger>
    {#snippet child({ props })}
      {@render children({ props })}
    {/snippet}
  </Popover.Trigger>
  <Popover.Content side="top" align="end" class="w-48 p-1">
    {#each items as item (item.tab)}
      <a
        href="/settings?tab={item.tab}"
        use:link
        onclick={() => (open = false)}
        class="flex items-center gap-2 rounded-md px-2 py-1.5 text-sm hover:bg-accent hover:text-accent-foreground"
      >
        <item.icon class="size-4" />
        <span>{item.label}</span>
      </a>
    {/each}
  </Popover.Content>
</Popover.Root>
