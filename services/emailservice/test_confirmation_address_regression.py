#!/usr/bin/env python3

import unittest
from pathlib import Path

from jinja2 import Environment, FileSystemLoader, select_autoescape

import demo_pb2


class ConfirmationAddressRegressionTest(unittest.TestCase):
    def test_confirmation_email_renders_protobuf_address_fields(self) -> None:
        templates_dir = Path(__file__).resolve().parent / "templates"
        env = Environment(
            loader=FileSystemLoader(str(templates_dir)),
            autoescape=select_autoescape(["html", "xml"]),
        )
        template = env.get_template("confirmation.html")

        order = demo_pb2.OrderResult(
            order_id="test-order",
            shipping_tracking_id="AB-12-3456789",
            shipping_cost=demo_pb2.Money(currency_code="USD", units=8, nanos=990_000_000),
        )
        order.shipping_address.street_address = "1600 Amphitheatre Parkway"
        order.shipping_address.city = "Mountain View"
        order.shipping_address.state = "CA"
        order.shipping_address.country = "United States"
        order.shipping_address.zip_code = 94043

        rendered = template.render(order=order)

        self.assertIn("1600 Amphitheatre Parkway", rendered)
        self.assertIn("Mountain View", rendered)
        self.assertIn("CA", rendered)
        self.assertIn("94043", rendered)
        self.assertIn("United States", rendered)
        self.assertNotIn("street_address_1", rendered)
        self.assertNotIn("street_address_2", rendered)


if __name__ == "__main__":
    unittest.main()
