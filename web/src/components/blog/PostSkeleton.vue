<template>
  <div class="card" aria-hidden="true">
    <div class="card-body p-lg-20 pb-lg-0 placeholder-glow">
      <div class="d-flex flex-column flex-xl-row">
        <div class="flex-lg-row-fluid me-xl-15">
          <div class="d-flex gap-5 mb-8">
            <span v-for="item in 3" :key="item" class="placeholder rounded skeleton-meta"></span>
          </div>
          <span class="placeholder col-7 d-block skeleton-title mb-8"></span>
          <span class="placeholder rounded skeleton-cover d-block mb-10"></span>
          <div class="mb-14">
            <span v-for="item in 6" :key="item" class="placeholder d-block mb-4"
              :class="item % 3 === 0 ? 'col-8' : 'col-12'"></span>
          </div>
          <div class="border border-dashed card-rounded p-8 d-flex align-items-center mb-14">
            <span class="placeholder rounded-circle skeleton-author me-7"></span>
            <div class="flex-grow-1">
              <span class="placeholder col-12 d-block mb-3"></span>
              <span class="placeholder col-9 d-block"></span>
            </div>
          </div>
        </div>
        <div class="flex-column flex-lg-row-auto w-100 w-xl-300px mb-10">
          <span class="placeholder col-5 d-block mb-7"></span>
          <span class="placeholder rounded skeleton-search d-block mb-16"></span>
          <span class="placeholder col-6 d-block mb-7"></span>
          <span v-for="item in 4" :key="item" class="placeholder col-12 d-block mb-5"></span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.placeholder { background-color: var(--kt-gray-300); }
.skeleton-meta { width: 130px; height: 18px; }
.skeleton-title { height: 30px; }
.skeleton-cover { width: 100%; height: 350px; }
.skeleton-author { width: 70px; height: 70px; flex: 0 0 70px; }
.skeleton-search { width: 100%; height: 44px; }
@media (max-width: 767.98px) {
  .skeleton-cover { height: 230px; }
  .skeleton-meta { width: 82px; }
}
</style>
