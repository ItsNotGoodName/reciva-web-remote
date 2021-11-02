<template >
	<div
		v-if="radioUUID && radioConnected"
		class="gap-1 grid lg:grid-cols-4 md:grid-cols-3 sm:grid-cols-2"
	>
		<loading-button
			:on-click="() => playPreset(preset.number)"
			class="hover:bg-gray-300 rounded p-1"
			v-bind:class="{ 'bg-gray-200': preset.number == radio.preset }"
			v-for="preset in radio.presets"
			:key="preset.number"
		>{{ preset.name }}</loading-button>
	</div>
</template>

<script>
import { mapState } from 'vuex';
import LoadingButton from './LoadingButton.vue'
import api from '../api';

export default {
	components: {
		LoadingButton
	},
	computed: {
		...mapState([
			'radio',
			'radioUUID',
			'radioConnecting',
			'radioConnected',
		])
	},
	methods: {
		playPreset(num) {
			return api.updateRadio(this.radio.uuid, { preset: num })
		},
	},

}
</script>

<style scoped>
</style>