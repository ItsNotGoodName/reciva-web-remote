<template>
  <div class="border-gray-100">
    <div
      class="p-2 border-2 rounded-t border-yellow-600 bg-yellow-600 text-white"
    >
      <strong class="px-2">My Streams</strong>
    </div>
    <div class="table border-l-2 w-full border-r-2 border-gray-100">
      <div class="table-row-group" :key="stream.id" v-for="stream in streams">
        <button class="w-full text-cell text-left" :title="stream.name">
          {{ stream.name }}
        </button>
      </div>
    </div>
    <div class="p-2 flex flex-row-reverse rounded-b bg-gray-100">
      <loading-button
        :on-click="loadStreams"
        class="rounded-r w-24 btn btn-warning"
        >Refresh</loading-button
      >
      <button class="rounded-l w-38 btn btn-success">Add a Stream</button>
    </div>
  </div>
</template>

<script>
import { mapActions, mapState } from "vuex";
import LoadingButton from "./LoadingButton.vue";

export default {
  components: { LoadingButton },
  computed: {
    ...mapState(["streams", "presets"]),
  },
  methods: {
    ...mapActions(["loadStreams"]),
  },
};
</script>

<style scoped>
</style>