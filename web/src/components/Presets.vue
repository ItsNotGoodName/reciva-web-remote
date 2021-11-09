<template>
  <div>
    <div class="p-2 border-2 rounded-t border-gray-600 bg-gray-600 text-white">
      <strong class="px-2">My Presets</strong>
    </div>
    <div class="border-l-2 border-b-2 border-r-2">
      <table class="w-full">
        <thead>
          <tr>
            <th>New Name</th>
            <th>URL</th>
            <th>New URL</th>
          </tr>
        </thead>
        <tbody>
          <tr
            class="hover:bg-gray-100"
            @click="openPreset(p.url)"
            :key="p.url"
            v-for="p in presets"
          >
            <td class="p-2">{{ p.newName }}</td>
            <td class="p-2">
              {{ p.url }}
            </td>
            <td class="p-2">{{ p.newUrl }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="p-2 flex flex-row-reverse rounded bg-light">
      <loading-button
        :on-click="readPresets"
        class="h-10 rounded w-24 btn btn-warning"
        >Refresh</loading-button
      >
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
      presets: (state) => state.p.presets,
    }),
  },
  methods: {
    ...mapActions(["openPreset", "readPresets", "clearPreset"]),
  },
};
</script>

<style scoped>
</style>