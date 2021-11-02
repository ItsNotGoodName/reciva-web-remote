<script>
import Alerts from './components/Alerts.vue'
import Player from './components/Player.vue'
import Presets from './components/Presets.vue'

export default {
	name: "App",
	components: {
		Alerts,
		Player,
		Presets,
	},
	mounted() {
		this.$store.dispatch("init")
			.catch((err) => this.$store.dispatch("addNotification", { "type": "error", "message": err }))
	},
}
</script>

<template >
	<div class="container mx-auto px-1 h-screen">
		<Player class="sticky top-0 mx-auto border-l-2 border-b-2 border-r-2 rounded-b-xl p-2 max-w-3xl" />
		<Presets class="my-1" />
		<Alerts />
	</div>
</template>

<style>
</style>
