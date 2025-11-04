<template>
    <!--begin::Section-->
    <div class="mb-17">
        <!--begin::Content-->
        <div class="d-flex flex-stack mb-5">
            <!--begin::Title-->
            <h3 class="text-dark">Hottest Bundles</h3>
            <!--end::Title-->
            <!--begin::Link-->
            <a href="#" class="fs-6 fw-semibold link-primary">View All Offers</a>
            <!--end::Link-->
        </div>
        <!--end::Content-->
        <!--begin::Separator-->
        <div class="separator separator-dashed mb-9"></div>
        <!--end::Separator-->
        <!--begin::Row-->
        <div class="row g-10">
            <!--begin::Col-->
            <div class="col-md-4">
                <!--begin::Hot sales post-->
                <div class="card-xl-stretch me-md-6">
                    <!--begin::Overlay-->
                    <BlogPhoto :image="getAssetPath('/media/stock/600x400/img-23.jpg')"/>
                    <!--end::Overlay-->
                    <!--begin::Body-->
                    <div class="mt-5">
                        <!--begin::Title-->
                        <a href="#" class="fs-4 text-dark fw-bold text-hover-primary text-dark lh-base">25 Products
                            Mega Bundle with 50% off discount amazing</a>
                        <!--end::Title-->
                        <!--begin::Text-->
                        <div class="fw-semibold fs-5 text-gray-600 text-dark mt-3">We’ve been focused on making a
                            the from also not been eye</div>
                        <!--end::Text-->
                        <!--begin::Text-->
                        <div class="fs-6 fw-bold mt-5 d-flex flex-stack">
                            <!--begin::Label-->
                            <span class="badge border border-dashed fs-2 fw-bold text-dark p-2">
                                <span class="fs-6 fw-semibold text-gray-400">$</span>28</span>
                            <!--end::Label-->
                            <!--begin::Action-->
                            <a href="#" class="btn btn-sm btn-primary">Purchase</a>
                            <!--end::Action-->
                        </div>
                        <!--end::Text-->
                    </div>
                    <!--end::Body-->
                </div>
                <!--end::Hot sales post-->
            </div>
            <!--end::Col-->
            <!--begin::Col-->
            <div class="col-md-4">
                <!--begin::Hot sales post-->
                <div class="card-xl-stretch mx-md-3">
                    <!--begin::Overlay-->
                    <BlogPhoto :image="getAssetPath('/media/stock/600x600/img-14.jpg')"/>
                    <!--end::Overlay-->
                    <!--begin::Body-->
                    <div class="mt-5">
                        <!--begin::Title-->
                        <a href="#" class="fs-4 text-dark fw-bold text-hover-primary text-dark lh-base">25 Products
                            Mega Bundle with 50% off discount amazing</a>
                        <!--end::Title-->
                        <!--begin::Text-->
                        <div class="fw-semibold fs-5 text-gray-600 text-dark mt-3">We’ve been focused on making a
                            the from also not been eye</div>
                        <!--end::Text-->
                        <!--begin::Text-->
                        <div class="fs-6 fw-bold mt-5 d-flex flex-stack">
                            <!--begin::Label-->
                            <span class="badge border border-dashed fs-2 fw-bold text-dark p-2">
                                <span class="fs-6 fw-semibold text-gray-400">$</span>27</span>
                            <!--end::Label-->
                            <!--begin::Action-->
                            <a href="#" class="btn btn-sm btn-primary">Purchase</a>
                            <!--end::Action-->
                        </div>
                        <!--end::Text-->
                    </div>
                    <!--end::Body-->
                </div>
                <!--end::Hot sales post-->
            </div>
            <!--end::Col-->
            <!--begin::Col-->
            <div class="col-md-4">
                <!--begin::Hot sales post-->
                <div class="card-xl-stretch ms-md-6">
                    <!--begin::Overlay-->
                    <BlogPhoto :image="getAssetPath('/media/stock/600x400/img-71.jpg')"/>
                    <!--end::Overlay-->
                    <!--begin::Body-->
                    <div class="mt-5">
                        <!--begin::Title-->
                        <a href="#" class="fs-4 text-dark fw-bold text-hover-primary text-dark lh-base">25 Products
                            Mega Bundle with 50% off discount amazing</a>
                        <!--end::Title-->
                        <!--begin::Text-->
                        <div class="fw-semibold fs-5 text-gray-600 text-dark mt-3">We’ve been focused on making a
                            the from also not been eye</div>
                        <!--end::Text-->
                        <!--begin::Text-->
                        <div class="fs-6 fw-bold mt-5 d-flex flex-stack">
                            <!--begin::Label-->
                            <span class="badge border border-dashed fs-2 fw-bold text-dark p-2">
                                <span class="fs-6 fw-semibold text-gray-400">$</span>25</span>
                            <!--end::Label-->
                            <!--begin::Action-->
                            <a href="#" class="btn btn-sm btn-primary">Purchase</a>
                            <!--end::Action-->
                        </div>
                        <!--end::Text-->
                    </div>
                    <!--end::Body-->
                </div>
                <!--end::Hot sales post-->
            </div>
            <!--end::Col-->
        </div>
        <!--end::Row-->
    </div>
    <!--end::Section-->
</template>
<script lang="ts">
import BlogPhoto from "@/components/blog/BlogPhoto.vue"
import { getAssetPath } from "@/core/helpers/assets";
export default {
    name: "section3",
    components: {
        BlogPhoto,
    },
    props: {
        widgetClasses: String,
    },
    setup() {
        return {
            getAssetPath,
        }
    }
}
</script>