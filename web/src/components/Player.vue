<template>
	<div class="bg-white space-y-2">
		<div class="space-y-2 sm:space-y-0 sm:flex">
			<div class="flex gap-2">
				<button class="p-1 bg-yellow-200 hover:bg-yellow-300 rounded" @click="discoverRadios">Discover</button>
				<select class="flex-1 h-8 bg-gray-200 rounded" v-model="radioUUID">
					<option :value="null" disabled value>Select Radio</option>
					<option v-bind:value="uuid" v-for="name,uuid in radios">{{ name }}</option>
				</select>
				<button
					v-if="radio"
					class="p-1 bg-gray-200 hover:bg-gray-300 rounded"
					@click="renewRadio"
				>Refresh</button>
			</div>
			<div class="ml-auto flex gap-2" v-if="radio">
				<div class="flex gap-2 flex-grow">
					<VolumeOffIcon v-if="radio.isMuted" class="w-8" />
					<VolumeUpIcon v-else class="w-8" />
					<button class="flex-grow my-auto h-8 text-left" @click="refreshRadioVolume">{{ radio.volume }}%</button>
				</div>
				<button class="w-8 border-2 hover:bg-gray-300 rounded" @click="decreaseRadioVolume">
					<ChevronDownIcon />
				</button>
				<button class="w-8 border-2 hover:bg-gray-300 rounded" @click="increaseRadioVolume">
					<ChevronUpIcon />
				</button>
				<button
					v-if="radio.power"
					class="bg-green-200 hover:bg-green-300 rounded w-12"
					@click="toggleRadioPower"
				>ON</button>
				<button v-else class="bg-red-200 hover:bg-red-300 rounded w-12" @click="toggleRadioPower">OFF</button>
			</div>
		</div>
		<div class="flex space-x-2" v-if="radio">
			<PlayIcon v-if="radio.state == 'Playing'" class="w-8" />
			<StopIcon v-else-if="radio.state == 'Stopped'" class="w-8" />
			<RefreshIcon v-else class="w-8" />
			<span class="my-auto truncate">{{ radio.title }}</span>
		</div>
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
		discoverRadios() {
			api.discoverRadios()
		},
		renewRadio() {
			api.renewRadio(this.radio.uuid)
		},
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
