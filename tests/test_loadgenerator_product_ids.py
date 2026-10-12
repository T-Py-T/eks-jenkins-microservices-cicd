import ast
import json
import unittest
from pathlib import Path


REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
LOCUSTFILE = REPOSITORY_ROOT / "services" / "loadgenerator" / "locustfile.py"
CATALOG = REPOSITORY_ROOT / "services" / "productcatalogservice" / "products.json"


class LoadGeneratorCatalogTests(unittest.TestCase):
    def test_locust_product_ids_exist_in_catalog(self) -> None:
        module = ast.parse(LOCUSTFILE.read_text())
        products = None
        for node in module.body:
            if isinstance(node, ast.Assign):
                for target in node.targets:
                    if isinstance(target, ast.Name) and target.id == "products":
                        products = ast.literal_eval(node.value)
                        break
        self.assertIsNotNone(products, "locustfile must define a products list")

        catalog_ids = {item["id"] for item in json.loads(CATALOG.read_text())["products"]}
        missing = set(products) - catalog_ids
        self.assertFalse(missing, f"locust products missing from catalog: {sorted(missing)}")

    def test_locust_tasks_cover_user_journey(self) -> None:
        module = ast.parse(LOCUSTFILE.read_text())
        task_names = set()
        for node in module.body:
            if isinstance(node, ast.ClassDef) and node.name == "UserBehavior":
                for child in node.body:
                    if isinstance(child, ast.Assign) and any(
                        isinstance(target, ast.Name) and target.id == "tasks"
                        for target in child.targets
                    ):
                        for key in child.value.keys:
                            if isinstance(key, ast.Name):
                                task_names.add(key.id)
        expected = {"index", "setCurrency", "browseProduct", "addToCart", "viewCart", "checkout"}
        self.assertTrue(expected.issubset(task_names), f"missing tasks: {expected - task_names}")


if __name__ == "__main__":
    unittest.main()
