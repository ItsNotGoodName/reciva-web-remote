<script>
import { mapActions } from 'vuex';
import Player from './components/Player.vue'

export default {
	name: "App",
	components: {
		Player
	},
	mounted() {
		this.loadAll()
	},
	computed: {
		state() {
			return this.$store.state
		},
	},
	methods: {
		...mapActions([
			'loadAll'
		])
	}
}
</script>

<template >
	<div class="container mx-auto px-2">
		<Player class="sticky top-0 mx-auto border-b-2" />
		<div>Lorem</div>
		<div>Lorem</div>
		<div>Lorem</div>
		<div>Lorem</div>
		<div>Lorem</div>
	</div>
</template>

<style>
</style>
