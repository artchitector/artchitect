# export NODE_OPTIONS=--openssl-legacy-provider
npm install
npm run build
pm2 start
pm2 ls
