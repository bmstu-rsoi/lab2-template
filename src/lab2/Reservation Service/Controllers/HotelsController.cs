using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using Reservation_Service.DTO;
using System.Text.Json;
using System.Text.Json.Serialization;

namespace Reservation_Service
{
    [Route("/")]
    [ApiController]
    public class HotelController : ControllerBase
    {
        private readonly ILogger<ReservationController> _logger;
        private readonly ReservationDBContext _hotelsContext;

        public HotelController(ILogger<ReservationController> logger, ReservationDBContext hotelsContext)
        {
            _logger = logger;
            _hotelsContext = hotelsContext;
        }

        [HttpGet("api/v1/hotels")]
        public async Task<ActionResult<PaginationResponse<IEnumerable<Hotels>>>> GetAllHotels([FromQuery] int? page,
        [FromQuery] int? size)
        {
            _logger.LogInformation("Request to hotels");

            var query = _hotelsContext.Hotels.AsNoTracking().AsQueryable();
            
            var total = await query.CountAsync();

            if (page.HasValue && size.HasValue)
            {
                query = query.OrderBy(l => l.Id).Skip((page.Value - 1) * size.Value).Take(size.Value);
            }

            var hotels = await query.ToListAsync();

            var response = new PaginationResponse<IEnumerable<Hotels>>()
            {
                Page = page.Value,
                PageSize = size.Value,
                Items = hotels,
                TotalElements = total
            };

            return response;
        }

        [HttpGet("api/v1/hotels/{hotelId}")]
        public async Task<ActionResult<Hotels>> GetHotelById([FromRoute] int hotelId)
        {
            var lib = await _hotelsContext.Hotels.FirstOrDefaultAsync(l => l.Id.Equals(hotelId));
            return lib;
        }

        [HttpGet("api/v1/hotels/byUid")]
        public async Task<ActionResult<Hotels>> GetHotelByUid([FromBody] Guid hotelId)
        {
            var lib = await _hotelsContext.Hotels.FirstOrDefaultAsync(l => l.HotelUid.Equals(hotelId));
            return lib;
        }
    }
}
