import ast
import re
import unittest
from pathlib import Path


REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
LOCUSTFILE = REPOSITORY_ROOT / "services" / "loadgenerator" / "locustfile.py"


def load_product_catalog() -> list[str]:
    module = ast.parse(LOCUSTFILE.read_text())
    for node in module.body:
        if not isinstance(node, ast.Assign):
            continue
        for target in node.targets:
            if isinstance(target, ast.Name) and target.id == "products":
                value = ast.literal_eval(node.value)
                if isinstance(value, list):
                    return value
    raise AssertionError("products catalog not found in locustfile")


class LoadGeneratorCatalogTests(unittest.TestCase):
    def test_product_catalog_is_non_empty_and_well_formed(self) -> None:
        products = load_product_catalog()
        self.assertGreaterEqual(len(products), 5)
        pattern = re.compile(r"^[A-Z0-9]{10}$")
        for product_id in products:
            self.assertRegex(product_id, pattern)

    def test_checkout_task_posts_expected_form_fields(self) -> None:
        source = LOCUSTFILE.read_text()
        required = {
            "email",
            "street_address",
            "zip_code",
            "city",
            "state",
            "country",
            "credit_card_number",
            "credit_card_expiration_month",
            "credit_card_expiration_year",
            "credit_card_cvv",
        }
        for field in required:
            self.assertIn(f"'{field}'", source)


if __name__ == "__main__":
    unittest.main()
