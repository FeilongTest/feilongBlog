<template>
    <div class="card">
    <!--begin::Body-->
    <div class="card-body p-lg-20 pb-lg-0">
        <!--begin::Layout-->
        <div class="d-flex flex-column flex-xl-row">
            <!--begin::Content-->
            <div class="flex-lg-row-fluid me-xl-15">
                <!--begin::Post content-->
                <div class="mb-17">
                    <!--begin::Wrapper-->
                    <div class="mb-8">
                        <!--begin::Info-->
                        <div class="d-flex flex-wrap mb-6">
                            <!--begin::Item-->
                            <div class="me-9 my-1">
                                <!--begin::Icon-->
                                <span class="svg-icon svg-icon-primary me-1 fs-2">
                                    <inline-svg
                                        :src="
                                        getAssetPath('/media/icons/duotune/general/gen016.svg')
                                        "
                                    />
                                </span>
                                <!--end::Icon-->
                                <!--begin::Label-->
                                <span class="fw-bold text-gray-400">{{ moment.unix(Number(articleInfo.ctime)).format("YYYY年MM月DD日 HH时mm分") }}</span>
                                <!--end::Label-->
                            </div>
                            <!--end::Item-->
                            <!--begin::Item-->
                            <div class="me-9 my-1">
                                <!--begin::Icon-->
                                <span class="svg-icon svg-icon-primary me-1 fs-2">
                                    <inline-svg
                                        :src="getAssetPath('/media/icons/duotune/art/art005.svg')"
                                    />
                                </span>
                                <!--end::Icon-->
                                <!--begin::Label-->
                                <span class="fw-bold text-gray-400">{{ moment.unix(Number(articleInfo.editTime)).format("YYYY年MM月DD日 HH时mm分") }}</span>
                                <!--begin::Label-->
                            </div>
                            <!--end::Item-->
                            <!--begin::Item-->
                            <div class="me-9 my-1">
                                <!--begin::Icon-->
                                <span class="svg-icon svg-icon-primary me-1 fs-2">
                                    <!-- /media/icons/duotune/communication/com003.svg评论图标 --> 
                                    <inline-svg
                                        :src="getAssetPath('/media/icons/duotune/coding/cod006.svg')"
                                    />
                                </span>
                                <!--end::Icon-->
                                <!--begin::Label-->
                                <span class="fw-bold text-gray-400">{{ articleInfo.view }}</span>
                                <!--end::Label-->
                            </div>
                            <!--end::Item-->
                        </div>
                        <!--end::Info-->
                        <!--begin::Title-->
                        <span class="text-dark text-hover-primary fs-2 fw-bold">{{ articleInfo.title }}
                        </span>
                        <!--end::Title-->
                        <!--begin::Container-->
                        <div v-if="articleInfo.type !== 0 && articleInfo.pic" class="overlay mt-8">
                            <!--begin::Image-->
                            <div class="bgi-no-repeat bgi-position-center bgi-size-cover card-rounded min-h-350px" :style="{'backgroundImage':'url('+articleInfo.pic+')'}"></div>
                            <!--end::Image-->
                        </div>
                        <!--end::Container-->
                    </div>
                    <!--end::Wrapper-->
                    <!--begin::Description-->
                    <div class="fs-5 fw-semibold text-gray-600">
                        <div class="content mb-8" v-html="sanitizedContent"></div>
                    </div>
                    <!--end::Description-->
                    <!--begin::Block-->
                    <div class="d-flex align-items-center border-1 border-dashed card-rounded p-5 p-lg-10 mb-14">
                        <!--begin::Section-->
                        <div class="text-center flex-shrink-0 me-7 me-lg-13">
                            <!--begin::Avatar-->
                            <div class="symbol symbol-70px symbol-circle mb-2">
                                <img :src="getAssetPath('/media/avatars/300-1.jpg')" class="" alt="" />
                            </div>
                            <!--end::Avatar-->
                            <!--begin::Info-->
                            <div class="mb-0">
                                <span class="text-gray-700 fw-bold">飞龙</span>
                                <span class="text-gray-400 fs-7 fw-semibold d-block mt-1">feilong</span>
                            </div>
                            <!--end::Info-->
                        </div>
                        <!--end::Section-->
                        <!--begin::Text-->
                        <div class="mb-0 fs-6">
                            <div class="text-muted fw-semibold lh-lg mb-2">免责声明:
                            本站提供的一切软件、教程和内容信息仅限用于学习和研究目的;不得将上述内容用于商业或者非法用途,否则,一切后果请用户自负。本站信息来自网络收集整理,版权争议与本站无关。您必须在下载后的24个小时之内,从您的电脑或手机中彻底删除上述内容。如果您喜欢该程序和内容，请支持正版，购买注册，得到更好的正版服务。我们非常重视版权问题，如有侵权请邮件与我们联系处理。敬请谅解！</div>
                            <!-- <a href="../../demo1/dist/pages/user-profile/overview.html" class="fw-semibold link-primary">作者资料</a> -->
                        </div>
                        <!--end::Text-->
                    </div>
                    <!--end::Block-->
                </div>
                <!--end::Post content-->

                <!--begin::Comment Section-->
                <div class="mt-15">
                    <CommentSection v-if="articleInfo.ID" :article-id="articleInfo.ID" />
                </div>
                <!--end::Comment Section-->
            </div>
            <!--end::Content-->
            <!--begin::Sidebar-->
            <div class="flex-column flex-lg-row-auto w-100 w-xl-300px mb-10">
                <!--begin::Search blog-->
                <div class="mb-16">
                    <h4 class="text-dark mb-7">搜索文章</h4>
                    <!--begin::Input group-->
                    <div class="position-relative">
                        <span 
                            class="svg-icon svg-icon-2 svg-icon-1 svg-icon-gray-500 position-absolute top-50 translate-middle-y ms-3"
                        >
                            <inline-svg
                                :src="getAssetPath('/media/icons/duotune/general/gen021.svg')"
                            />
                        </span>
                        <input 
                            type="text" 
                            class="form-control form-control-solid ps-10" 
                            name="search" 
                            v-model="searchKeyword"
                            @keyup.enter="handleSearch"
                            placeholder="输入文章标题关键字搜索" 
                        />
                        <button 
                            v-if="searchKeyword"
                            @click="clearSearch"
                            class="btn btn-sm btn-icon btn-active-light-primary position-absolute end-0 top-50 translate-middle-y me-2"
                            type="button"
                        >
                            <i class="bi bi-x-lg fs-6"></i>
                        </button>
                    </div>
                    <!--end::Input group-->
                </div>
                <!--end::Search blog-->
                <!--begin::Catigories-->
                <div class="mb-16">
                    <h4 class="text-dark mb-7">相关分类</h4>
                    
                    <!--begin::Item-->
                    <template v-for="item in articleSummary" :key="item.fid">
                        <div class="d-flex flex-stack fw-semibold fs-5 text-muted mb-4">
                            <!--begin::Text-->
                            <router-link class="text-muted text-hover-primary pe-2" :to="{ name: 'category', params: { id: item.fid } }">
                                <span class="badge badge-light-warning fw-bold my-2 me-3">推荐板块</span>{{item.name}}
                            </router-link>
                            <!--end::Text-->
                            <!--begin::Number-->
                            <div class="m-0">
                                共
                                <span class="badge badge-light-info fw-bold my-2">{{item.total}}</span>
                                篇
                            </div>
                            <!--end::Number-->
                        </div>
                    </template>
                    <!--end::Item-->
                </div>
                <!--end::Catigories-->
                <!--begin::Recent posts-->
                <div class="m-0">
                    <h4 class="text-dark mb-7">最新发布</h4>
                    <template v-for="item in articleList" :key="item.ID">
                        <!--begin::Item-->
                        <div class="d-flex mb-7">
                            <!--begin::Symbol-->
                            <div class="symbol symbol-60px symbol-2by3 me-4">
                                <div class="symbol-label" :style="{'backgroundImage':'url('+item.pic+')'}"></div>
                            </div>
                            <!--end::Symbol-->
                            <!--begin::Title-->
                            <div class="m-0 align-self-center" >
                                <span @click="gotoContent('/content/'+ item.fid + '/'+ item.ID)" class="text-dark fw-bold text-hover-primary fs-6 m-width=100 pe-4">{{ item.title }}</span>
                            </div>
                            <!--end::Title-->
                        </div>
                        <!--end::Item-->
                    </template>
                </div>
                <!--end::Recent posts-->
            </div>
            <!--end::Sidebar-->
        </div>
        <!--end::Layout-->
    </div>
    <!--end::Body-->
</div>
</template>
<script lang="ts">
import { getAssetPath } from "@/core/helpers/assets";
import { useRouter } from "vue-router";
import { computed, ref,onMounted,watch } from "vue"
import service from "@/utils/request";
import type { ArticleList,ArticelSummary } from "@/core/blog/ArticleTypes";
import dayjs from "dayjs";
import DOMPurify from "dompurify";
import CommentSection from "@/components/blog/CommentSection.vue";

//弹窗
import Swal from "sweetalert2";

 

export default {
    name: "widget-1",
    components: {
        CommentSection,
    },
    props: {
        widgetClasses: String,
    },
    setup() {
        const router = useRouter();
        const articleInfo = ref<ArticleList>({} as ArticleList);
        const articleSummary = ref<Array<ArticelSummary>>([]);
        const articleList = ref<Array<ArticleList>>([]);
        const searchKeyword = ref("");
		const sanitizedContent = computed(() => DOMPurify.sanitize(articleInfo.value.content || ""));


        onMounted(()=>{
            const currentRoute = router.currentRoute.value;
            getArticle(Number(currentRoute.params.id));
            getArticleList();
        })


        // 监听当前路由
        watch(
            () => router.currentRoute.value,
            (newValue: any) => {
                if(newValue.name == "blog-content"){
                    // 如果 URL 中有 query 参数，清理它
                    if (newValue.query && Object.keys(newValue.query).length > 0) {
                        router.replace({
                            name: "blog-content",
                            params: newValue.params
                        });
                    }
                    getArticle(Number(newValue.params.id))
                }
            },
        )




    
        const getArticle = async(id:number) => {
            service.post("/base/getArticle",{id:id})
            .then((res:any) => {
                if(res.code == 7){
                    Swal.fire({
                        title:"出错了",
                        icon:"error",
                        text: "没有找到文章信息,即将返回主页面!",
                        buttonsStyling: false,
                        confirmButtonText: "Return Home",
                        heightAuto: false,
                        customClass: {
                            confirmButton: "btn fw-semobold btn-light-primary",
                        },
                    }).then(() => {
                        // Go to page after successfully login
                        router.push({ name: "blog-home" });
                    });
                    return;
                }
                articleInfo.value = res.data;
                // 动态更新页面标题
                const appName = import.meta.env.VITE_APP_NAME || "Feilong'S Blog";
                document.title = `${articleInfo.value.title} - ${appName}`;
                getSummary(articleInfo.value.fid);
            })
        }

        const getSummary = async(fid:number) => {
            service.post("/base/getSummary",{id:fid})
            .then((res:any) => {
                articleSummary.value = res.data
            })
        }

        const getArticleList = async() => {
            service({
                url: "/base/getArticleList",
                method: "get",
                params: {
                    page: 1,
                    pageSize: 5,
                }
            }).then(({data}) => {
                if(data !== undefined){
                    articleList.value = data.list;
                }    
            })
        }

        const gotoContent = (path:string) => {
            setTimeout(()=>{
                router.replace(path);
            })
        }

        const handleSearch = () => {
            if (searchKeyword.value.trim()) {
                router.push({ 
                    name: "blog-home", 
                    query: { keyword: searchKeyword.value.trim() } 
                });
            }
        }

        const clearSearch = () => {
            searchKeyword.value = "";
            router.push({ name: "blog-home" });
        }

        return {
            getAssetPath,
            articleInfo,
            articleSummary,
            articleList,
            moment: dayjs,
			sanitizedContent,
            gotoContent,
            searchKeyword,
            handleSearch,
            clearSearch,
        }
    }
}
</script>

<style scoped lang="scss">
// 因为生成出来的 有p标签或多层标签包裹,img就在深层里面 :deep() 穿透解决
.content{
    :deep(img){
        max-width: 100%;
        width: 100%;
        height: 100%;
    }
}

</style>
