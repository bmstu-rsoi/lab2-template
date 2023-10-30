using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection.Metadata.Ecma335;
using System.Threading.Tasks;

namespace Loyalty_Service.Controllers
{
    [Route("/")]
    [ApiController]
    public class LoyaltyController : ControllerBase
    {
        private readonly ILogger<LoyaltyController> _logger;
        private readonly LoyaltyDBContext _loyaltyContext;

        public LoyaltyController(ILogger<LoyaltyController> logger, LoyaltyDBContext loyaltyContext)
        {
            _logger = logger;
            _loyaltyContext = loyaltyContext;
        }

        [IgnoreAntiforgeryToken]
        [HttpGet("manage/health")]
        public async Task<ActionResult> HealthCheck()
        {
            return Ok();
        }

        [HttpGet("api/v1/loyalty")]
        public async Task<ActionResult<Loyalty>> GetByUsername(
            [FromHeader(Name = "X-User-Name")] string xUserName)
        {
            if (string.IsNullOrWhiteSpace(xUserName))
            {
                return BadRequest();

            }

            var response = await _loyaltyContext.Loyalty.FirstOrDefaultAsync(r => r.Username.Equals(xUserName));
            return response;
        }

        [HttpDelete("api/v1/loyalty")]
        public async Task<ActionResult<Loyalty>> DecreaseByUsername(
            [FromHeader(Name = "X-User-Name")] string xUserName)
        {
            var res = await _loyaltyContext.Loyalty
                .FirstOrDefaultAsync(r => r.Username.Equals(xUserName));
            if (res.ReservationCount <= 0)
            {
                return res;
            }
            res.ReservationCount--;
            res.Status = res.ReservationCount switch
            {
                >20 => LoyaltyStatuses.GOLD,
                >10 => LoyaltyStatuses.SILVER,
                _ => LoyaltyStatuses.BRONZE
            };

            res.Discount = res.ReservationCount switch
            {
                > 20 => 10,
                > 10 => 7,
                _ => 5
            };

            await _loyaltyContext.SaveChangesAsync();
            return res;
        }


        [HttpPut("api/v1/loyalty")]
        public async Task<ActionResult<Loyalty>> IncreaseByUsername(
            [FromHeader(Name = "X-User-Name")] string xUserName)
        {
            var res = await _loyaltyContext.Loyalty
                .FirstOrDefaultAsync(r => r.Username.Equals(xUserName));
            res.ReservationCount++;
            res.Status = res.ReservationCount switch
            {
                > 20 => LoyaltyStatuses.GOLD,
                > 10 => LoyaltyStatuses.SILVER,
                _ => LoyaltyStatuses.BRONZE
            };

            res.Discount = res.ReservationCount switch
            {
                > 20 => 10,
                > 10 => 7,
                _ => 5
            };

            await _loyaltyContext.SaveChangesAsync();
            return res;
        }
    }
}
