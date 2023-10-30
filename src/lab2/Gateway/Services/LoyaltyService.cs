using System.Net.Http;
using System;
using Gateway.Models;
using System.Collections.Generic;
using System.Net.Http.Json;
using System.Threading.Tasks;
using Gateway.Controllers;
using Microsoft.Extensions.Logging;

namespace Gateway.Services
{
    public class LoyaltyService
    {
        private readonly HttpClient _httpClient;
        
        public LoyaltyService()
        {
            _httpClient = new HttpClient();
            _httpClient.BaseAddress = new Uri("http://loyalty:8050/");
        }

        public async Task<Loyalty?> GetLoyaltyByUsernameAsync(string username)
        {
            using var req = new HttpRequestMessage(HttpMethod.Get, "api/v1/loyalty");
            req.Headers.Add("X-User-Name", username);
            using var res = await _httpClient.SendAsync(req);
            var response = await res.Content.ReadFromJsonAsync<Loyalty>();
            return response;
        }

        public async Task<Loyalty?> PutLoyaltyByUsernameAsync(string username)
        {
            using var req = new HttpRequestMessage(HttpMethod.Put, "api/v1/loyalty");
            req.Headers.Add("X-User-Name", username);
            using var res = await _httpClient.SendAsync(req);
            var response = await res.Content.ReadFromJsonAsync<Loyalty>();
            return response;
        }

        public async Task<Loyalty?> DeleteLoyaltyByUsernameAsync(string username)
        {
            using var req = new HttpRequestMessage(HttpMethod.Delete, "api/v1/loyalty");
            req.Headers.Add("X-User-Name", username);
            using var res = await _httpClient.SendAsync(req);
            var response = await res.Content.ReadFromJsonAsync<Loyalty>();
            return response;
        }
    }
}
