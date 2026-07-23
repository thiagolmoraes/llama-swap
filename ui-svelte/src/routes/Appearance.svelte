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
  <div class="mt-4 mb-4">
    <h3 class="text-lg font-semibold">Appearance</h3>
  </div>

  <div class="rounded-lg border p-4 space-y-3 max-w-md">
    <div class="flex items-center justify-between gap-4">
      <span class="text-sm">Theme</span>
      <Select.Root
        type="single"
        value={$themeName}
        onValueChange={(v) => v && themeName.set(v as typeof $themeName)}
      >
        <Select.Trigger class="w-40">{themeLabel}</Select.Trigger>
        <Select.Content>
          {#each themes as theme (theme.value)}
            <Select.Item value={theme.value}>{theme.label}</Select.Item>
          {/each}
        </Select.Content>
      </Select.Root>
    </div>
    <div class="flex items-center justify-between gap-4">
      <span class="text-sm">Mode</span>
      <Select.Root
        type="single"
        value={$themeMode}
        onValueChange={(v) => v && themeMode.set(v as ThemeMode)}
      >
        <Select.Trigger class="w-40">{modeLabel}</Select.Trigger>
        <Select.Content>
          {#each modes as mode (mode.value)}
            <Select.Item value={mode.value}>{mode.label}</Select.Item>
          {/each}
        </Select.Content>
      </Select.Root>
    </div>
  </div>
</div>
