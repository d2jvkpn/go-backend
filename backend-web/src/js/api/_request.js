import { request } from '@/js/api'

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
