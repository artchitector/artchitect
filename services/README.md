### Artchitect services

**Main services**
- **soul** - core service, running on server with GPU. Soul handles creation process and service operations.
- **gate** - api-gateway service. Connected with soul via database and redis.
- **vision** - nuxt.js (vue.js) frontend service, SSR. connected to gate via http-api and websocket.

**Util services**
- **saver** - api to save and get images (for image-storage-servers)