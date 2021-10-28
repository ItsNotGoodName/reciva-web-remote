<template>
	<div class="bg-white space-y-2">
		<div class="space-y-2 sm:space-y-0 sm:flex">
			<div class="flex gap-2">
				<button class="p-1 bg-yellow-200 hover:bg-yellow-300 rounded w-" @click="discoverRadios">
					<div v-if="discover" class="w-6 h-6 my-auto border-2 border-blue-600 rounded-full loader" />
					<div v-else>Discover</div>
				</button>
				<select class="flex-1 h-8 bg-gray-200 rounded" v-model="radioUUID">
					<option :value="null" disabled value>Select Radio</option>
					<option v-bind:value="uuid" v-for="name,uuid in radios">{{ name }}</option>
				</select>
				<button v-if="radio" class="p-1 bg-gray-200 hover:bg-gray-300 rounded" @click="renewRadio">
					<div v-if="renew" class="w-6 h-6 my-auto border-2 border-blue-600 rounded-full loader" />
					<div v-else>Refresh</div>
				</button>
			</div>
			<div class="ml-auto flex gap-2" v-if="radio">
				<div class="flex gap-2 flex-grow">
					<VolumeOffIcon v-if="radio.isMuted" class="w-8" />
					<VolumeUpIcon v-else class="w-8" />
					<button class="flex-grow my-auto h-8 text-left" @click="refreshRadioVolume">{{ radio.volume }}%</button>
				</div>
				<button class="w-8 border-2 hover:bg-gray-300 rounded" @click="decreaseRadioVolume">
					<div
						v-if="decreaseVolume"
						class="w-6 h-6 my-auto border-2 border-blue-600 rounded-full loader"
					/>
					<ChevronDownIcon v-else />
				</button>
				<button class="w-8 border-2 hover:bg-gray-300 rounded" @click="increaseRadioVolume">
					<div
						v-if="increaseVolume"
						class="w-6 h-6 my-auto border-2 border-blue-600 rounded-full loader"
					/>
					<ChevronUpIcon v-else />
				</button>
				<button
					v-if="radio.power"
					class="bg-green-200 hover:bg-green-300 rounded w-12"
					@click="toggleRadioPower"
				>
					<div v-if="power" class="w-6 h-6 m-auto border-2 border-blue-600 rounded-full loader" />
					<div v-else>ON</div>
				</button>
				<button v-else class="bg-red-200 hover:bg-red-300 rounded w-12" @click="toggleRadioPower">
					<div v-if="power" class="w-6 h-6 m-auto border-2 border-blue-600 rounded-full loader" />
					<div v-else>OFF</div>
				</button>
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
import { mapState, mapActions, mapMutations } from 'vuex';

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
	data() {
		return {
			discover: false,
			power: false,
			renew: false,
			increaseVolume: false,
			decreaseVolume: false,
		}
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
		...mapActions([
			'loadRadios'
		]),
		...mapMutations([
			'SET_RADIO_POWER'
		]),
		discoverRadios() {
			if (this.discover) {
				return
			}
			this.discover = true
			api.discoverRadios()
				.then(() => {
					this.loadRadios()
					this.discover = false
				})
				.catch(() => this.discover = false)
		},
		renewRadio() {
			if (this.renew) {
				return
			}
			this.renew = true
			api.renewRadio(this.radio.uuid)
				.then(() => this.renew = false)
				.catch(() => this.renew = false)
		},
		toggleRadioPower() {
			if (this.power) {
				return
			}
			this.power = true
			let newPower = !this.radio.power
			api.updateRadio(this.radio.uuid, { power: newPower })
				.then(() => {
					this.SET_RADIO_POWER(newPower)
					this.power = false
				}).catch(() => this.power = false)
		},
		refreshRadioVolume() {
			api.refreshRadioVolume(this.radio.uuid)
		},
		increaseRadioVolume() {
			if (this.increaseVolume) {
				return
			}
			this.increaseVolume = true
			api.updateRadio(this.radio.uuid, { volume: this.radio.volume + 5 })
				.then(() => this.increaseVolume = false)
				.catch(() => this.increaseVolume = false)
		},
		decreaseRadioVolume() {
			if (this.decreaseVolume) {
				return
			}
			this.decreaseVolume = true
			api.updateRadio(this.radio.uuid, { volume: this.radio.volume - 5 })
				.then(() => this.decreaseVolume = false)
				.catch(() => this.decreaseVolume = false)
		}
	},
}
</script>

<style scoped>
</style>
