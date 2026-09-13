#!/usr/bin/env bash

set -euo pipefail

repository_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
container_engine=${CONTAINER_ENGINE:-podman}
image_prefix=${SMOKE_IMAGE_PREFIX:-localhost/eks-}
image_tag=${SMOKE_IMAGE_TAG:-smoke}
run_id="eks-smoke-$$"
network="$run_id"
artifacts=$(mktemp -d)
containers=()

services=(
  adservice cartservice checkoutservice currencyservice emailservice frontend
  loadgenerator paymentservice productcatalogservice recommendationservice
  shippingservice
)

cleanup() {
  if ((${#containers[@]})); then
    "$container_engine" rm --force "${containers[@]}" >/dev/null 2>&1 || true
  fi
  "$container_engine" network rm "$network" >/dev/null 2>&1 || true
  rm -rf "$artifacts"
}
trap cleanup EXIT

image_reference() {
  printf '%s%s:%s' "$image_prefix" "$1" "$image_tag"
}

build_images() {
  local service context
  for service in "${services[@]}"; do
    context="$repository_root/services/$service"
    if [[ "$service" == cartservice ]]; then
      context="$repository_root/services/cartservice/src"
    fi
    "$container_engine" build --tag "$(image_reference "$service")" "$context"
  done
}

start_container() {
  local service=$1
  shift
  local name="$run_id-$service"
  "$container_engine" run --detach --rm \
    --name "$name" \
    --network "$network" \
    --network-alias "$service" \
    "$@" \
    "$(image_reference "$service")" >/dev/null
  containers+=("$name")
}

assert_containers_running() {
  local name running
  for name in "${containers[@]}"; do
    running=$("$container_engine" inspect --format '{{.State.Running}}' "$name")
    if [[ "$running" != true ]]; then
      "$container_engine" logs "$name" >&2 || true
      printf 'container stopped before the storefront was ready: %s\n' "$name" >&2
      return 1
    fi
  done
}

wait_for_storefront() {
  local url=$1
  local attempt
  for attempt in {1..45}; do
    if curl --fail --silent --show-error \
      --header 'Cookie: shop_session-id=integration-smoke' \
      "$url/_healthz" >/dev/null 2>&1; then
      return 0
    fi
    assert_containers_running
    sleep 1
  done
  printf 'storefront did not become ready: %s\n' "$url" >&2
  return 1
}

run_checkout() {
  local product_id=0PUK6V6EV0
  local response="$artifacts/checkout.html"

  curl --fail --silent --show-error \
    --header 'Cookie: shop_session-id=integration-smoke' \
    --data-urlencode "product_id=$product_id" \
    --data-urlencode 'quantity=1' \
    "$storefront_url/cart" >/dev/null

  curl --fail --silent --show-error \
    --header 'Cookie: shop_session-id=integration-smoke' \
    --data-urlencode 'email=someone@example.com' \
    --data-urlencode 'street_address=1600 Amphitheatre Parkway' \
    --data-urlencode 'zip_code=94043' \
    --data-urlencode 'city=Mountain View' \
    --data-urlencode 'state=CA' \
    --data-urlencode 'country=United States' \
    --data-urlencode 'credit_card_number=4432-8015-6152-0454' \
    --data-urlencode 'credit_card_expiration_month=1' \
    --data-urlencode 'credit_card_expiration_year=2039' \
    --data-urlencode 'credit_card_cvv=672' \
    "$storefront_url/cart/checkout" >"$response"

  grep -q 'Your order is complete!' "$response"
}

run_load_traffic() {
  local output="$artifacts/loadgenerator.log"
  "$container_engine" run --rm \
    --network "$network" \
    --env FRONTEND_ADDR=frontend:8080 \
    --env USERS=2 \
    --env RATE=1 \
    --env LOCUST_RUN_TIME=20s \
    --env LOCUST_EXIT_CODE_ON_ERROR=1 \
    "$(image_reference loadgenerator)" | tee "$output"

  if ! grep -Eq 'Aggregated[[:space:]]+[1-9][0-9]*[[:space:]]+0\(0\.00%\)' "$output"; then
    printf 'load generator did not report a nonzero, zero-failure aggregate\n' >&2
    return 1
  fi
}

build_images
"$container_engine" network create "$network" >/dev/null

redis_image=$(awk '/image: redis:/ {print $2; exit}' "$repository_root/deploy/eks/deployment-service.yml")
redis_name="$run_id-redis-cart"
"$container_engine" run --detach --rm \
  --name "$redis_name" \
  --network "$network" \
  --network-alias redis-cart \
  "$redis_image" >/dev/null
containers+=("$redis_name")

start_container productcatalogservice --env PORT=3550 --env DISABLE_PROFILER=1
start_container currencyservice --env PORT=7000 --env DISABLE_PROFILER=1
start_container cartservice --env REDIS_ADDR=redis-cart:6379
start_container emailservice --env PORT=5000 --env DISABLE_PROFILER=1
start_container paymentservice --env PORT=50051 --env DISABLE_PROFILER=1
start_container shippingservice --env PORT=50051 --env DISABLE_PROFILER=1
start_container recommendationservice \
  --env PORT=8080 \
  --env PRODUCT_CATALOG_SERVICE_ADDR=productcatalogservice:3550 \
  --env DISABLE_PROFILER=1
start_container adservice --env PORT=9555
start_container checkoutservice \
  --env PORT=5050 \
  --env PRODUCT_CATALOG_SERVICE_ADDR=productcatalogservice:3550 \
  --env SHIPPING_SERVICE_ADDR=shippingservice:50051 \
  --env PAYMENT_SERVICE_ADDR=paymentservice:50051 \
  --env EMAIL_SERVICE_ADDR=emailservice:5000 \
  --env CURRENCY_SERVICE_ADDR=currencyservice:7000 \
  --env CART_SERVICE_ADDR=cartservice:7070
start_container frontend \
  --publish 127.0.0.1::8080 \
  --env PORT=8080 \
  --env PRODUCT_CATALOG_SERVICE_ADDR=productcatalogservice:3550 \
  --env CURRENCY_SERVICE_ADDR=currencyservice:7000 \
  --env CART_SERVICE_ADDR=cartservice:7070 \
  --env RECOMMENDATION_SERVICE_ADDR=recommendationservice:8080 \
  --env SHIPPING_SERVICE_ADDR=shippingservice:50051 \
  --env CHECKOUT_SERVICE_ADDR=checkoutservice:5050 \
  --env AD_SERVICE_ADDR=adservice:9555 \
  --env ENABLE_PROFILER=0 \
  --env ENV_PLATFORM=local

frontend_name="$run_id-frontend"
frontend_port=$("$container_engine" port "$frontend_name" 8080/tcp | awk -F: 'NR == 1 {print $NF}')
storefront_url="http://127.0.0.1:$frontend_port"
wait_for_storefront "$storefront_url"

curl --fail --silent --show-error \
  --header 'Cookie: shop_session-id=integration-smoke' \
  "$storefront_url/" >"$artifacts/storefront.html"
grep -q 'Sunglasses' "$artifacts/storefront.html"

run_checkout
run_load_traffic
printf 'storefront, checkout, and load traffic passed at %s\n' "$storefront_url"
