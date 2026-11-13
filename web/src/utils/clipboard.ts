/** 返回是否复制成功，拒绝权限和非安全上下文均有明确结果。 */
export async function copyText(text: string): Promise<boolean> {
  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(text);
      return true;
    }
  } catch { /* 尝试浏览器兼容方案。 */ }
  const previous = document.activeElement as HTMLElement | null;
  const input = document.createElement("textarea");
  input.value = text;
  input.setAttribute("readonly", "");
  input.style.cssText = "position:fixed;left:-9999px;top:0";
  document.body.appendChild(input);
  input.select();
  try { return document.execCommand("copy"); }
  catch { return false; }
  finally { input.remove(); previous?.focus({ preventScroll: true }); }
}
