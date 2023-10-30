FROM mcr.microsoft.com/dotnet/aspnet:7.0 AS base
WORKDIR /app
EXPOSE 8060

FROM mcr.microsoft.com/dotnet/sdk:7.0 AS build
WORKDIR /src
COPY ["Payment Service.csproj", "."]
RUN dotnet restore "./Payment Service.csproj"
COPY . .
WORKDIR "/src/."
RUN dotnet build "Payment Service.csproj" -c Release -o /app/build

FROM build AS publish
RUN dotnet publish "Payment Service.csproj" -c Release -o /app/publish

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "Payment Service.dll", "--urls=http://0.0.0.0:8060"]