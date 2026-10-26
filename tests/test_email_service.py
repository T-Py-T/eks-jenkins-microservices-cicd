import os
import sys
import unittest
from pathlib import Path
from types import SimpleNamespace
from unittest import mock

from jinja2 import Environment, FileSystemLoader, select_autoescape


REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
EMAIL_SERVICE_ROOT = REPOSITORY_ROOT / "services" / "emailservice"


class EmailConfirmationTemplateTests(unittest.TestCase):
    def setUp(self) -> None:
        self.env = Environment(
            loader=FileSystemLoader(EMAIL_SERVICE_ROOT / "templates"),
            autoescape=select_autoescape(["html", "xml"]),
        )
        self.template = self.env.get_template("confirmation.html")

    def test_renders_order_metadata_and_line_items(self) -> None:
        order = SimpleNamespace(
            order_id="order-123",
            shipping_tracking_id="track-456",
            shipping_cost=SimpleNamespace(units=12, nanos=500_000_000, currency_code="USD"),
            shipping_address=SimpleNamespace(
                street_address_1="1600 Amphitheatre Parkway",
                street_address_2="",
                city="Mountain View",
                country="United States",
                zip_code="94043",
            ),
            items=[
                SimpleNamespace(
                    item=SimpleNamespace(product_id="OLJCESPC7Z", quantity=2),
                    cost=SimpleNamespace(units=15, nanos=990_000_000, currency_code="USD"),
                )
            ],
        )

        rendered = self.template.render(order=order)

        self.assertIn("order-123", rendered)
        self.assertIn("track-456", rendered)
        self.assertIn("OLJCESPC7Z", rendered)
        self.assertIn("United States", rendered)


class EmailServiceBehaviorTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        os.chdir(EMAIL_SERVICE_ROOT)
        sys.path.insert(0, str(EMAIL_SERVICE_ROOT))
        from grpc_health.v1 import health_pb2

        import demo_pb2  # noqa: F401
        import email_server  # noqa: F401

        cls.demo_pb2 = demo_pb2
        cls.email_server = email_server
        cls.health_pb2 = health_pb2

    @classmethod
    def tearDownClass(cls) -> None:
        if str(EMAIL_SERVICE_ROOT) in sys.path:
            sys.path.remove(str(EMAIL_SERVICE_ROOT))

    def test_health_check_reports_serving(self) -> None:
        service = self.email_server.BaseEmailService()
        response = service.Check(request=None, context=None)
        self.assertEqual(
            response.status,
            self.health_pb2.HealthCheckResponse.SERVING,
        )

    def test_dummy_send_order_confirmation_is_successful(self) -> None:
        service = self.email_server.DummyEmailService()
        request = self.demo_pb2.SendOrderConfirmationRequest(
            email="buyer@example.com",
            order=self.demo_pb2.OrderResult(order_id="abc"),
        )
        context = mock.Mock()

        response = service.SendOrderConfirmation(request, context)

        self.assertIsInstance(response, self.demo_pb2.Empty)
        context.set_code.assert_not_called()


if __name__ == "__main__":
    unittest.main()
