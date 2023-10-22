package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"lab2/src/gateway-service/internal/models"
)

func GetPrivilegeShortInfo(bonusServiceAddress, username string) (*models.Privilege, error) {
	requestURL := fmt.Sprintf("%s/api/v1/bonus/%s", bonusServiceAddress, username)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Println("Failed to create an http request")
		return nil, err
	}

	client := &http.Client{
		Timeout: 10 * time.Minute,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed request to privilege service: %s", err)
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			fmt.Println("Failed to close response body")
		}
	}(res.Body)

	privilege := &models.Privilege{}
	if res.StatusCode != http.StatusNotFound {
		if err = json.NewDecoder(res.Body).Decode(privilege); err != nil {
			return nil, fmt.Errorf("Failed to decode response: %s", err)
		}
	}

	return privilege, nil
}

func GetPrivilegeHistory(bonusServiceAddress string, privilegeID int) (*[]models.PrivilegeHistory, error) {
	requestURL := fmt.Sprintf("%s/api/v1/bonus/history/%d", bonusServiceAddress, privilegeID)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Println("Failed to create an http request")
		return nil, err
	}

	client := &http.Client{
		Timeout: 10 * time.Minute,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed request to privilege service: %s", err)
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			fmt.Println("Failed to close response body")
		}
	}(res.Body)

	privilegeHistory := &[]models.PrivilegeHistory{}
	if res.StatusCode != http.StatusNotFound {
		if err = json.NewDecoder(res.Body).Decode(privilegeHistory); err != nil {
			return nil, fmt.Errorf("Failed to decode response: %s", err)
		}
	}

	return privilegeHistory, nil
}

func CreatePrivilegeHistoryRecord(bonusServiceAddress, uid, date, optype string, ID, diff int) error {
	requestURL := fmt.Sprintf("%s/api/v1/bonus", bonusServiceAddress)

	record := &models.PrivilegeHistory{
		PrivilegeID:   ID,
		TicketUID:     uid,
		Date:          date,
		BalanceDiff:   diff,
		OperationType: optype,
	}

	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("encoding error: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(data))
	if err != nil {
		fmt.Println("Failed to create an http request")
		return err
	}

	client := &http.Client{
		Timeout: 10 * time.Minute,
	}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed request to privilege service: %s", err)
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			fmt.Println("Failed to close response body")
		}
	}(res.Body)

	return nil
}

func CreatePrivilege(bonusServiceAddress, username string, balance int) error {
	requestURL := fmt.Sprintf("%s/api/v1/bonus/privilege", bonusServiceAddress)

	record := &models.Privilege{
		Username: username,
		Balance:  balance,
	}

	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("encoding error: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(data))
	if err != nil {
		fmt.Println("Failed to create an http request")
		return err
	}

	client := &http.Client{
		Timeout: 10 * time.Minute,
	}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed request to privilege service: %s", err)
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			fmt.Println("Failed to close response body")
		}
	}(res.Body)

	return nil
}
