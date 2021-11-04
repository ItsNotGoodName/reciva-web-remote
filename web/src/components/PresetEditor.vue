<template>
  <div class="sm:flex gap-2">
    <Streams class="flex-grow" />
    <edit-stream class="flex-auto" />
  </div>
</template>

<script>
import Streams from "./Streams.vue";
import EditStream from "./EditStream.vue";

export default {
  components: {
    Streams,
    EditStream,
  },
};
</script>

<style>
</style>