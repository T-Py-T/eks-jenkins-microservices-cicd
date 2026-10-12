/*
 * Copyright 2018 Google LLC.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/**
 * Helper function that handles decimal/fractional carrying for Money amounts.
 */
function carry (amount) {
  const fractionSize = Math.pow(10, 9);
  amount.nanos += (amount.units % 1) * fractionSize;
  amount.units = Math.floor(amount.units) + Math.floor(amount.nanos / fractionSize);
  amount.nanos = amount.nanos % fractionSize;
  return amount;
}

function assertCurrencyCode (rates, code) {
  if (rates[code] === undefined) {
    throw new Error(`unknown currency code: ${code}`);
  }
  const rate = Number(rates[code]);
  if (!Number.isFinite(rate) || rate <= 0) {
    throw new Error(`invalid conversion rate for ${code}`);
  }
  return rate;
}

/**
 * Converts a Money amount from one currency to another using ECB-style rate table.
 */
function convertMoney (from, toCode, rates) {
  const fromRate = assertCurrencyCode(rates, from.currency_code);
  const toRate = assertCurrencyCode(rates, toCode);

  const euros = carry({
    units: from.units / fromRate,
    nanos: from.nanos / fromRate
  });

  euros.nanos = Math.round(euros.nanos);

  const result = carry({
    units: euros.units * toRate,
    nanos: euros.nanos * toRate
  });

  result.units = Math.floor(result.units);
  result.nanos = Math.floor(result.nanos);
  result.currency_code = toCode;
  return result;
}

function getSupportedCurrencyCodes (rates) {
  return Object.keys(rates);
}

module.exports = {
  carry,
  convertMoney,
  getSupportedCurrencyCodes,
};
