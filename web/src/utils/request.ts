import axios from 'axios' // 引入axios
import JwtService from "@/core/services/JwtService";
import { ElMessage, ElMessageBox } from 'element-plus'

// 确保环境变量有值，如果没有则使用默认值
const baseURL = import.meta.env.VITE_BLOG_API_URL || '/blog';

// 开发环境调试信息
if (import.meta.env.DEV) {
  console.log('API Base URL:', baseURL);
  console.log('Environment variables:', {
    VITE_BLOG_API_URL: import.meta.env.VITE_BLOG_API_URL,
    MODE: import.meta.env.MODE,
    DEV: import.meta.env.DEV
  });
}

const service = axios.create({
  baseURL: baseURL,   
  // baseURL: getBaseUrl(),
  timeout: 99999
})

// http request 拦截器
service.interceptors.request.use(
  config => {
    // 开发环境打印请求信息
    if (import.meta.env.DEV) {
      console.log('Request URL:', (config.baseURL || '') + (config.url || ''));
      console.log('Request Method:', config.method);
    }
    config.headers = {
      'Content-Type': 'application/json',
      'x-token': `${JwtService.getToken()}`,
      ...config.headers
    }
    return config
  },
  error => {
    ElMessage({
      showClose: true,
      message: error,
      type: 'error'
    })
    return error
  }
)

// http response 拦截器
service.interceptors.response.use(
  response => {    
    if (response.data.code === 0 || response.headers.success === 'true') {
      if (response.headers.msg) {
        response.data.msg = decodeURI(response.headers.msg)
      }
      return response.data
    } else {
      ElMessage({
        showClose: true,
        message: response.data.msg,
        type: 'error'
      })
      return response.data.msg ? response.data : response
    }
  },
  error => {
    // 开发环境打印错误信息
    if (import.meta.env.DEV) {
      console.error('Request Error:', error);
      console.error('Error Message:', error.message);
      console.error('Error Code:', error.code);
      if (error.config) {
        console.error('Failed URL:', (error.config.baseURL || '') + (error.config.url || ''));
      }
    }
    if (!error.response) {
      ElMessageBox.confirm(`
        <p>检测到请求错误</p>
        <p>${error.message || error}</p>
        <p>请检查：</p>
        <ul>
          <li>后端服务是否运行在 http://localhost:8889</li>
          <li>网络连接是否正常</li>
          <li>代理配置是否正确</li>
        </ul>
        `, '请求报错', {
        dangerouslyUseHTMLString: true,
        distinguishCancelAndClose: true,
        confirmButtonText: '稍后重试',
        cancelButtonText: '取消'
      })
      return
    }

    switch (error.response.status) {
      case 500:
        ElMessageBox.confirm(`
        <p>检测到接口错误${error}</p>
        <p>错误码<span style="color:red"> 500 </span>：此类错误内容常见于后台panic，请先查看后台日志，如果影响您正常使用可强制登出清理缓存</p>
        `, '接口报错', {
          dangerouslyUseHTMLString: true,
          distinguishCancelAndClose: true,
          confirmButtonText: '清理缓存',
          cancelButtonText: '取消'
        })
          .then(() => {
          })
        break
      case 404:
        ElMessageBox.confirm(`
          <p>检测到接口错误${error}</p>
          <p>错误码<span style="color:red"> 404 </span>：此类错误多为接口未注册（或未重启）或者请求路径（方法）与api路径（方法）不符--如果为自动化代码请检查是否存在空格</p>
          `, '接口报错', {
          dangerouslyUseHTMLString: true,
          distinguishCancelAndClose: true,
          confirmButtonText: '我知道了',
          cancelButtonText: '取消'
        })
        break
    }

    return error
  }
)

export default service