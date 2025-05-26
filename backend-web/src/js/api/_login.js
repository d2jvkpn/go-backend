import { ElMessage } from 'element-plus'


export function login(data, callback, onError) {
  if (!onError) {
    onError = () => {}
  }

  fetch(`${import.meta.env.VITE_API_URL}/api/v1/open/account/login?_platform=web`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  .then((response) => {
    if (response.status == 200) {
      return response.json()
    }

    if (response.status > 500) {
      ElMessage.error(`Internal Server Error: ${response.status}`)
    } else if (response.status == 400) {
      ElMessage.error(`Bad Request: 400`)
    } else {
      console.log(`!!! Response Error: status=${response.status}`)
    }

    return Promise.reject(new Error(`status=${response.status}`));
    // throw new Error(`status=${response.status}`);
  })
  .then((res) => {
    if (res.code == "ok") {
      callback(res.data);
    } else {
      ElMessage.warn(res.msg)
      onError()
    }
  })
  .catch(err => {
    if (err instanceof TypeError && err.message.startsWith("NetworkError")) {
      ElMessage.error("Network error");
      console.log(`!!! Network error: ${err}`)
    } else if (err instanceof SyntaxError) {
      ElMessage.warn(`SyntaxError error`);
      console.log(`!!! SyntaxError error: ${err}`)
    } else {
      console.log(`!!! Unknown error: ${err}`)
    }

    onError()
  })
  // .finally( () => { TODO })
}

export async function login2(data, callback, onError) {
  if (!onError) {
    onError = () => {}
  }

  try {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/api/v1/open/account/login?platform=web`);
    if (response.status !== 200) {
      throw new Error(`status=${response.status}`);
    }

    const res = await response.json();
    callback(res.data)
  } catch (err) {
    console.error(`!!! Requets Error: ${err}`);
  }
}

export function logout(callback, onError) {
  if (!onError) {
    onError = () => {}
  }

  fetch(`${import.meta.env.VITE_API_URL}/api/v1/auth/account/logout`, {
    method: 'POST',
    headers: { 'Authorization': `Bearer ${localStorage.getItem("token")}` },
  })
  .then((response) => {
    if (response.status == 200) {
      return response.json()
    }

    if (response.status > 500) {
      ElMessage.error(`Internal Server Error: ${response.status}`)
    } else if (response.status == 400) {
      ElMessage.error(`Bad Request: 400`)
    } else {
      console.log(`!!! Response Error: status=${response.status}`)
    }

    return Promise.reject(new Error(`status=${response.status}`));
    // throw new Error(`status=${response.status}`);
  })
  .then((res) => {
    if (res.code == "ok") {
      callback();
    } else {
      ElMessage.warn(res.msg)
      onError()
    }
  })
  .catch(err => {
    if (err instanceof TypeError && err.message.startsWith("NetworkError")) {
      ElMessage.error("Network error");
      console.log(`!!! Network error: ${err}`)
    } else if (err instanceof SyntaxError) {
      ElMessage.warn(`SyntaxError error`);
      console.log(`!!! SyntaxError error: ${err}`)
    } else {
      console.log(`!!! Unknown error: ${err}`)
    }

    onError()
  })
  // .finally( () => { TODO })
}
