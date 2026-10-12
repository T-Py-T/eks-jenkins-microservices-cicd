const assert = require('node:assert/strict');
const test = require('node:test');

const createLogger = require('../logger');

test('createLogger returns a pino logger with the given name', () => {
  const logger = createLogger('paymentservice-test');
  assert.equal(logger.bindings().name, 'paymentservice-test');
});

test('createLogger child inherits service name binding', () => {
  const logger = createLogger('paymentservice-parent');
  const child = logger.child({ component: 'charge' });
  assert.equal(child.bindings().name, 'paymentservice-parent');
  assert.equal(child.bindings().component, 'charge');
});
