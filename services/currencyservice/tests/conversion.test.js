const assert = require('node:assert/strict');
const path = require('node:path');
const test = require('node:test');

const rates = require(path.join(__dirname, '..', 'data', 'currency_conversion.json'));

function carry(amount) {
  const fractionSize = Math.pow(10, 9);
  amount.nanos += (amount.units % 1) * fractionSize;
  amount.units = Math.floor(amount.units) + Math.floor(amount.nanos / fractionSize);
  amount.nanos = amount.nanos % fractionSize;
  return amount;
}

function convert(from, toCode) {
  const euros = carry({
    units: from.units / rates[from.currency_code],
    nanos: from.nanos / rates[from.currency_code],
  });
  euros.nanos = Math.round(euros.nanos);

  const result = carry({
    units: euros.units * rates[toCode],
    nanos: euros.nanos * rates[toCode],
  });
  result.units = Math.floor(result.units);
  result.nanos = Math.floor(result.nanos);
  result.currency_code = toCode;
  return result;
}

test('currency data includes EUR baseline and major codes', () => {
  assert.equal(rates.EUR, '1.0');
  assert.ok(rates.USD);
  assert.ok(rates.GBP);
  assert.ok(rates.JPY);
});

test('USD to EUR conversion uses ECB rates', () => {
  const usd = { currency_code: 'USD', units: 11, nanos: 305_000_000 };
  const eur = convert(usd, 'EUR');
  assert.equal(eur.currency_code, 'EUR');
  assert.ok(eur.units >= 9 && eur.units <= 11);
});

test('EUR to USD round trip stays near original amount', () => {
  const original = { currency_code: 'EUR', units: 50, nanos: 0 };
  const usd = convert(original, 'USD');
  const back = convert(usd, 'EUR');
  assert.equal(back.currency_code, 'EUR');
  assert.ok(Math.abs(back.units - original.units) <= 1);
});
