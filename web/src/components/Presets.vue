<template>
  <div>
    <div class="p-2 border-2 rounded-t border-gray-600 bg-gray-600 text-white">
      <strong class="px-2">My Presets</strong>
    </div>
    <div class="rounded-b border-l-2 border-b-2 border-r-2">
      <table class="w-full">
        <thead>
          <tr>
            <th>URI</th>
            <th>Stream Name</th>
          </tr>
        </thead>
        <tbody>
          <tr class="hover:bg-gray-100" :key="p.uri" v-for="p in presets">
            <td class="p-2">
              {{ p.uri }}
            </td>
            <td class="p-2">{{ streams[p.sid] && streams[p.sid].name }}</td>
          </tr>
        </tbody>
      </table>
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