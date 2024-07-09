docker build -t localhost:5000/shisha-server-amd64:latest -f Dockerfile-server --platform linux/amd64 .

docker buildx create --use --name multi-builder --platform linux/arm64,linux/amd64

docker buildx build \
 --tag truebad0ur/shisha-server:latest \
 --platform linux/amd64,linux/arm64 \
 --builder multi-builder --push \
 -f Dockerfile-server .