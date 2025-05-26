import { ElMessage } from 'element-plus'

export function hello() {
  ElMessage.info("Hello!")
}

/*
Promise.all([
  fetch('/api/users'),
  fetch('/api/products')
]).then(([usersRes, productsRes]) => {
  callback()
});
*/

// import { hello } from "@/js/_hello.js"
// hello()
