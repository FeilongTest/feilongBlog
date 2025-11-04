<template>
  <div
    class="engage-toolbar d-flex position-fixed px-5 fw-bold zindex-2 top-50 end-0 transform-90 mt-20 gap-2"
  >
    <!--begin::Links drawer toggle-->
    <button
      id="kt_engage_demos_toggle"
      class="engage-demos-toggle btn btn-flex h-35px bg-body btn-color-gray-700 btn-active-color-gray-900 shadow-sm fs-6 px-4 rounded-top-0"
      title="友情链接"
      data-bs-toggle="tooltip"
      data-bs-placement="left"
      data-bs-dismiss="click"
      data-bs-trigger="hover"
    >
      <i class="bi bi-link-45deg fs-4"></i>
      <span id="kt_engage_demos_label">友情链接</span>
    </button>
    <!--end::Links drawer toggle-->

  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "kt-toolbar-buttons",
  components: {},
});
</script>
