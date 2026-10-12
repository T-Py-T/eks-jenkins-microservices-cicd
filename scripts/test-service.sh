#!/usr/bin/env bash

set -euo pipefail

repository_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
service=${1:-}

case "$service" in
  adservice)
    (cd "$repository_root/services/adservice" && ./gradlew --no-daemon build)
    ;;
  cartservice)
    dotnet test "$repository_root/services/cartservice/tests/cartservice.tests.csproj" --configuration Release
    ;;
  checkoutservice|frontend|productcatalogservice|shippingservice)
    (cd "$repository_root/services/$service" && go test ./...)
    ;;
  currencyservice)
    (cd "$repository_root/services/currencyservice" && npm ci && npm audit --omit=dev && npm test)
    node "$repository_root/tests/node_service_smoke.js" \
      "$repository_root/services/currencyservice" server.js \
      "CurrencyService gRPC server started" 17000
    ;;
  paymentservice)
    (cd "$repository_root/services/paymentservice" && npm ci && npm audit --omit=dev && npm test)
    node "$repository_root/tests/node_service_smoke.js" \
      "$repository_root/services/paymentservice" index.js \
      "PaymentService gRPC server started" 15000
    ;;
  emailservice|loadgenerator)
    python3 -m compileall -q "$repository_root/services/$service"
    python3 -m pip_audit -r "$repository_root/services/$service/requirements.txt"
    ;;
  recommendationservice)
    python3 -m compileall -q "$repository_root/services/$service"
    python3 -m pip_audit -r "$repository_root/services/$service/requirements.txt"
    (cd "$repository_root/services/recommendationservice" && python3 -m unittest test_recommendation_logic)
    ;;
  *)
    printf 'unknown service: %s\n' "$service" >&2
    exit 2
    ;;
esac
