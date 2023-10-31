using System.Net.Http;
using System;
using Gateway.Models;
using System.Net.Http.Json;
using System.Threading.Tasks;

namespace Gateway.Services
{
    public class PaymentService
    {
        private readonly HttpClient _httpClient;

        public PaymentService()
        {
            _httpClient = new HttpClient();
            _httpClient.BaseAddress = new Uri("http://payment:8060/");
        }

        public async Task<Payment?> GetPaymentByUidAsync(Guid paymentUid)
        {
            using var req = new HttpRequestMessage(HttpMethod.Get, $"api/v1/payments/{paymentUid}");
            using var res = await _httpClient.SendAsync(req);
            var response = await res.Content.ReadFromJsonAsync<Payment>();
            return response;
        }

        public async Task<Payment?> DeletePaymentByUidAsync(Guid paymentUid)
        {
            using var req = new HttpRequestMessage(HttpMethod.Delete, $"api/v1/payments/{paymentUid}");
            using var res = await _httpClient.SendAsync(req);
            var response = await res.Content.ReadFromJsonAsync<Payment>();
            return response;
        }

        public async Task<Payment?> CreatePaymentAsync(Payment request)
        {
            using var req = new HttpRequestMessage(HttpMethod.Post, "api/v1/payments");
            req.Content = JsonContent.Create(request, typeof(Payment));
            using var res = await _httpClient.SendAsync(req);
            var response = await res.Content.ReadFromJsonAsync<Payment>();
            return response;
        }
    }
}
