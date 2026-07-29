# 网络唤醒
前端基于vue3+tailwind css+pinia制作<br>
后端基于go1.26，支持unix socket+baseurl访问<br>
构建步骤
```
cd wol
npm install
npm run build
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o wol
```
启动环境
```
export WOL_SOCKET="path-to-app.sock"
export WOL_BASE_URL="/example"
./wol -d path-to-storage-devices.json
```
