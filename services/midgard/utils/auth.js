export default function loggedIn() {
  if (!process.client) {
    return false
  } else {
    return localStorage.getItem("token")
  }
}
