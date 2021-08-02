<script>
	import { getNotificationsContext } from 'svelte-notifications';
	const { addNotification } = getNotificationsContext();

	export const prerender = true;

	let loading = false;
	let url = "";
	let pages = 1;
	let mode = "recent"
	let filterNoPP = false

	async function checkURL() {
		const res = await fetch(`https://recentbeat.com/check/?url=${encodeURIComponent(url)}`);

		if (!res.ok) {
			return false;
		}

		const data = await res.json();

		return data.status == "ok"
	}

	async function handleBPList(e, oneclick = false) {
		loading = true;
		let ok = await checkURL();
		if (!ok) {
			addNotification({
				text: 'Invalid URL given',
				position: 'top-center',
			})
			loading = false;
			return
		}

		window.location = `/playlist/?url=${encodeURIComponent(url)}&mode=${mode}&pages=${pages}${filterNoPP ? "&filterNoPP=true" : ""}${oneclick ? "&oneclick=true" : ""}`;
		loading = false;
	}

	async function handleOneCLick(e) {
		return handleBPList(e,true)
	}
</script>

<svelte:head>
	<title>Recent Beater</title>
</svelte:head>

<style>
	.about {
		color: #fff;
	}
</style>

<section>
	<div class="row justify-content-center my-4">
		<div class="col-6 col-sm-12 about text-center">
			Recent Beater will give you a a Beat Saber playlist of the recent plays of a given ScoreSaber profile.
		</div>
	</div>
	<div class="row justify-content-center my-4">
		<div class="col-12 col-md-6">
			<input class="form-control form-control-lg" type="text" placeholder="https://scoresaber.com/u/765611979825216..." aria-label=".form-control-lg" bind:value={url}>
		</div>
	</div>
	<div class="row justify-content-center my-4">
		<div class="col-6 col-md-2">
			<label for="type">Type</label>
			<select class="form-select" id="type" aria-label="select type" bind:value={mode}>
				<option value="recent" selected>Recent Scores</option>
				<option value="top">Top Scores</option>
			</select>
		</div>
		<div class="col-4 col-md-1">
			<label for="type">Pages</label>
			<select class="form-select" aria-label="select type" bind:value={pages} >
				<option value="1" selected>1</option>
				{#each Array(9) as _, row}
					<option value={row+2}>{row+2}</option>
				{/each}
			</select>
		</div>
		<div class="col-2 col-md-1">
			<label>&nbsp;</label>
			<div class="form-check form-switch">
				<input class="form-check-input" type="checkbox" id="filterNoPP" bind:checked={filterNoPP}>
				<label class="form-check-label" for="filterNoPP">Filter songs with no PP</label>
			</div>
		</div>
	</div>

	<div class="row justify-content-center my-4">
		<div class="col-6 text-center">
			<button on:click={handleBPList} disabled={loading} type="button" class="btn btn-lg btn-primary mx-1"><i class="bi bi-download"></i> .bplist</button>
			<button on:click={handleOneCLick} disabled={loading}  type="button" class="btn btn-lg btn-danger mx-1"><i class="bi bi-cloud-download-fill"></i> OneClick</button>
		</div>
	</div>
</section>
