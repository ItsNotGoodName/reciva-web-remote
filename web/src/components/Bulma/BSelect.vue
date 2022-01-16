<template>
  <div class="select" :class="{ 'is-loading': loading }">
    <select
      :value="modelValue"
      @input="$emit('update:modelValue', $event.target.value)"
    >
      <slot></slot>
    </select>
  </div>
</template>

<script>
export default {
  props: {
    modelValue: {
      type: String,
      default: "",
    },
    loading: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss" scoped></style>
