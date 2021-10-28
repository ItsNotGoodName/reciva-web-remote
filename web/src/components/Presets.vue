<template>
	<div v-if="radio" class="gap-1 grid">
		<loading-button
			:on-click="() => playPreset(preset.number)"
			class="hover:bg-gray-200 rounded p-1"
			v-bind:class="{ 'bg-gray-200': preset.url == radio.url }"
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
			'radio'
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