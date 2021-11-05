<template>
  <div :class="{ invisible: !showStream }">
    <div class="p-2 bg-info rounded-t">
      <span v-if="isStreamEdit">
        <strong class="px-2">Stream Details</strong>
      </span>
      <span v-else>
        <strong class="px-2">Add Stream</strong>
      </span>
    </div>
    <div class="flex border-gray border-l-2 border-r-2">
      <form class="w-full mx-auto bg-white rounded p-2">
        <div class="mb-4">
          <label class="block text-gray-700 text-sm font-bold mb-2" for="name">
            Stream Name
          </label>
          <input
            v-model="streamName"
            class="
              appearance-none
              border
              rounded
              w-full
              py-2
              px-3
              text-gray-700
              leading-tight
            "
            id="name"
            type="text"
            placeholder="Stream Name"
          />
        </div>
        <div class="mb-4">
          <label
            class="block text-gray-700 text-sm font-bold mb-2"
            for="content"
          >
            Stream Content
          </label>
          <input
            v-model="streamContent"
            class="
              appearance-none
              border
              rounded
              w-full
              py-2
              px-3
              text-gray-700
              leading-tight
            "
            id="content"
            type="text"
            placeholder="Stream Content"
          />
        </div>
      </form>
    </div>
    <div class="p-2 flex flex-row-reverse rounded-b bg-light">
      <template v-if="isStreamEdit">
        <loading-button
          :on-click="deleteStream"
          class="btn btn-danger rounded-r"
          >Delete</loading-button
        >
        <button @click="hideStream()" class="btn btn-secondary">Close</button>
        <loading-button
          :on-click="updateStream"
          class="btn text-white btn-success rounded-l"
        >
          Save
        </loading-button>
      </template>
      <template v-else>
        <button @click="hideStream()" class="btn btn-secondary rounded-r w-18">
          Cancel
        </button>
        <loading-button
          :on-click="newStream"
          class="btn text-white btn-success rounded-l w-16"
          >Add</loading-button
        >
      </template>
    </div>
  </div>
</template>

<script>
import { mapActions, mapState } from "vuex";

import LoadingButton from "./LoadingButton.vue";

export default {
  components: {
    LoadingButton,
  },
  computed: {
    ...mapState(["isStreamEdit", "showStream"]),
    streamName: {
      get() {
        return this.$store.state.stream.name;
      },
      set(name) {
        this.$store.commit("SET_STREAM_NAME", name);
      },
    },
    streamContent: {
      get() {
        return this.$store.state.stream.content;
      },
      set(content) {
        this.$store.commit("SET_STREAM_CONTENT", content);
      },
    },
    streamURI: {
      get() {
        return this.$store.state.stream.uri;
      },
      set(uri) {
        this.$store.commit("SET_STREAM_URI", uri);
      },
    },
  },
  methods: {
    ...mapActions([
      "loadStreams",
      "hideStream",
      "updateStream",
      "newStream",
      "deleteStream",
    ]),
  },
};
</script>

<style scoped>
</style>