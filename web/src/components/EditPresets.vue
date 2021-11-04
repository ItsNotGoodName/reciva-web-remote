<template>
    <div class="border-gray-100">
        <div class="p-2 border-2 rounded-t border-yellow-400 bg-yellow-400 text-white">
            <strong class="px-2">My Streams</strong>
        </div>
        <div class="table border-l-2 w-full border-r-2 border-gray-100">
            <div class="table-row-group" :key="stream.id" v-for="stream in streams">
                <button
                    class="w-full text-cell p-2 text-left hover:bg-gray-100"
                    :title="stream.name"
                >{{ stream.name }}</button>
            </div>
        </div>
        <div class="p-2 flex flex-row-reverse rounded-b bg-gray-100">
            <loading-button
                :on-click="loadStreams"
                class="w-20 p-2 rounded-md text-white bg-yellow-400"
            >Refresh</loading-button>
            <button class="p-2 rounded-md text-white bg-green-400 mr-2">Add a Stream</button>
        </div>
    </div>
</template>

<script>
import { mapActions, mapState } from 'vuex';
import LoadingButton from './LoadingButton.vue';

export default {
    components: { LoadingButton },
    computed: {
        ...mapState([
            'streams',
            'presets'
        ])
    },
    methods: {
        ...mapActions([
            'loadStreams',
        ])
    }
}
</script>

<style scoped>
</style>