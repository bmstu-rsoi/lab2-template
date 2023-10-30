using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace Gateway.DTO
{
    public class PaginationResponse<T>
    {
        public int Page { get; set; }
        public int PageSize { get; set; }
        public int TotalElements { get; set; }
        public T Items { get; set; }
    }
}
