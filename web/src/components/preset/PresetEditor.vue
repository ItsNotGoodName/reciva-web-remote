<template>
  <div class="p-fluid">
    <div class="card">
      <DataTable
        :value="presets"
        editMode="cell"
        @cell-edit-complete="onCellEditComplete"
        class="editable-cells-table"
        responsiveLayout="scroll"
      >
        <Column
          v-for="col of columns"
          :field="col.field"
          :header="col.header"
          :key="col.field"
          style="width: 25%"
        >
          <template v-if="col.field != 'url'" #editor="{ data, field }">
            <InputText v-model="data[field]" autofocus />
          </template>
        </Column>
      </DataTable>
    </div>
  </div>
</template>

<script>
import { mapActions, mapState } from "vuex";

import Column from "primevue/column";
import DataTable from "primevue/datatable";
import InputText from "primevue/inputtext";

export default {
  components: {
    Column,
    DataTable,
    InputText,
  },
  data() {
    return {
      columns: null,
    };
  },
  created() {
    this.columns = [
      {
        field: "url",
        header: "URL",
      },
      {
        field: "newName",
        header: "New Name",
      },
      {
        field: "newUrl",
        header: "New URL",
      },
    ];
  },
  methods: {
    ...mapActions(["updatePreset"]),
    onCellEditComplete(event) {
      let { data, newValue, field } = event;
      data[field] = newValue;
      this.updatePreset(data)
        .then(() => {
          this.$toast.add({
            severity: "success",
            summary: "Preset Updated",
            life: 1000,
          });
        })
        .catch((err) => {
          this.$toast.add({
            severity: "error",
            summary: "Preset Update Failed",
            detail: err,
            life: 3000,
          });
        });
    },
  },
  computed: {
    ...mapState({
      presets: (state) => state.p.presets,
    }),
  },
};
</script>

<style scoped>
::v-deep(.editable-cells-table td.p-cell-editing) {
  padding-top: 0;
  padding-bottom: 0;
}
</style>