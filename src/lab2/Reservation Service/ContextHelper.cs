using System.Linq;
using System;
using System.Threading.Tasks;


namespace Reservation_Service
{
    public class ContextHelper
    {
        public static async Task Seed(ReservationDBContext context)
        {
            if (!context.Hotels.Any())
            {
                var lib = new Hotels()
                {
                    Id = 1,
                    HotelUid = Guid.Parse("049161bb-badd-4fa8-9d90-87c9a82b0668"),
                    Name = "Ararat Park Hyatt Moscow",
                    Country = "Россия",
                    City = "Москва",
                    Address = "Неглинная ул., 4",
                    Stars = 5,
                    Price = 10000
                };
                await context.Hotels.AddAsync(lib);
                await context.SaveChangesAsync();
            }
        }
    }
}
