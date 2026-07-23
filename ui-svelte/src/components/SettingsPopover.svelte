<script lang="ts">
  import * as Popover from "$lib/components/ui/popover/index.js";
  import * as Tabs from "$lib/components/ui/tabs/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import * as Select from "$lib/components/ui/select/index.js";
  import { connectionState, themeName, themeMode, themes, type ThemeMode } from "../stores/theme";
  import { versionInfo, fetchHFTokenStatus, saveHFToken, clearHFToken } from "../stores/api";

  let { children }: { children: import("svelte").Snippet<[{ props: Record<string, unknown> }]> } = $props();

  const modes: { value: ThemeMode; label: string }[] = [
    { value: "light", label: "Light" },
    { value: "dark", label: "Dark" },
    { value: "system", label: "System" },
  ];

  let themeLabel = $derived(themes.find((t) => t.value === $themeName)?.label ?? "Default");
  let modeLabel = $derived(modes.find((m) => m.value === $themeMode)?.label ?? "System");

  let open = $state(false);
  let hfTokenConfigured = $state(false);
  let hfTokenInput = $state("");
  let hfTokenStatus = $state<"idle" | "saving" | "saved" | "error">("idle");
  let hfTokenError = $state("");

  async function loadHFTokenStatus(): Promise<void> {
    try {
      hfTokenConfigured = await fetchHFTokenStatus();
    } catch (error) {
      console.error(error);
    }
  }

  async function handleSaveHFToken(): Promise<void> {
    if (!hfTokenInput.trim()) return;
    hfTokenStatus = "saving";
    hfTokenError = "";
    try {
      await saveHFToken(hfTokenInput.trim());
      hfTokenInput = "";
      hfTokenConfigured = true;
      hfTokenStatus = "saved";
    } catch (error) {
      hfTokenStatus = "error";
      hfTokenError = error instanceof Error ? error.message : "Failed to save token";
    }
  }

  async function handleClearHFToken(): Promise<void> {
    hfTokenStatus = "saving";
    hfTokenError = "";
    try {
      await clearHFToken();
      hfTokenConfigured = false;
      hfTokenStatus = "idle";
    } catch (error) {
      hfTokenStatus = "error";
      hfTokenError = error instanceof Error ? error.message : "Failed to clear token";
    }
  }

  $effect(() => {
    if (open) {
      hfTokenStatus = "idle";
      hfTokenError = "";
      void loadHFTokenStatus();
    }
  });
</script>

<Popover.Root bind:open>
  <Popover.Trigger>
    {#snippet child({ props })}
      {@render children({ props })}
    {/snippet}
  </Popover.Trigger>
  <Popover.Content side="top" align="end" class="p-0">
    <Tabs.Root value="integrations" class="w-full">
      <Tabs.List class="w-full rounded-t-lg rounded-b-none px-2 pt-2">
        <Tabs.Trigger value="integrations">Integrations</Tabs.Trigger>
        <Tabs.Trigger value="appearance">Appearance</Tabs.Trigger>
        <Tabs.Trigger value="about">About</Tabs.Trigger>
      </Tabs.List>

      <Tabs.Content value="integrations" class="p-4 space-y-3">
        <div class="space-y-1">
          <h4 class="text-sm font-semibold">Hugging Face</h4>
          <p class="text-xs text-muted-foreground">
            {hfTokenConfigured ? "Token configured." : "No token configured."}
          </p>
        </div>
        <div class="space-y-1.5">
          <Label for="hf-token-input">API Token</Label>
          <Input
            id="hf-token-input"
            type="password"
            placeholder={hfTokenConfigured ? "•••••••••••••••• (replace)" : "hf_..."}
            bind:value={hfTokenInput}
            disabled={hfTokenStatus === "saving"}
          />
        </div>
        {#if hfTokenStatus === "error"}
          <p class="text-xs text-destructive">{hfTokenError}</p>
        {/if}
        <div class="flex items-center justify-between gap-2">
          <Button
            size="sm"
            onclick={handleSaveHFToken}
            disabled={!hfTokenInput.trim() || hfTokenStatus === "saving"}
          >
            {hfTokenStatus === "saving" ? "Saving…" : "Save"}
          </Button>
          {#if hfTokenConfigured}
            <Button
              size="sm"
              variant="ghost"
              onclick={handleClearHFToken}
              disabled={hfTokenStatus === "saving"}
            >
              Remove
            </Button>
          {/if}
        </div>
      </Tabs.Content>

      <Tabs.Content value="appearance" class="p-4 space-y-3">
        <div class="flex items-center justify-between gap-4">
          <span class="text-sm">Theme</span>
          <Select.Root
            type="single"
            value={$themeName}
            onValueChange={(v) => v && themeName.set(v as typeof $themeName)}
          >
            <Select.Trigger class="w-36">{themeLabel}</Select.Trigger>
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
            <Select.Trigger class="w-36">{modeLabel}</Select.Trigger>
            <Select.Content>
              {#each modes as mode (mode.value)}
                <Select.Item value={mode.value}>{mode.label}</Select.Item>
              {/each}
            </Select.Content>
          </Select.Root>
        </div>
      </Tabs.Content>

      <Tabs.Content value="about" class="p-4">
        <dl class="text-sm space-y-1">
          <div class="flex justify-between gap-4">
            <dt class="text-muted-foreground">Event Stream</dt>
            <dd class="font-medium">{$connectionState ?? "unknown"}</dd>
          </div>
          <div class="flex justify-between gap-4">
            <dt class="text-muted-foreground">Version</dt>
            <dd class="font-medium">{$versionInfo?.version ?? "unknown"}</dd>
          </div>
          <div class="flex justify-between gap-4">
            <dt class="text-muted-foreground">Commit Hash</dt>
            <dd class="font-medium">{$versionInfo?.commit?.substring(0, 7) ?? "unknown"}</dd>
          </div>
          <div class="flex justify-between gap-4">
            <dt class="text-muted-foreground">Build Date</dt>
            <dd class="font-medium">{$versionInfo?.build_date ?? "unknown"}</dd>
          </div>
        </dl>
      </Tabs.Content>
    </Tabs.Root>
  </Popover.Content>
</Popover.Root>
