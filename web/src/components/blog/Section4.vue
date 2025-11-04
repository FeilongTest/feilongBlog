<template>
    <!--begin::latest instagram-->
    <div class="">
        <!--begin::Section-->
        <div class="m-0">
            <!--begin::Content-->
            <div class="d-flex flex-stack">
                <!--begin::Title-->
                <h3 class="text-dark">Latest Instagram Posts</h3>
                <!--end::Title-->
                <!--begin::Link-->
                <a href="#" class="fs-6 fw-semibold link-primary">View Instagram</a>
                <!--end::Link-->
            </div>
            <!--end::Content-->
            <!--begin::Separator-->
            <div class="separator separator-dashed border-gray-300 mb-9 mt-5"></div>
            <!--end::Separator-->
        </div>
        <!--end::Section-->
        <!--begin::Row-->
        <div class="row g-10 row-cols-2 row-cols-lg-5">
            <!--begin::Col-->
            <div class="col">
                <!--begin::Overlay-->
                <BlogPhoto :image="getAssetPath('/media/stock/900x600/16.jpg')"/>
                <!--end::Overlay-->
            </div>
            <!--end::Col-->
            <!--begin::Col-->
            <div class="col">
                <!--begin::Overlay-->
                <BlogPhoto :image="getAssetPath('/media/stock/900x600/13.jpg')"/>
                <!--end::Overlay-->
            </div>
            <!--end::Col-->
            <!--begin::Col-->
            <div class="col">
                <!--begin::Overlay-->
                <BlogPhoto :image="getAssetPath('/media/stock/900x600/19.jpg')"/>
                <!--end::Overlay-->
            </div>
            <!--end::Col-->
            <!--begin::Col-->
            <div class="col">
                <!--begin::Overlay-->
                <BlogPhoto :image="getAssetPath('/media/stock/900x600/15.jpg')"/>
                <!--end::Overlay-->
            </div>
            <!--end::Col-->
            <!--begin::Col-->
            <div class="col">
                <!--begin::Overlay-->
                <BlogPhoto :image="getAssetPath('/media/stock/900x600/12.jpg')"/>
                <!--end::Overlay-->
            </div>
            <!--end::Col-->
        </div>
        <!--begin::Row-->
    </div>
    <!--end::latest instagram-->
</template>
<script lang="ts">
import BlogPhoto from "@/components/blog/BlogPhoto.vue"
import { getAssetPath } from "@/core/helpers/assets";
export default {
    name: "section4",
    components: {
        BlogPhoto,
    },
    setup() {
        return {
            getAssetPath,
        }
    }
}
</script>