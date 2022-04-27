<script setup lang="ts">
import { PropType, ref } from 'vue';

import DButton from './DaisyUI/DButton.vue';
import DDropdownButton from './DaisyUI/DDropdownButton.vue';

const { audiosources, audiosource, setAudiosource } = defineProps({
  audiosources: {
    type: Array as PropType<string[]>,
    default: [],
  },
  audiosource: {
    type: String,
    default: '',
  },
  setAudiosource: {
    type: Function as PropType<(value: string) => Promise<void>>,
    default: () => { return Promise.resolve(); },
  }
});

const loading = ref(false);

const onClick = (a: string) => {
  if (loading.value) {
    return;
  }

  loading.value = true;
  setAudiosource(a).finally(() => {
    loading.value = false;
  });
}
</script>

<template>
  <div class="dropdown dropdown-top dropdown-end">
    <d-dropdown-button :class="{ 'btn-secondary': audiosource }" aria-label="Audiosource">
      <v-icon name="fa-itunes-note" />
    </d-dropdown-button>
    <ul tabindex="0" class="menu card card-compact bg-base-200 p-2 menu-compact dropdown-content space-y-2 mb-2 w-52">
      <d-button :loading="loading" :key="a" v-for="a in audiosources" :class="{ 'btn-secondary': a == audiosource }"
        @click="() => onClick(a)">
        {{ a }}
      </d-button>
    </ul>
  </div>
</template>

<style>
</style>

