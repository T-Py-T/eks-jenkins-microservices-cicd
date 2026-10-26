import json
import unittest
from decimal import Decimal
from pathlib import Path


REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
CURRENCY_DATA_PATH = (
    REPOSITORY_ROOT / "services" / "currencyservice" / "data" / "currency_conversion.json"
)


class CurrencyConversionDataTests(unittest.TestCase):
    def setUp(self) -> None:
        self.rates = json.loads(CURRENCY_DATA_PATH.read_text())

    def test_euro_is_base_currency(self) -> None:
        self.assertEqual(self.rates["EUR"], "1.0")

    def test_rates_are_positive_decimals(self) -> None:
        for code, raw_rate in self.rates.items():
            rate = Decimal(raw_rate)
            self.assertGreater(rate, 0, msg=f"{code} must be positive")
            self.assertIsInstance(raw_rate, str)

    def test_frontend_whitelist_currencies_are_present(self) -> None:
        required = {"USD", "EUR", "CAD", "JPY", "GBP", "TRY"}
        missing = required - set(self.rates)
        self.assertFalse(missing, f"missing conversion rates for {sorted(missing)}")


if __name__ == "__main__":
    unittest.main()
