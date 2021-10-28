<template>
	<button @click="handleClick">
		<div v-if="loading" class="w-5 h-5 m-auto border-2 border-blue-600 rounded-full loader" />
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
				.then((() => this.loading = false))
				.catch((() => this.loading = false))
		},
	},
}
</script>

<style scoped>
@keyframes loader-rotate {
	0% {
		transform: rotate(0);
	}
	100% {
		transform: rotate(360deg);
	}
}
.loader {
	border-right-color: transparent;
	animation: loader-rotate 1s linear infinite;
}
</style>