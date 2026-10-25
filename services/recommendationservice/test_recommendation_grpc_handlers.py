#!/usr/bin/env python3
# Copyright 2018 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import unittest

import demo_pb2
from grpc_health.v1 import health_pb2

import recommendation_server


class _StubProductCatalog:
    def __init__(self, product_ids: list[str]) -> None:
        self._product_ids = product_ids

    def ListProducts(self, _request: demo_pb2.Empty) -> demo_pb2.ListProductsResponse:
        response = demo_pb2.ListProductsResponse()
        for product_id in self._product_ids:
            product = response.products.add()
            product.id = product_id
            product.name = product_id
        return response


class RecommendationGrpcHandlerTests(unittest.TestCase):
    def setUp(self) -> None:
        self._previous_stub = getattr(recommendation_server, "product_catalog_stub", None)
        recommendation_server.product_catalog_stub = _StubProductCatalog(
            ["cat-a", "cat-b", "cat-c", "cat-d", "cat-e", "cat-f", "cat-g"]
        )
        self.service = recommendation_server.RecommendationService()

    def tearDown(self) -> None:
        if self._previous_stub is None:
            recommendation_server.__dict__.pop("product_catalog_stub", None)
        else:
            recommendation_server.product_catalog_stub = self._previous_stub

    def test_check_reports_serving(self) -> None:
        resp = self.service.Check(None, None)
        self.assertEqual(resp.status, health_pb2.HealthCheckResponse.SERVING)

    def test_list_recommendations_excludes_request_product_ids(self) -> None:
        request = demo_pb2.ListRecommendationsRequest(
            user_id="user-42",
            product_ids=["cat-a", "cat-b"],
        )
        response = self.service.ListRecommendations(request, None)

        self.assertLessEqual(len(response.product_ids), 5)
        for excluded in ("cat-a", "cat-b"):
            self.assertNotIn(excluded, response.product_ids)

    def test_list_recommendations_returns_only_catalog_ids(self) -> None:
        request = demo_pb2.ListRecommendationsRequest(user_id="user-42", product_ids=[])
        response = self.service.ListRecommendations(request, None)

        allowed = {"cat-a", "cat-b", "cat-c", "cat-d", "cat-e", "cat-f", "cat-g"}
        for product_id in response.product_ids:
            self.assertIn(product_id, allowed)


if __name__ == "__main__":
    unittest.main()
