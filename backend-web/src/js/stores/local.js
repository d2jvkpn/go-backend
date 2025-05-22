export function setAccount(data) {
  localStorage.setItem("id", data.id);
  localStorage.setItem("firstname", data.firstname);
  localStorage.setItem("lastname", data.lastname);
  localStorage.setItem("token", data.token);
  localStorage.setItem("level", data.level);
}

export function getAccount() {
  return {
    id: localStorage.getItem("id"),
    firstname: localStorage.getItem("firstname"),
    lastname: localStorage.getItem("lastname"),
    token: localStorage.getItem("token"),
    level: localStorage.getItem("level"),
  }
}

export function clearAccount() {
  //localStorage.clear()
  localStorage.removeItem("id");
  localStorage.removeItem("firstname");
  localStorage.removeItem("lastname");
  localStorage.removeItem("token");
  localStorage.removeItem("level");
}
