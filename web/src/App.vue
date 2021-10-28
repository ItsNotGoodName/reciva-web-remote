<script>
import { mapActions } from 'vuex';
import Player from './components/Player.vue'
import Presets from './components/Presets.vue'

export default {
	name: "App",
	components: {
		Player,
		Presets
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
		<Player
			class="sticky top-0 mx-auto border-l-2 border-b-2 border-r-2 rounded-b-xl p-2 max-w-3xl mb-1"
		/>
		<Presets />
	</div>
</template>

<style>
</style>
