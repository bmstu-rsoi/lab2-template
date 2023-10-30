using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;

namespace Reservation_Service
{
    [Route("/")]
    [ApiController]
    public class ReservationController : ControllerBase
    {
        private readonly ILogger<ReservationController> _logger;
        private readonly ReservationDBContext _reservationsContext;

        public ReservationController(ILogger<ReservationController> logger, ReservationDBContext reservationsContext)
        {
            _logger = logger;
            _reservationsContext = reservationsContext;
        }

        [HttpGet("manage/health")]
        public async Task<ActionResult> HealthCheck()
        {
            return Ok();
        }

        [HttpGet("api/v1/reservations")]
        public async Task<ActionResult<IEnumerable<Reservation>>> GetByUsername(
            [FromHeader(Name = "X-User-Name")] string xUserName)
        {
            if (string.IsNullOrWhiteSpace(xUserName))
            {
                return BadRequest();

            }

            var query = _reservationsContext.Reservations.AsNoTracking().AsQueryable();
            query = query.Where(r => r.Username.Equals(xUserName));
            var response = await query.ToListAsync();
            return response;
        }

        [HttpGet("api/v1/reservations/{reservationUid:guid}")]
        public async Task<ActionResult<Reservation?>> GetByUid([FromRoute] Guid reservationUid)
        {
            var reservation = await _reservationsContext.Reservations.AsNoTracking()
                .FirstOrDefaultAsync(r => r.ReservationUid.Equals(reservationUid));

            return reservation;
        }

        [HttpDelete("api/v1/reservations/{reservationUid}")]
        public async Task<ActionResult<Reservation?>> DeleteByUid([FromRoute] Guid reservationUid)
        {
            var res = await _reservationsContext.Reservations
                .FirstOrDefaultAsync(r => r.ReservationUid.Equals(reservationUid));
            res.Status = ReservationStatuses.CANCELED;
            //res.Username = reservation.Username;
            //res.PaymentUid = reservation.PaymentUid;
            //res.ReservationUid = reservation.ReservationUid;
            //res.StartDate = reservation.StartDate;
            //res.EndDate = reservation.EndDate;
            //res.HotelId = reservation.HotelId;
            //res.Id = reservation.Id;
            await _reservationsContext.SaveChangesAsync();
            return res;
        }


        [HttpPost("api/v1/reservations")]
        public async Task<ActionResult<Reservation>> CreateReservation([FromHeader(Name = "X-User-Name")] string xUserName,
            [FromBody] Reservation request)
        {
            if (string.IsNullOrWhiteSpace(xUserName))
            {
                return BadRequest();
            }

            var newReservation = new Reservation()
            {
                Status = ReservationStatuses.PAID,
                Username = xUserName,
                PaymentUid = request.PaymentUid,
                HotelId = request.HotelId,
                ReservationUid = Guid.NewGuid(),
                //StartDate = DateOnly.FromDateTime(DateTime.Now),
                StartDate = request.StartDate,
                EndDate = request.EndDate,
                Id = request.Id,
            };
            await _reservationsContext.Reservations.AddAsync(newReservation);
            await _reservationsContext.SaveChangesAsync();
            _logger.LogInformation("\nReservation:" + "\n" + newReservation.Status + "\n" + newReservation.ReservationUid + "\n" + newReservation.StartDate + "\n" + newReservation.EndDate + "\n" + newReservation.Id + "\n");
            return newReservation;
        }
    }
}
