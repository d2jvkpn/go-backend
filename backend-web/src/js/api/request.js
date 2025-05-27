import axios from 'axios'
import qs from 'qs'
import { ElMessage } from 'element-plus'

import { ApiError } from "../types/errors.js"
import stores from "../stores";


const request = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  paramsSerializer: params => qs.stringify(params, { arrayFormat: 'repeat' }),
  timeout: 3000,
})

request.interceptors.request.use(
  config => {
    /*
    config.params = {
      ...config.params,
      _platform: 'web',
      // _version: '1.0'
    }
    */

    let version = localStorage.getItem("version");
    let env = localStorage.getItem("env");
    config.headers['x-client'] = `app=backend-web; version=${version}; env=${env}`

    /*
    if (['post', 'put', 'patch'].includes(config.method?.toLowerCase())) {
      config.headers['Content-Type'] = 'application/json'
    }
    */

    console.debug(`--> ${config.method}@${config.url}, ${JSON.stringify(config.params)}`)

    const token = localStorage.getItem("token")
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    // return Promise.reject(new Error(`!!! no token`))
    // return Promise.reject(new ApiError("", "no token", {}))

    return config
  },
  err => {
    // return Promise.reject(err)
    console.log(`!!! axios interceptors error: ${err.message}`);
  }
)

function handleStatus(err) {
  const status = err.response?.status || 0;
  const data = err.response?.data || {};

  if (status >= 500) {
    ElMessage.error(`Internal Server Error ${status}`)
    return
  }

  switch (status) {
  case 400:
    ElMessage.error(`Bad Request ${status}: ${data.msg}`)
    break;
  case 401:
    ElMessage.error(`Please login again ${status}`)
    stores.clearAccount();
    window.location.href = `${import.meta.env.VITE_BASE_PATH}/login`;
    break;
  case 403:
    ElMessage.error(`Permission denied ${status}: ${data.msg}`)
    break;
  case 404:
    ElMessage.error(`Not route ${status}`);
    break;
  case 409:
    ElMessage.error(`Biz error ${status}: ${data.msg}`);
    break
  case 429:
    ElMessage.error(`Too many requests ${status}`);
    break;
  default:
    ElMessage.error(`Error ${status}: ${err.message}, code=${data.code}, kind=${data.kind}`);
  }

  return
};

request.interceptors.response.use(
  response => {
    if (response.status == 200 && response.data.code === "ok") {
      return response.data.data;
    }

   return Promise.reject(new ApiError(
      response.data.code, response.data.msg,
      { kind: response.data.kind, status: response.status, requestId: response.requestId },
   ))
  },
  err => { // always throw an ApiError
    if (err instanceof ApiError) {
      const { status, data: res } = err.response;
      console.log(`!!! ApiError: status=${status}, data=${JSON.stringify(res)}`);
    } else if (axios.isAxiosError(err) && !err.response) {
      /*
      if (err.code === 'ECONNABORTED') {
        console.error(`!!! 请求超时: ${err.code}, ${err.message}`)
      } else if (err.message === 'Network Error') {
        console.error(`!!! 网络错误, 请检查你的连接: ${err.code}, ${err.message}`)
      } else {
        console.error(`!!! 未知网络错误: ${err.code}, ${err.message}`)
      }
      */
      // both code and kind are empty, "Network Error"...
      err = new ApiError("", err.message, { status: 0, kind: "", raw: err })
      ElMessage.error(`!!! AxiosError request error: ${err.message}`);
    } else if (axios.isAxiosError(err)) {
      handleStatus(err);
      const { status, data: res } = err.response;
      console.log(`!!! GAxiosError respoonse error: status=${status}, data=${JSON.stringify(res)}`);

      err = new ApiError(
        res.code, res.msg,
        { status: status, kind: res.kind, requestId: res.requestId, raw: err },
      )
    } else { // Unexpected
      console.log(`!!! UnknownError: ${err}`)
      err = new ApiError("", `Unknown Error: ${err.msg}`, { status: err.response.status, kind: "", raw: err })
    }
    return Promise.reject(err)
  }
)

export default request
