<script lang="ts">
  import { Button } from "$lib/components/ui/button/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import { fetchHFTokenStatus, saveHFToken, clearHFToken } from "../stores/api";

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

  void loadHFTokenStatus();
</script>

<div class="p-2">
  <div class="mt-4 mb-6">
    <p class="font-mono text-[10px] tracking-[0.15em] uppercase text-primary/70">SET.01 — Integrations</p>
    <h3 class="text-lg font-semibold">Integrations</h3>
  </div>

  <div class="max-w-md rounded-lg border border-l-2 border-l-primary/40 p-4 space-y-4">
    <div class="flex items-start justify-between gap-4">
      <div class="space-y-1">
        <h4 class="text-sm font-semibold text-muted-foreground">Hugging Face</h4>
        <p class="text-xs text-muted-foreground">API token for model downloads and gated repos.</p>
      </div>
      <div class="flex items-center gap-1.5 pt-0.5">
        <span
          class="size-1.5 rounded-full {hfTokenConfigured ? 'bg-success' : 'bg-muted-foreground/40'}"
        ></span>
        <span class="font-mono text-[10px] uppercase tracking-wide text-muted-foreground">
          {hfTokenConfigured ? "configured" : "unset"}
        </span>
      </div>
    </div>

    <div class="space-y-1.5">
      <Label for="hf-token-input" class="font-mono text-xs">API_TOKEN</Label>
      <Input
        id="hf-token-input"
        type="password"
        placeholder={hfTokenConfigured ? "•••••••••••••••• (replace)" : "hf_..."}
        bind:value={hfTokenInput}
        disabled={hfTokenStatus === "saving"}
        class="font-mono text-sm"
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
        <Button size="sm" variant="ghost" onclick={handleClearHFToken} disabled={hfTokenStatus === "saving"}>
          Remove
        </Button>
      {/if}
    </div>
  </div>
</div>
