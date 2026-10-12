// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

using System;
using System.Threading.Tasks;
using Grpc.Core;
using Grpc.Net.Client;
using Hipstershop;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.TestHost;
using Microsoft.Extensions.Hosting;
using Xunit;
using static Hipstershop.CartService;

namespace cartservice.tests
{
    public class AddItemQuantityValidationRegressionTests
    {
        private readonly IHostBuilder _host = new HostBuilder().ConfigureWebHost(webBuilder =>
        {
            webBuilder
                .UseStartup<Startup>()
                .UseTestServer();
        });

        [Fact]
        public async Task AddItem_ZeroQuantity_ReturnsInvalidArgument()
        {
            using var server = await _host.StartAsync();
            var httpClient = server.GetTestClient();
            var channel = GrpcChannel.ForAddress(httpClient.BaseAddress, new GrpcChannelOptions
            {
                HttpClient = httpClient
            });
            var cartClient = new CartServiceClient(channel);

            var rpc = async () => await cartClient.AddItemAsync(new AddItemRequest
            {
                UserId = Guid.NewGuid().ToString(),
                Item = new CartItem { ProductId = "0PUK6V6EV0", Quantity = 0 }
            });

            var ex = await Assert.ThrowsAsync<RpcException>(rpc);
            Assert.Equal(StatusCode.InvalidArgument, ex.StatusCode);
        }

        [Fact]
        public async Task AddItem_NegativeQuantity_ReturnsInvalidArgument()
        {
            using var server = await _host.StartAsync();
            var httpClient = server.GetTestClient();
            var channel = GrpcChannel.ForAddress(httpClient.BaseAddress, new GrpcChannelOptions
            {
                HttpClient = httpClient
            });
            var cartClient = new CartServiceClient(channel);

            var rpc = async () => await cartClient.AddItemAsync(new AddItemRequest
            {
                UserId = Guid.NewGuid().ToString(),
                Item = new CartItem { ProductId = "0PUK6V6EV0", Quantity = -3 }
            });

            var ex = await Assert.ThrowsAsync<RpcException>(rpc);
            Assert.Equal(StatusCode.InvalidArgument, ex.StatusCode);
        }
    }
}
