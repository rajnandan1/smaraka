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
  import BrowserBack from "$lib/BrowserBack.svelte";
  let ghAllChecked = true;
  let uploading = false;
  let importInProgress = "";
  let logos = [
    {
      name: "Mozilla Firefox",
      image: `${base}/firefox.png`,
      selected: false,
      guide: "https://support.mozilla.org/en-US/kb/export-firefox-bookmarks-to-backup-or-transfer"
    },
    {
      name: "Google Chrome",
      image: `${base}/chrome.png`,
      selected: false,
      guide: "https://support.google.com/chrome/answer/96816?hl=en"
    },
    {
      name: "Apple Safari",
      image: `${base}/safari.png`,
      selected: false,
      guide: "https://www.ionos.com/digitalguide/websites/web-development/export-safari-bookmarks/"
    },
    {
      name: "Pocket",
      image: `${base}/pocket.png`,
      selected: false,
      guide: "https://support.mozilla.org/en-US/kb/exporting-your-pocket-list"
    },
    {
      name: "Raindrop",
      image: `${base}/raindrop.png`,
      selected: false,
      guide: "https://help.raindrop.io/export"
    }
  ];
  onMount(() => {
    //get hash parameter
  });
  let urls = [];

  function initURLs(res) {
    return res
      .map((url) => {
        url.checked = true;
        return url;
      })
      .reverse();
  }
  let file = null;
  function submitFile() {
    if (file) {
      uploading = true;
      uploadFirefoxFile(file)
        .then((res) => {
          urls = initURLs(res);
          uploading = false;
        })
        .catch((err) => {
          console.log(err);
          toast.error("Failed to fetch URLs from file");
          uploading = false;
        });
    }
  }
  const handleFileChange = (event) => {
    file = event.target.files[0];
    submitFile();
  };

  function selectUnselectAllGH() {
    urls = urls.map((gh) => {
      gh.checked = !ghAllChecked;
      return gh;
    });
  }
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
  function clearFile() {
    file = null;
    urls = [];
    document.getElementById("fileupload").value = "";
  }
  function handleBrowserClick(event) {
    let logo = event.detail;
    window.open(logo.guide, "_blank");
  }
</script>

<div class="container mx-auto w-full max-w-2xl">
  <BrowserBack />
  <h1 class="text-lg font-medium leading-10">Import From Other Apps</h1>
  <p class="text-sm text-secondary-foreground">Learn how to export from popular browsers and services</p>
  <div class="mt-4">
    <LogoBar {logos} on:click={handleBrowserClick} />
  </div>

  <div class="my-4 grid w-full grid-cols-4 items-center gap-1.5">
    <div class="relative col-span-3">
      <Label for="fileupload">Select HTML File</Label>
      <Input
        on:change={handleFileChange}
        id="fileupload"
        type="file"
        class="mt-1 cursor-pointer hover:bg-secondary"
        accept=".html"
        disabled={uploading}
      />
      {#if uploading}
        <LoaderCircle class="absolute right-3 top-10  h-4 w-4 animate-spin" />
      {:else if !uploading && !!file && urls.length > 0}
        <Button variant="ghost" on:click={clearFile} size="icon" class="absolute right-3 top-10 ml-2 h-4 w-4">
          <X class="h-4 w-4" />
        </Button>
      {/if}
    </div>
  </div>
  {#if urls.length > 0}
    <div class="mt-4">
      <Card.Root>
        <Card.Content>
          <div class="mb-4 mt-4 flex">
            <Checkbox
              onCheckedChange={selectUnselectAllGH}
              id="select-all-gh"
              disabled={importInProgress != ""}
              bind:checked={ghAllChecked}
              aria-labelledby="select-all-gh-label"
            />
            <Label id="select-all-gh-label" for="select-all-gh" class="text-md ml-2 font-normal leading-5">
              Fetched {urls.length} urls for
            </Label>
          </div>
          <hr />
          <div class="h-[30vh] overflow-y-auto">
            {#each urls as url, i}
              <div class="border-b py-2">
                <p class="overflow-hidden text-ellipsis whitespace-nowrap">
                  <Checkbox id="url-{i}" bind:checked={url.checked} disabled={importInProgress != ""} class="" />
                  <label for="url-{i}" class="ml-2 cursor-pointer text-foreground hover:text-accent-foreground"
                    >{url.name}</label
                  >
                </p>
                <p class="overflow-hidden text-ellipsis whitespace-nowrap">
                  <span class="pl-8 text-sm text-secondary-foreground">{url.url}</span>
                </p>
              </div>
            {/each}
          </div>

          <div class="py-4">
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
        </Card.Content>
      </Card.Root>
    </div>
  {/if}
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
