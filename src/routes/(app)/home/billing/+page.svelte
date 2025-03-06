<script>
	import BrowserBack from "$lib/BrowserBack.svelte";
	import { getMySubscription } from "$lib/api";
	import { onMount } from "svelte";
	import { toast } from "svelte-sonner";
	import Loader from "$lib/loader.svelte";
	import { format, formatDistance, addDays } from "date-fns";
	import * as Alert from "$lib/components/ui/alert";
	import { Button } from "$lib/components/ui/button";

	let subscription = null;
	let loading = true;
	let isFreeTrialExpired = false;
	onMount(async () => {
		loading = true;
		subscription = await getMySubscription();
		loading = false;

		let freeTrialEndDate = new Date(subscription.created_at);
		//if more than 14 days since free trial started then isFreeTrialExpired = true
		let dateDiff = new Date() - freeTrialEndDate;
		if (dateDiff > 14 * 24 * 60 * 60 * 1000) {
			isFreeTrialExpired = true;
		}
	});
	let redirect = false;
	function buy() {
		redirect = true;
		location.href = "/api/billing/buy";
	}
</script>

<div class="container mx-auto w-full max-w-2xl">
	<BrowserBack />
	<h2 class="text-lg font-medium leading-10">Billing</h2>
	<p class="mb-4 text-sm text-secondary-foreground">Manage your access to OkBookmarks</p>
	{#if loading}
		<div class="my-4 flex justify-center">
			<Loader />
		</div>
	{/if}
	{#if !loading && !!subscription && subscription.current_state == "INACTIVE"}
		<div class="mb-4 rounded-lg bg-primary-foreground bg-opacity-10">
			<Alert.Root>
				<Alert.Title>You are yet to purchase OkBookmarks Access</Alert.Title>
				<Alert.Description class="pt-2">
					{#if isFreeTrialExpired}
						<p class="text-md">
							Your free access to okbookmarks has expired and it has been {formatDistance(
								new Date(),
								new Date(subscription.created_at)
							)} since then. Please consider purchasing lifetime access to OkBookmarks.
						</p>
					{:else}
						<p class="text-md">
							Get lifetime access to OkBookmarks by purchasing it for $14.99. You can
							keep on using OkBookmarks free till {format(
								addDays(new Date(subscription.created_at), 14),
								"do MMMM, yyyy"
							)}({formatDistance(
								addDays(new Date(subscription.created_at), 14),
								new Date()
							)} remaning)
						</p>
					{/if}
				</Alert.Description>
			</Alert.Root>
			<Button class="mt-4 w-full" disabled={redirect} on:click={buy}>
				Buy Lifetime Access for $14.99
			</Button>
		</div>
	{:else if !loading && !!subscription && subscription.current_state == "ACTIVE"}
		<div class="mb-4 rounded-lg bg-primary-foreground bg-opacity-10">
			<Alert.Root>
				<Alert.Title>You have lifetime access to OkBookmarks</Alert.Title>
				<Alert.Description class="pt-2">
					<p class="text-md">
						You have lifetime access to OkBookmarks. Enjoy the service and let us know
						if you have any feedback.
					</p>
					<p class="text-md">
						For queries/refund/receipt/feedback, please reach out to us at
						support@okbookmarks.com
					</p>
				</Alert.Description>
			</Alert.Root>
		</div>
	{/if}
</div>
