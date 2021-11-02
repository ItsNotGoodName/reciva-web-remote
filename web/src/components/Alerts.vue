<template>
	<div class="fixed bottom-0 left-0 right-0 flex-grow mx-8 mb-8 space-y-2">
		<ErrorAlert :error="v" v-for="v  in notifications" />
	</div>
</template>

<script>
import { mapActions, mapState } from 'vuex';
import ErrorAlert from './ErrorAlert.vue';

export default {
	components: {
		ErrorAlert
	},
	computed: {
		...mapState(['notifications'])
	},
	methods: {
		...mapActions([
			'loadAll'
		])
	}
}
</script>

<style scoped>
</style>