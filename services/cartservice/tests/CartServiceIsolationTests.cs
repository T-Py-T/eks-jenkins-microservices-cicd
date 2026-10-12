// Copyright 2018 Google LLC
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
    public class CartServiceIsolationTests
    {
        private readonly IHostBuilder _host;

        public CartServiceIsolationTests()
        {
            _host = new HostBuilder().ConfigureWebHost(webBuilder =>
            {
                webBuilder
                    .UseStartup<Startup>()
                    .UseTestServer();
            });
        }

        [Fact]
        public async Task CartsAreIsolatedPerUser()
        {
            using var server = await _host.StartAsync();
            var httpClient = server.GetTestClient();
            var channel = GrpcChannel.ForAddress(httpClient.BaseAddress, new GrpcChannelOptions
            {
                HttpClient = httpClient
            });
            var client = new CartServiceClient(channel);

            var userA = Guid.NewGuid().ToString();
            var userB = Guid.NewGuid().ToString();

            await client.AddItemAsync(new AddItemRequest
            {
                UserId = userA,
                Item = new CartItem { ProductId = "sku-a", Quantity = 2 }
            });

            var cartB = await client.GetCartAsync(new GetCartRequest { UserId = userB });
            Assert.Empty(cartB.Items);

            var cartA = await client.GetCartAsync(new GetCartRequest { UserId = userA });
            Assert.Single(cartA.Items);
            Assert.Equal("sku-a", cartA.Items[0].ProductId);
            Assert.Equal(2, cartA.Items[0].Quantity);

            await client.EmptyCartAsync(new EmptyCartRequest { UserId = userA });
        }

        [Fact]
        public async Task AddingDistinctProductsAccumulatesLineItems()
        {
            using var server = await _host.StartAsync();
            var httpClient = server.GetTestClient();
            var channel = GrpcChannel.ForAddress(httpClient.BaseAddress, new GrpcChannelOptions
            {
                HttpClient = httpClient
            });
            var client = new CartServiceClient(channel);
            var userId = Guid.NewGuid().ToString();

            await client.AddItemAsync(new AddItemRequest
            {
                UserId = userId,
                Item = new CartItem { ProductId = "one", Quantity = 1 }
            });
            await client.AddItemAsync(new AddItemRequest
            {
                UserId = userId,
                Item = new CartItem { ProductId = "two", Quantity = 3 }
            });

            var cart = await client.GetCartAsync(new GetCartRequest { UserId = userId });
            Assert.Equal(2, cart.Items.Count);
            Assert.Equal(4, cart.Items[0].Quantity + cart.Items[1].Quantity);

            await client.EmptyCartAsync(new EmptyCartRequest { UserId = userId });
        }
    }
}
