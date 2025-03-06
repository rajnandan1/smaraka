<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import {
    ChevronRight,
    ChevronLeft,
    ArrowRight,
    QrCode,
    SquareArrowOutUpRight,
    LoaderCircle,
    X,
    CornerDownLeft,
    Plus,
    Star,
    Github,
    Globe,
    Trash2,
    Link
  } from "lucide-svelte";
  import * as Card from "$lib/components/ui/card";
  import { Checkbox } from "$lib/components/ui/checkbox";
  import { base } from "$app/paths";
  import { onMount, onDestroy } from "svelte";
  import { fetchBookmarks, callSearch, addNewBookmarkURL, deleteBulkBookmarks, countBookmarks } from "$lib/api";
  import Bookmark from "$lib/components/bookmarkCard.svelte";
  import BookmarkTable from "$lib/components/bookmarkTable.svelte";
  import JobQueue from "$lib/components/JobQueue.svelte";
  import Loader from "$lib/loader.svelte";
  import * as Dialog from "$lib/components/ui/dialog";
  import { Input } from "$lib/components/ui/input";
  import * as Drawer from "$lib/components/ui/drawer";
  import autoAnimate from "@formkit/auto-animate";

  export let data;

  let listElement;
  let limit = 25;
  let listEnd = false;
  let searchView = false;
  let searchedText = "";
  let nextID = "z";
  let firstID = "z";
  let pageNumber = 1;
  let pageIndex = [];

  function addKeyboardShortcuts(e) {
    if (e.key == "f" && e.shiftKey && e.metaKey) {
      e.preventDefault();
      document.getElementById("search-box").focus();
    }
  }

  function loadBookMarks(direction) {
    if (loading || searchView) return;
    loading = true;
    nextID = pageIndex[pageIndex.length - 1];

    fetchBookmarks(limit, nextID, "").then(
      (res) => {
        if (!!res.data) {
          loading = false;
          listEnd = res.is_last;
          bookmarks = res.data;
          nextID = res.next_id;
          setPageNumbers();
        }
      },
      function (err) {
        loading = false;
      }
    );
  }

  onMount(() => {
    loadFirstPage();
    document.addEventListener("keydown", addKeyboardShortcuts);
  });
  onDestroy(() => {
    // if (!!document) {
    // 	document.removeEventListener("keydown", addKeyboardShortcuts);
    // }
  });

  let searchQuery = "";

  let bookmarks = [];
  let loading = false;
  let deleteLoading = false;
  let mdSearchOpen = false;

  function closeSearchView() {
    searchView = false;
    searchQuery = "";
    bookmarks = [];
    nextID = "z";
    listEnd = false;
    loadBookMarks();
  }

  function doSearch() {
    mdSearchOpen = false;
    if (searchQuery === "") {
      closeSearchView();
      bookmarks = [];
      searchedText = "";
      nextID = "z";
      listEnd = false;
      loadBookMarks();
    } else {
      bookmarks = [];
      searchedText = searchQuery;
      nextID = "z";
      searchView = true;
      listEnd = false;
      loading = true;
      callSearch(searchQuery).then((res) => {
        bookmarks = res;
        loading = false;
      });
    }
  }

  function newDataIsHere(event) {
    bookmarks = [event.detail, ...bookmarks];
    setTimeout(() => {
      window.scrollTo({
        top: 0,
        behavior: "smooth"
      });
    }, 100);
  }
  function openNewBookmarkModal() {
    const eventAwesome = new CustomEvent("openNewBookmarkModal", {
      bubbles: true
    });
    document.dispatchEvent(eventAwesome);
  }
  function openSearchModal() {
    mdSearchOpen = true;
  }
  function onOpenChange(e) {
    mdSearchOpen = !mdSearchOpen;
  }
  let showQrCodeImage = false;
  let qrBookmark;
  function showQR(e) {
    qrBookmark = e.detail;
    showQrCodeImage = true;
  }
  function toggleShowQrCodeImage(e) {
    showQrCodeImage = !showQrCodeImage;
  }
  function nextPage() {
    pageIndex.push(nextID);
    pageIndex = [...new Set(pageIndex)];
    loadBookMarks("forward");
  }
  function previousPage() {
    //remove 1 elements from last from pageIndex and reassign
    pageIndex = pageIndex.slice(0, -1);

    loadBookMarks("backward");
  }
  function loadFirstPage() {
    nextID = "z";
    pageIndex = [];
    nextPage();
  }
  let selectAll = false;
  function checkUncheck() {
    bookmarks = bookmarks.map((bookmark) => {
      bookmark.checked = !selectAll;
      return bookmark;
    });
  }
  function doDelete() {
    deleteLoading = true;
    let bookmarkIds = bookmarks
      .filter((bookmark) => bookmark.checked)
      .map((bookmark) => bookmark.organization_relation_id);
    deleteBulkBookmarks(bookmarkIds).then(
      (res) => {
        bookmarks = bookmarks.filter((bookmark) => !bookmark.checked);
        deleteLoading = false;
        loadBookMarks("forward");
        selectAll = false;
      },
      function (err) {
        toast.error("Could not delete bookmarks");
        deleteLoading = false;
      }
    );
  }
  function checkUncheckBookmark(e) {
    let organization_relation_id = e.detail.organization_relation_id;
    let checked = e.detail.checked;
    bookmarks = bookmarks.map((bookmark) => {
      if (bookmark.organization_relation_id == organization_relation_id) {
        bookmark.checked = checked;
      }
      return bookmark;
    });
  }
  let showingFrom = 0;
  let showingTo = 0;
  let totalActiveBookmarks = 0;
  function setPageNumbers() {
    showingFrom = (pageIndex.length - 1) * limit + 1;
    showingTo = (pageIndex.length - 1) * limit + Math.min(limit, bookmarks.length);
    countBookmarks().then(
      (res) => {
        totalActiveBookmarks = res.count;
      },
      function (err) {
        totalActiveBookmarks = 0;
      }
    );
  }
</script>

<svelte:window on:newBookmarkAdded={newDataIsHere} on:openSearchModal={openSearchModal} />
<div class="container max-w-6xl p-4 pb-32">
  <JobQueue on:fetchNewBookmarks={loadFirstPage} />
  <div class="p-table grid grid-cols-1 overflow-hidden rounded-md border">
    <div class="p-row relative grid h-20 grid-cols-12 border-b bg-card pl-10 pr-4 pt-1.5">
      <div class="absolute left-3.5 top-7">
        <Checkbox
          bind:checked={selectAll}
          onCheckedChange={checkUncheck}
          class="border-foreground shadow-lg"
          id="selectAll"
        />
      </div>
      <div class="col-span-8 pt-[1px]">
        <div class="flex" use:autoAnimate>
          {#if bookmarks.filter((bookmark) => bookmark.checked).length > 0}
            <Button variant="ghost" size="icon" class="mt-3" on:click={doDelete} disabled={deleteLoading}>
              {#if deleteLoading}
                <LoaderCircle class="h-4 w-4 animate-spin" />
              {:else}
                <Trash2 class="h-4 w-4" />
              {/if}
            </Button>
          {/if}

          <form
            class="relative mt-4 hidden w-1/2 max-w-4xl justify-center rounded-md p-2 md:bottom-4 md:flex"
            on:submit|preventDefault={doSearch}
          >
            <Input
              type="text"
              name="q"
              class="searchh-input h-12 pr-8 focus-visible:ring-0 focus-visible:ring-secondary-foreground focus-visible:ring-opacity-50 focus-visible:ring-offset-0 focus-visible:ring-offset-transparent"
              bind:value={searchQuery}
              aria-label="Search for bookmarks"
              id="search-box"
              inputmode="search"
              autocapitalize="off"
              placeholder="Search for bookmarks / ⇧⌘F"
            />

            <Button
              size="icon"
              variant="ghost"
              type="submit"
              class="search-input-submit absolute  right-3 top-3 rounded-sm opacity-0 "
            >
              <ArrowRight class="h-4 w-4" />
            </Button>
          </form>
          {#if loading === true}
            <div class="mt-5">
              <LoaderCircle class="h-6 w-6 animate-spin" />
            </div>
          {/if}
        </div>
      </div>
      <div class="col-span-4 text-right">
        <div class="mt-6 flex justify-end">
          <Button
            variant="ghost"
            size="icon"
            class="h-6 w-6"
            disabled={pageIndex.length == 1 || loading}
            on:click={previousPage}
          >
            <ChevronLeft class="h-4 w-4" />
          </Button>
          <Button variant="ghost" class="h-6 text-xs" disabled={loading}>
            <span>
              {showingFrom} - {showingTo} of {totalActiveBookmarks.toLocaleString()}
            </span>
          </Button>
          <Button variant="ghost" size="icon" class="h-6 w-6" on:click={nextPage} disabled={listEnd || loading}>
            <ChevronRight class="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
    <div use:autoAnimate>
      {#each bookmarks as bookmark}
        <BookmarkTable {bookmark} on:checkUncheck={checkUncheckBookmark} on:showQR={showQR} />
      {/each}
    </div>
    {#if loading === false && bookmarks.length === 0 && listEnd === true}
      <div class="p-row relative grid grid-cols-12 border-b">
        <div class="relative col-span-12 mt-1 py-4">
          <p class="hidden text-center text-8xl text-secondary-foreground opacity-30 md:block">
            ⇧⌘U<br /> to add a bookmark
          </p>

          <div class="mx-auto grid grid-cols-2 gap-2 md:mt-8 md:max-w-2xl md:grid-cols-4">
            <a
              class="cursor-pointer rounded-md py-4 text-center hover:bg-card hover:shadow-sm"
              on:click={openNewBookmarkModal}
            >
              <Plus class="mx-auto mb-4 h-4 w-4" />
              New URL
            </a>
            <a class="cursor-pointer rounded-md py-4 text-center hover:bg-card hover:shadow-sm" href="{base}/github">
              <Github class="mx-auto mb-4 h-4 w-4" />
              Starred Repos
            </a>
            <a class="cursor-pointer rounded-md py-4 text-center hover:bg-card hover:shadow-sm" href="{base}/browsers">
              <Globe class="mx-auto mb-4 h-4 w-4" />
              Bookmarks
            </a>
            <a class="cursor-pointer rounded-md py-4 text-center hover:bg-card hover:shadow-sm" href="{base}/urls">
              <Link class="mx-auto mb-4 h-4 w-4" />
              Bulk URLs
            </a>
          </div>
        </div>
      </div>
    {/if}
  </div>
</div>

<Dialog.Root open={mdSearchOpen} {onOpenChange}>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Search for something</Dialog.Title>
      <Dialog.Description>
        <form class="relative my-4" on:submit|preventDefault={doSearch}>
          <Input
            type="search"
            class="h-12 rounded-sm pr-12"
            bind:value={searchQuery}
            aria-label="Search for bookmarks"
            placeholder="Search like a pro "
            autocapitalize="off"
            inputmode="search"
          />

          <Button size="icon" type="submit" class=" absolute  right-1 top-1 rounded-sm transition-all">
            <ArrowRight class="h-4 w-4" />
          </Button>
        </form>
      </Dialog.Description>
    </Dialog.Header>
  </Dialog.Content>
</Dialog.Root>

<Dialog.Root open={showQrCodeImage} onOpenChange={toggleShowQrCodeImage}>
  <Dialog.Content class="max-w-[305px] pt-3">
    <Dialog.Header>
      <Dialog.Description>
        <p class=" mb-2 text-center text-lg text-card-foreground">Scan QR code</p>
        <p class="text-md relative mb-4 text-ellipsis pl-8">
          <img src={qrBookmark.image_small} class="absolute left-0 top-1 inline-block h-5 w-5 rounded-sm" alt="" />
          {qrBookmark.title}
        </p>
        <div class="flex justify-center">
          <img src={qrBookmark.qrCodeImage} alt="" />
        </div>
      </Dialog.Description>
    </Dialog.Header>
  </Dialog.Content>
</Dialog.Root>
