<template>
  <div v-if="radioReady" class="flex">
    <Button
      class="flex-shrink-0 p-button-rounded pl-1 mr-2"
      v-if="radio.state == 'Playing'"
      icon="pi pi-play"
    />
    <Button
      class="flex-shrink-0 p-button-warning p-button-rounded mr-2"
      v-else-if="radio.state == 'Connecting'"
      icon="pi pi-spin pi-spinner"
    />
    <Button
      class="flex-shrink-0 p-button-danger p-button-rounded mr-2"
      v-else-if="radio.state == 'Stopped'"
      icon="pi pi-stop"
    />
    <Button @click="toggle" class="h-full flex-grow-1">
      {{ radio.title }}
    </Button>
    <OverlayPanel @show="readPreset" ref="op">
      <table>
        <caption>
          Preset Information
        </caption>
        <tbody>
          <tr>
            <td>
              <Badge class="h-full w-full" severity="warning">Metadata</Badge>
            </td>
            <td>{{ radio.metadata }}</td>
          </tr>
          <tr>
            <td><Badge class="h-full w-full">URL</Badge></td>
            <td>{{ radio.url }}</td>
          </tr>
          <tr>
            <td>
              <Badge class="h-full w-full flex-nowrap" severity="success"
                >New URL</Badge
              >
            </td>
            <td>
              <Skeleton v-if="loading" />
              <span v-else>
                {{ preset.newUrl }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </OverlayPanel>
  </div>
</template>

<script>
import Badge from "primevue/badge";
import Button from "primevue/button";
import OverlayPanel from "primevue/overlaypanel";
import Skeleton from "primevue/skeleton";

import { mapGetters, mapState } from "vuex";
export default {
  components: {
    Badge,
    Button,
    OverlayPanel,
    Skeleton,
  },
  data() {
    return {
      loading: false,
      reload: false,
    };
  },
  computed: {
    ...mapState({
      radio: (state) => state.radio,
      preset: (state) => state.p.preset,
    }),
    ...mapGetters(["radioReady"]),
  },
  methods: {
    toggle(event) {
      this.$refs.op.toggle(event);
    },
    readPreset() {
      if (this.loading) {
        this.reload = true;
        return;
      }
      this.loading = true;
      this.reload = false;

      this.$store.dispatch("readPreset").finally(() => {
        this.loading = false;
        if (this.reload) {
          this.$nextTick(this.readPreset);
        }
      });
    },
  },
};
</script>

<style scoped>
</style>