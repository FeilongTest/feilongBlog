<template>
  <div class="card">
    <div class="card-body p-lg-20">
      <div class="mw-700px">
        <span class="text-primary fw-bold fs-7 text-uppercase">About</span>
        <h1 class="text-gray-900 fw-bold fs-2x mt-3 mb-6">关于 {{ site.siteName }}</h1>
        <p v-if="site.siteSlogan" class="text-gray-600 fs-5 mb-6">{{ site.siteSlogan }}</p>
        <p class="text-gray-700 fs-5 lh-lg mb-10">{{ site.description }}</p>

        <div class="d-flex align-items-center border border-dashed card-rounded p-6 mb-10">
          <div class="symbol symbol-70px symbol-circle me-6 flex-shrink-0">
            <AuthorAvatar class="about-avatar" :src="profile.pic" :alt="authorName" />
          </div>
          <div>
            <div class="text-gray-900 fw-bold fs-4">{{ authorName }}</div>
            <p v-if="profile.bio" class="text-gray-600 mb-1 mt-2">{{ profile.bio }}</p>
            <a v-if="profile.email" class="text-primary" :href="`mailto:${profile.email}`">{{ profile.email }}</a>
          </div>
        </div>

        <div class="d-flex flex-wrap gap-3">
          <router-link to="/" class="btn btn-sm btn-primary">返回首页</router-link>
          <router-link to="/archive" class="btn btn-sm btn-light">文章归档</router-link>
          <a href="/rss.xml" class="btn btn-sm btn-light">RSS 订阅</a>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from "vue";
import { absoluteUrl, applySeo } from "@/core/seo";
import AuthorAvatar from "@/components/common/AuthorAvatar.vue";
import { resolveMediaUrl } from "@/utils/media";
import { useAuthStore, type PublicProfile } from "@/stores/auth";
import { useSiteStore } from "@/stores/site";

const siteStore = useSiteStore();
const authStore = useAuthStore();
const site = computed(() => siteStore.setting);
const profile = computed<PublicProfile>(() => authStore.publicProfile || { trueName: "", email: "", pic: "", bio: "" });
const authorName = computed(() => profile.value.trueName || site.value.author || site.value.siteName);
onMounted(async () => {
  const setting = await siteStore.loadSiteSetting();
  const contact = await authStore.refreshPublicProfile().catch(() => undefined);
  const name = contact?.trueName || setting.author || setting.siteName;
  const description = contact?.bio || setting.description;
  const canonical = absoluteUrl("/about");
  const avatar = resolveMediaUrl(contact?.pic);
  applySeo({
    title: `关于 - ${setting.siteName}`,
    description,
    keywords: ["关于", name, setting.siteName, setting.keywords].filter(Boolean).join(","),
    canonical,
    siteName: setting.siteName,
    type: "website",
    image: avatar || undefined,
    jsonLd: [
      {
        "@context": "https://schema.org",
        "@type": "AboutPage",
        name: `关于 ${setting.siteName}`,
        url: canonical,
        description: setting.description,
        mainEntity: {
          "@type": "Person",
          name,
          url: canonical,
          description: contact?.bio || undefined,
          email: contact?.email || undefined,
          image: avatar || undefined,
        },
      },
    ],
  });
});
</script>

<style scoped>
.about-avatar {
  width: 70px;
  height: 70px;
  object-fit: cover;
  border-radius: 50%;
}
</style>
