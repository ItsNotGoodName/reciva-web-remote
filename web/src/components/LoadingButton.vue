<template>
	<button @click="handleClick">
		<div v-if="loading" class="w-6 h-6 m-auto border-2 border-blue-600 rounded-full loader" />
		<slot v-else></slot>
	</button>
</template>

<script>
export default {
	props: {
		onClick: {
			type: Function,
			default: () => Promise.resolve()
		},
	},
	data() {
		return {
			loading: false
		}
	}, methods: {
		handleClick() {
			if (this.loading) {
				return
			}
			this.loading = true
			this.onClick()
				.then((() => { this.loading = false }).bind(this))
				.catch((() => { this.loading = false }).bind(this))
		},
	},
}
</script>

<style scoped>
</style>