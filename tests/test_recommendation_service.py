import os
import sys
import unittest
from pathlib import Path
from types import SimpleNamespace
from unittest import mock

REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
RECOMMENDATION_ROOT = REPOSITORY_ROOT / "services" / "recommendationservice"


class RecommendationServiceLogicTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        os.chdir(RECOMMENDATION_ROOT)
        sys.path.insert(0, str(RECOMMENDATION_ROOT))
        from grpc_health.v1 import health_pb2

        import demo_pb2  # noqa: F401
        import recommendation_server as rs

        cls.health_pb2 = health_pb2
        cls.rs = rs
        cls.RecommendationService = rs.RecommendationService

    def setUp(self) -> None:
        catalog = mock.Mock()
        products = [
            SimpleNamespace(id="OLJCESPC7Z"),
            SimpleNamespace(id="66VCHSJNUP"),
            SimpleNamespace(id="1YMWWN1N4O"),
            SimpleNamespace(id="L9ECAV7KIM"),
            SimpleNamespace(id="2ZYFJ3GM2N"),
            SimpleNamespace(id="0PUK6V6EV0"),
        ]
        catalog.ListProducts.return_value = SimpleNamespace(products=products)
        self.rs.product_catalog_stub = catalog
        self.service = self.RecommendationService()

    def test_list_recommendations_excludes_requested_products(self) -> None:
        import demo_pb2

        request = demo_pb2.ListRecommendationsRequest(
            user_id="user-1",
            product_ids=["OLJCESPC7Z", "66VCHSJNUP"],
        )
        response = self.service.ListRecommendations(request, None)

        self.assertTrue(response.product_ids)
        self.assertLessEqual(len(response.product_ids), 5)
        for product_id in response.product_ids:
            self.assertNotIn(product_id, request.product_ids)

    def test_health_check_reports_serving(self) -> None:
        response = self.service.Check(None, None)
        self.assertEqual(
            response.status,
            self.health_pb2.HealthCheckResponse.SERVING,
        )


if __name__ == "__main__":
    unittest.main()
