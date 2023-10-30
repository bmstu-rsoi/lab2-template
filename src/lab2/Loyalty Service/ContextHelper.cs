using System.Linq;
using System.Threading.Tasks;

namespace Loyalty_Service
{
    public class ContextHelper
    {
        public static async Task Seed(LoyaltyDBContext context)
        {
            if (!context.Loyalty.Any())
            {
                var lib = new Loyalty()
                {
                    Id = 1,
                    Username = "Test Max",
                    ReservationCount = 25,
                    Status = "GOLD",
                    Discount = 10
                };
                await context.Loyalty.AddAsync(lib);
                await context.SaveChangesAsync();
            }
        }
    }
}
