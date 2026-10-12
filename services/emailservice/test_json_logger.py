#!/usr/bin/env python3
#
# Copyright 2026 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import io
import json
import logging
import unittest

from logger import CustomJsonFormatter, getJSONLogger


class EmailServiceLoggerTest(unittest.TestCase):
    def test_custom_json_formatter_uppercases_severity(self):
        formatter = CustomJsonFormatter('%(timestamp)s %(severity)s %(name)s %(message)s')
        record = logging.LogRecord(
            name='emailservice-test',
            level=logging.WARNING,
            pathname=__file__,
            lineno=1,
            msg='payment failed',
            args=(),
            exc_info=None,
        )
        out = io.StringIO()
        handler = logging.StreamHandler(out)
        handler.setFormatter(formatter)
        logger = logging.getLogger('emailservice.formatter.test')
        logger.handlers = []
        logger.addHandler(handler)
        logger.setLevel(logging.WARNING)
        logger.warning('payment failed')

        payload = json.loads(out.getvalue().strip())
        self.assertEqual(payload['severity'], 'WARNING')
        self.assertEqual(payload['message'], 'payment failed')
        self.assertIn('timestamp', payload)

    def test_get_json_logger_is_named_and_non_propagating(self):
        logger = getJSONLogger('emailservice-unit-test')
        self.assertEqual(logger.name, 'emailservice-unit-test')
        self.assertFalse(logger.propagate)
        self.assertEqual(logger.level, logging.INFO)


if __name__ == '__main__':
    unittest.main()
