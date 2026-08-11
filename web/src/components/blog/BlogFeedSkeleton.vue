<template>
  <div class="card feed-skeleton" aria-hidden="true">
    <div class="card-body pb-5 placeholder-glow">
      <div class="d-flex align-items-center mb-5">
        <span class="placeholder rounded-circle skeleton-avatar me-5"></span>
        <div class="flex-grow-1">
          <span class="placeholder col-7 d-block mb-3"></span>
          <span class="placeholder col-4 d-block"></span>
        </div>
      </div>
      <span class="placeholder skeleton-media rounded d-block mb-5"></span>
      <span class="placeholder col-12 d-block mb-3"></span>
      <span class="placeholder col-10 d-block mb-5"></span>
      <div class="d-flex gap-4">
        <span v-for="item in 3" :key="item" class="placeholder skeleton-action rounded"></span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.feed-skeleton { overflow: hidden; }
.placeholder { background-color: var(--kt-gray-300); }
.skeleton-avatar { width: 45px; height: 45px; flex: 0 0 45px; }
.skeleton-media { width: 100%; height: 220px; }
.skeleton-action { width: 72px; height: 34px; }
</style>
