<script>
  import { onMount } from "svelte";
  import * as Select from "$lib/components/ui/select";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import { Button } from "$lib/components/ui/button";
  import { uploadFirefoxFile, bulkAddUrls } from "$lib/api";
  import * as Card from "$lib/components/ui/card";
  import { Checkbox } from "$lib/components/ui/checkbox";
  import { LoaderCircle, ListChecks, ArrowRight, X } from "lucide-svelte";
  import { confettiEffectStar } from "$lib/utils";
  import * as Dialog from "$lib/components/ui/dialog";
  import { toast } from "svelte-sonner";
  import LogoBar from "$lib/components/logo-bar.svelte";
  import { base } from "$app/paths";
  import { Textarea } from "$lib/components/ui/textarea";
  import BrowserBack from "$lib/BrowserBack.svelte";

  let importInProgress = "";

  onMount(() => {
    //get hash parameter
  });
  let urls = [];

  function beginImport() {
    importInProgress = "ONGOING";
    let urlsOnly = [];
    for (let i = 0; i < urls.length; i++) {
      if (urls[i].checked) {
        urlsOnly.push(urls[i].url);
      }
    }

    bulkAddUrls(urlsOnly, "reverse").then(
      function () {
        confettiEffectStar();
        setTimeout(function () {
          importInProgress = "SUCCESS";
        }, 300);
      },
      function (err) {
        console.log(err);
        importInProgress = "FAILED";
      }
    );
  }
</script>

<div class="container mx-auto w-full max-w-2xl">
  <BrowserBack />
  <h1 class="text-lg font-medium leading-10">Import URLs in Bulk</h1>
  <p class="text-sm text-secondary-foreground">
    Add multiple URLs in the textbox below. Each URL should be on a new line
  </p>

  <div class="mt-4 grid w-full gap-1.5">
    <Label for="message-2" class="mb-1">URLs</Label>
    <Textarea
      placeholder="https://google.com
https://example.com"
      id="message-2"
      class="overflow-auto whitespace-pre"
      rows="16"
      on:change={(e) => {
        urls = e.target.value
          .split("\n")
          .map((url) => {
            return { url: url.trim(), checked: true };
          })
          .filter((url) => url.url != "" && url.url.startsWith("http"));
      }}
    />
    <p class="text-sm text-muted-foreground">Please make sure each URL is on a new line</p>
    {#if urls.length > 0}
      <div class="mt-2">
        {#if importInProgress == "" || importInProgress == "FAILED"}
          <Button
            class="shadow-bk w-full"
            on:click={beginImport}
            disabled={urls.filter((gh) => gh.checked).length == 0}
          >
            Import {urls.filter((gh) => gh.checked).length}
            {urls.filter((gh) => gh.checked).length > 1 ? "URLs" : "URL"}
          </Button>
        {:else if importInProgress == "ONGOING"}
          <Button class="w-full" disabled>
            Importing {urls.filter((gh) => gh.checked).length} URLs
            <LoaderCircle class="ml-2 h-4 w-4 animate-spin" />
          </Button>
        {:else if importInProgress == "SUCCESS"}
          <Button class="w-full" href="{base}/">
            <ListChecks class="mr-2 h-4 w-4" />
            See Imported URLs
          </Button>
        {/if}
      </div>
    {/if}
  </div>
</div>
<Dialog.Root open={importInProgress == "SUCCESS"}>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Import Successful</Dialog.Title>
    </Dialog.Header>
    <div>
      <p class="text-lg">However</p>
      <p class="text-sm text-secondary-foreground">
        We have started importing the URLs in the background. Once a URL is imported, you will see it the main home
        page. Depending on the number of URLs, this might take a few minutes.
      </p>
      <Button class="mt-4 w-full" href="{base}/home">
        View Bookmarks
        <ArrowRight class="ml-2 h-4 w-4" />
      </Button>
    </div>
  </Dialog.Content>
</Dialog.Root>
