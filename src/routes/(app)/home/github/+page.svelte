<script>
  import Header from "$lib/Header.svelte";
  import { page } from "$app/stores";
  import logo from "$lib/images/svelte-logo.svg";
  import { Button } from "$lib/components/ui/button";
  import confetti from "canvas-confetti";
  import * as Dialog from "$lib/components/ui/dialog";
  import { base } from "$app/paths";
  import {
    Plus,
    ChevronDown,
    Github,
    Star,
    LoaderCircle,
    Check,
    ListChecks,
    ArrowRight,
    ArrowLeft
  } from "lucide-svelte";
  import { Input } from "$lib/components/ui/input";
  import { createEventDispatcher } from "svelte";
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu";
  import * as Avatar from "$lib/components/ui/avatar";
  import * as Drawer from "$lib/components/ui/drawer";
  import { Checkbox } from "$lib/components/ui/checkbox";
  import { Label } from "$lib/components/ui/label";
  import * as Card from "$lib/components/ui/card";
  import { fetchGithubStarredRepos, bulkAddUrls } from "$lib/api";
  import BrowserBack from "$lib/BrowserBack.svelte";
  let ghUsername = "";
  let importInProgress = "";
  let fetchInProgress = "";
  let currentImportProgress = 0;

  let ghData = [];

  function fetchGHRepos() {
    fetchInProgress = "ONGOING";
    ghData = [];
    importInProgress = "";
    fetchGithubStarredRepos(ghUsername)
      .then((res) => {
        ghData = res;
        fetchInProgress = "COMPLETE";
        initghData();
      })
      .catch((err) => {
        console.log(err);
        fetchInProgress = "FAILED";
      });
  }

  var defaults = {
    spread: 360,
    ticks: 50,
    gravity: 0,
    decay: 0.94,
    startVelocity: 30,
    colors: ["FFE400", "FFBD00", "E89400", "FFCA6C", "FDFFB8"]
  };

  function shoot() {
    confetti({
      ...defaults,
      particleCount: 40,
      scalar: 1.2,
      shapes: ["star"]
    });

    confetti({
      ...defaults,
      particleCount: 10,
      scalar: 0.75,
      shapes: ["circle"]
    });
  }

  function initghData() {
    ghData = ghData.map((gh) => {
      gh.checked = true;
      gh.status = "READY";
      return gh;
    });
  }

  function selectUnselectAllGH() {
    ghData = ghData.map((gh) => {
      gh.checked = !ghAllChecked;
      return gh;
    });
  }

  function beginImport() {
    importInProgress = "ONGOING";
    let urls = [];
    for (let i = 0; i < ghData.length; i++) {
      if (ghData[i].checked) {
        urls.push(ghData[i].url);
      }
    }

    bulkAddUrls(urls, "reverse").then(
      function () {
        importInProgress = "SUCCESS";
        setTimeout(shoot, 0);
        setTimeout(shoot, 200);
        setTimeout(shoot, 400);
      },
      function (err) {
        console.log(err);
        importInProgress = "FAILED";
      }
    );
  }

  let ghAllChecked = true;
</script>

<div class="container mx-auto w-full max-w-2xl">
  <BrowserBack />
  <h1 class="text-lg font-medium leading-10">Import Github Stars</h1>
  <p class="text-sm text-secondary-foreground">Import all the starred github repositories using username</p>
  <div class="mt-4">
    <form class="relative flex items-center gap-2" on:submit|preventDefault={fetchGHRepos}>
      <Input
        placeholder="username"
        bind:value={ghUsername}
        class="z-2 text-md pl-[35px] md:pl-[158px]"
        id="gh-un"
        autocapitalize="off"
        required
      />
      <Label id="gh-un-label" for="gh-un" class="text-md z-1  absolute left-3 top-2  font-normal  ">
        <span class="hidden md:block">https://github.com/</span>
        <span class="md:hidden"><Github class="inline h-4 w-4" />/</span>
      </Label>
      <Button type="submit" disabled={fetchInProgress == "ONGOING"}>
        Fetch
        {#if fetchInProgress == "ONGOING"}
          <LoaderCircle class="ml-2 h-4 w-4 animate-spin" />
        {/if}
      </Button>
    </form>
    {#if fetchInProgress == "FAILED"}
      <p class="mt-2 text-sm text-red-500">Could not fetch repositories</p>
    {:else if fetchInProgress == "ONGOING"}
      <p class="mt-2 text-sm text-secondary-foreground">This usually takes some time...</p>
    {/if}
  </div>
  {#if fetchInProgress == "COMPLETE"}
    <div class="mt-4">
      <Card.Root>
        <Card.Content>
          <div class="mb-2 mt-4 flex">
            <Checkbox
              onCheckedChange={selectUnselectAllGH}
              id="select-all-gh"
              disabled={importInProgress != ""}
              bind:checked={ghAllChecked}
              aria-labelledby="select-all-gh-label"
            />
            <Label id="select-all-gh-label" for="select-all-gh" class="text-md ml-2 font-normal leading-5">
              Fetched {ghData.length} repositories starred by <i>{ghUsername}</i>
            </Label>
          </div>
          <div class="h-[35vh] overflow-y-auto">
            {#each ghData as gh, i}
              <div class="gh-repo relative {i < ghData.length - 1 ? 'border-b' : ''} py-4 pl-7">
                <Checkbox bind:checked={gh.checked} disabled={importInProgress != ""} class="absolute left-0 top-5" />

                <p class="overflow-hidden text-ellipsis text-nowrap">
                  <a href={gh.url} target="_blank">{gh.title}</a>
                </p>
                <p class="overflow-hidden text-ellipsis text-nowrap text-sm text-secondary-foreground">
                  <Star class="inline-block h-3 w-3 text-yellow-600" />
                  {gh.stars.toLocaleString()}.
                  {gh.excerpt}
                </p>
                <div class="absolute right-2 top-5">
                  {#if gh.status === "PROCESSING"}
                    <LoaderCircle class="mr-2 h-4 w-4 animate-spin" />
                  {:else if gh.status === "COMPLETE"}
                    <Check class="mr-2 h-4 w-4 text-green-500" />
                  {:else if gh.status === "FAILED"}
                    <Button variant="outline" size="sm">Retry</Button>
                  {/if}
                </div>
              </div>
            {/each}
          </div>
          <div class="py-4">
            {#if importInProgress == "" || importInProgress == "FAILED"}
              <Button class="w-full" on:click={beginImport}>
                Import {ghData.filter((gh) => gh.checked).length}
                {ghData.filter((gh) => gh.checked).length > 1 ? "repositories" : "repository"}
              </Button>
            {:else if importInProgress == "ONGOING"}
              <Button class="w-full" disabled>
                Importing {ghData.filter((gh) => gh.checked).length} repositories
                <LoaderCircle class="ml-2 h-4 w-4 animate-spin" />
              </Button>
            {:else if importInProgress == "SUCCESS"}
              <Button class="w-full" href="{base}/">
                <ListChecks class="mr-2 h-4 w-4" />
                See Imported
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
      <p class="text-lg">Repositories Under Process</p>
      <p class="text-sm text-secondary-foreground">
        We have started importing the repositories. Once a repository is imported, you will see it the main home page.
        Depending on the number of repositories, this might take a few minutes.
      </p>
      <Button class="mt-4 w-full" href="{base}/home">
        View Bookmarks
        <ArrowRight class="ml-2 h-4 w-4" />
      </Button>
    </div>
  </Dialog.Content>
</Dialog.Root>
