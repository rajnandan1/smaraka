<script>
	import { page } from "$app/stores";
	import { onMount } from "svelte";
	import { getBookmarkById, updateBookmark, refetchBookmark, deleteBookmarkById } from "$lib/api";
	import { toast } from "svelte-sonner";
	import Loader from "$lib/loader.svelte";
	import { ServerCrash, ExternalLink, Trash, ListRestart, LoaderCircle } from "lucide-svelte";
	import * as Card from "$lib/components/ui/card";
	import { Input } from "$lib/components/ui/input";
	import { Label } from "$lib/components/ui/label";
	import { Textarea } from "$lib/components/ui/textarea";
	import { Button } from "$lib/components/ui/button";
	import BrowserBack from "$lib/BrowserBack.svelte";
	let bookmark = {};
	let pageState = "";
	let interactDisabled = false;
	function loadBookMark() {
		pageState = "LOADING";

		//get hash param
		const hash = window.location.hash;
		const bookmarkID = hash.substring(1);
		getBookmarkById(bookmarkID).then(
			(res) => {
				// res.status = "PENDING";
				bookmark = res;
				if (bookmark.status == "COMPLETE") {
					interactDisabled = false;
				} else {
					interactDisabled = true;
				}
				pageState = "LOADED";
			},
			function (err) {
				toast.error("Bookmark not found");
				pageState = "ERROR";
			}
		);
	}

	function fetchBookmarkURL() {
		interactDisabled = true;
		bookmark.status = "PENDING";
		refetchBookmark(bookmark.id).then(
			(res) => {
				bookmark = res;
				if (bookmark.status == "COMPLETE") {
					interactDisabled = false;
				} else {
					interactDisabled = true;
				}
			},
			function (err) {
				toast.error("Could not fetch bookmark");
				interactDisabled = false;
			}
		);
	}

	function deleteBookmark() {
		if (!confirm("Are you sure you want to delete this bookmark?")) {
			return;
		}
		interactDisabled = true;
		deleteBookmarkById(bookmark.id).then(
			(res) => {
				toast.success("Bookmark deleted");
				interactDisabled = false;
				location.href = "/app";
			},
			function (err) {
				toast.error("Could not delete bookmark");
				interactDisabled = false;
			}
		);
	}

	onMount(() => {
		loadBookMark();
	});
</script>

<div class="container max-w-6xl p-4 pb-32">
	<BrowserBack />
	{#if pageState === "LOADING"}
		<div class="absolute left-1/2 mt-10 -translate-x-1/2"><Loader /></div>
	{:else if pageState === "ERROR"}
		<div class="mt-10 text-center">
			<h1 class="text-center text-2xl font-bold">
				<ServerCrash class="mx-auto h-12 w-12" />
			</h1>
			<h1 class="mt-6 text-2xl font-bold">Bookmark not found</h1>
		</div>
	{:else if pageState === "LOADED"}
		<div class="grid grid-cols-1">
			<Card.Root>
				<Card.Header>
					<Card.Title class="relative">
						<img src={bookmark.image_small} alt="" class="left-2 block h-6 w-6" />
						<p class="mt-4">
							<Button
								variant="ghost"
								href={bookmark.url}
								target="_blank"
								class=" max-w-full justify-start   overflow-hidden text-ellipsis whitespace-nowrap px-0 text-xl hover:bg-transparent"
							>
								{bookmark.url}
								<ExternalLink class="ml-2 h-4 w-4" />
							</Button>
						</p>

						<Button
							variant="outline"
							on:click={fetchBookmarkURL}
							class="absolute right-2 top-0"
							disabled={interactDisabled}
						>
							{#if bookmark.status == "COMPLETE"}
								<ListRestart class="h-4 w-4" />
							{:else}
								<LoaderCircle class="h-4 w-4 animate-spin" />
							{/if}
						</Button>
					</Card.Title>
				</Card.Header>
				<Card.Content class="">
					<div class="flex w-full flex-col gap-1.5">
						<Label for="bk_title">Title</Label>

						<p class="text-sm text-secondary-foreground">
							{bookmark.title}
						</p>
					</div>
					<div class="mt-8 flex w-full flex-col gap-1.5">
						<Label for="bk_excerpt">Excerpt</Label>
						<p class="text-sm text-secondary-foreground">
							{bookmark.excerpt}
						</p>
					</div>
					<div class="mt-8 flex w-full flex-col gap-1.5">
						<Label for="bk_text">Text</Label>
						<pre class="overflow-x-auto text-sm text-secondary-foreground">
							{bookmark.full_text}</pre>
					</div>
				</Card.Content>
			</Card.Root>
		</div>
	{/if}
</div>
