const assert = require('node:assert/strict');
const test = require('node:test');

const {
  carry,
  convertMoney,
  getSupportedCurrencyCodes,
} = require('../currency');

const SAMPLE_RATES = {
  EUR: '1.0',
  USD: '1.1305',
  CHF: '1.1360',
};

test('carry normalizes fractional units into nanos', () => {
  const amount = carry({ units: 1.5, nanos: 500_000_000 });
  assert.equal(amount.units, 2);
  assert.equal(amount.nanos, 0);
});

test('getSupportedCurrencyCodes lists rate table keys', () => {
  assert.deepEqual(
    getSupportedCurrencyCodes(SAMPLE_RATES).sort(),
    ['CHF', 'EUR', 'USD']
  );
});

test('convertMoney converts CHF to EUR using ECB rates', () => {
  const result = convertMoney(
    { currency_code: 'CHF', units: 300, nanos: 0 },
    'EUR',
    SAMPLE_RATES
  );
  assert.equal(result.currency_code, 'EUR');
  assert.equal(result.units, 264);
  assert.equal(result.nanos, 84_507_042);
});

test('convertMoney rejects unknown source currency', () => {
  assert.throws(
    () => convertMoney(
      { currency_code: 'ZZZ', units: 1, nanos: 0 },
      'EUR',
      SAMPLE_RATES
    ),
    /unknown currency code: ZZZ/
  );
});

test('convertMoney rejects unknown target currency', () => {
  assert.throws(
    () => convertMoney(
      { currency_code: 'EUR', units: 1, nanos: 0 },
      'ZZZ',
      SAMPLE_RATES
    ),
    /unknown currency code: ZZZ/
  );
});
