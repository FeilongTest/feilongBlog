<template>
    <a class="d-block overlay" data-fslightbox="lightbox-hot-sales"
        :href="getImage">
        <!--begin::Image-->
        <div class="overlay-wrapper bgi-no-repeat bgi-position-center bgi-size-cover card-rounded min-h-175px"
            :style="{ 'backgroundImage': 'url(' + getImage + ')' }"></div>
        <!--end::Image-->
        <!--begin::Action-->
        <div class="overlay-layer card-rounded bg-dark bg-opacity-25">
            <i class="bi bi-eye fs-2x text-white" />
        </div>
        <!--end::Action-->
    </a>
</template>
<script lang="ts">
import { getAssetPath } from "@/core/helpers/assets";
import { computed } from "@vue/reactivity";
export default {
    name: "BlogPhoto",
    props: {
        image: String,
    },
    setup(props) {
        const getImage = computed(() => {
           return props.image?props.image:''
        })
        return {
            getAssetPath,
            getImage
        }
    }
}
</script>