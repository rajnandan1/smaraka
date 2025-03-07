<script>
  import { page } from "$app/stores";
  import logo from "$lib/images/svelte-logo.svg";
  import { Button } from "$lib/components/ui/button";
  import { base } from "$app/paths";
  import { toast } from "svelte-sonner";
  import * as Tooltip from "$lib/components/ui/tooltip";
  //import env
  import { PUBLIC_SMARAKA_FRONT_BASE } from "$env/static/public";

  import {
    Plus,
    ChevronDown,
    Github,
    House,
    Coffee,
    Receipt,
    BookLock,
    X,
    MessageCircleQuestion,
    Star,
    LoaderCircle,
    Check,
    ArrowRight,
    Globe,
    Link,
    Sun,
    Moon,
    LogOut,
    Search,
    EllipsisVertical,
    FileDown,
    FileClock
  } from "lucide-svelte";
  import { Input } from "$lib/components/ui/input";
  import { toggleMode } from "mode-watcher";
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu";
  import * as Avatar from "$lib/components/ui/avatar";
  import * as Drawer from "$lib/components/ui/drawer";
  import { Checkbox } from "$lib/components/ui/checkbox";
  import { Label } from "$lib/components/ui/label";
  import { addNewBookmarkURL, getMySubscription } from "$lib/api";
  import * as Dialog from "$lib/components/ui/dialog";
  import { onMount } from "svelte";
  import { afterNavigate } from "$app/navigation";
  import { format, formatDistance, addDays } from "date-fns";

  let githubDrawer = false;
  let quickAddOpen = false;
  let newBookmarkURL = "";
  let addingNewLoader = false;
  let totalBookmarks = {};

  function addKeyboardShortcuts(e) {
    if (e.key == "u" && e.shiftKey && e.metaKey) {
      e.preventDefault();
      if (!quickAddOpen) {
        quickAddOpen = true;
      }
    }
  }
  let pageID = $page.route.id;
  afterNavigate(() => {
    pageID = $page.route.id;
  });

  let userName = $page.data?.user?.name || "User";
  let userEmail = $page.data?.user?.email || "user@email.com";
  let subscription;
  let isFreeTrialExpired = false;
  onMount(() => {
    document.addEventListener("keydown", addKeyboardShortcuts);

    return () => {
      document.removeEventListener("keydown", addKeyboardShortcuts);
    };
  });

  async function addURLSingle() {
    addingNewLoader = true;
    if (!!!newBookmarkURL) {
      toast.error("URL cannot be empty");
      addingNewLoader = false;
      return;
    }
    try {
      let newData = await addNewBookmarkURL(newBookmarkURL);
      newData.visible = true;
      newBookmarkURL = "";
      quickAddOpen = false;
      const eventAwesome = new CustomEvent("newBookmarkAdded", {
        bubbles: true,
        detail: newData
      });
      document.dispatchEvent(eventAwesome);
    } catch (error) {
      toast.error(error.response.data?.message || "An error occurred");
    } finally {
      addingNewLoader = false;
    }
  }
  function onOpenChange(e) {
    quickAddOpen = !quickAddOpen;
  }
  function openSearchModal(e) {
    const eventAwesome = new CustomEvent("openSearchModal", {
      bubbles: true
    });
    document.dispatchEvent(eventAwesome);
  }
</script>

<svelte:window
  on:openNewBookmarkModal={(e) => {
    quickAddOpen = true;
  }}
/>

<div class="sticky top-0 z-40 mb-4 w-full backdrop-blur">
  <header class="container flex w-full max-w-6xl flex-wrap p-4 text-sm sm:flex-nowrap sm:justify-start">
    <nav class="mx-auto flex w-full basis-full flex-wrap items-center justify-between">
      <a
        class="flex-none text-xl font-medium focus:opacity-80 focus:outline-none dark:text-white sm:order-1"
        href="{base}/home"
      >
        <img src="{base}/smaraka.png" class="inline h-12" alt="" />
        <span class="hidden md:inline-block">Smaraka</span>
      </a>
      <div class="flex items-center gap-x-2 sm:order-3">
        <div class="hidden overflow-hidden rounded-md bg-primary md:block">
          <Tooltip.Root openDelay="100">
            <Tooltip.Trigger>
              <Button
                on:click={(e) => {
                  quickAddOpen = true;
                }}
                class=" rounded-none   border-r-2 border-accent  text-primary-foreground md:pl-3"
              >
                <Plus class="h-4 w-4 md:mr-2" />

                <span class="hidden md:inline-block">Add New Bookmark</span>
              </Button>
            </Tooltip.Trigger>
            <Tooltip.Content class="text-xs font-medium">⇧⌘U</Tooltip.Content>
          </Tooltip.Root>
          <DropdownMenu.Root>
            <DropdownMenu.Trigger>
              <Button size="icon" class="-ml-1 rounded-none text-primary-foreground">
                <ChevronDown class="h-4 w-4" />
              </Button>
            </DropdownMenu.Trigger>
            <DropdownMenu.Content>
              <DropdownMenu.Group>
                <DropdownMenu.Label>Import URLs</DropdownMenu.Label>
                <DropdownMenu.Separator />
                <DropdownMenu.Item class="cursor-pointer" href="{base}/home/github">
                  <Github class="mr-2 h-4 w-4" />
                  Github Starred Repos
                </DropdownMenu.Item>
                <DropdownMenu.Item class="cursor-pointer" href="{base}/home/browsers">
                  <Globe class="mr-2 h-4 w-4" /> Import Bookmarks
                </DropdownMenu.Item>
                <DropdownMenu.Item class="cursor-pointer" href="{base}/home/urls">
                  <Link class="mr-2 h-4 w-4" /> Bulk URLs
                </DropdownMenu.Item>
                <DropdownMenu.Item class="cursor-pointer" href="{base}/home/schedules">
                  <FileClock class="mr-2 h-4 w-4" /> Schedules
                </DropdownMenu.Item>
              </DropdownMenu.Group>
            </DropdownMenu.Content>
          </DropdownMenu.Root>
        </div>
        <div class="ml-3 mt-1 flex rounded-full border bg-card p-0">
          <DropdownMenu.Root>
            <DropdownMenu.Trigger>
              <Avatar.Root class=" mt-0 h-8 w-8 rounded-full ">
                <Avatar.Image src="{base}/avatar.png" alt="name" />
                <Avatar.Fallback>
                  {userEmail.charAt(0).toUpperCase()}
                </Avatar.Fallback>
              </Avatar.Root>
            </DropdownMenu.Trigger>
            <DropdownMenu.Content>
              <DropdownMenu.Group>
                <DropdownMenu.Item class="cursor-pointer" on:click={toggleMode}>
                  <Sun class="mr-2 h-4 w-4 rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
                  <Moon class="absolute mr-2 h-4 w-4 rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
                  Toggle theme
                </DropdownMenu.Item>
                <DropdownMenu.Item class="cursor-pointer" href="/api/ui/url/bookmarks-export">
                  <FileDown class="mr-2 h-4 w-4" />
                  Export Data
                </DropdownMenu.Item>
                <DropdownMenu.Item class="cursor-pointer" href="{base}/secrets">
                  <BookLock class="mr-2 h-4 w-4" />
                  Secrets
                </DropdownMenu.Item>
                <DropdownMenu.Item
                  target="_blank"
                  class="cursor-pointer"
                  href="https://chromewebstore.google.com/detail/okbookmarks/mokcokjphlglbpnlhanebmalbofjifdf"
                >
                  <img src="{base}/chrome.png" class="mr-2 h-4 w-4" />
                  Chrome Extension
                </DropdownMenu.Item>
                <DropdownMenu.Item
                  target="_blank"
                  class="cursor-pointer"
                  href="https://addons.mozilla.org/en-US/firefox/addon/okbookmarks/"
                >
                  <img src="{base}/firefox.png" class="mr-2 h-4 w-4" />
                  Firefox Add-on
                </DropdownMenu.Item>
                <DropdownMenu.Item class="cursor-pointer" target="_blank" href="https://discord.gg/32aHaqR7Z5">
                  <MessageCircleQuestion class="mr-2 h-4 w-4" />
                  Help & Support
                </DropdownMenu.Item>
                <DropdownMenu.Item class="cursor-pointer" href="https://github.com/rajnandan1/smaraka">
                  <Github class="mr-2 h-4 w-4" />
                  Github
                </DropdownMenu.Item>
                <DropdownMenu.Item class="cursor-pointer" href="https://github.com/rajnandan1/smaraka">
                  <Coffee class="mr-2 h-4 w-4" />
                  Buy me a coffee
                </DropdownMenu.Item>
                <DropdownMenu.Separator class="border" />
                <DropdownMenu.Item class="cursor-pointer" href="/api/ui/logout">
                  <LogOut class="mr-2 h-4 w-4" />
                  Logout
                </DropdownMenu.Item>
              </DropdownMenu.Group>
            </DropdownMenu.Content>
          </DropdownMenu.Root>
        </div>
      </div>
    </nav>
  </header>
</div>
<div class="fixed bottom-0 left-0 z-50 w-full bg-primary px-4 pb-4 pt-2 md:hidden">
  <div class="grid grid-cols-3">
    <div>
      {#if pageID == "/"}
        <Button size="icon" class="text-left" on:click={openSearchModal}>
          <Search class="h-6 w-6" />
        </Button>
      {:else}
        <Button size="icon" class="text-left" href={base}>
          <House class="h-6 w-6" />
        </Button>
      {/if}
    </div>
    <div class="text-center">
      <Button
        size="icon"
        class="text-left"
        on:click={(e) => {
          quickAddOpen = true;
        }}
      >
        <Plus class="h-6 w-6" />
      </Button>
    </div>
    <div class="text-right">
      <DropdownMenu.Root>
        <DropdownMenu.Trigger>
          <Button size="icon" class="text-left">
            <EllipsisVertical class="h-6 w-6" />
          </Button>
        </DropdownMenu.Trigger>
        <DropdownMenu.Content>
          <DropdownMenu.Group>
            <DropdownMenu.Label>Import URLs</DropdownMenu.Label>
            <DropdownMenu.Separator />
            <DropdownMenu.Item class="cursor-pointer" href="{base}/github">
              <Github class="mr-2 h-4 w-4" />
              Github Starred Repos
            </DropdownMenu.Item>
            <DropdownMenu.Item class="cursor-pointer" href="{base}/browsers">
              <Globe class="mr-2 h-4 w-4" /> Import Bookmarks
            </DropdownMenu.Item>
            <DropdownMenu.Item class="cursor-pointer" href="{base}/urls">
              <Link class="mr-2 h-4 w-4" /> Bulk URLs
            </DropdownMenu.Item>
          </DropdownMenu.Group>
        </DropdownMenu.Content>
      </DropdownMenu.Root>
    </div>
  </div>
</div>
<Dialog.Root open={quickAddOpen} {onOpenChange}>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Bookmark a URL</Dialog.Title>
      <Dialog.Description>
        <form class="relative my-4" on:submit|preventDefault={addURLSingle}>
          <Input
            type="url"
            class="h-12 rounded-sm pr-12"
            bind:value={newBookmarkURL}
            placeholder="https://example.com/blog"
          />

          <Button
            size="icon"
            type="submit"
            disabled={addingNewLoader || newBookmarkURL == ""}
            class=" absolute  right-1 top-1 rounded-sm transition-all"
          >
            {#if addingNewLoader}
              <LoaderCircle class="h-4 w-4 animate-spin" />
            {:else}
              <ArrowRight class="h-4 w-4" />
            {/if}
          </Button>
        </form>
      </Dialog.Description>
    </Dialog.Header>
  </Dialog.Content>
</Dialog.Root>
