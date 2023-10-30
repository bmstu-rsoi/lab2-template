using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.EntityFrameworkCore;

namespace Loyalty_Service
{
    public partial class LoyaltyDBContext : DbContext
    {
        public LoyaltyDBContext()
        {
            //Database.EnsureCreated();
        }

        public LoyaltyDBContext(DbContextOptions<LoyaltyDBContext> options)
            : base(options)
        {
            //Database.EnsureCreated();
        }

        public virtual DbSet<Loyalty> Loyalty { get; set; }

        protected override void OnConfiguring(DbContextOptionsBuilder optionsBuilder)
        {
            if (!optionsBuilder.IsConfigured)
            {
                //var databaseHost = Environment.GetEnvironmentVariable("DB_HOST");
                //var databasePort = Environment.GetEnvironmentVariable("DB_PORT");
                //var database = Environment.GetEnvironmentVariable("DATABASE");
                //var username = Environment.GetEnvironmentVariable("USERNAME");
                //var password = Environment.GetEnvironmentVariable("PASSWORD");
                optionsBuilder.UseNpgsql(
                    "Host=postgres;Port=5432;Database=loyalties;Username=postgres;Password=postgres");
            }
        }

        protected override void OnModelCreating(ModelBuilder modelBuilder)
        {
            modelBuilder.HasAnnotation("Relational:Collation", "Russian_Russia.1251");

            modelBuilder.Entity<Loyalty>(entity =>
            {
                entity.ToTable("loyalty");

                entity.Property(e => e.Id)
                    .HasColumnName("id");

                entity.Property(e => e.Username)
                        .HasMaxLength(80)
                        .HasColumnName("username");

                entity.Property(e => e.ReservationCount).HasColumnName("reservation_count");

                entity.Property(e => e.Status)
                .HasMaxLength(80)
                .HasColumnName("status")
                .HasDefaultValueSql("'BRONZE'::character varying");

                entity.Property(e => e.Discount).HasColumnName("discount").HasDefaultValueSql("5::character varying");
            });

            OnModelCreatingPartial(modelBuilder);
        }

        partial void OnModelCreatingPartial(ModelBuilder modelBuilder);
    }
}
