<template>
    <div class="card" :class="widgetClasses">
        <div class="card-body p-lg-20">
            <Section1/>
            <Section2/>
            <Section3/>
            <Section4/>
        </div>
    </div>
</template>
<script lang="ts">
import { getAssetPath } from "@/core/helpers/assets";
import Section1 from "@/components/blog/Section1.vue";
import Section2 from "@/components/blog/Section2.vue";
import Section3 from "@/components/blog/Section3.vue";
import Section4 from "@/components/blog/Section4.vue";

export default {
    name: "blog-home",
    props: {
        widgetClasses: String,
    },
    components: {
        Section1,
        Section2,
        Section3,
        Section4
    },
    setup() {
        return {
            getAssetPath,
        }
    }
}
</script>
