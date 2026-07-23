<script lang="ts">
  import * as Select from "$lib/components/ui/select/index.js";
  import { themeName, themeMode, themes, type ThemeMode } from "../stores/theme";

  const modes: { value: ThemeMode; label: string }[] = [
    { value: "light", label: "Light" },
    { value: "dark", label: "Dark" },
    { value: "system", label: "System" },
  ];

  let themeLabel = $derived(themes.find((t) => t.value === $themeName)?.label ?? "Default");
  let modeLabel = $derived(modes.find((m) => m.value === $themeMode)?.label ?? "System");
</script>

<div class="p-2">
  <div class="mt-4 mb-6">
    <p class="font-mono text-[10px] tracking-[0.15em] uppercase text-primary/70">SET.02 — Appearance</p>
    <h3 class="text-lg font-semibold">Appearance</h3>
  </div>

  <div class="max-w-md rounded-lg border border-l-2 border-l-primary/40 p-4 space-y-1">
    <div class="flex items-center justify-between gap-4 py-2">
      <div>
        <span class="text-sm">Theme</span>
        <p class="font-mono text-[10px] uppercase tracking-wide text-muted-foreground">color palette</p>
      </div>
      <Select.Root
        type="single"
        value={$themeName}
        onValueChange={(v) => v && themeName.set(v as typeof $themeName)}
      >
        <Select.Trigger class="w-40 font-mono text-xs">{themeLabel}</Select.Trigger>
        <Select.Content>
          {#each themes as theme (theme.value)}
            <Select.Item value={theme.value} class="font-mono text-xs">{theme.label}</Select.Item>
          {/each}
        </Select.Content>
      </Select.Root>
    </div>
    <div class="h-px bg-border"></div>
    <div class="flex items-center justify-between gap-4 py-2">
      <div>
        <span class="text-sm">Mode</span>
        <p class="font-mono text-[10px] uppercase tracking-wide text-muted-foreground">light / dark</p>
      </div>
      <Select.Root
        type="single"
        value={$themeMode}
        onValueChange={(v) => v && themeMode.set(v as ThemeMode)}
      >
        <Select.Trigger class="w-40 font-mono text-xs">{modeLabel}</Select.Trigger>
        <Select.Content>
          {#each modes as mode (mode.value)}
            <Select.Item value={mode.value} class="font-mono text-xs">{mode.label}</Select.Item>
          {/each}
        </Select.Content>
      </Select.Root>
    </div>
  </div>
</div>
