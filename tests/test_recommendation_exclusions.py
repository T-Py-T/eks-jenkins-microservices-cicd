import random
import unittest


def list_recommendation_ids(all_product_ids, cart_product_ids, max_responses=5):
    """Pure-python mirror of recommendationservice list logic."""
    filtered_products = list(set(all_product_ids) - set(cart_product_ids))
    num_products = len(filtered_products)
    num_return = min(max_responses, num_products)
    if num_return == 0:
        return []
    indices = random.sample(range(num_products), num_return)
    return [filtered_products[index] for index in indices]


class RecommendationExclusionTests(unittest.TestCase):
    def test_excludes_products_already_in_cart(self) -> None:
        random.seed(0)
        catalog = ["a", "b", "c", "d", "e"]
        cart = ["a", "c"]
        recommendations = list_recommendation_ids(catalog, cart, max_responses=5)
        self.assertTrue(set(recommendations).isdisjoint(cart))
        self.assertEqual(len(recommendations), 3)

    def test_caps_results_at_max_responses(self) -> None:
        random.seed(1)
        catalog = [f"p{i}" for i in range(20)]
        recommendations = list_recommendation_ids(catalog, [], max_responses=5)
        self.assertEqual(len(recommendations), 5)

    def test_returns_empty_when_cart_contains_entire_catalog(self) -> None:
        catalog = ["only-item"]
        recommendations = list_recommendation_ids(catalog, catalog, max_responses=5)
        self.assertEqual(recommendations, [])


if __name__ == "__main__":
    unittest.main()
