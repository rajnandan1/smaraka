<script>
	import BrowserBack from "$lib/BrowserBack.svelte";
	import { onMount } from "svelte";
	import { getExtensionSecrets, createNewSecret, deactivateSecret } from "$lib/api";
	import * as Alert from "$lib/components/ui/alert";
	import { Button } from "$lib/components/ui/button";
	import { format, formatDistance } from "date-fns";
	import { LoaderCircle, ListChecks, X, Copy } from "lucide-svelte";
	import * as Dialog from "$lib/components/ui/dialog";
	import { toast } from "svelte-sonner";
	import Loader from "$lib/loader.svelte";
	import { Input } from "$lib/components/ui/input";
	import autoAnimate from "@formkit/auto-animate";
	let secrets = [];
	let loading = true;
	let secretName = "";
	let creatingSecret = false;

	let newSecretValue = "";
	onMount(async () => {
		loading = true;
		secrets = await getExtensionSecrets();
		loading = false;
	});
	function createSecret(e) {
		e.preventDefault;
		creatingSecret = true;
		createNewSecret(secretName).then(
			async (res) => {
				secrets = await getExtensionSecrets();
				creatingSecret = false;
				secretName = "";
				newSecretValue = res.secret_value;
			},
			function (err) {
				toast.error(err.message);
				creatingSecret = false;
			}
		);
	}
	function resetForm() {
		secretName = "";
		newSecretValue = "";
		creatingSecret = false;
	}
	function deactivate(i) {
		if (
			confirm(
				"Are you sure you want to deactivate this secret: " + secrets[i].secret_name + "?"
			)
		) {
			secrets[i].deactivating = true;
			deactivateSecret(secrets[i].id).then(
				async (res) => {
					secrets = await getExtensionSecrets();
				},
				function (err) {
					toast.error(err.message);
				}
			);
		}
	}
</script>

<div class="container mx-auto w-full max-w-2xl">
	<BrowserBack />
	<h1 class="text-lg font-medium leading-10">Keys and Secret management</h1>
	<p class="mb-4 text-sm text-secondary-foreground">Manage your keys and secrets here</p>

	<Dialog.Root onOpenChange={resetForm}>
		<Dialog.Trigger>
			<Button variant="outline" class="mb-4">Generate new secret</Button>
		</Dialog.Trigger>
		<Dialog.Content>
			<Dialog.Header>
				<Dialog.Title>Create a New Secret</Dialog.Title>
				<Dialog.Description>Generate a new secret to use in extension</Dialog.Description>
			</Dialog.Header>
			<form class="grid grid-cols-4 gap-2" on:submit={createSecret}>
				<div class="col-span-3">
					<Input
						bind:value={secretName}
						type="text"
						placeholder="Enter a name for your secret"
						class="max-w-xs"
						required
					/>
				</div>
				<div class="col-span-1">
					<Button class="mb-4 w-full" type="submit" disabled={creatingSecret}>
						{#if creatingSecret}
							<LoaderCircle class="h-6 w-6 animate-spin" />
						{:else}
							Create
						{/if}
					</Button>
				</div>
				<div class="col-span-4">
					{#if newSecretValue}
						<Alert.Root>
							<Alert.Title class="text-sm">
								Copy Paste the secret. It won't be shown again
							</Alert.Title>
							<Alert.Description>
								<pre
									class="overflow-auto rounded-md bg-yellow-200 px-1 py-2 text-gray-600">{newSecretValue}</pre>
								<Button
									variant="outline"
									class="mt-2 h-8 text-xs"
									on:click={() => {
										navigator.clipboard.writeText(newSecretValue);
										toast.success("Copied to clipboard");
									}}
								>
									Copy
								</Button>
							</Alert.Description>
						</Alert.Root>
					{/if}
				</div>
			</form>
		</Dialog.Content>
	</Dialog.Root>
	{#if loading}
		<div class="my-4 flex justify-center">
			<Loader />
		</div>
	{/if}
	{#if loading === false && secrets.length === 0}
		<div class="flex justify-center">
			<p class="text-secondary-foreground">No secrets found</p>
		</div>
	{:else}
		<div class="grid grid-cols-4 gap-2 text-sm" use:autoAnimate>
			<div class="col-span-1">
				<span class="font-semibold">Name</span>
			</div>
			<div class="col-span-1">
				<span class="font-semibold">Last Used</span>
			</div>
			<div class="col-span-1">
				<span class="font-semibold">Created at</span>
			</div>
			<div class="col-span-1"></div>
			{#each secrets as secret, i}
				<div class="col-span-4 border"></div>
				<div class="col-span-1 pt-[2px]">
					<span class="font-normal">{secret.secret_name}</span>
				</div>
				<div class="col-span-1 pt-[2px]">
					<span class="font-normal">
						{formatDistance(new Date(), new Date(secret.last_used_at))}
					</span>
				</div>
				<div class="col-span-1 pt-[2px]">
					<span class="font-senormalmibold"
						>{format(new Date(secret.created_at), "do MMM yyyy")}</span
					>
				</div>
				<div class="col-span-1 justify-end text-right">
					<Button
						variant="destructive"
						size="icon"
						disabled={!!secret.deactivating}
						on:click={() => deactivate(i)}
						class="float-end flex h-6 w-6 items-center"
					>
						<X class="h-4 w-4" />
					</Button>
				</div>
			{/each}
		</div>
	{/if}
</div>
