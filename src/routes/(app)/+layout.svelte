<script>
  import "../../app.css";
  import "../../mybookmark.css";
  import { Toaster } from "$lib/components/ui/sonner";
  import { ModeWatcher } from "mode-watcher";
  import Header from "$lib/Header.svelte";
  import { toast } from "svelte-sonner";
  import { getUser } from "$lib/api";

  let user = {};

  function noauth() {
    toast("You are not authenticated", {
      type: "error",
      duration: 5000
    });
    setTimeout(() => {
      window.location.href = "/api/ui/logout";
    }, 3000);
  }

  export async function load() {
    try {
      let user = await getUser();
      return {
        user: user
      };
    } catch (error) {
      noauth();
    }
  }
</script>

<ModeWatcher />
<svelte:window on:noauth={noauth} />
<div class="app min-h-screen bg-background">
  <Toaster position="top-center" richColors />
  <Header></Header>
  <main>
    <slot></slot>
  </main>
</div>
