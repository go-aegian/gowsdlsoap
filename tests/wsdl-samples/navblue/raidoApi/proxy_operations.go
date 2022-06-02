// Code generated by gowsdlsoap DO NOT EDIT.

package raidoApi

import (
	"context"
	"github.com/go-aegian/gowsdlsoap/proxy"
)

var envelopeXmlns map[string]string = map[string]string{

	"http": "http://schemas.xmlsoap.org/wsdl/http/",

	"s": "http://www.w3.org/2001/XMLSchema",

	"soap": "http://schemas.xmlsoap.org/wsdl/soap/",

	"soap12": "http://schemas.xmlsoap.org/wsdl/soap12/",

	"tns": "http://raido.aviolinx.com/api/",

	"wsdl": "http://schemas.xmlsoap.org/wsdl/",
}

type RaidoAPISoap interface {
	Ping(request *Ping) (*PingResponse, error)

	PingContext(ctx context.Context, request *Ping) (*PingResponse, error)

	GetVersion(request *GetVersion) (*GetVersionResponse, error)

	GetVersionContext(ctx context.Context, request *GetVersion) (*GetVersionResponse, error)

	AuthenticateUser(request *AuthenticateUser) (*AuthenticateUserResponse, error)

	AuthenticateUserContext(ctx context.Context, request *AuthenticateUser) (*AuthenticateUserResponse, error)

	GetSchedules(request *GetSchedules) (*GetSchedulesResponse, error)

	GetSchedulesContext(ctx context.Context, request *GetSchedules) (*GetSchedulesResponse, error)

	GetFlights(request *GetFlights) (*GetFlightsResponse, error)

	GetFlightsContext(ctx context.Context, request *GetFlights) (*GetFlightsResponse, error)

	GetMaintenances(request *GetMaintenances) (*GetMaintenancesResponse, error)

	GetMaintenancesContext(ctx context.Context, request *GetMaintenances) (*GetMaintenancesResponse, error)

	GetCrews(request *GetCrews) (*GetCrewsResponse, error)

	GetCrewsContext(ctx context.Context, request *GetCrews) (*GetCrewsResponse, error)

	GetUsers(request *GetUsers) (*GetUsersResponse, error)

	GetUsersContext(ctx context.Context, request *GetUsers) (*GetUsersResponse, error)

	GetRosters(request *GetRosters) (*GetRostersResponse, error)

	GetRostersContext(ctx context.Context, request *GetRosters) (*GetRostersResponse, error)

	GetPairings(request *GetPairings) (*GetPairingsResponse, error)

	GetPairingsContext(ctx context.Context, request *GetPairings) (*GetPairingsResponse, error)

	GetAircrafts(request *GetAircrafts) (*GetAircraftsResponse, error)

	GetAircraftsContext(ctx context.Context, request *GetAircrafts) (*GetAircraftsResponse, error)

	GetAccumulatedValues(request *GetAccumulatedValues) (*GetAccumulatedValuesResponse, error)

	GetAccumulatedValuesContext(ctx context.Context, request *GetAccumulatedValues) (*GetAccumulatedValuesResponse, error)

	GetCrewRevisions(request *GetCrewRevisions) (*GetCrewRevisionsResponse, error)

	GetCrewRevisionsContext(ctx context.Context, request *GetCrewRevisions) (*GetCrewRevisionsResponse, error)

	GetAirports(request *GetAirports) (*GetAirportsResponse, error)

	GetAirportsContext(ctx context.Context, request *GetAirports) (*GetAirportsResponse, error)

	GetHotelBookings(request *GetHotelBookings) (*GetHotelBookingsResponse, error)

	GetHotelBookingsContext(ctx context.Context, request *GetHotelBookings) (*GetHotelBookingsResponse, error)

	SetCrewRevision(request *SetCrewRevision) (*SetCrewRevisionResponse, error)

	SetCrewRevisionContext(ctx context.Context, request *SetCrewRevision) (*SetCrewRevisionResponse, error)

	GetConfigurationData(request *GetConfigurationData) (*GetConfigurationDataResponse, error)

	GetConfigurationDataContext(ctx context.Context, request *GetConfigurationData) (*GetConfigurationDataResponse, error)

	GetRosterTransactions(request *GetRosterTransactions) (*GetRosterTransactionsResponse, error)

	GetRosterTransactionsContext(ctx context.Context, request *GetRosterTransactions) (*GetRosterTransactionsResponse, error)

	SetRosterDesignator(request *SetRosterDesignator) (*SetRosterDesignatorResponse, error)

	SetRosterDesignatorContext(ctx context.Context, request *SetRosterDesignator) (*SetRosterDesignatorResponse, error)

	SetFlightData(request *SetFlightData) (*SetFlightDataResponse, error)

	SetFlightDataContext(ctx context.Context, request *SetFlightData) (*SetFlightDataResponse, error)

	SetRosterData(request *SetRosterData) (*SetRosterDataResponse, error)

	SetRosterDataContext(ctx context.Context, request *SetRosterData) (*SetRosterDataResponse, error)

	SetCrewDocument(request *SetCrewDocument) (*SetCrewDocumentResponse, error)

	SetCrewDocumentContext(ctx context.Context, request *SetCrewDocument) (*SetCrewDocumentResponse, error)

	SetRosters(request *SetRosters) (*SetRostersResponse, error)

	SetRostersContext(ctx context.Context, request *SetRosters) (*SetRostersResponse, error)

	SetRoster(request *SetRoster) (*SetRosterResponse, error)

	SetRosterContext(ctx context.Context, request *SetRoster) (*SetRosterResponse, error)

	SetCrew(request *SetCrew) (*SetCrewResponse, error)

	SetCrewContext(ctx context.Context, request *SetCrew) (*SetCrewResponse, error)

	SetUser(request *SetUser) (*SetUserResponse, error)

	SetUserContext(ctx context.Context, request *SetUser) (*SetUserResponse, error)

	SetMaintenance(request *SetMaintenance) (*SetMaintenanceResponse, error)

	SetMaintenanceContext(ctx context.Context, request *SetMaintenance) (*SetMaintenanceResponse, error)

	DeleteMaintenance(request *DeleteMaintenance) (*DeleteMaintenanceResponse, error)

	DeleteMaintenanceContext(ctx context.Context, request *DeleteMaintenance) (*DeleteMaintenanceResponse, error)

	SetAircraftData(request *SetAircraftData) (*SetAircraftDataResponse, error)

	SetAircraftDataContext(ctx context.Context, request *SetAircraftData) (*SetAircraftDataResponse, error)

	DeleteRosters(request *DeleteRosters) (*DeleteRostersResponse, error)

	DeleteRostersContext(ctx context.Context, request *DeleteRosters) (*DeleteRostersResponse, error)

	SetExternalCrew(request *SetExternalCrew) (*SetExternalCrewResponse, error)

	SetExternalCrewContext(ctx context.Context, request *SetExternalCrew) (*SetExternalCrewResponse, error)
}

type raidoAPISoap struct {
	client *proxy.Client
}

func NewRaidoAPISoap(client *proxy.Client) RaidoAPISoap {
	client.SetXmlns(envelopeXmlns)
	return &raidoAPISoap{client: client}
}

func (service *raidoAPISoap) PingContext(ctx context.Context, request *Ping) (*PingResponse, error) {
	response := new(PingResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/Ping", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) Ping(request *Ping) (*PingResponse, error) {
	return service.PingContext(context.Background(), request)
}

func (service *raidoAPISoap) GetVersionContext(ctx context.Context, request *GetVersion) (*GetVersionResponse, error) {
	response := new(GetVersionResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/GetVersion", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) GetVersion(request *GetVersion) (*GetVersionResponse, error) {
	return service.GetVersionContext(context.Background(), request)
}

func (service *raidoAPISoap) AuthenticateUserContext(ctx context.Context, request *AuthenticateUser) (*AuthenticateUserResponse, error) {
	response := new(AuthenticateUserResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/AuthenticateUser", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) AuthenticateUser(request *AuthenticateUser) (*AuthenticateUserResponse, error) {
	return service.AuthenticateUserContext(context.Background(), request)
}

func (service *raidoAPISoap) GetSchedulesContext(ctx context.Context, request *GetSchedules) (*GetSchedulesResponse, error) {
	response := new(GetSchedulesResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/GetSchedules", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) GetSchedules(request *GetSchedules) (*GetSchedulesResponse, error) {
	return service.GetSchedulesContext(context.Background(), request)
}

func (service *raidoAPISoap) GetFlightsContext(ctx context.Context, request *GetFlights) (*GetFlightsResponse, error) {
	response := new(GetFlightsResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/GetFlights", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) GetFlights(request *GetFlights) (*GetFlightsResponse, error) {
	return service.GetFlightsContext(context.Background(), request)
}

func (service *raidoAPISoap) GetMaintenancesContext(ctx context.Context, request *GetMaintenances) (*GetMaintenancesResponse, error) {
	response := new(GetMaintenancesResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/GetMaintenances", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) GetMaintenances(request *GetMaintenances) (*GetMaintenancesResponse, error) {
	return service.GetMaintenancesContext(context.Background(), request)
}

func (service *raidoAPISoap) GetCrewsContext(ctx context.Context, request *GetCrews) (*GetCrewsResponse, error) {
	response := new(GetCrewsResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/GetCrews", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) GetCrews(request *GetCrews) (*GetCrewsResponse, error) {
	return service.GetCrewsContext(context.Background(), request)
}

func (service *raidoAPISoap) GetUsersContext(ctx context.Context, request *GetUsers) (*GetUsersResponse, error) {
	response := new(GetUsersResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/GetUsers", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) GetUsers(request *GetUsers) (*GetUsersResponse, error) {
	return service.GetUsersContext(context.Background(), request)
}

func (service *raidoAPISoap) GetRostersContext(ctx context.Context, request *GetRosters) (*GetRostersResponse, error) {
	response := new(GetRostersResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/GetRosters", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) GetRosters(request *GetRosters) (*GetRostersResponse, error) {
	return service.GetRostersContext(context.Background(), request)
}

func (service *raidoAPISoap) GetPairingsContext(ctx context.Context, request *GetPairings) (*GetPairingsResponse, error) {
	response := new(GetPairingsResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/GetPairings", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) GetPairings(request *GetPairings) (*GetPairingsResponse, error) {
	return service.GetPairingsContext(context.Background(), request)
}

func (service *raidoAPISoap) GetAircraftsContext(ctx context.Context, request *GetAircrafts) (*GetAircraftsResponse, error) {
	response := new(GetAircraftsResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/GetAircrafts", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) GetAircrafts(request *GetAircrafts) (*GetAircraftsResponse, error) {
	return service.GetAircraftsContext(context.Background(), request)
}

func (service *raidoAPISoap) GetAccumulatedValuesContext(ctx context.Context, request *GetAccumulatedValues) (*GetAccumulatedValuesResponse, error) {
	response := new(GetAccumulatedValuesResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/GetAccumulatedValues", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) GetAccumulatedValues(request *GetAccumulatedValues) (*GetAccumulatedValuesResponse, error) {
	return service.GetAccumulatedValuesContext(context.Background(), request)
}

func (service *raidoAPISoap) GetCrewRevisionsContext(ctx context.Context, request *GetCrewRevisions) (*GetCrewRevisionsResponse, error) {
	response := new(GetCrewRevisionsResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/GetCrewRevisions", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) GetCrewRevisions(request *GetCrewRevisions) (*GetCrewRevisionsResponse, error) {
	return service.GetCrewRevisionsContext(context.Background(), request)
}

func (service *raidoAPISoap) GetAirportsContext(ctx context.Context, request *GetAirports) (*GetAirportsResponse, error) {
	response := new(GetAirportsResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/GetAirports", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) GetAirports(request *GetAirports) (*GetAirportsResponse, error) {
	return service.GetAirportsContext(context.Background(), request)
}

func (service *raidoAPISoap) GetHotelBookingsContext(ctx context.Context, request *GetHotelBookings) (*GetHotelBookingsResponse, error) {
	response := new(GetHotelBookingsResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/GetHotelBookings", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) GetHotelBookings(request *GetHotelBookings) (*GetHotelBookingsResponse, error) {
	return service.GetHotelBookingsContext(context.Background(), request)
}

func (service *raidoAPISoap) SetCrewRevisionContext(ctx context.Context, request *SetCrewRevision) (*SetCrewRevisionResponse, error) {
	response := new(SetCrewRevisionResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/SetCrewRevision", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) SetCrewRevision(request *SetCrewRevision) (*SetCrewRevisionResponse, error) {
	return service.SetCrewRevisionContext(context.Background(), request)
}

func (service *raidoAPISoap) GetConfigurationDataContext(ctx context.Context, request *GetConfigurationData) (*GetConfigurationDataResponse, error) {
	response := new(GetConfigurationDataResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/GetConfigurationData", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) GetConfigurationData(request *GetConfigurationData) (*GetConfigurationDataResponse, error) {
	return service.GetConfigurationDataContext(context.Background(), request)
}

func (service *raidoAPISoap) GetRosterTransactionsContext(ctx context.Context, request *GetRosterTransactions) (*GetRosterTransactionsResponse, error) {
	response := new(GetRosterTransactionsResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/GetRosterTransactions", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) GetRosterTransactions(request *GetRosterTransactions) (*GetRosterTransactionsResponse, error) {
	return service.GetRosterTransactionsContext(context.Background(), request)
}

func (service *raidoAPISoap) SetRosterDesignatorContext(ctx context.Context, request *SetRosterDesignator) (*SetRosterDesignatorResponse, error) {
	response := new(SetRosterDesignatorResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/SetRosterDesignator", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) SetRosterDesignator(request *SetRosterDesignator) (*SetRosterDesignatorResponse, error) {
	return service.SetRosterDesignatorContext(context.Background(), request)
}

func (service *raidoAPISoap) SetFlightDataContext(ctx context.Context, request *SetFlightData) (*SetFlightDataResponse, error) {
	response := new(SetFlightDataResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/SetFlightData", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) SetFlightData(request *SetFlightData) (*SetFlightDataResponse, error) {
	return service.SetFlightDataContext(context.Background(), request)
}

func (service *raidoAPISoap) SetRosterDataContext(ctx context.Context, request *SetRosterData) (*SetRosterDataResponse, error) {
	response := new(SetRosterDataResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/SetRosterData", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) SetRosterData(request *SetRosterData) (*SetRosterDataResponse, error) {
	return service.SetRosterDataContext(context.Background(), request)
}

func (service *raidoAPISoap) SetCrewDocumentContext(ctx context.Context, request *SetCrewDocument) (*SetCrewDocumentResponse, error) {
	response := new(SetCrewDocumentResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/SetCrewDocument", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) SetCrewDocument(request *SetCrewDocument) (*SetCrewDocumentResponse, error) {
	return service.SetCrewDocumentContext(context.Background(), request)
}

func (service *raidoAPISoap) SetRostersContext(ctx context.Context, request *SetRosters) (*SetRostersResponse, error) {
	response := new(SetRostersResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/SetRosters", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) SetRosters(request *SetRosters) (*SetRostersResponse, error) {
	return service.SetRostersContext(context.Background(), request)
}

func (service *raidoAPISoap) SetRosterContext(ctx context.Context, request *SetRoster) (*SetRosterResponse, error) {
	response := new(SetRosterResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/SetRoster", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) SetRoster(request *SetRoster) (*SetRosterResponse, error) {
	return service.SetRosterContext(context.Background(), request)
}

func (service *raidoAPISoap) SetCrewContext(ctx context.Context, request *SetCrew) (*SetCrewResponse, error) {
	response := new(SetCrewResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/SetCrew", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) SetCrew(request *SetCrew) (*SetCrewResponse, error) {
	return service.SetCrewContext(context.Background(), request)
}

func (service *raidoAPISoap) SetUserContext(ctx context.Context, request *SetUser) (*SetUserResponse, error) {
	response := new(SetUserResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/SetUser", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) SetUser(request *SetUser) (*SetUserResponse, error) {
	return service.SetUserContext(context.Background(), request)
}

func (service *raidoAPISoap) SetMaintenanceContext(ctx context.Context, request *SetMaintenance) (*SetMaintenanceResponse, error) {
	response := new(SetMaintenanceResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/SetMaintenance", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) SetMaintenance(request *SetMaintenance) (*SetMaintenanceResponse, error) {
	return service.SetMaintenanceContext(context.Background(), request)
}

func (service *raidoAPISoap) DeleteMaintenanceContext(ctx context.Context, request *DeleteMaintenance) (*DeleteMaintenanceResponse, error) {
	response := new(DeleteMaintenanceResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/DeleteMaintenance", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) DeleteMaintenance(request *DeleteMaintenance) (*DeleteMaintenanceResponse, error) {
	return service.DeleteMaintenanceContext(context.Background(), request)
}

func (service *raidoAPISoap) SetAircraftDataContext(ctx context.Context, request *SetAircraftData) (*SetAircraftDataResponse, error) {
	response := new(SetAircraftDataResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/SetAircraftData", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) SetAircraftData(request *SetAircraftData) (*SetAircraftDataResponse, error) {
	return service.SetAircraftDataContext(context.Background(), request)
}

func (service *raidoAPISoap) DeleteRostersContext(ctx context.Context, request *DeleteRosters) (*DeleteRostersResponse, error) {
	response := new(DeleteRostersResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/DeleteRosters", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) DeleteRosters(request *DeleteRosters) (*DeleteRostersResponse, error) {
	return service.DeleteRostersContext(context.Background(), request)
}

func (service *raidoAPISoap) SetExternalCrewContext(ctx context.Context, request *SetExternalCrew) (*SetExternalCrewResponse, error) {
	response := new(SetExternalCrewResponse)
	err := service.client.CallContext(ctx, "http://raido.aviolinx.com/api/SetExternalCrew", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *raidoAPISoap) SetExternalCrew(request *SetExternalCrew) (*SetExternalCrewResponse, error) {
	return service.SetExternalCrewContext(context.Background(), request)
}
