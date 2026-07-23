<script lang="ts">
  import { connectionState } from "../stores/theme";
  import { versionInfo } from "../stores/api";

  let streamConnected = $derived($connectionState === "connected");
</script>

<div class="p-2">
  <div class="mt-4 mb-6">
    <p class="font-mono text-[10px] tracking-[0.15em] uppercase text-primary/70">SET.03 — About</p>
    <h3 class="text-lg font-semibold">About</h3>
  </div>

  <div class="max-w-md rounded-lg border border-l-2 border-l-primary/40 p-4 space-y-1">
    <div class="flex items-center justify-between gap-4 py-1.5">
      <dt class="font-mono text-[11px] uppercase tracking-wide text-muted-foreground">event_stream</dt>
      <dd class="flex items-center gap-1.5">
        <span class="size-1.5 rounded-full {streamConnected ? 'bg-success' : 'bg-muted-foreground/40'}"></span>
        <span class="font-mono text-xs font-medium">{$connectionState ?? "unknown"}</span>
      </dd>
    </div>
    <div class="h-px bg-border"></div>
    <div class="flex items-center justify-between gap-4 py-1.5">
      <dt class="font-mono text-[11px] uppercase tracking-wide text-muted-foreground">version</dt>
      <dd class="font-mono text-xs font-medium">{$versionInfo?.version ?? "unknown"}</dd>
    </div>
    <div class="h-px bg-border"></div>
    <div class="flex items-center justify-between gap-4 py-1.5">
      <dt class="font-mono text-[11px] uppercase tracking-wide text-muted-foreground">commit</dt>
      <dd class="font-mono text-xs font-medium">{$versionInfo?.commit?.substring(0, 7) ?? "unknown"}</dd>
    </div>
    <div class="h-px bg-border"></div>
    <div class="flex items-center justify-between gap-4 py-1.5">
      <dt class="font-mono text-[11px] uppercase tracking-wide text-muted-foreground">build_date</dt>
      <dd class="font-mono text-xs font-medium">{$versionInfo?.build_date ?? "unknown"}</dd>
    </div>
  </div>
</div>
