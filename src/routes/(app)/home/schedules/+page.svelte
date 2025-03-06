<script>
  import BrowserBack from "$lib/BrowserBack.svelte";
  import { onMount } from "svelte";
  import { updateSchedule, getSchedules } from "$lib/api";
  import * as Alert from "$lib/components/ui/alert";
  import { Button } from "$lib/components/ui/button";
  import { Switch } from "$lib/components/ui/switch";

  import { format, formatDistance } from "date-fns";
  import { base } from "$app/paths";
  import { LoaderCircle, ListChecks, X, Copy } from "lucide-svelte";
  import * as Dialog from "$lib/components/ui/dialog";
  import { toast } from "svelte-sonner";
  import Loader from "$lib/loader.svelte";
  import { Input } from "$lib/components/ui/input";
  import autoAnimate from "@formkit/auto-animate";
  import * as Card from "$lib/components/ui/card";
  let loading = true;
  let schedules = [];

  onMount(async () => {
    loading = true;
    schedules = await getSchedules();
    loading = false;
    schedules = schedules.map((sc) => {
      sc.isOn = sc.status === "active";
      return sc;
    });
  });

  function updateScheduleStatus(e, sc) {
    sc.isOn = e;
    updateSchedule(sc.schedule_id, sc.isOn ? "active" : "inactive", sc.interval);
  }
</script>

<div class="container mx-auto w-full max-w-2xl">
  <BrowserBack />
  <h1 class="text-lg font-medium leading-10">Manage Schedules</h1>
  <p class="mb-4 text-sm text-secondary-foreground">
    Schedules fetches urls from a given page and stores it for you at a defined interval. You can enable or disable the
    schedule as per your need.
  </p>

  {#if loading}
    <div class="my-4 flex justify-center">
      <Loader />
    </div>
  {:else}
    <div class="mt-4">
      <Card.Root>
        <Card.Content class="pb-0">
          <div class="">
            {#each schedules as sc, i}
              <div class="gh-repo relative {i < schedules.length - 1 ? 'border-b' : ''} py-4 pr-16">
                <p class="overflow-hidden text-ellipsis text-nowrap text-foreground">
                  {sc.schedule_name}
                  <span class="mx-1 rounded-sm border border-primary px-1 text-xs"
                    >{sc.interval} {sc.interval > 1 ? "days" : "day"}</span
                  >
                </p>
                <p
                  class="overflow-hidden text-ellipsis text-nowrap text-sm text-secondary-foreground"
                  title={sc.schedule_description}
                >
                  {sc.schedule_description}
                </p>
                <p class="overflow-hidden text-ellipsis text-nowrap text-sm text-secondary-foreground"></p>
                <div class="absolute right-2 top-6">
                  <Switch
                    bind:checked={sc.isOn}
                    onCheckedChange={(e) => {
                      updateScheduleStatus(e, sc);
                    }}
                  />
                </div>
              </div>
            {/each}
          </div>
        </Card.Content>
      </Card.Root>
    </div>
  {/if}
</div>
