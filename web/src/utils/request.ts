import axios, { type AxiosError } from "axios";
import { ElMessage } from "element-plus";
import JwtService from "@/core/services/JwtService";

const service = axios.create({
  baseURL: import.meta.env.VITE_BLOG_API_URL || "/blog",
  timeout: 15_000,
});

service.interceptors.request.use((config) => {
  const token = JwtService.getToken();
  if (token) config.headers.set("x-token", token);
  return config;
});

service.interceptors.response.use(
  (response) => {
    const newToken = response.headers["new-token"];
    if (newToken) JwtService.saveToken(newToken);

    if (response.data.code === 0 || response.headers.success === "true") {
      if (response.headers.msg) response.data.msg = decodeURI(response.headers.msg);
      return response.data;
    }

    if (response.data?.data?.reload) JwtService.destroyToken();
    const message = response.data?.msg || "请求失败";
    ElMessage.error(message);
    return Promise.reject(new Error(message));
  },
  (error: AxiosError<{ msg?: string }>) => {
    if (error.response?.status === 401) JwtService.destroyToken();
    const status = error.response?.status;
    const message = error.response?.data?.msg
      || (status ? `请求失败（HTTP ${status}）` : "无法连接服务器，请检查网络或后端服务");
    ElMessage.error(message);
    return Promise.reject(error);
  },
);

export default service;
