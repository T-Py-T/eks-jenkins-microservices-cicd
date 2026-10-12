#!/usr/bin/env python3
# Copyright 2018 Google LLC
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

import unittest
from unittest import mock

import grpc
from google.api_core.exceptions import GoogleAPICallError
from grpc_health.v1 import health_pb2
from jinja2 import TemplateError

import demo_pb2
from email_server import BaseEmailService, DummyEmailService, EmailService


class EmailGrpcHandlerTests(unittest.TestCase):
    def test_base_email_check_reports_serving(self) -> None:
        svc = BaseEmailService()
        resp = svc.Check(None, None)
        self.assertEqual(resp.status, health_pb2.HealthCheckResponse.SERVING)

    def test_dummy_send_order_confirmation_returns_empty(self) -> None:
        svc = DummyEmailService()
        ctx = mock.Mock()
        out = svc.SendOrderConfirmation(
            demo_pb2.SendOrderConfirmationRequest(email="shopper@example.com"),
            ctx,
        )
        self.assertIsInstance(out, demo_pb2.Empty)
        ctx.set_code.assert_not_called()

    def test_send_order_confirmation_template_error_sets_internal(self) -> None:
        svc = EmailService.__new__(EmailService)
        ctx = mock.Mock()
        request = demo_pb2.SendOrderConfirmationRequest(
            email="shopper@example.com",
            order=demo_pb2.OrderResult(order_id="order-1"),
        )

        with mock.patch("email_server.template.render", side_effect=TemplateError("bad template")):
            out = svc.SendOrderConfirmation(request, ctx)

        self.assertIsInstance(out, demo_pb2.Empty)
        ctx.set_code.assert_called_once_with(grpc.StatusCode.INTERNAL)
        ctx.set_details.assert_called_once()

    def test_send_order_confirmation_api_error_sets_internal(self) -> None:
        svc = EmailService.__new__(EmailService)
        svc.client = object()
        ctx = mock.Mock()
        request = demo_pb2.SendOrderConfirmationRequest(
            email="shopper@example.com",
            order=demo_pb2.OrderResult(order_id="order-2"),
        )

        with mock.patch("email_server.template.render", return_value="<p>ok</p>"), mock.patch(
            "email_server.EmailService.send_email",
            side_effect=GoogleAPICallError("mail down"),
        ):
            out = svc.SendOrderConfirmation(request, ctx)

        self.assertIsInstance(out, demo_pb2.Empty)
        ctx.set_code.assert_called_once_with(grpc.StatusCode.INTERNAL)
        ctx.set_details.assert_called_once()


if __name__ == "__main__":
    unittest.main()
