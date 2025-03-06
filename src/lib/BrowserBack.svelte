<script>
	import { ArrowLeft } from "lucide-svelte";
	import { Button } from "$lib/components/ui/button";
	function goBack() {
		window.history.back();
	}
</script>

<Button variant="ghost" class="  -ml-5  hover:bg-transparent" on:click={goBack}>
	<ArrowLeft class="h-5 w-5" />
</Button>
