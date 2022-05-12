package aclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

//GetAllAvengers returns list of Avengers
func (c Client) GetAllAvengers() ([]Avenger, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/avengers/getAllAvengers", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	var avengers []Avenger
	err = json.Unmarshal(body, &avengers)
	if err != nil {
		return nil, err
	}

	return avengers, nil
}

//GetAvengerByID returns a specific Avenger
func (c Client) GetAvengerById(id string) (*Avenger, error) {
	avengers, err := c.GetAllAvengers()
	if err != nil {
		return nil, err
	}

	for _, avenger := range avengers {
		if avenger.ID == id {
			return &avenger, nil
		}
	}
	return nil, nil
}

//GetAvengerByName returns a specific Avenger
func (c Client) GetAvengerByName(name string) (*Avenger, error) {
	avengers, err := c.GetAllAvengers()
	if err != nil {
		return nil, err
	}

	for _, avenger := range avengers {
		if avenger.Name == name {
			return &avenger, nil
		}
	}
	return nil, nil
}

//CreateAvenger will create an Avenger
func (c *Client) CreateAvenger(avenger Avenger) (*Avenger, error) {
	avg, err := json.Marshal(avenger)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/avengers/createNewAvenger", c.HostURL), strings.NewReader(string(avg)))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var insertedID InsertedResult
	err = json.Unmarshal(body, &insertedID)
	if err != nil {
		return nil, err
	}
	avenger.ID = insertedID.InsertedID

	return &avenger, nil
}

//UpdateAvengerByName will update an Avenger
func (c *Client) UpdateAvengerByName(avenger Avenger) (*UpdateResult, error) {
	avg, err := json.Marshal(avenger)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/avengers/updateAvengerByName", c.HostURL), strings.NewReader(string(avg)))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	var updateResult UpdateResult
	err = json.Unmarshal(body, &updateResult)
	if err != nil {
		return nil, err
	}

	return &updateResult, nil
}

//DeleteAvengerByName will delete an Avenger
func (c *Client) DeleteAvengerByName(avengerName string) (*DeleteResult, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/avengers/deleteAvengerByName", c.HostURL), http.NoBody)

	params := url.Values{}
	params.Add("name", avengerName)
	req.URL.RawQuery = params.Encode()

	if err != nil {
		return nil, err
	}

	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	var deleteResult DeleteResult
	err = json.Unmarshal(body, &deleteResult)
	if err != nil {
		return nil, err
	}

	return &deleteResult, nil
}
