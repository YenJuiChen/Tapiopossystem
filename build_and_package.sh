#!/bin/bash

set -e

# === 設定參數 ===
IMAGE_NAME="hcj-fdg-pos"
IMAGE_TAG="latest"
OUTPUT_DIR="deploy_package"
TAR_NAME="hcj-fdg-pos-image.tar"
REMOTE_USER="ec2-user"
REMOTE_HOST="ec2-13-208-243-217.ap-northeast-3.compute.amazonaws.com"
PEM_PATH="/Users/linminze/fcj-fdg-pos.pem"

# === 清理舊資料 ===
rm -rf ${OUTPUT_DIR}
mkdir -p ${OUTPUT_DIR}

# === 清除舊 image 和 builder cache（選擇性）===
echo "🧹 清除舊的 Docker 緩存（如需）..."
docker builder prune -f || true
docker rmi ${IMAGE_NAME}:${IMAGE_TAG} || true

# === 建立映像檔 ===
echo "🚧 建立 Docker 映像（平台：linux/amd64）..."
docker build --platform linux/amd64 -t ${IMAGE_NAME}:${IMAGE_TAG} .

# === 匯出為 .tar 檔 ===
echo "📦 匯出映像為 tar 檔..."
docker save -o ${OUTPUT_DIR}/${TAR_NAME} ${IMAGE_NAME}:${IMAGE_TAG}

# === 複製部署檔案 ===
echo "📁 複製部署相關檔案..."
cp docker-compose.yml ${OUTPUT_DIR}/

# === 壓縮部署包 ===
echo "🎁 壓縮 deploy_package..."
tar -czf ${OUTPUT_DIR}.tar.gz -C ${OUTPUT_DIR} .

# === 傳送到遠端 ===
echo "🚀 傳送 deploy_package 到遠端機器..."
scp -i "${PEM_PATH}" ${OUTPUT_DIR}.tar.gz ${REMOTE_USER}@${REMOTE_HOST}:~

echo "✅ 傳送完成！請登入遠端執行 ./deploy.sh 完成部署"
