<template >
	<div
		v-if="radioUUID && radioConnected"
		class="gap-1 grid lg:grid-cols-4 md:grid-cols-3 sm:grid-cols-2"
	>
		<loading-button
			:on-click="() => playRadioPreset(preset.number)"
			class="hover:bg-blue-300 hover:text-white rounded p-1"
			v-bind:class="{ 'bg-blue-500 text-white': preset.number == radio.preset }"
			v-for="preset in radio.presets"
			:key="preset.number"
		>{{ preset.name }}</loading-button>
	</div>
</template>

<script>
import { mapState, mapActions } from 'vuex';

import LoadingButton from './LoadingButton.vue'

export default {
	components: {
		LoadingButton
	},
	computed: {
		...mapState([
			'radio',
			'radioConnected',
			'radioConnecting',
			'radioUUID',
		])
	},
	methods: {
		...mapActions(['playRadioPreset']),
	},
}
</script>

<style scoped>
</style>