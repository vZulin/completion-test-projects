<script setup lang="ts">
import { computed, ref } from "vue";

const count = ref(0);
const doubled = computed(() => count.value * 2);

function increment(): void {
  count.value += 1;
}
</script>

<template>
  <button type="button" @click="increment">
    Count: {{ count }} / Doubled: {{ doubled }}
  </button>
</template>
