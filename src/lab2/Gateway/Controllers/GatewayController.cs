using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using Gateway.Services;
using Gateway.DTO;
using Gateway.Models;
using System.Collections.Generic;
using System.Threading.Tasks;
using System.Linq;
using System;

namespace Gateway.Controllers
{
    [ApiController]
    [Route("/api/v1/")]
    public class GatewayController : ControllerBase
    {
        private readonly ILogger<GatewayController> _logger;
        private readonly ReservationService _reservationsService;
        private readonly PaymentService _paymentsService;
        private readonly LoyaltyService _loyaltyService;

        public GatewayController(ILogger<GatewayController> logger, ReservationService reservationsService,
            PaymentService paymentsService, LoyaltyService loyaltyService)
        {
            _logger = logger;
            _reservationsService = reservationsService;
            _paymentsService = paymentsService;
            _loyaltyService = loyaltyService;
        }

        /// <summary>
        /// Проверяет жив ли ещё сервис
        /// </summary>
        /// <returns>Если сервис жив возвращает 200 </returns>
        /// <response code="200" cref="Person">Работает</response>
        [IgnoreAntiforgeryToken]
        [HttpGet("/manage/health")]
        public async Task<ActionResult> HealthCheck()
        {
            return Ok();
        }

    
        [HttpGet("hotels")]
        public async Task<PaginationResponse<IEnumerable<Hotels>>?> GetAllHotels(
        [FromQuery] int? page,
        [FromQuery] int? size)
        {
            var response = await _reservationsService.GetHotelsAsync(page, size);
            return response;
        }

        /// <summary>
        /// Возвращает информацию о пользователе
        /// </summary>
        /// <returns>Запись о пользователе и его статус лояльности  </returns>
        /// <response code="200" cref="Person">Работает</response>
        [HttpGet("me")]
        public async Task<UserInfoResponse?> GetUserInfoByUsername(
        [FromHeader(Name = "X-User-Name")] string xUserName)
        {
            var reservations = await _reservationsService.GetReservationsByUsernameAsync(xUserName);
            if (reservations == null || !reservations.Any())
            {
                return null;
            }

            var response = new UserInfoResponse();
            response.Reservations = new List<UserReservationInfo>();
            var tasks = reservations.Select(reservation => Task.Run(async () =>
            {
                var hotel = await _reservationsService.GetHotelsByIdAsync(reservation.HotelId);
                var payment = await _paymentsService.GetPaymentByUidAsync(reservation.PaymentUid);

                response.Reservations.Add(new UserReservationInfo()
                {
                    ReservationUid = reservation.ReservationUid,
                    Status = reservation.Status,
                    StartDate = DateOnly.FromDateTime(reservation.StartDate),
                    EndDate = DateOnly.FromDateTime(reservation.EndDate),
                    Hotel = new HotelInfo()
                    {
                        HotelUid = hotel.HotelUid,
                        Name = hotel.Name,
                        FullAddress = hotel.Country + ", " + hotel.City + ", " + hotel.Address,
                        Stars = hotel.Stars,
                    },
                    Payment = new PaymentInfo()
                    {
                        Status = payment.Status,
                        Price = payment.Price,
                    },
                });
            }));

            await Task.WhenAll(tasks);

            var loyalty = await _loyaltyService.GetLoyaltyByUsernameAsync(xUserName);

            response.Loyalty = new LoyaltyInfo()
            {
                Status = loyalty.Status,
                Discount = loyalty.Discount,
            };

            return response;
        }

        /// <summary>
        /// Получить информацию о бронировании
        /// </summary>
        /// <returns>Записи о всех бронях с отелями и оплатами </returns>
        /// <response code="200" cref="Person">Работает</response>
        [HttpGet("reservations")]
        public async Task<List<UserReservationInfo>?> GetReservationsInfoByUsername(
        [FromHeader(Name = "X-User-Name")] string xUserName)
        {
            var reservations = await _reservationsService.GetReservationsByUsernameAsync(xUserName);
            if (reservations == null || !reservations.Any())
            {
                return null;
            }

            var response = new List<UserReservationInfo>();
            var tasks = reservations.Select(reservation => Task.Run(async () =>
            {
                var hotel = await _reservationsService.GetHotelsByIdAsync(reservation.HotelId);
                var payment = await _paymentsService.GetPaymentByUidAsync(reservation.PaymentUid);
                _logger.LogInformation("\nHotel adress: " + hotel.Country + hotel.City + hotel.Address + "\n");

                response.Add(new UserReservationInfo()
                {
                    ReservationUid = reservation.ReservationUid,
                    Status = reservation.Status,
                    StartDate = DateOnly.FromDateTime(reservation.StartDate),
                    EndDate = DateOnly.FromDateTime(reservation.EndDate),
                    Hotel = new HotelInfo()
                    {
                        HotelUid = hotel.HotelUid,
                        Name = hotel.Name,
                        FullAddress = hotel.Country + ", " + hotel.City + ", " + hotel.Address,
                        Stars = hotel.Stars,
                    },
                    Payment = new PaymentInfo()
                    {
                        Status = payment.Status,
                        Price = payment.Price,
                    },
                });
            }));

            
            await Task.WhenAll(tasks);

            return response;
        }

        [HttpGet("reservations/{reservationsUid}")]
        public async Task<ActionResult<UserReservationInfo>?> GetReservationsInfoByUsername(
        [FromRoute] Guid reservationsUid,
        [FromHeader(Name = "X-User-Name")] string xUserName)
        {
            var reservation = await _reservationsService.GetReservationsByUidAsync(reservationsUid);
            if (reservation == null)
            {
                return BadRequest();
            }

            if (!reservation.Username.Equals(xUserName))
            {
                return BadRequest();
            }

            var hotel = await _reservationsService.GetHotelsByIdAsync(reservation.HotelId);
            var payment = await _paymentsService.GetPaymentByUidAsync(reservation.PaymentUid);

            var response = new UserReservationInfo()
            {
                ReservationUid = reservation.ReservationUid,
                Status = reservation.Status,
                StartDate = DateOnly.FromDateTime(reservation.StartDate),
                EndDate = DateOnly.FromDateTime(reservation.EndDate),
                Hotel = new HotelInfo()
                {
                    HotelUid = hotel.HotelUid,
                    Name = hotel.Name,
                    FullAddress = hotel.Country + ", " + hotel.City + ", " + hotel.Address,
                    Stars = hotel.Stars,
                },
                Payment = new PaymentInfo()
                {
                    Status = payment.Status,
                    Price = payment.Price,
                },
            };

            _logger.LogInformation("\nHotel adress: " + response.Hotel.FullAddress + "\n");


            return response;
        }

        [HttpPost("reservations")]
        public async Task<ActionResult<CreateReservationResponse?>> CreateReservation(
        [FromHeader(Name = "X-User-Name")] string xUserName,
        [FromBody] CreateReservationRequest request)
        {

            var hotel = await _reservationsService.GetHotelsByUidAsync(request.HotelUid);

            if (hotel == null)
            {
                return BadRequest(null);
            }

            int sum = ((request.EndDate - request.StartDate).Days) * hotel.Price;

            var loyalty = await _loyaltyService.GetLoyaltyByUsernameAsync(xUserName);

            _logger.LogInformation("\nPrice: " + sum + "\n");

            if (loyalty == null)
            {
                sum -= sum * 5 / 100;
            }
            else
            {
                sum -= sum * loyalty.Discount / 100;
            }

            _logger.LogInformation("\nPrice after discount: " + sum + "\n");

            Payment paymentRequest = new Payment()
            {
                Price = sum,
            };

            var payment = await _paymentsService.CreatePaymentAsync(paymentRequest);

            if (payment == null)
            {
                return BadRequest(null);
            }

            loyalty = await _loyaltyService.PutLoyaltyByUsernameAsync(xUserName);

            Reservation reservationRequest = new Reservation()
            {
                PaymentUid = payment.PaymentUid,
                HotelId = hotel.Id,
                EndDate = request.EndDate,
                StartDate = request.StartDate,
            };

            var reservation = await _reservationsService.CreateReservationAsync(xUserName, reservationRequest);

            var reservationResponse = new CreateReservationResponse()
            {
                ReservationUid = reservation.ReservationUid,
                HotelUid = hotel.HotelUid,
                //StartDate = DateOnly.FromDateTime(reservation.StartDate),
                //EndDate = DateOnly.FromDateTime(reservation.EndDate),
                StartDate = DateOnly.FromDateTime(reservation.StartDate),
                EndDate = DateOnly.FromDateTime(reservation.EndDate),
                Status = reservation.Status,
                Discount = loyalty.Discount,
                Payment = new PaymentInfo()
                {
                    Status = payment.Status,
                    Price = payment.Price,
                }
            };

            return Ok(reservationResponse);
        }

        [HttpDelete("reservations/{reservationsUid}")]
        public async Task<ActionResult> DeleteReservationsByUid(
        [FromRoute] Guid reservationsUid,
        [FromHeader(Name = "X-User-Name")] string xUserName)
        {
            var reservation = await _reservationsService.GetReservationsByUidAsync(reservationsUid);
            if (reservation == null)
            {
                return BadRequest();
            }

            var updateReservationTask = _reservationsService.DeleteReservationAsync(reservation.ReservationUid);

            var updatePaymentTask = _paymentsService.DeletePaymentByUidAsync(reservation.PaymentUid);

            var updateLoyaltyTask = _loyaltyService.DeleteLoyaltyByUsernameAsync(xUserName);

            await updateReservationTask;
            await updatePaymentTask;
            await updateLoyaltyTask;

            return Ok(null);
        }

        [HttpGet("loyalty")]
        public async Task<LoyaltyInfoResponse?> GetLoyaltyInfoByUsername(
        [FromHeader(Name = "X-User-Name")] string xUserName)
        {
            var loyalty = await _loyaltyService.GetLoyaltyByUsernameAsync(xUserName);

            var response = new LoyaltyInfoResponse()
            {
                Status = loyalty.Status,
                Discount = loyalty.Discount,
                ReservationCount = loyalty.ReservationCount,
            };

            return response;
        }


    }
}
