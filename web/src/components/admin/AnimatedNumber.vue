<template>
  <span>{{ displayedValue.toLocaleString("zh-CN") }}</span>
</template>

<script setup lang="ts">
import { onBeforeUnmount, ref, watch } from "vue";

const props = withDefaults(defineProps<{ value: number; duration?: number }>(), {
  duration: 650,
});

const displayedValue = ref(0);
let animationFrame = 0;

const animateTo = (target: number) => {
  cancelAnimationFrame(animationFrame);
  const safeTarget = Number.isFinite(target) ? Math.max(0, target) : 0;

  if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) {
    displayedValue.value = safeTarget;
    return;
  }

  const startValue = displayedValue.value;
  const startedAt = performance.now();
  const tick = (now: number) => {
    const progress = Math.min((now - startedAt) / props.duration, 1);
    const eased = 1 - Math.pow(1 - progress, 3);
    displayedValue.value = Math.round(startValue + (safeTarget - startValue) * eased);
    if (progress < 1) animationFrame = requestAnimationFrame(tick);
  };
  animationFrame = requestAnimationFrame(tick);
};

watch(() => props.value, animateTo, { immediate: true });
onBeforeUnmount(() => cancelAnimationFrame(animationFrame));
</script>
