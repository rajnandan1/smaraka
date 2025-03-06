<script>
  import {
    ChevronRight,
    ArrowRight,
    QrCode,
    SquareArrowOutUpRight,
    LoaderCircle,
    X,
    Copy,
    Settings
  } from "lucide-svelte";
  import * as Card from "$lib/components/ui/card";
  import { Button } from "$lib/components/ui/button";
  import QRCode from "qrcode";
  import { toast } from "svelte-sonner";
  import { slide } from "svelte/transition";
  import { onMount } from "svelte";
  import * as Popover from "$lib/components/ui/popover";
  import { base } from "$app/paths";
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu";
  import { deleteBookmarkById } from "$lib/api";
  import { Checkbox } from "$lib/components/ui/checkbox";
  import { createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();

  export let bookmark;
  bookmark.visible = true;
  let fallbackImage = "https://placehold.co/40";
  let showQr = false;
  let qrCodeImage = "";
  let qrStatus = "NONE";

  const handleImageError = (ev) => {
    //set parent background color
    if (!!ev.target.parentElement.style.backgroundColor) {
      return;
    }
    ev.target.parentElement.style.cssText = `
			background-color: ${bookmark.accent_color};
			box-shadow: rgba(149, 157, 165, 0.2) 0px 8px 24px;
		`;
    ev.target.remove();
  };
  async function generateQR() {
    qrStatus = "LOADING";
    try {
      qrCodeImage = await QRCode.toDataURL(bookmark.url, {
        errorCorrectionLevel: "H",
        width: 220,
        margin: 1
      });
      bookmark.qrCodeImage = qrCodeImage;
      const eventAwesome = new CustomEvent("showQR", {
        bubbles: true,
        detail: bookmark
      });
      document.dispatchEvent(eventAwesome);
      dispatch("showQR", bookmark);
    } catch (err) {
      toast.error("Could not show QR code");
    } finally {
      qrStatus = "LOADED";
    }
  }
  function fadeSlide(node, options) {
    const slideTrans = slide(node, options);
    return {
      duration: options.duration,
      css: (t) => `
				${slideTrans.css(t)}
			`
    };
  }

  onMount(() => {});

  function copyURL() {
    navigator.clipboard.writeText(bookmark.url);
    toast.success("URL copied to clipboard");
  }

  function doDelete() {
    if (!confirm("Are you sure you want to delete this bookmark?")) {
      return;
    }

    deleteBookmarkById(bookmark.id).then(
      (res) => {
        toast.success("Bookmark deleted");
        bookmark.visible = false;
      },
      function (err) {
        toast.error("Could not delete bookmark");
      }
    );
  }
  function checkUncheck() {
    dispatch("checkUncheck", {
      organization_relation_id: bookmark.organization_relation_id,
      checked: !bookmark.checked
    });
  }
</script>

{#if bookmark.organization_url_status == "ACTIVE"}
  <div class=" p-row relative grid grid-cols-12 border-b py-2 pl-24 pr-4 hover:bg-secondary">
    <div class="absolute left-3.5 top-2.5">
      <Checkbox
        class="border-foreground  shadow-lg"
        bind:checked={bookmark.checked}
        onCheckedChange={checkUncheck}
        id={bookmark.id}
      />
    </div>
    <div class="absolute left-14 top-2.5">
      <div class="h-5 w-5 overflow-hidden rounded-sm shadow-sm">
        {#if !!bookmark.image_small}
          <img src={bookmark.image_small} class=" h-5 w-5" on:error={handleImageError} />
        {:else}
          <div class="h-5 w-5" style="background-color: {bookmark.accent_color}"></div>
        {/if}
      </div>
    </div>
    <div class="col-span-12 pt-[1px] md:col-span-10">
      <p class=" h-6 overflow-hidden text-ellipsis whitespace-nowrap">
        {bookmark.title}
        <span class="pl-2 text-sm text-secondary-foreground">{bookmark.excerpt}</span>
      </p>
    </div>
    <div class="col-span-12 -ml-4 mt-4 pl-2 pt-[1px] md:col-span-2 md:ml-0 md:mt-0 md:justify-end md:text-right">
      <Button class="h-6 w-6" href={bookmark.url} target="_blank" variant="ghost" size="icon">
        <SquareArrowOutUpRight class="h-4 w-4 text-secondary-foreground" />
      </Button>
      <Button
        class="h-0 w-0 md:h-6 md:w-6 "
        variant="ghost"
        size="icon"
        disabled={qrStatus == "LOADING"}
        on:click={generateQR}
      >
        {#if qrStatus == "LOADING"}
          <LoaderCircle class="h-4 w-4 animate-spin text-secondary-foreground" />
        {:else}
          <QrCode class="h-4 w-4 text-secondary-foreground" />
        {/if}
      </Button>
      <Button class="h-6 w-6" variant="ghost" size="icon" on:click={copyURL}>
        <Copy class="h-4 w-4 text-secondary-foreground" />
      </Button>
      <Button
        class="h-6 w-6"
        variant="ghost"
        size="icon"
        href="{base}/home/bookmark#{bookmark.organization_relation_id}"
      >
        <Settings class="h-4 w-4 text-secondary-foreground" />
      </Button>
    </div>
  </div>
{/if}
