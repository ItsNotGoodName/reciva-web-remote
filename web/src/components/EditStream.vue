<template>
  <div :class="{ hidden: !streamShow }">
    <div class="p-2 bg-info rounded-t">
      <span v-if="stream.id">
        <strong class="px-2">Stream Details</strong>
      </span>
      <span v-else>
        <strong class="px-2">Add Stream</strong>
      </span>
    </div>
    <div class="flex border-l-2 border-r-2 b">
      <form class="w-full mx-auto bg-white rounded p-2">
        <div class="mb-4">
          <label class="block text-sm font-bold mb-2" for="name">
            Stream Name
          </label>
          <input
            v-model="streamName"
            :class="{ 'border-yellow-500': stream.nameChanged }"
            class="border rounded w-full py-2 px-3"
            :disabled="streamLoading"
            id="name"
            type="text"
            placeholder="Stream Name"
          />
        </div>
        <div class="mb-4">
          <label class="block text-sm font-bold mb-2" for="content">
            Stream Content
          </label>
          <input
            v-model="streamContent"
            :class="{ 'border-yellow-500': stream.contentChanged }"
            class="border rounded w-full py-2 px-3"
            :disabled="streamLoading"
            id="content"
            type="text"
            placeholder="Stream Content"
          />
        </div>
      </form>
    </div>
    <div class="p-2 flex flex-row-reverse rounded-b bg-light">
      <template v-if="stream.id">
        <template v-if="stream.deleteConfirm">
          <button @click="TOGGLE_DELETE_CONFIRM" class="btn btn-secondary">
            Cancel
          </button>
          <loading-button
            :on-click="deleteStream"
            class="btn btn-danger rounded-r w-26"
            >Confirm Delete</loading-button
          >
        </template>
        <template v-else>
          <button
            @click="TOGGLE_DELETE_CONFIRM"
            class="btn btn-danger rounded-r"
          >
            Delete
          </button>
          <button @click="closeStream()" class="btn btn-secondary">
            Close
          </button>
          <loading-button
            :on-click="updateStream"
            class="btn text-white btn-success rounded-l w-16"
          >
            Save
          </loading-button>
        </template>
      </template>
      <template v-else>
        <button @click="closeStream()" class="btn btn-secondary rounded-r">
          Cancel
        </button>
        <loading-button
          :on-click="createStream"
          class="btn text-white btn-success rounded-l w-16"
          >Add</loading-button
        >
      </template>
    </div>
  </div>
</template>

<script>
import { mapMutations, mapActions, mapState } from "vuex";

import LoadingButton from "./LoadingButton.vue";

export default {
  components: {
    LoadingButton,
  },
  computed: {
    ...mapState({
      streamShow: (state) => state.p.streamShow,
      stream: (state) => state.p.stream,
      streamLoading: (state) => state.p.streamLoading,
    }),
    streamName: {
      get() {
        return this.$store.state.p.stream.name;
      },
      set(name) {
        this.$store.commit("SET_STREAM_NAME", name);
      },
    },
    streamContent: {
      get() {
        return this.$store.state.p.stream.content;
      },
      set(content) {
        this.$store.commit("SET_STREAM_CONTENT", content);
      },
    },
  },
  methods: {
    ...mapMutations(["TOGGLE_DELETE_CONFIRM"]),
    ...mapActions([
      "loadStreams",
      "closeStream",
      "updateStream",
      "createStream",
      "deleteStream",
    ]),
  },
};
</script>

<style scoped>
</style>