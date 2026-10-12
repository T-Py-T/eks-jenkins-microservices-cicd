const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const test = require('node:test');

const ratesPath = path.join(__dirname, '..', 'data', 'currency_conversion.json');
const rates = JSON.parse(fs.readFileSync(ratesPath, 'utf8'));

test('ECB rate file uses EUR as the 1.0 baseline', () => {
  assert.equal(rates.EUR, '1.0');
});

test('every currency rate is a positive numeric string', () => {
  for (const [code, value] of Object.entries(rates)) {
    assert.match(code, /^[A-Z]{3}$/, `unexpected currency code ${code}`);
    const n = Number(value);
    assert.ok(Number.isFinite(n) && n > 0, `${code} rate must be positive, got ${value}`);
  }
});

test('major checkout currencies are present', () => {
  for (const code of ['USD', 'GBP', 'JPY', 'CAD']) {
    assert.ok(rates[code], `missing rate for ${code}`);
  }
});
