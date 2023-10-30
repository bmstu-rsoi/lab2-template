using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace Reservation_Service
{
    /// <summary>
    /// Запись о человеке
    /// </summary>
    public class Reservation
    {
        public int Id { get; set; }
        public Guid ReservationUid { get; set; }
        public string Username { get; set; } = null!;
        public Guid PaymentUid { get; set; }
        public int? HotelId { get; set; }
        public string Status { get; set; } = null!;
        public DateTime StartDate { get; set; }
        public DateTime EndDate { get; set; }

        public virtual Hotels? Hotel { get; set; }
    }
}
