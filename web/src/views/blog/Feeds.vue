<template>
    <!--begin::Row-->
    <div class="row g-5 g-xl-8">

        <!-- 第一列双数 -->
        <!--begin::Col-->
        <div class="col-xl-6">
            <div v-for="(article, index) in articleList">
                <blogFeeds v-if="index % 2 === 0" widget-classes="mb-5 mb-xl-8" :article="article">
                </blogFeeds>
            </div>
        </div>
        <!--end::Col-->

        <!-- 第一列单数 -->
        <!--begin::Col-->
        <div class="col-xl-6">
            <div v-for="(article, index) in articleList">
                <blogFeeds v-if="index % 2 === 1" widget-classes="mb-5 mb-xl-8" :article="article">
                </blogFeeds>
            </div>
        </div>
        <!--end::Col-->
    </div>
    <!--end::Row-->

    <!--begin::Pagination-->
    <div class="d-flex flex-stack flex-wrap pt-10">
        <div class="fs-6 fw-semobold text-gray-700">
            共 {{ total }} 篇文章
        </div>
        <TablePagination v-if="total > 1" :total-pages="Math.ceil(total / 10)" :total="total" :per-page="10"
            :current-page="page" @page-change="pageChange" />

    </div>


    <!--end::Pagination-->
</template>
<script setup lang="ts">
import blogFeeds from "@/components/blog/BlogFeeds.vue"
import { ref, onMounted, watch } from "vue"
import service from "@/utils/request";
import TablePagination from "@/components/kt-datatable/table-partials/table-content/table-footer/TablePagination.vue";
import { useRouter } from "vue-router";
import type { ArticleList } from "@/core/blog/ArticleTypes"



const page = ref(1);
const total = ref(0);
const pageSize = ref(10);
const articleList = ref<Array<ArticleList>>([]);
const fid = ref(0);
const keyword = ref("");
const router = useRouter();


onMounted(() => {
    //判断当前是否有分类打开
    const currentRoute = router.currentRoute.value;
    if (Object.keys(currentRoute.params).length !== 0) {
        fid.value = Number(currentRoute.params.id);
    }
    // 获取搜索关键字
    if (currentRoute.query.keyword) {
        keyword.value = String(currentRoute.query.keyword);
    }
    getArticleList();
});

// 监听当前路由
watch(
    () => router.currentRoute.value,
    (newValue: any, oldValue: any) => {        
        if(newValue.name == "category"){
            page.value = 1;
            total.value = 0;            
            fid.value = newValue.params.id;
            keyword.value = "";
            // 如果 URL 中有 query 参数，清理它
            if (newValue.query && Object.keys(newValue.query).length > 0) {
                router.replace({
                    name: "category",
                    params: { id: newValue.params.id }
                });
            }
            getArticleList();
        }else if(newValue.name == "blog-home"){            
            page.value = 1;
            total.value = 0;            
            fid.value = 0;
            keyword.value = newValue.query.keyword ? String(newValue.query.keyword) : "";
            getArticleList();
        } else if(newValue.name == "blog-content") {
            // 从搜索页面跳转到文章详情页时，清理 query 参数
            if (newValue.query && Object.keys(newValue.query).length > 0) {
                router.replace({
                    name: "blog-content",
                    params: newValue.params
                });
            }
        }
    },
)


const getArticleList = () => {
    service({
        url: "/base/getArticleList",
        method: "get",
        params: {
            page: page.value,
            pageSize: pageSize.value,
            fid: fid.value,
            keyword: keyword.value || undefined
        }
    }).then(({data}) => {
        if(data !== undefined){
            total.value = data.total;
            articleList.value = data.list;
        }    
    })
}

const pageChange = (newPage: number) => {
    page.value = newPage;
    getArticleList();
}

</script>
