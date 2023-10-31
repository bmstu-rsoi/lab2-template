using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;

namespace Payment_Service
{
    [Route("/")]
    [ApiController]
    public class PaymentController : ControllerBase
    {
        private readonly ILogger<PaymentController> _logger;
        private readonly PaymentDBContext _paymentContext;

        public PaymentController(ILogger<PaymentController> logger, PaymentDBContext paymentContext)
        {
            _logger = logger;
            _paymentContext = paymentContext;
        }

        [HttpGet("manage/health")]
        public async Task<ActionResult> HealthCheck()
        {
            return Ok();
        }

        [HttpGet("api/v1/payments/{paymentUid}")]
        public async Task<ActionResult<Payment>> GetByUid([FromRoute] Guid paymentUid)
        {
            var reservation = await _paymentContext.Payments.AsNoTracking()
                .FirstOrDefaultAsync(r => r.PaymentUid.Equals(paymentUid));

            return reservation;
        }

        [HttpDelete("api/v1/payments/{paymentUid}")]
        public async Task<ActionResult<Payment?>> UpdateByUid([FromRoute] Guid paymentUid)
        {
            var res = await _paymentContext.Payments
                .FirstOrDefaultAsync(r => r.PaymentUid.Equals(paymentUid));
            res.Status = PaymentStatuses.CANCELED;
            await _paymentContext.SaveChangesAsync();
            return res;
        }


        [HttpPost("api/v1/payments")]
        public async Task<ActionResult<Payment>> CreatePayment(
            [FromBody] Payment request)
        {
            if (request == null)
            {
                return BadRequest();
            }

            var newReservation = new Payment()
            {
                Status = PaymentStatuses.PAID,
                PaymentUid = Guid.NewGuid(),
                Price = request.Price,
                
            };
            await _paymentContext.Payments.AddAsync(newReservation);
            await _paymentContext.SaveChangesAsync();

            return newReservation;
        }
    }
}
