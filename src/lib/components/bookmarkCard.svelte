<script>
	import {
		ChevronRight,
		ArrowRight,
		QrCode,
		SquareArrowOutUpRight,
		X,
		Copy,
		LoaderCircle,
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

	export let bookmark;
	bookmark.visible = true;
	let showQr = false;
	let qrCodeImage = "";
	let qrStatus = "NONE";

	const handleImageError = (ev) => {
		//set parent background color
		ev.target.parentElement.style.cssText = `
			background-color: ${bookmark.accent_color};
			border-radius: 0.5rem;
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
		} catch (err) {
			toast.error("Could not show QR code");
		} finally {
			qrStatus = "LOADED";
		}
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
</script>

{#if bookmark.visible}
	<div class="col-span-1" id="bookmark-{bookmark.id}">
		<Card.Root class="bgg-card relative h-60 overflow-hidden bg-opacity-50">
			<Card.Header
				class="cover-image relative h-28"
				style="background-image: linear-gradient(0deg, rgba(0, 0, 0, .15) 0%, transparent 100%),url({bookmark.image_large}); background-color: {bookmark.accent_color};"
			>
				<div class="absolute left-3 top-3 h-6 w-6 rounded-s">
					{#if !!bookmark.image_small}
						<img
							src={bookmark.image_small}
							class=" h-6 w-6"
							alt=""
							on:error={handleImageError}
						/>
					{:else}
						<div
							class=" h-6 w-6 rounded-sm shadow-sm"
							style="background-color: {bookmark.accent_color}"
						></div>
					{/if}
				</div>
			</Card.Header>

			<Card.Content class="p-0">
				<p class="h-14 overflow-hidden text-ellipsis p-4 text-sm font-medium">
					{bookmark.title}
				</p>
			</Card.Content>
			<Card.Footer class="mt-2 justify-between p-4 pr-2">
				<div class="col-span-1">
					<Button
						class="hidden h-6 w-6 pl-1 text-center md:inline-block"
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
					<Button class="h-6 w-6 " variant="ghost" size="icon" on:click={copyURL}>
						<Copy class="h-4 w-4 text-secondary-foreground" />
					</Button>
					<DropdownMenu.Root>
						<DropdownMenu.Trigger
							><Button class="h-6 w-6" variant="ghost" size="icon">
								<Settings class="h-4 w-4 text-secondary-foreground" />
							</Button></DropdownMenu.Trigger
						>
						<DropdownMenu.Content>
							<DropdownMenu.Group>
								<DropdownMenu.Item
									class="cursor-pointer"
									href="{base}/bookmark#{bookmark.id}"
								>
									Configure
								</DropdownMenu.Item>
								<DropdownMenu.Separator />
								<DropdownMenu.Item class="cursor-pointer" on:click={doDelete}
									>Delete</DropdownMenu.Item
								>
							</DropdownMenu.Group>
						</DropdownMenu.Content>
					</DropdownMenu.Root>
				</div>
				<div class="col-span-1">
					<Button
						href={bookmark.url}
						target="_blank"
						variant="outline"
						size="sm"
						class="text-xs font-medium uppercase text-secondary-foreground"
					>
						Open
						<SquareArrowOutUpRight class="ml-2 h-4 w-4" />
					</Button>
				</div>
			</Card.Footer>
		</Card.Root>
	</div>
{/if}
