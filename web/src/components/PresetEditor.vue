<template>
  <div class="sm:flex gap-2">
    <Streams class="flex-grow mb-2" />
    <div class="flex-auto">
      <edit-stream class="mb-2" />
      <Presets />
    </div>
  </div>
</template>

<script>
import Streams from "./Streams.vue";
import EditStream from "./EditStream.vue";
import Presets from "./Presets.vue";

export default {
  components: {
    Streams,
    EditStream,
    Presets,
  },
};
</script>

<style>
</style>