<template>
  <div ref="editorElement" class="markdown-editor"></div>
</template>

<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from "vue";
import Vditor from "vditor";
import "vditor/dist/index.css";
import service from "@/utils/request";
import { resolveMediaUrl } from "@/utils/media";
import JwtService from "@/core/services/JwtService";

const props = defineProps<{ modelValue: string }>();
const emit = defineEmits<{ "update:modelValue": [value: string] }>();
const editorElement = ref<HTMLElement>();
let editor: Vditor | undefined;
let applyingExternalValue = false;
let ready = false;

const resolveUploadedUrl = resolveMediaUrl;

const uploadImages = async (files: File[]) => {
  for (const file of files) {
    const formData = new FormData();
    formData.append("file", file);
    const res: any = await service.post("/admin/file/upload", formData, {
      headers: { "x-token": JwtService.getToken() || "" },
    });
    if (res.code !== 0) return res.msg || "图片上传失败";
    const url = resolveUploadedUrl(res.data.url);
    editor?.insertValue(`![${file.name}](${url})`);
  }
  return null;
};

onMounted(async () => {
  await nextTick();
  if (!editorElement.value) return;
  editor = new Vditor(editorElement.value, {
    cdn: `${window.location.origin}/vendor/vditor`,
    after: () => { ready = true; editor?.setValue(props.modelValue || "", true); },
    height: 500,
    mode: "ir",
    value: props.modelValue || "",
    cache: { enable: false },
    placeholder: "使用 Markdown 编写文章，支持代码块、表格和图片……",
    counter: { enable: true },
    outline: { enable: true, position: "left" },
    preview: {
      delay: 150,
      hljs: { enable: true, lineNumber: true, style: "github" },
    },
    upload: {
      accept: "image/*",
      multiple: true,
      handler: uploadImages,
    },
    input: (value) => {
      if (!applyingExternalValue) emit("update:modelValue", value);
    },
  });
});

watch(
  () => props.modelValue,
  (value) => {
    if (!editor || !ready || editor.getValue() === value) return;
    applyingExternalValue = true;
    editor.setValue(value || "", true);
    applyingExternalValue = false;
  },
);

onBeforeUnmount(() => {
  editor?.destroy();
  editor = undefined;
});
</script>

<style scoped>
.markdown-editor {
  min-height: 500px;
  border-radius: 0.475rem;
  overflow: hidden;
}
</style>
