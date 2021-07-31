<script>
	import { Styles } from 'sveltestrap';
	import Notifications from 'svelte-notifications';

	import Header from '$lib/header/Header.svelte';

	import '../app.css';
</script>

<Styles />

<Header />

<main>
	<Notifications>
		<slot />
	</Notifications>
</main>


<footer class="text-center text-lg-start bg-light text-muted mt-auto">
	<div class="text-center p-4" style="background-color: rgba(0, 0, 0, 0.05);">
		&copy; Recent Beater - Fork me on <a href="https://github.com/meyskens/recent-beater">GitHub</a> - <a href="https://github.com/sponsors/meyskens">Sponsor this site</a>
	</div>
</footer>
