using System;
using System.Linq;
using Microsoft.EntityFrameworkCore;
#nullable disable

namespace Reservation_Service
{
    public partial class ReservationDBContext : DbContext
    {
            public ReservationDBContext()
            {
                Database.EnsureCreated();
            }

            public ReservationDBContext(DbContextOptions<ReservationDBContext> options)
                : base(options)
            {
                Database.EnsureCreated();
            }

            public virtual DbSet<Reservation> Reservations { get; set; } = null!;
            public virtual DbSet<Hotels> Hotels { get; set; } = null!;

            protected override void OnConfiguring(DbContextOptionsBuilder optionsBuilder)
            {
                if (!optionsBuilder.IsConfigured)
                {
                optionsBuilder.UseNpgsql(
                    "Host=postgres;Port=5432;Database=reservations;Username=postgres;Password=postgres");
            }
            }

            protected override void OnModelCreating(ModelBuilder modelBuilder)
            {



                modelBuilder.HasAnnotation("Relational:Collation", "Russian_Russia.1251");




                modelBuilder.Entity<Reservation>(entity =>
                {
                    //entity.HasKey(h => new { h.HotelId });

                    entity.ToTable("reservation");

                    entity.HasIndex(e => e.ReservationUid, "reservation_reservation_uid_key")
                    .IsUnique();

                    entity.Property(e => e.Id)
                        .HasColumnName("id");

                    entity.Property(e => e.ReservationUid).HasColumnName("reservation_uid");

                    entity.Property(e => e.Username)
                    .HasMaxLength(80)
                    .HasColumnName("username");

                    entity.Property(e => e.PaymentUid).HasColumnName("payment_uid");

                    entity.Property(e => e.HotelId).HasColumnName("hotel_id");

                    entity.Property(e => e.Status)
                    .HasMaxLength(20)
                    .HasColumnName("status");

                    entity.Property(e => e.StartDate)
                    .HasColumnType("timestamp without time zone")
                    .HasColumnName("start_date");

                    entity.Property(e => e.EndDate)
                    .HasColumnType("timestamp without time zone")
                    .HasColumnName("end_data");

                    entity.HasOne(d => d.Hotel)
                    .WithMany()
                    .HasForeignKey(d => d.HotelId)
                    .HasConstraintName("reservation_hotel_id_fkey");

                });

                modelBuilder.Entity<Hotels>(entity =>
                {
                entity.ToTable("hotels");

                entity.HasIndex(e => e.HotelUid, "hotels_hotel_uid_key")
                    .IsUnique();

                entity.Property(e => e.Id)
                    .ValueGeneratedNever()
                    .HasColumnName("id");

                entity.Property(e => e.HotelUid).HasColumnName("hotel_uid");

                    entity.Property(e => e.Name)
                        .HasMaxLength(255)
                        .HasColumnName("name");

                    entity.Property(e => e.Country)
                        .HasMaxLength(80)
                        .HasColumnName("country");

                    entity.Property(e => e.City)
                        .HasMaxLength(80)
                        .HasColumnName("city");

                    entity.Property(e => e.Address)
                    .HasMaxLength(255)
                    .HasColumnName("address");

                entity.Property(e => e.Stars).HasColumnName("stars");

                entity.Property(e => e.Price).HasColumnName("price");

                });

            OnModelCreatingPartial(modelBuilder);
            }

            partial void OnModelCreatingPartial(ModelBuilder modelBuilder);
    }
}
