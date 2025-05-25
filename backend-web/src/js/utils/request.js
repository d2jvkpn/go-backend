import axios from 'axios'
import qs from 'qs'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'

import { ApiError } from "./errors.js"
import { clearAccount } from "../stores/storage.js";


const router = useRouter();

const service = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  paramsSerializer: params => qs.stringify(params, { arrayFormat: 'repeat' }),
  timeout: 5000,
})

service.interceptors.request.use(
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
      clearAccount();
      router.push('/login');
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
      err = new ApiError("", `Unknown Error: ${err.msg}`, { status: err.response.status, kind: "", raw: err })
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

    /*
    if (error.response) {
      const { status, data } = error.response
      if (status === 401) {
        ElMessage.error('Authentication failed, please login again')
        clearAccount()
        // go to login page
      } else if (status === 400) {
        ElMessage.error(data.message || 'Invalid request')
      } else {
        ElMessage.error(data.message || 'Failed to change password')
      }
    } else if (error.request) {
      ElMessage.error('Network error, please try again later')
    } else {
      if (error.message) {
        ElMessage.error(error.message)
      }
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
