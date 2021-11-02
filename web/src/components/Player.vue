<template>
	<div class="bg-white space-y-2">
		<div class="space-y-2 sm:space-y-0 sm:flex">
			<div class="flex gap-2">
				<loading-button
					class="w-20 p-1 bg-yellow-200 hover:bg-yellow-300 rounded"
					:on-click="discoverRadios"
				>Discover</loading-button>
				<select class="flex-1 h-8 bg-gray-200 rounded" v-model="radioUUID">
					<option :value="null" disabled value>Select Radio</option>
					<option v-bind:value="uuid" v-for="name,uuid in radios">{{ name }}</option>
				</select>
				<loading-button
					v-if="radioUUID"
					class="w-16 p-1 rounded bg-gray-200 hover:bg-gray-300"
					:on-click="refreshRadio"
				>Refresh</loading-button>
			</div>
			<div class="ml-auto flex gap-2" v-if="radioConnected">
				<div class="flex gap-2 flex-grow">
					<VolumeOffIcon v-if="radio.isMuted" class="w-8 h-8" />
					<VolumeUpIcon v-else class="w-8 h-8" />
					<loading-button
						class="w-10 h-8 sm:text-center text-left flex-grow"
						:on-click="refreshRadioVolume"
					>{{ radio.volume }}%</loading-button>
				</div>
				<loading-button class="w-16 hover:bg-gray-200 rounded" :on-click="decreaseRadioVolume">
					<ChevronDownIcon class="w-16 h-8 border-2" />
				</loading-button>
				<loading-button class="w-16 hover:bg-gray-200 rounded" :on-click="increaseRadioVolume">
					<ChevronUpIcon class="w-16 h-8 border-2" />
				</loading-button>
				<loading-button
					v-if="radio.power"
					class="bg-green-200 hover:bg-green-300 rounded w-16"
					:on-click="toggleRadioPower"
				>ON</loading-button>
				<loading-button
					v-else
					class="bg-red-200 hover:bg-red-300 rounded w-16"
					:on-click="toggleRadioPower"
				>OFF</loading-button>
			</div>
		</div>
		<div class="flex space-x-2" v-if="radioConnected">
			<div class="w-8 h-8 flex-shrink-0">
				<PlayIcon v-if="radio.state == 'Playing'" />
				<StopIcon v-else-if="radio.state == 'Stopped'" />
				<RefreshIcon v-else />
			</div>
			<div class="my-auto truncate space-x-2">
				<span class="px-2 py-0.5 text-base rounded-full text-white bg-blue-500">{{ radio.title }}</span>
				<span v-if="radio.metadata">{{ radio.metadata }}</span>
			</div>
		</div>
	</div>
</template>

<script>
import { mapState, mapActions, mapMutations } from 'vuex';

import PlayIcon from "@heroicons/vue/solid/PlayIcon"
import StopIcon from "@heroicons/vue/solid/StopIcon"
import RefreshIcon from "@heroicons/vue/solid/RefreshIcon"
import VolumeOffIcon from "@heroicons/vue/solid/VolumeOffIcon"
import VolumeUpIcon from "@heroicons/vue/solid/VolumeUpIcon"
import ChevronUpIcon from "@heroicons/vue/solid/ChevronUpIcon"
import ChevronDownIcon from "@heroicons/vue/solid/ChevronDownIcon"
import LoadingButton from "./LoadingButton.vue"
import api from '../api';

export default {
	components: {
		PlayIcon,
		StopIcon,
		RefreshIcon,
		VolumeOffIcon,
		VolumeUpIcon,
		ChevronUpIcon,
		ChevronDownIcon,
		LoadingButton
	},
	computed: {
		...mapState([
			'radio',
			'radioUUID',
			'radioConnected',
			'radioConnecting',
			'radios',
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
		...mapActions([
			'refreshRadio',
			'discoverRadios'
		]),
		...mapMutations([
			'SET_RADIO_POWER',
		]),
		toggleRadioPower() {
			let newPower = !this.radio.power
			return api.updateRadio(this.radio.uuid, { power: newPower })
				.then(() => this.SET_RADIO_POWER(newPower))
		},
		refreshRadioVolume() {
			return api.refreshRadioVolume(this.radio.uuid)
		},
		increaseRadioVolume() {
			return api.updateRadio(this.radio.uuid, { volume: this.radio.volume + 5 })
		},
		decreaseRadioVolume() {
			return api.updateRadio(this.radio.uuid, { volume: this.radio.volume - 5 })
		}
	},
}
</script>

<style scoped>
</style>


