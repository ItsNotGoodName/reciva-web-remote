<template>
	<div class="bg-white flex flex-wrap gap-2">
		<template v-if="radio">
			<div class="flex space-x-2 flex-grow">
				<button
					v-if="radio.power"
					class="bg-green-200 hover:bg-green-400 rounded w-10"
					@click="toggleRadioPower"
				>ON</button>
				<button v-else class="bg-red-200 hover:bg-red-400 rounded w-10" @click="toggleRadioPower">OFF</button>
				<PlayIcon v-if="radio.state == 'Playing'" class="w-8" />
				<StopIcon v-else-if="radio.state == 'Stopped'" class="w-8" />
				<RefreshIcon v-else class="w-8" />
				<div class="my-auto">{{ radio.title }}</div>
			</div>
			<div class="flex">
				<VolumeOffIcon v-if="radio.isMuted" class="w-8" />
				<VolumeUpIcon v-else class="w-8" />
				<button class="ml-2 w-8 hover:bg-gray-200 rounded" @click="decreaseRadioVolume">
					<ChevronDownIcon />
				</button>
				<button
					class="my-auto p-1 hover:bg-gray-200 rounded w-10"
					@click="refreshRadioVolume"
				>{{ radio.volume }}%</button>
				<button class="w-8 hover:bg-gray-200 rounded" @click="increaseRadioVolume">
					<ChevronUpIcon />
				</button>
			</div>
		</template>
		<select class="ml-auto h-8" v-model="radioUUID">
			<option :value="null" disabled value>Select Radio</option>
			<option v-bind:value="uuid" v-for="name,uuid in radios">{{ name }}</option>
		</select>
	</div>
</template>

<script>
import { mapState } from 'vuex';

import PlayIcon from "@heroicons/vue/solid/PlayIcon"
import StopIcon from "@heroicons/vue/solid/StopIcon"
import RefreshIcon from "@heroicons/vue/solid/RefreshIcon"
import VolumeOffIcon from "@heroicons/vue/solid/VolumeOffIcon"
import VolumeUpIcon from "@heroicons/vue/solid/VolumeUpIcon"
import ChevronUpIcon from "@heroicons/vue/solid/ChevronUpIcon"
import ChevronDownIcon from "@heroicons/vue/solid/ChevronDownIcon"
import api from '../api';

export default {
	components: {
		PlayIcon,
		StopIcon,
		RefreshIcon,
		VolumeOffIcon,
		VolumeUpIcon,
		ChevronUpIcon,
		ChevronDownIcon
	},
	computed: {
		...mapState([
			'radio',
			'radios'
		]),
		radioUUID: {
			get() {
				return this.$store.state.radioUUID
			},
			set(uuid) {
				this.$store.dispatch("setRadioUUID", uuid)
			}
		}
	},
	methods: {
		toggleRadioPower() {
			api.updateRadio(this.radio.uuid, { power: !this.radio.power })
		},
		refreshRadioVolume() {
			api.refreshRadioVolume(this.radio.uuid)
		},
		increaseRadioVolume() {
			api.updateRadio(this.radio.uuid, { volume: this.radio.volume + 5 })
		},
		decreaseRadioVolume() {
			api.updateRadio(this.radio.uuid, { volume: this.radio.volume - 5 })
		}
	},
}
</script>

<style scoped>
</style>