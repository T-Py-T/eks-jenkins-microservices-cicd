const assert = require('node:assert/strict');
const test = require('node:test');

const charge = require('../charge');

function requestWithCard(number) {
  return {
    amount: { currency_code: 'USD', units: 25, nanos: 0 },
    credit_card: {
      credit_card_number: number,
      credit_card_cvv: 123,
      credit_card_expiration_month: 12,
      credit_card_expiration_year: new Date().getFullYear() + 1,
    },
  };
}

test('returns a transaction id for a valid Visa card', () => {
  const result = charge(requestWithCard('4111111111111111'));
  assert.match(result.transaction_id, /^[0-9a-f-]{36}$/);
});

test('rejects an invalid card number', () => {
  assert.throws(() => charge(requestWithCard('4111111111111112')), /invalid/i);
});
