<script>
	import { LoaderCircle } from "lucide-svelte";
	import { getJobQueueStatus } from "$lib/api";
	import * as Alert from "$lib/components/ui/alert";
	import { onMount, onDestroy } from "svelte";
	import { createEventDispatcher } from "svelte";

	const dispatch = createEventDispatcher();

	let jobQueueStatus = {};
	let interval = null;
	let show = false;
	let remaining = 0;
	const intervalMs = 5000;
	let counter = 0;
	async function pollAPIProcess() {
		jobQueueStatus = await getJobQueueStatus();
		//count all keys for total
		remaining = jobQueueStatus["PENDING"] || 0;
	}
	//create a promise that will wait for 5 secs
	async function waitTill(ms) {
		return new Promise((resolve) => {
			setTimeout(() => {
				resolve();
			}, ms);
		});
	}

	function pollStatus() {
		interval = setInterval(async () => {
			await pollAPIProcess();
			if (remaining == 0) {
				stopPoller();
			} else {
				counter++;
				if (remaining > 0 && counter % 3 == 0) {
					dispatch("fetchNewBookmarks");
				}
			}
		}, intervalMs);
	}

	onMount(async () => {
		await pollAPIProcess();
		await waitTill(intervalMs);
		pollStatus();
	});

	function stopPoller() {
		if (counter > 0) {
			dispatch("fetchNewBookmarks");
		}

		clearInterval(interval);
	}

	onDestroy(() => {
		stopPoller();
	});
</script>

{#if remaining > 0}
	<div class="mb-4">
		<Alert.Root>
			<LoaderCircle class="h-4 w-4 animate-spin" />
			<Alert.Title>Indexing your bookmarks!</Alert.Title>
			<Alert.Description>
				We are indexing your bookmarks for search. {remaining} URLs remaining.
			</Alert.Description>
		</Alert.Root>
	</div>
{/if}
