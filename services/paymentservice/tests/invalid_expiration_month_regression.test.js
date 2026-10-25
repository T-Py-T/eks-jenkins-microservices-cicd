const assert = require('node:assert/strict');
const test = require('node:test');

const charge = require('../charge');

function requestWithExpirationMonth(month) {
  return {
    amount: { currency_code: 'USD', units: 25, nanos: 0 },
    credit_card: {
      credit_card_number: '4111111111111111',
      credit_card_cvv: 123,
      credit_card_expiration_month: month,
      credit_card_expiration_year: new Date().getFullYear() + 2,
    },
  };
}

test('rejects expiration month 0', () => {
  assert.throws(() => charge(requestWithExpirationMonth(0)), /invalid/i);
});

test('rejects expiration month 13', () => {
  assert.throws(() => charge(requestWithExpirationMonth(13)), /invalid/i);
});

test('accepts expiration month 12 for a valid card', () => {
  const result = charge(requestWithExpirationMonth(12));
  assert.match(result.transaction_id, /^[0-9a-f-]{36}$/);
});
