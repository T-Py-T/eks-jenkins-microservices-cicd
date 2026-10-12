const assert = require('node:assert/strict');
const test = require('node:test');

const charge = require('../charge');

function requestWithCard(overrides = {}) {
  const year = new Date().getFullYear() + 2;
  return {
    amount: { currency_code: 'USD', units: 99, nanos: 0 },
    credit_card: {
      credit_card_number: '4111111111111111',
      credit_card_cvv: 123,
      credit_card_expiration_month: 12,
      credit_card_expiration_year: year,
      ...overrides,
    },
  };
}

test('accepts a valid MasterCard number', () => {
  const result = charge(
    requestWithCard({ credit_card_number: '5555555555554444' })
  );
  assert.match(result.transaction_id, /^[0-9a-f-]{36}$/);
});

test('rejects expired cards', () => {
  assert.throws(
    () =>
      charge(
        requestWithCard({
          credit_card_expiration_month: 1,
          credit_card_expiration_year: 2000,
        })
      ),
    /expired/i
  );
});

test('rejects unsupported card brands', () => {
  assert.throws(
    () => charge(requestWithCard({ credit_card_number: '378282246310005' })),
    /cannot process/i
  );
});
