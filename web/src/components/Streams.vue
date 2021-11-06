<template>
  <div>
    <div
      class="p-2 border-2 rounded-t border-yellow-600 bg-yellow-600 text-white"
    >
      <strong class="px-2">My Streams</strong>
    </div>
    <div class="table border-gray border-l-2 border-r-2 w-full border-gray-100">
      <div class="table-row-group" :key="id" v-for="(name, id) in streams">
        <loading-button
          :on-click="() => readStream(id)"
          className="border-blue-600"
          class="h-10 w-full text-cell text-left btn"
          :class="[id == stream.id ? 'btn-white-on' : 'btn-white']"
          :title="name"
        >
          {{ name }}
        </loading-button>
      </div>
    </div>
    <div class="p-2 flex flex-row-reverse rounded-b bg-light">
      <loading-button
        :on-click="readStreams"
        class="rounded-r w-24 btn btn-warning"
        >Refresh</loading-button
      >
      <button @click="addStream()" class="rounded-l w-38 btn btn-success">
        Add a Stream
      </button>
    </div>
  </div>
</template>

<script>
import { mapActions, mapState } from "vuex";
import LoadingButton from "./LoadingButton.vue";

export default {
  components: { LoadingButton },
  computed: {
    ...mapState({
      streams: (state) => state.p.streams,
      presets: (state) => state.p.presets,
      stream: (state) => state.p.stream,
    }),
  },
  methods: {
    ...mapActions(["readStreams", "readStream", "addStream"]),
  },
};
</script>

<style scoped>
</style>