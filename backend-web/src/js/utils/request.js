import axios from 'axios'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

import { ApiError } from "./errors.js"

const service = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 5000,
})

service.interceptors.request.use(
  config => {
    config.params = {
      ...config.params,
      _platform: 'web',
      // _version: '1.0'
    }

    /*
    if (['post', 'put', 'patch'].includes(config.method?.toLowerCase())) {
      config.headers['Content-Type'] = 'application/json'
    }
    */

    const token = localStorage.getItem("token")

    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }

    return config
  },
  err => {
    // return Promise.reject(err)
    console.log(`!!! axios interceptors error: ${err.message}`);
  }
)

function handleStatus(err) {
  const status = err.response?.status;
  const data = err.response?.data || {};

  if (status >= 500) {
    ElMessage.error(`Internal Server Error ${status}`)
    return
  }

  switch(status) {
    case 400:
      ElMessage.error(`Bad Request ${status}: ${data.msg}`)
      break;
    case 401:
      ElMessage.error(`Please login again ${status}`)
      localStorage.clear()
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
      return
  }

  return
};

service.interceptors.response.use(
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
      console.log(`!!! Got an ApiError: status=${err.response.status}, data=${JSON.stringify(err.response.data)}`);
    } if (axios.isAxiosError(err) && !err.response) { // Unexpected
      // both code and kind are empty, "Network Error"...
      err = new ApiError("", err.message, { status: 0, kind: "", raw: err })
    } else if(axios.isAxiosError(err)) {
      handleStatus(err);
      console.log(`!!! Got an AxiosError: status=${err.response.status}, data=${JSON.stringify(err.response.data)}`);

      let data = err.response.data;
      let details = { status: err.response.status, kind: data.kind, requestId: data.requestId, raw: err };

      err = new ApiError(data.code, data.msg, details)
    } else { // Unexpected
      console.log(`!!! Got UnknownError: ${err}`)
      err = new ApiError("", `Unknown Error: ${err.message}`, { status: err.response.status, kind: "", raw: err })
    }

    /*
    if (err instanceof TypeError && err.message.startsWith("NetworkError")) {
      ElMessage.error("Network error");
      err = new ApiError("NetworkError", "network error", { raw: err });
    } else if (err instanceof TypeError) {
      err = new ApiError("TypeError", "type error", { raw: err })
    } else if (err instanceof SyntaxError) {
      ElMessage.warn(`SyntaxError`);
      err = new ApiError("SyntaxError", "syntax error", { raw: err })
    }
    */

    return Promise.reject(err)
  }
)

export { service }

/*
//
import request from '@/utils/request'

export function createUser(data) {
  service.post('/api/users',
    { name: 'John', age: 30 },   // json body
    { params: { debug: true } }, // query parameters
  )
}

createUser({ firstname: "Jane", "lastname": "Doe", password: "123abcABC" })
.then(data => {
  console.log('--> created user:', data)
})
.catch(err => {
  console.error('!!! failed to create user:', err)
})

//
export default {
  user: {
    login(data) {
      return service.post('/auth/login', data)
    },
    getInfo() {
      return service.get('/auth/info')
    },
    logout() {
      return service.post('/auth/logout')
    },
  },

  product: {
    list(params) {
      return service.get('/products', { params })
    },
    detail(id) {
      return service.get(`/products/${id}`)
    },
    create(data) {
      return service.post('/products', data)
    }
  }
}
*/
