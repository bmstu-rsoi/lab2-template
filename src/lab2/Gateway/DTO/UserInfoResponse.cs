using Gateway.Models;
using System;
using System.Collections;
using System.Collections.Generic;

namespace Gateway.DTO
{
    public class UserInfoResponse
    {
        public List<UserReservationInfo> Reservations { get; set; } = null!;
        public LoyaltyInfo Loyalty { get; set; } = null!;
    }


    public class UserReservationInfo
    {
        public Guid ReservationUid { get; set; }
        public HotelInfo Hotel { get; set; }
        public string Status { get; set; } = null!;
        public DateOnly StartDate { get; set; }
        public DateOnly EndDate { get; set; }
        public PaymentInfo Payment { get; set; }
    }

    public class LoyaltyInfo
    {
        public string Status { get; set; } = null!;
        public int Discount { get; set; }
    }

    public class PaymentInfo
    {
        public string Status { get; set; } = null!;
        public int Price { get; set; }
    }

    public class HotelInfo
    {
        public Guid HotelUid { get; set; }
        public string Name { get; set; } = null!;
        public string FullAddress { get; set; } = null!;
        public int? Stars { get; set; }
    }

    public class LoyaltyInfoResponse
    {
        public string Status { get; set; } = null!;
        public int Discount { get; set; }
        public int ReservationCount { get; set; }
    }

    public class CreateReservationRequest
    {
        public Guid HotelUid { get; set; }
        public DateTime StartDate { get; set; }
        public DateTime EndDate { get; set; }
    }

    public class CreateReservationResponse
    {
        public Guid ReservationUid { get; set; }
        public Guid HotelUid { get; set; }
        public string Status { get; set; } = null!;
        public int Discount { get; set; }
        public DateOnly StartDate { get; set; }
        public DateOnly EndDate { get; set; }
        public PaymentInfo Payment { get; set; }
    }
}

