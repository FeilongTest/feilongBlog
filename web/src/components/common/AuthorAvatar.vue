<template>
  <img :src="avatarUrl" :alt="alt || authorName" @error="failed = true" />
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useAuthStore } from "@/stores/auth";
import { getAssetPath } from "@/core/helpers/assets";
import { resolveMediaUrl } from "@/utils/media";

const props = defineProps<{ src?: string; alt?: string }>();
const auth = useAuthStore();
const failed = ref(false);
const source = computed(() => resolveMediaUrl(props.src ?? auth.publicProfile?.pic));
const authorName = computed(() => auth.publicProfile?.trueName || "作者头像");
const avatarUrl = computed(() => !failed.value && source.value ? source.value : getAssetPath("/media/avatars/300-1.jpg"));
watch(source, () => { failed.value = false; });
onMounted(() => { if (props.src === undefined) void auth.refreshPublicProfile().catch(() => undefined); });
</script>
