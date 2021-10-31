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
					v-if="radioUUID || radioConnected"
					class="w-16 p-1 bg-gray-200 hover:bg-gray-300 rounded"
					:on-click="refreshRadio"
				>Refresh</loading-button>
			</div>
			<div class="ml-auto flex gap-2" v-if="radioConnected">
				<div class="flex gap-2 flex-grow">
					<VolumeOffIcon v-if="radio.isMuted" class="w-8 h-8" />
					<VolumeUpIcon v-else class="w-8 h-8" />
					<button class="h-8 text-left flex-grow" @click="refreshRadioVolume">{{ radio.volume }}%</button>
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
			<PlayIcon v-if="radio.state == 'Playing'" class="w-8 h-8" />
			<StopIcon v-else-if="radio.state == 'Stopped'" class="w-8 h-8" />
			<RefreshIcon v-else class="w-8 h-8" />
			<span class="my-auto truncate">
				{{ radio.title }}
				<span v-if="radio.metadata">| {{ radio.metadata }}</span>
			</span>
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
			'loadRadios',
			'refreshRadio'
		]),
		...mapMutations([
			'SET_RADIO_POWER'
		]),
		discoverRadios() {
			return api.discoverRadios().then(() => this.loadRadios)
		},
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

