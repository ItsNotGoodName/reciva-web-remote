<script>
import { mapActions } from 'vuex';

export default {
	name: "App",
	components: { },
	mounted() {
		this.loadAll()
	},
	computed: {
		state() {
			return this.$store.state
		}
	},
	methods: {
		...mapActions([
			'loadAll'
		])
	}
}
</script>

<template>
	<div class="container mx-auto px-2">
		<pre><code>{{ JSON.stringify(state, null, 2) }}</code></pre>
	</div>
</template>

<style>
</style>
