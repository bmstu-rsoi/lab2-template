FROM mcr.microsoft.com/dotnet/aspnet:7.0 AS base
WORKDIR /app
EXPOSE 8050

FROM mcr.microsoft.com/dotnet/sdk:7.0 AS build
WORKDIR /src
COPY ["Loyalty Service.csproj", "."]
RUN dotnet restore "./Loyalty Service.csproj"
COPY . .
WORKDIR "/src/."
RUN dotnet build "Loyalty Service.csproj" -c Release -o /app/build

FROM build AS publish
RUN dotnet publish "Loyalty Service.csproj" -c Release -o /app/publish

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "Loyalty Service.dll", "--urls=http://0.0.0.0:8050"]