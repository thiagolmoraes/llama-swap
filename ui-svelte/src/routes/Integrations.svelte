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
  <div class="mt-4 mb-4">
    <h3 class="text-lg font-semibold">Integrations</h3>
  </div>

  <div class="rounded-lg border p-4 space-y-3 max-w-md">
    <div class="space-y-1">
      <h4 class="text-sm font-semibold text-muted-foreground">Hugging Face</h4>
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
        <Button size="sm" variant="ghost" onclick={handleClearHFToken} disabled={hfTokenStatus === "saving"}>
          Remove
        </Button>
      {/if}
    </div>
  </div>
</div>
