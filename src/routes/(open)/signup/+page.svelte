<script lang="ts">
  import { Input } from "$lib/components/ui/input/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import * as Card from "$lib/components/ui/card";
  import { Button } from "$lib/components/ui/button/index.js";
  import { base } from "$app/paths";
  import { LoaderCircle } from "lucide-svelte";
  import { onMount } from "svelte";

  let email = "";
  let password = "";
  let name = "";
  let error = "";
  let loading = false;

  let errorMsg = "";

  onMount(() => {
    const urlParams = new URLSearchParams(window.location.search);
    errorMsg = urlParams.get("error");
  });
</script>

<div class="relative grid h-screen grid-cols-2 bg-background">
  <div class="absolute left-0 top-8 flex w-full justify-between px-8">
    <div>
      <a
        class="flex-none text-xl font-semibold focus:opacity-80 focus:outline-none dark:text-white sm:order-1"
        href={base}
      >
        <img src="{base}/smaraka.png" class="inline h-12" alt="" />
        <span class="inline-block">Smaraka</span>
      </a>
    </div>

    <div>
      <Button href="{base}/login" variant="secondary" class="font-normal text-white">Login</Button>
    </div>
  </div>
  <div
    class="hidden h-screen bg-card md:flex"
    style="background-image: url({base}/road.svg);background-repeat: no-repeat;background-size: 80%; background-position: bottom;"
  ></div>
  <div class="col-span-2 px-4 md:col-span-1">
    <Card.Root class="mx-auto mt-32 max-w-md border-none bg-transparent md:mt-48">
      <Card.Header class="text-center">
        <Card.Title>Create an account for Smaraka</Card.Title>
        <Card.Description>Welcome! Create an account to save your bookmarks.</Card.Description>
      </Card.Header>
      <Card.Content>
        {#if !!errorMsg}
          <p class="my-4 rounded-md border border-destructive p-2 text-xs font-medium text-destructive">
            Error: {errorMsg}
          </p>
        {/if}
        <form method="post" action="/api/ui/user/sign-up">
          <div class="flex w-full max-w-sm flex-col gap-1.5">
            <Label for="name-2">Your Name</Label>
            <Input name="name" bind:value={name} type="text" id="name-2" placeholder="Raj Nandan Sharma" required />
            <p class="text-xs text-muted-foreground">What should we call you?</p>
          </div>
          <div class="mt-4 flex w-full max-w-sm flex-col gap-1.5">
            <Label for="email-2">Email</Label>
            <Input name="email" bind:value={email} type="email" id="email-2" placeholder="raj@example.com" required />
            <p class="text-xs text-muted-foreground">What is your email address?</p>
          </div>
          <div class="mt-4 flex w-full max-w-sm flex-col gap-1.5">
            <Label for="pass-2">Password</Label>
            <Input name="password" bind:value={password} type="password" id="pass-2" placeholder="Password" required />
            <p class="text-xs text-muted-foreground">Create your password.</p>
          </div>
          <div class="mt-4 flex w-full max-w-sm flex-col justify-end gap-1.5">
            {#if !!error}
              <p class="text-left text-xs font-medium text-red-500">{error}</p>
            {/if}
            <Button type="submit" disabled={loading}>
              Create Account
              {#if loading}
                <LoaderCircle class="ml-2 h-4 w-4 animate-spin" />
              {/if}
            </Button>
          </div>
        </form>
      </Card.Content>
    </Card.Root>
  </div>
</div>
