<template>
	<transition name="fade" v-on:afterLeave="afterLeave">
		<div
			@click="this.show = false"
			class="flex w-full max-w-sm mx-auto overflow-hidden bg-white rounded-lg shadow-md dark:bg-gray-800 cursor-pointer"
			v-show="show"
		>
			<div class="flex items-center justify-center w-12 bg-red-500">
				<svg
					class="w-6 h-6 text-white fill-current"
					viewBox="0 0 40 40"
					xmlns="http://www.w3.org/2000/svg"
				>
					<path
						d="M20 3.36667C10.8167 3.36667 3.3667 10.8167 3.3667 20C3.3667 29.1833 10.8167 36.6333 20 36.6333C29.1834 36.6333 36.6334 29.1833 36.6334 20C36.6334 10.8167 29.1834 3.36667 20 3.36667ZM19.1334 33.3333V22.9H13.3334L21.6667 6.66667V17.1H27.25L19.1334 33.3333Z"
					/>
				</svg>
			</div>

			<div class="px-4 py-2 -mx-3">
				<div class="mx-3">
					<span class="font-semibold text-red-500 dark:text-red-400">Error</span>
					<p class="text-sm text-gray-600 dark:text-gray-200">{{ this.error.message }}</p>
				</div>
			</div>
		</div>
	</transition>
</template>

<script>
import { mapMutations } from "vuex"

export default {
	mounted() {
		this.$nextTick(function () {
			this.show = true
			this.error.timeout && setTimeout((() => {
				this.show = false
			}).bind(this), this.error.timeout)
		})
	},
	data() {
		return {
			show: false
		}
	},
	methods: {
		afterLeave() {
			this.DELETE_NOTIFICATION(0)
		},
		...mapMutations(['DELETE_NOTIFICATION'])
	},
	props: {
		error: {
			type: Object
		},
	},
}
</script>

<style>
.fade-enter-active,
.fade-leave-active {
	@apply transition;
	@apply transform;
	@apply duration-300;
}

.fade-enter-from {
	@apply translate-y-full;
}

.fade-leave-to {
	@apply translate-x-full;
	@apply opacity-0;
}
</style>