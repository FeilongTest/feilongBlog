<template>
  <nav
    aria-label="文章分页"
    class="col-sm-12 col-md-7 d-flex align-items-center justify-content-center justify-content-md-end"
  >
    <div class="dataTables_paginate paging_simple_numbers">
      <ul class="pagination">
        <li
          class="paginate_button page-item"
          :class="{ disabled: isInFirstPage }"
          :style="{ cursor: !isInFirstPage ? 'pointer' : 'auto' }"
        >
          <button
            type="button"
            class="page-link"
            aria-label="第一页"
            :disabled="isInFirstPage"
            @click="onClickFirstPage"
          >
            <span class="svg-icon">
              <inline-svg
                :src="getAssetPath('/media/icons/duotune/arrows/arr079.svg')"
              />
            </span>
          </button>
        </li>

        <li
          class="paginate_button page-item"
          :class="{ disabled: isInFirstPage }"
          :style="{ cursor: !isInFirstPage ? 'pointer' : 'auto' }"
        >
          <button
            type="button"
            class="page-link"
            aria-label="上一页"
            :disabled="isInFirstPage"
            @click="onClickPreviousPage"
          >
            <span class="svg-icon">
              <inline-svg
                :src="getAssetPath('/media/icons/duotune/arrows/arr074.svg')"
              />
            </span>
          </button>
        </li>

        <li
          v-for="(page, i) in pages"
          class="paginate_button page-item"
          :class="{
            active: isPageActive(page.name),
          }"
          :style="{ cursor: !page.isDisabled ? 'pointer' : 'auto' }"
          :key="i"
        >
          <button
            type="button"
            class="page-link"
            :aria-label="`第 ${page.name} 页`"
            :aria-current="isPageActive(page.name) ? 'page' : undefined"
            @click="onClickPage(page.name)"
          >
            {{ page.name }}
          </button>
        </li>

        <li
          class="paginate_button page-item"
          :class="{ disabled: isInLastPage }"
          :style="{ cursor: !isInLastPage ? 'pointer' : 'auto' }"
        >
          <button
            type="button"
            class="paginate_button page-link"
            aria-label="下一页"
            :disabled="isInLastPage"
            @click="onClickNextPage"
          >
            <span class="svg-icon">
              <inline-svg
                :src="getAssetPath('/media/icons/duotune/arrows/arr071.svg')"
              />
            </span>
          </button>
        </li>

        <li
          class="paginate_button page-item"
          :class="{ disabled: isInLastPage }"
          :style="{ cursor: !isInLastPage ? 'pointer' : 'auto' }"
        >
          <button
            type="button"
            class="paginate_button page-link"
            aria-label="最后一页"
            :disabled="isInLastPage"
            @click="onClickLastPage"
          >
            <span class="svg-icon">
              <inline-svg
                :src="getAssetPath('/media/icons/duotune/arrows/arr080.svg')"
              />
            </span>
          </button>
        </li>
      </ul>
    </div>
  </nav>
</template>

<script lang="ts">
import { getAssetPath } from "@/core/helpers/assets";
import { defineComponent, computed } from "vue";

export default defineComponent({
  name: "table-pagination",
  props: {
    maxVisibleButtons: {
      type: Number,
      required: false,
      default: 5,
    },
    totalPages: {
      type: Number,
      required: true,
    },
    total: {
      type: Number,
      required: true,
    },
    perPage: {
      type: Number,
      required: true,
    },
    currentPage: {
      type: Number,
      required: true,
    },
  },
  emits: ["page-change"],
  setup(props, { emit }) {
    const startPage = computed(() => {
      if (
        props.totalPages < props.maxVisibleButtons ||
        props.currentPage === 1 ||
        props.currentPage <= Math.floor(props.maxVisibleButtons / 2) ||
        (props.currentPage + 2 > props.totalPages &&
          props.totalPages === props.maxVisibleButtons)
      ) {
        return 1;
      }

      if (props.currentPage + 2 > props.totalPages) {
        return props.totalPages - props.maxVisibleButtons + 1;
      }

      return props.currentPage - 2;
    });

    const endPage = computed(() => {
      return Math.min(
        startPage.value + props.maxVisibleButtons - 1,
        props.totalPages
      );
    });

    const pages = computed(() => {
      const range: Array<{
        name: number;
        isDisabled: boolean;
      }> = [];

      for (let i = startPage.value; i <= endPage.value; i += 1) {
        range.push({
          name: i,
          isDisabled: i === props.currentPage,
        });
      }

      return range;
    });

    const isInFirstPage = computed(() => {
      return props.currentPage <= 1;
    });
    const isInLastPage = computed(() => {
      return props.currentPage >= props.totalPages;
    });

    const onClickFirstPage = () => {
      if (!isInFirstPage.value) emit("page-change", 1);
    };
    const onClickPreviousPage = () => {
      if (!isInFirstPage.value) emit("page-change", props.currentPage - 1);
    };
    const onClickPage = (page: number) => {
      emit("page-change", page);
    };
    const onClickNextPage = () => {
      if (!isInLastPage.value) emit("page-change", props.currentPage + 1);
    };
    const onClickLastPage = () => {
      if (!isInLastPage.value) emit("page-change", props.totalPages);
    };
    const isPageActive = (page: number) => {
      return props.currentPage === page;
    };

    return {
      startPage,
      endPage,
      pages,
      isInFirstPage,
      isInLastPage,
      onClickFirstPage,
      onClickPreviousPage,
      onClickPage,
      onClickNextPage,
      onClickLastPage,
      isPageActive,
      getAssetPath,
    };
  },
});
</script>
