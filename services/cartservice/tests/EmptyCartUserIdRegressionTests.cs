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
using Grpc.Net.Client;
using Hipstershop;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.TestHost;
using Microsoft.Extensions.Hosting;
using Xunit;
using static Hipstershop.CartService;

namespace cartservice.tests
{
    public class EmptyCartUserIdRegressionTests
    {
        private readonly IHostBuilder _host;

        public EmptyCartUserIdRegressionTests()
        {
            _host = new HostBuilder().ConfigureWebHost(webBuilder =>
            {
                webBuilder
                    .UseStartup<Startup>()
                    .UseTestServer();
            });
        }

        [Fact]
        public async Task EmptyCart_PreservesUserIdOnStoredCart()
        {
            using var server = await _host.StartAsync();
            var httpClient = server.GetTestClient();

            string userId = Guid.NewGuid().ToString();

            var channel = GrpcChannel.ForAddress(httpClient.BaseAddress, new GrpcChannelOptions
            {
                HttpClient = httpClient
            });

            var client = new CartServiceClient(channel);

            await client.AddItemAsync(new AddItemRequest
            {
                UserId = userId,
                Item = new CartItem { ProductId = "1", Quantity = 1 }
            });

            await client.EmptyCartAsync(new EmptyCartRequest { UserId = userId });

            var cart = await client.GetCartAsync(new GetCartRequest { UserId = userId });

            Assert.NotNull(cart);
            Assert.Equal(userId, cart.UserId);
            Assert.Empty(cart.Items);
        }
    }
}
