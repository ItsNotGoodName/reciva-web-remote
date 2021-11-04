<template>
	<div class="bg-white space-y-2">
		<div class="space-y-2 sm:space-y-0 sm:flex">
			<div class="flex gap-2">
				<loading-button
					class="w-20 p-1 bg-yellow-200 hover:bg-yellow-300 rounded"
					:on-click="discoverRadios"
				>Discover</loading-button>
				<select class="flex-1 h-8 bg-gray-200 rounded" v-model="radioUUID">
					<option :value="null" disabled>Select Radio</option>
					<option :key="uuid" v-bind:value="uuid" v-for="name,uuid in radios">{{ name }}</option>
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
import { mapState, mapActions, } from 'vuex';

import ChevronDownIcon from "@heroicons/vue/solid/ChevronDownIcon"
import ChevronUpIcon from "@heroicons/vue/solid/ChevronUpIcon"
import PlayIcon from "@heroicons/vue/solid/PlayIcon"
import RefreshIcon from "@heroicons/vue/solid/RefreshIcon"
import StopIcon from "@heroicons/vue/solid/StopIcon"
import VolumeOffIcon from "@heroicons/vue/solid/VolumeOffIcon"
import VolumeUpIcon from "@heroicons/vue/solid/VolumeUpIcon"

import LoadingButton from "./LoadingButton.vue"

export default {
	components: {
		ChevronDownIcon,
		ChevronUpIcon,
		PlayIcon,
		RefreshIcon,
		StopIcon,
		VolumeOffIcon,
		VolumeUpIcon,
		LoadingButton,
	},
	computed: {
		...mapState([
			'radio',
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
			'decreaseRadioVolume',
			'discoverRadios',
			'increaseRadioVolume',
			'refreshRadio',
			'refreshRadioVolume',
			'toggleRadioPower',
		]),
	},
}
</script>

<style scoped>
</style>


