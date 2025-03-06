<script>
  import "../../app.css";
  import "../../mybookmark.css";
  import { Toaster } from "$lib/components/ui/sonner";
  import { ModeWatcher } from "mode-watcher";
  import Header from "$lib/Header.svelte";
  import { toast } from "svelte-sonner";
</script>

<ModeWatcher />
<div class="app bg-background">
  <main>
    <slot></slot>
  </main>
</div>
