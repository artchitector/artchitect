export default function ({$axios, store}) {
  $axios.onRequest((req) => {
    if (process.server) {
      return
    }
    req.headers.common['Authorization'] = localStorage.getItem("token")
  })
}
