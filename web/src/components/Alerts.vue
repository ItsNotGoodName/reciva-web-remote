<template>
	<div class="fixed bottom-0 left-0 right-0 flex-grow mx-8 mb-8 space-y-2">
		<TransitionGroup>
			<Alert :key="k" :notification="v" v-for="v, k  in notifications" />
		</TransitionGroup>
	</div>
</template>

<script>
import { mapState } from 'vuex';

import Alert from './Alert.vue';

export default {
	components: {
		Alert,
	},
	computed: {
		...mapState(['notifications']),
	},
}
</script>

<style scoped>
.v-enter-active,
.v-leave-active {
	@apply transition;
	@apply transform;
	@apply duration-300;
}

.v-enter-from {
	@apply translate-y-full;
}

.v-leave-to {
	@apply translate-x-full;
	@apply opacity-0;
}
</style>