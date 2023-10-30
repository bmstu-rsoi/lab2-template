FROM mcr.microsoft.com/dotnet/aspnet:7.0 AS base
WORKDIR /app
EXPOSE 8070

FROM mcr.microsoft.com/dotnet/sdk:7.0 AS build
WORKDIR /src
COPY ["Reservation Service.csproj", "."]
RUN dotnet restore "./Reservation Service.csproj"
COPY . .
WORKDIR "/src/."
RUN dotnet build "Reservation Service.csproj" -c Release -o /app/build

FROM build AS publish
RUN dotnet publish "Reservation Service.csproj" -c Release -o /app/publish

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "Reservation Service.dll", "--urls=http://0.0.0.0:8070"]