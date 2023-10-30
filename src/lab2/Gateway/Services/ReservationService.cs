using System.Collections.Generic;
using System.Net.Http.Json;
using System.Net.Http;
using System.Threading.Tasks;
using System;
using Gateway.DTO;
using Gateway.Utils;
using Gateway.Models;

namespace Gateway.Services
{
    public class ReservationService
    {
        private readonly HttpClient _httpClient;

        public ReservationService()
        {
            _httpClient = new HttpClient();
            _httpClient.BaseAddress = new Uri("http://reservation:8070/");
        }

        public async Task<PaginationResponse<IEnumerable<Hotels>>?> GetHotelsAsync(int? page,
        int? size)
        {
            using var req = new HttpRequestMessage(HttpMethod.Get, $"api/v1/hotels?page={page}&size={size}");
            using var res = await _httpClient.SendAsync(req);
            var response = await res.Content.ReadFromJsonAsync<PaginationResponse<IEnumerable<Hotels>>>();
            return response;
        }

        public async Task<Hotels?> GetHotelsByIdAsync(int? id)
        {
            using var req = new HttpRequestMessage(HttpMethod.Get, $"api/v1/hotels/{id}");
            using var res = await _httpClient.SendAsync(req);
            var response = await res.Content.ReadFromJsonAsync<Hotels>();
            return response;
        }

        public async Task<Hotels?> GetHotelsByUidAsync(Guid? id)
        {
            using var req = new HttpRequestMessage(HttpMethod.Get, $"api/v1/hotels/byUid");
            req.Content = JsonContent.Create(id, typeof(Guid?));
            using var res = await _httpClient.SendAsync(req);
            var response = await res.Content.ReadFromJsonAsync<Hotels>();
            return response;
        }


        public async Task<IEnumerable<Reservation>?> GetReservationsByUsernameAsync(string username)
        {
            using var req = new HttpRequestMessage(HttpMethod.Get, "api/v1/reservations");
            req.Headers.Add("X-User-Name", username);
            using var res = await _httpClient.SendAsync(req);
            var response = await res.Content.ReadFromJsonAsync<IEnumerable<Reservation>>();
            return response;
        }

        public async Task<Reservation?> GetReservationsByUidAsync(Guid reservationUid)
        {
            using var req = new HttpRequestMessage(HttpMethod.Get, $"api/v1/reservations/{reservationUid}");
            using var res = await _httpClient.SendAsync(req);
            var response = await res.Content.ReadFromJsonAsync<Reservation>();
            return response;
        }

        public async Task<Reservation?> CreateReservationAsync(string username, Reservation request)
        {
            using var req = new HttpRequestMessage(HttpMethod.Post, "api/v1/reservations");
            req.Headers.Add("X-User-Name", username);
            req.Content = JsonContent.Create(request, typeof(Reservation));
            using var res = await _httpClient.SendAsync(req);
            res.EnsureSuccessStatusCode();
            var response = await res.Content.ReadFromJsonAsync<Reservation>();
            return response;
        }

        

        public async Task<Reservation?> DeleteReservationAsync(Guid reservationUid)
        {
            using var req = new HttpRequestMessage(HttpMethod.Delete, $"api/v1/reservations/{reservationUid}");
            //req.Content = JsonContent.Create(reservation, typeof(Reservation));
            using var res = await _httpClient.SendAsync(req);
            var response = await res.Content.ReadFromJsonAsync<Reservation>();
            return response;
        }
    }
}
