<template>
  <div class="py-2 border-b border-slate-100 last:border-b-0 md:border-b-0">
    <p class="text-xs font-medium text-slate-500 uppercase tracking-wider">{{ label }}</p>
    <p class="mt-1 text-md sm:text-lg text-slate-700 font-semibold">{{ value }}</p>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  label: string;
  value: string | number | undefined;
}>();
</script>
