<template>
  <div
    ref="editorElement"
    class="markdown-editor"
    @input="flush"
    @compositionend="flush"
    @focusout="flush"
  ></div>
</template>

<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from "vue";
import Vditor from "vditor";
import { ElNotification } from "element-plus";
import "vditor/dist/index.css";
import service from "@/utils/request";
import { resolveMediaUrl } from "@/utils/media";

const props = defineProps<{ modelValue: string }>();
const emit = defineEmits<{
  "update:modelValue": [value: string];
  uploading: [value: boolean];
}>();
const editorElement = ref<HTMLElement>();
let editor: Vditor | undefined;
let applyingExternalValue = false;
let ready = false;

// Vditor's input callback is debounced with its undo stack. Read its live
// document at editing/action boundaries instead of waiting to persist user input.
function flush() {
  if (editor && ready && !applyingExternalValue) {
    emit("update:modelValue", editor.getValue());
  }
}
defineExpose({ flush });

const resolveUploadedUrl = resolveMediaUrl;

let uploads = 0;
const uploadImages = async (files: File[]) => {
  if (files.some((file) => file.size > 10 * 1024 * 1024)) {
    ElNotification({
      title: "图片过大",
      message: "单张图片不能超过 10 MB",
      type: "warning",
    });
    return null;
  }
  uploads++;
  emit("uploading", true);
  try {
    for (const file of files) {
      const formData = new FormData();
      formData.append("file", file);
      const res = await service.post("/admin/file/upload", formData);
      const url = resolveUploadedUrl(res.data.url);
      editor?.insertValue(`![${file.name}](${url})`);
      flush();
    }
    return null;
  } catch {
    // The request interceptor reports the failure; release the editor controls.
    return null;
  } finally {
    uploads--;
    emit("uploading", uploads > 0);
  }
};

onMounted(async () => {
  await nextTick();
  if (!editorElement.value) return;
  editor = new Vditor(editorElement.value, {
    cdn: `${window.location.origin}/vendor/vditor`,
    after: () => {
      ready = true;
      editor?.setValue(props.modelValue || "", true);
      editorElement.value
        ?.querySelectorAll<HTMLElement>('[contenteditable="true"]')
        .forEach((element) => {
          element.setAttribute("role", "textbox");
          element.setAttribute("aria-label", "Markdown 正文");
          element.setAttribute("aria-multiline", "true");
        });
    },
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
  }
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
