import unittest
from pathlib import Path

from jinja2 import Environment, FileSystemLoader, select_autoescape


EMAIL_SERVICE_ROOT = Path(__file__).resolve().parents[1] / "services" / "emailservice"


class EmailConfirmationTemplateTests(unittest.TestCase):
    def setUp(self) -> None:
        env = Environment(
            loader=FileSystemLoader(EMAIL_SERVICE_ROOT / "templates"),
            autoescape=select_autoescape(["html", "xml"]),
        )
        self.template = env.get_template("confirmation.html")

    def test_renders_order_identifiers_and_totals(self) -> None:
        order = type(
            "Order",
            (),
            {
                "order_id": "order-123",
                "shipping_tracking_id": "AB-20123-10987654321",
                "shipping_cost": type(
                    "Money",
                    (),
                    {"units": 8, "nanos": 990000000, "currency_code": "USD"},
                )(),
                "shipping_address": type(
                    "Address",
                    (),
                    {
                        "street_address_1": "1600 Amphitheatre Parkway",
                        "street_address_2": "",
                        "city": "Mountain View",
                        "country": "United States",
                        "zip_code": "94043",
                    },
                )(),
                "items": [],
            },
        )()

        html = self.template.render(order=order)
        self.assertIn("order-123", html)
        self.assertIn("AB-20123-10987654321", html)
        self.assertIn("USD", html)

    def test_confirmation_template_escapes_html_injection(self) -> None:
        order = type(
            "Order",
            (),
            {
                "order_id": "<script>alert(1)</script>",
                "shipping_tracking_id": "safe-id",
                "shipping_cost": type(
                    "Money",
                    (),
                    {"units": 1, "nanos": 0, "currency_code": "USD"},
                )(),
                "shipping_address": type(
                    "Address",
                    (),
                    {
                        "street_address_1": "Main",
                        "street_address_2": "",
                        "city": "Town",
                        "country": "US",
                        "zip_code": "00000",
                    },
                )(),
                "items": [],
            },
        )()

        html = self.template.render(order=order)
        self.assertNotIn("<script>", html)
        self.assertIn("&lt;script&gt;alert(1)&lt;/script&gt;", html)


if __name__ == "__main__":
    unittest.main()
