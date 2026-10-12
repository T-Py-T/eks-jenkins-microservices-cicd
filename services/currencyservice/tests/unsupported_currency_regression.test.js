const assert = require('node:assert/strict');
const test = require('node:test');

const { assertCurrencySupported } = require('../convert_validation');

const rates = { USD: 1.08, EUR: 1 };

test('accepts a currency code present in the rate table', () => {
  assertCurrencySupported(rates, 'USD');
});

test('rejects an unknown source currency code', () => {
  assert.throws(
    () => assertCurrencySupported(rates, 'XXX'),
    /unsupported currency code: XXX/
  );
});

test('rejects a missing currency code', () => {
  assert.throws(
    () => assertCurrencySupported(rates, ''),
    /unsupported currency code/
  );
});
