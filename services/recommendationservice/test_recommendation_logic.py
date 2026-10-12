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

from recommendation_logic import pick_recommendation_ids


class RecommendationLogicTests(unittest.TestCase):
    def test_excludes_requested_products(self) -> None:
        catalog = ["a", "b", "c", "d", "e", "f"]
        picks = pick_recommendation_ids(
            catalog,
            ["a", "b"],
            max_responses=5,
            sample=lambda population, count: list(population)[:count],
        )
        self.assertEqual(sorted(picks), ["c", "d", "e", "f"])

    def test_caps_results_at_max_responses(self) -> None:
        catalog = [f"id-{index}" for index in range(10)]
        picks = pick_recommendation_ids(
            catalog,
            [],
            max_responses=3,
            sample=lambda population, count: list(population)[:count],
        )
        self.assertEqual(len(picks), 3)
        self.assertTrue(all(product_id in catalog for product_id in picks))

    def test_returns_empty_when_everything_is_excluded(self) -> None:
        catalog = ["only"]
        self.assertEqual(pick_recommendation_ids(catalog, ["only"]), [])

    def test_rejects_negative_max_responses(self) -> None:
        with self.assertRaises(ValueError):
            pick_recommendation_ids(["a"], [], max_responses=-1)


if __name__ == "__main__":
    unittest.main()
