<template>
  <main class="auth-shell">
    <section class="auth-form-side">
      <router-link :to="{ name: 'blog-home' }" class="mobile-brand">
        <span class="brand-dot"></span>
        <span>studio.</span>
      </router-link>

      <div class="auth-form-wrap">
        <router-view />
      </div>

      <footer class="auth-footer">© {{ currentYear }} 飞龙博客管理后台</footer>
    </section>

    <aside class="auth-brand-side">
      <div class="brand-glow brand-glow-one"></div>
      <div class="brand-glow brand-glow-two"></div>
      <div class="brand-grid"></div>

      <div class="brand-content">
        <router-link :to="{ name: 'blog-home' }" class="brand-lockup">
          <span class="brand-mark">F</span>
          <span>studio.</span>
        </router-link>

        <div class="brand-copy">
          <span class="brand-kicker">FEILONG'S BLOG</span>
          <h1>记录想法，<br />沉淀每一次创造。</h1>
          <p>在一个安静、清晰的空间里管理文章、动态、评论与站点数据。</p>
        </div>

        <div class="brand-features">
          <span><i class="bi bi-journal-text"></i> 内容管理</span>
          <span><i class="bi bi-bar-chart"></i> 数据统计</span>
          <span><i class="bi bi-cloud-arrow-up"></i> R2 存储</span>
        </div>
      </div>
    </aside>
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted } from "vue";
import LayoutService from "@/core/services/LayoutService";
import { useBodyStore } from "@/stores/body";

const bodyStore = useBodyStore();
const currentYear = computed(() => new Date().getFullYear());

onMounted(() => {
  LayoutService.emptyElementClassesAndAttributes(document.body);
  bodyStore.addBodyClassname("app-blank");
  bodyStore.addBodyClassname("bg-body");
});
</script>

<style scoped>
.auth-shell { min-height: 100vh; display: grid; grid-template-columns: minmax(420px, 46%) 1fr; background: #fff; }
.auth-form-side { position: relative; display: flex; min-width: 0; flex-direction: column; padding: 2.5rem clamp(2rem, 5vw, 6rem); }
.auth-form-wrap { display: flex; flex: 1; align-items: center; justify-content: center; }
.auth-footer { color: #99a1b7; font-size: .875rem; text-align: center; }
.mobile-brand { display: none; align-items: center; gap: .65rem; color: #181c32; font-size: 1.5rem; font-weight: 800; }
.brand-dot { width: .7rem; height: .7rem; border-radius: 50%; background: #009ef7; box-shadow: 0 0 0 .35rem rgba(0, 158, 247, .12); }
.auth-brand-side { position: relative; min-height: 100vh; overflow: hidden; color: #fff; background: linear-gradient(145deg, #111827 0%, #172554 55%, #312e81 100%); }
.brand-grid { position: absolute; inset: 0; opacity: .16; background-image: linear-gradient(rgba(255,255,255,.16) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,.16) 1px, transparent 1px); background-size: 54px 54px; mask-image: linear-gradient(to bottom, #000, transparent 88%); }
.brand-glow { position: absolute; border-radius: 50%; filter: blur(2px); }
.brand-glow-one { width: 28rem; height: 28rem; top: -12rem; right: -8rem; background: rgba(0, 158, 247, .3); }
.brand-glow-two { width: 22rem; height: 22rem; bottom: -10rem; left: -7rem; background: rgba(124, 58, 237, .34); }
.brand-content { position: relative; z-index: 1; display: flex; height: 100%; flex-direction: column; justify-content: space-between; padding: clamp(3rem, 6vw, 7rem); }
.brand-lockup { display: inline-flex; width: fit-content; align-items: center; gap: .85rem; color: #fff; font-size: 2rem; font-weight: 800; letter-spacing: -.04em; }
.brand-mark { display: grid; width: 2.75rem; height: 2.75rem; place-items: center; border: 1px solid rgba(255,255,255,.25); border-radius: .85rem; background: rgba(255,255,255,.1); backdrop-filter: blur(12px); font-size: 1.15rem; }
.brand-copy { max-width: 42rem; padding: 5rem 0; }
.brand-kicker { color: #7dd3fc; font-size: .78rem; font-weight: 700; letter-spacing: .22em; }
.brand-copy h1 { margin: 1.5rem 0; color: #fff; font-size: clamp(2.7rem, 5vw, 5.5rem); font-weight: 750; line-height: 1.08; letter-spacing: -.055em; }
.brand-copy p { max-width: 34rem; margin: 0; color: rgba(255,255,255,.68); font-size: 1.05rem; line-height: 1.9; }
.brand-features { display: flex; flex-wrap: wrap; gap: 1rem 1.75rem; color: rgba(255,255,255,.72); font-size: .9rem; }
.brand-features span { display: inline-flex; align-items: center; gap: .55rem; }
.brand-features i { color: #7dd3fc; }
@media (max-width: 991.98px) {
  .auth-shell { grid-template-columns: 1fr; }
  .auth-brand-side { display: none; }
  .auth-form-side { min-height: 100vh; padding: 2rem clamp(1.25rem, 7vw, 4rem); }
  .mobile-brand { display: inline-flex; }
}
</style>
