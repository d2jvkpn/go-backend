import axios from 'axios'

const service = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 5000,
})

service.interceptors.request.use(
  config => {
    config.params = {
      ...config.params,
      _platform: 'web',
      _version: '1.0'
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

service.interceptors.response.use(
  response => {
    if (response.data.code === "ok") {
      return response.data.data
    } else {
      return Promise.reject(response.data.msg)
    }
  },
  err => {
    return Promise.reject(err)
  }
)

export default service
