function setAccount(data) {
  localStorage.setItem("id", data.id);
  localStorage.setItem("firstname", data.firstname);
  localStorage.setItem("lastname", data.lastname);
  localStorage.setItem("level", data.level);
  localStorage.setItem("token", data.token);
  localStorage.setItem("expiresAt", data.expiresAt);
}

function getAccount() {
  return {
    id: localStorage.getItem("id"),
    firstname: localStorage.getItem("firstname"),
    lastname: localStorage.getItem("lastname"),
    level: localStorage.getItem("level"),
    token: localStorage.getItem("token"),
    expiresAt: Number(localStorage.getItem("expiresAt")),
  }
}

function clearAccount() {
  //localStorage.clear()
  localStorage.removeItem("id");
  localStorage.removeItem("firstname");
  localStorage.removeItem("lastname");
  localStorage.removeItem("level");
  localStorage.removeItem("token");
  localStorage.removeItem("expiresAt");
}

function checkIsLoggedIn() {
  if (!localStorage.getItem("token")) {
    return false;
  }

  const now = Date.now()/1000;
  const expiresAt = Number(localStorage.getItem("expiresAt"));

  if (expiresAt <= now) {
    return false;
  }

  return true;
}

export default {
  setAccount,
  getAccount,
  clearAccount,
  checkIsLoggedIn,
}
