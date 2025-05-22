export function setAccount(data) {
  localStorage.setItem("id", data.id);
  localStorage.setItem("firstname", data.firstname);
  localStorage.setItem("lastname", data.lastname);
  localStorage.setItem("level", data.level);
  localStorage.setItem("token", data.token);
  localStorage.setItem("expiresAt", data.expiresAt);
}

export function getAccount() {
  return {
    id: localStorage.getItem("id"),
    firstname: localStorage.getItem("firstname"),
    lastname: localStorage.getItem("lastname"),
    level: localStorage.getItem("level"),
    token: localStorage.getItem("token"),
    expiresAt: localStorage.getItem("expiresAt"),
  }
}

export function clearAccount() {
  //localStorage.clear()
  localStorage.removeItem("id");
  localStorage.removeItem("firstname");
  localStorage.removeItem("lastname");
  localStorage.removeItem("level");
  localStorage.removeItem("token");
  localStorage.removeItem("expiresAt");
}
