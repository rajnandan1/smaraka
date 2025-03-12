<script>
  import BrowserBack from "$lib/BrowserBack.svelte";
  import { onMount } from "svelte";
  import { updateSchedule, getSchedules, createSchedule, deleteSchedules, playSchedules } from "$lib/api";
  import * as Alert from "$lib/components/ui/alert";
  import { Button } from "$lib/components/ui/button";
  import { Switch } from "$lib/components/ui/switch";
  import * as Select from "$lib/components/ui/select";
  import { Label } from "$lib/components/ui/label";
  import { Checkbox } from "$lib/components/ui/checkbox";

  import { format, formatDistance } from "date-fns";
  import { base } from "$app/paths";
  import { LoaderCircle, ListChecks, X, Copy, Trash, Play } from "lucide-svelte";
  import * as Dialog from "$lib/components/ui/dialog";
  import { toast } from "svelte-sonner";
  import Loader from "$lib/loader.svelte";
  import { Input } from "$lib/components/ui/input";
  import autoAnimate from "@formkit/auto-animate";
  import * as Card from "$lib/components/ui/card";
  let loading = false;
  let schedules = [];
  let showCreateSchedule = false;

  async function fetchData(data) {
    if (!!!data) {
      loading = true;
      data = await getSchedules();
    }

    if (!!data) {
      schedules = data;
      loading = false;
      schedules = schedules.map((sc) => {
        sc.isOn = sc.schedule_status === "ACTIVE";
        sc.isSelected = false;
        return sc;
      });
    }
  }

  onMount(async () => {
    fetchData();
  });

  function updateScheduleStatus(e, sc) {
    sc.isOn = e;
    updateSchedule(sc.schedule_id, sc.isOn ? "ACTIVE" : "INACTIVE");
  }
  function onOpenChange(e) {
    showCreateSchedule = !showCreateSchedule;
  }

  let scheduleTypes = [
    { value: "GH_TRENDING", label: "Github Trending", description: "Get trending repositories from Github" },
    {
      value: "PH_LEADERBOARD",
      label: "Product Hunt Leaderboard",
      description: "Get trending products from Product Hunt"
    },
    {
      value: "GH_STARRED_REPO",
      label: "Github Starred Repos",
      description: "Get your starred repositories from Github"
    },
    { value: "HN_TRENDING", label: "Hacker New Trending", description: "Get trending news from Hacker News" }
  ];

  let schedulesIntervals = [
    { value: 1, label: "Every day" },
    { value: 7, label: "Every 7 days" },
    { value: 30, label: "Every 30 days" }
  ];

  let selectedType;
  let selectedInterval;
  let scheduleURL = "";
  let ghUsername = "";

  function scheduleTypeChange(e) {
    selectedType = e;
    prepareSchedule();
  }
  function scheduleIntervalChange(e) {
    selectedInterval = e;

    prepareSchedule();
  }

  function prepareSchedule() {
    scheduleURL = "";
    if (!!!selectedType || !!!selectedInterval) {
      return;
    }

    if (selectedType.value === "GH_TRENDING") {
      scheduleURL = "https://github.com/trending";
      if (selectedInterval.value === 1) {
        scheduleURL += "?since=daily";
      } else if (selectedInterval.value === 7) {
        scheduleURL += "?since=weekly";
      } else if (selectedInterval.value === 30) {
        scheduleURL += "?since=monthly";
      }
    } else if (selectedType.value === "PH_LEADERBOARD") {
      scheduleURL = "https://www.producthunt.com/leaderboard";
      if (selectedInterval.value === 1) {
        scheduleURL += "/daily";
      } else if (selectedInterval.value === 7) {
        scheduleURL += "/weekly";
      } else if (selectedInterval.value === 30) {
        scheduleURL += "/monthly";
      }
    } else if (selectedType.value === "GH_STARRED_REPO") {
      scheduleURL = `https://github.com/stars/${ghUsername}/repositories`;
    } else if (selectedType.value === "HN_TRENDING") {
      scheduleURL = `https://news.ycombinator.com/best?p=1`;
    }
  }

  let creatingSchedule = false;
  function createNewSchedule() {
    creatingSchedule = true;
    let createObj = {
      schedule_name: selectedType.label,
      schedule_description: scheduleTypes.find((st) => st.value === selectedType.value).description,
      interval_days: selectedInterval.value,
      schedule_url: scheduleURL,
      schedule_type: selectedType.value
    };

    createSchedule(createObj).then(
      (res) => {
        creatingSchedule = false;
        toast.success("Schedule created successfully");
        showCreateSchedule = false;
        fetchData(res);
      },
      (err) => {
        creatingSchedule = false;
        toast.error("Failed to create schedule. Check if the schedule already exists.");
      }
    );
  }
  let deletingSchedule = false;
  function doDeleteSchedules(e) {
    deletingSchedule = true;
    let selectedSchedules = schedules.filter((sc) => sc.isSelected);
    let scheduleIds = selectedSchedules.map((sc) => sc.schedule_id);
    deleteSchedules(scheduleIds).then(
      (res) => {
        toast.success("Schedules deleted successfully");
        fetchData(res);
        deletingSchedule = false;
      },
      (err) => {
        toast.error("Failed to delete schedules");
        deletingSchedule = false;
      }
    );
  }
  let playingSchedule = false;
  function doPlaySchedules(e) {
    playingSchedule = true;
    let selectedSchedules = schedules.filter((sc) => sc.isSelected);
    let scheduleIds = selectedSchedules.map((sc) => sc.schedule_id);
    playSchedules(scheduleIds).then(
      (res) => {
        toast.success("Schedules added successfully");
        playingSchedule = false;
      },
      (err) => {
        toast.error("Failed to delete schedules");
        playingSchedule = false;
      }
    );
  }
</script>

<div class="container mx-auto w-full max-w-2xl">
  <BrowserBack />
  <h1 class="text-lg font-medium leading-10">Manage Schedules</h1>
  <p class="mb-4 text-sm text-secondary-foreground">
    Schedules fetches urls from a given page at given interval and stores it for you so that you can search them later.
    You can enable or disable the schedule as per your need.
  </p>
  <div class="mt-4 flex justify-end gap-x-2">
    {#if schedules.filter((sc) => sc.isSelected).length > 0}
      <div class="flex w-2/3 justify-start gap-x-2">
        <Button class="bg-green-500 hover:bg-green-600" on:click={doPlaySchedules}>
          <Play class="mr-2 h-4 w-4" />
          Run
          {#if playingSchedule}
            <LoaderCircle class="ml-2 h-4 w-4 animate-spin" />
          {/if}
        </Button>
        <Button variant="ghost" class="text-destructive" on:click={doDeleteSchedules}>
          <Trash class="mr-2 h-4 w-4" />
          Delete
          {#if deletingSchedule}
            <LoaderCircle class="ml-2 h-4 w-4 animate-spin" />
          {/if}
        </Button>
      </div>
    {/if}
    <div class="flex w-1/3 justify-end">
      <Button
        variant="outline"
        on:click={(e) => {
          showCreateSchedule = true;
        }}>Add new schedule</Button
      >
    </div>
  </div>
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
              <div class="gh-repo relative {i < schedules.length - 1 ? 'border-b' : ''} relative py-4 pl-7 pr-16">
                <Checkbox
                  bind:checked={sc.isSelected}
                  onCheckedChange={(e) => {
                    sc.isSelected = e;
                  }}
                  class="absolute left-0 top-5"
                />
                <p class="overflow-hidden text-ellipsis text-nowrap text-foreground">
                  {sc.schedule_name}
                  <span class="mx-1 rounded-sm border border-primary px-1 text-xs">
                    {sc.interval_days}
                    {sc.interval_days > 1 ? "days" : "day"}
                  </span>
                </p>
                <div class="flex flex-col">
                  <p
                    class="overflow-hidden text-ellipsis text-nowrap text-sm text-muted-foreground"
                    title={sc.schedule_description}
                  >
                    {sc.schedule_description}
                  </p>
                  <p
                    class="overflow-hidden text-ellipsis text-nowrap text-sm text-secondary-foreground"
                    title={sc.schedule_url}
                  >
                    {sc.schedule_url}
                  </p>
                </div>

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
<Dialog.Root open={showCreateSchedule} {onOpenChange}>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Create a New Schedule</Dialog.Title>
      <Dialog.Description>
        <div class="mt-4 flex flex-col gap-y-4">
          <div>
            <Label for="sch_type">Schedule Type</Label>

            <Select.Root selected={selectedType} onSelectedChange={scheduleTypeChange}>
              <Select.Trigger id="sch_type" class="mt-2 w-[280px]">
                <Select.Value placeholder="Select Type" />
              </Select.Trigger>
              <Select.Content>
                {#each scheduleTypes as type}
                  <Select.Item value={type.value} label={type.label}>{type.label}</Select.Item>
                {/each}
              </Select.Content>
            </Select.Root>
          </div>
          <div>
            <Label for="sch_interval">Schedule Interval</Label>

            <Select.Root selected={selectedInterval} onSelectedChange={scheduleIntervalChange}>
              <Select.Trigger id="sch_interval" class="mt-2 w-[180px]">
                <Select.Value placeholder="Select Interval" />
              </Select.Trigger>
              <Select.Content>
                {#each schedulesIntervals as type}
                  <Select.Item value={type.value} label={type.label}>{type.label}</Select.Item>
                {/each}
              </Select.Content>
            </Select.Root>
          </div>
          {#if selectedType && selectedType.value === "GH_STARRED_REPO"}
            <div>
              <Label for="ghusername">Github Username</Label>
              <Input
                class="mt-2 w-[280px]"
                type="text"
                id="ghusername"
                bind:value={ghUsername}
                placeholder="github username"
              />
            </div>
          {/if}
          <div class="">
            <Button disabled={!scheduleURL || creatingSchedule} on:click={createNewSchedule}>
              Save Schedule
              {#if creatingSchedule}
                <LoaderCircle class="ml-2 h-4 w-4 animate-spin" />
              {/if}
            </Button>
          </div>
        </div>
      </Dialog.Description>
    </Dialog.Header>
  </Dialog.Content>
</Dialog.Root>
