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

from __future__ import annotations

import random
from typing import Callable, Sequence


def pick_recommendation_ids(
    catalog_product_ids: Sequence[str],
    excluded_product_ids: Sequence[str],
    max_responses: int = 5,
    *,
    sample: Callable[[range, int], list[int]] | None = None,
) -> list[str]:
    """Return up to max_responses product ids not already in excluded_product_ids."""
    if max_responses < 0:
        raise ValueError("max_responses must be non-negative")

    filtered_products = list(set(catalog_product_ids) - set(excluded_product_ids))
    num_products = len(filtered_products)
    num_return = min(max_responses, num_products)
    if num_return == 0:
        return []

    sampler = sample if sample is not None else random.sample
    indices = sampler(range(num_products), num_return)
    return [filtered_products[index] for index in indices]
