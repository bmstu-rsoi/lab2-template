using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace Loyalty_Service
{
    /// <summary>
    /// Запись о человеке
    /// </summary>
    public class Loyalty
    {
        public int Id { get; set; }
        public string Username { get; set; } = null!;
        public string Status { get; set; } = null!;
        public int ReservationCount { get; set; }
        public int Discount { get; set; }
    }
}
