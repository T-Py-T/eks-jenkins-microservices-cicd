const assert = require('node:assert/strict');
const test = require('node:test');

const charge = require('../charge');

function requestWithCard(overrides = {}) {
  const nextYear = new Date().getFullYear() + 2;
  return {
    amount: { currency_code: 'USD', units: 42, nanos: 500000000 },
    credit_card: {
      credit_card_number: '5555555555554444',
      credit_card_cvv: 123,
      credit_card_expiration_month: 12,
      credit_card_expiration_year: nextYear,
      ...overrides,
    },
  };
}

test('accepts a valid Mastercard', () => {
  const result = charge(requestWithCard({ credit_card_number: '5555555555554444' }));
  assert.match(result.transaction_id, /^[0-9a-f-]{36}$/);
});

test('rejects American Express cards', () => {
  assert.throws(
    () => charge(requestWithCard({ credit_card_number: '378282246310005' })),
    /cannot process/i,
  );
});

test('rejects expired cards', () => {
  assert.throws(
    () => charge(requestWithCard({
      credit_card_expiration_month: 1,
      credit_card_expiration_year: 2000,
    })),
    /expired/i,
  );
});
