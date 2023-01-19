package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Artists struct {
	Id            int
	Image         string
	Name          string
	Members       []string
	CreationDate  int
	FirstAlbum    string
	Locations     string
	ConcertDates  string
	Relations     string
	RelationsData Relations
	LocationsData Locations
}

type Relations struct {
	Id             int
	DatesLocations map[string][]string
}

type Locations struct {
	Id        int
	Locations []string
}

type ApiClient struct {
	Client http.Client
}

const (
	url = "https://groupietrackers.herokuapp.com/api/artists"
)

func NewClient() ApiClient {
	netClient := http.Client{
		Timeout: time.Second * 10,
	}
	api := ApiClient{Client: netClient}
	return api
}

func (a ApiClient) ConvertAllArtist() ([]Artists, int, error) {
	var allArtist []Artists

	res, err := a.Client.Get(url)
	if err != nil {
		return allArtist, 500, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return allArtist, 500, err
	}

	jsonErr := json.Unmarshal(body, &allArtist)
	if jsonErr != nil {
		return allArtist, 500, err
	}
	time.Sleep(5 * time.Millisecond)
	for i := 0; i < 52; i++ {
		go a.OneArtist(allArtist, i)
	}

	return allArtist, 0, err
}

func (a ApiClient) OneArtist(allArtist []Artists, i int) ([]Artists, int, error) {
	relation, err := a.Client.Get(allArtist[i].Relations)
	if err != nil {
		return allArtist, 404, err
	}
	relationBody, err := ioutil.ReadAll(relation.Body)
	if err != nil {
		return allArtist, 500, err
	}

	jsonErrRelation := json.Unmarshal(relationBody, &allArtist[i].RelationsData)
	if jsonErrRelation != nil {
		return allArtist, 500, err
	}

	location, err := a.Client.Get(allArtist[i].Locations)
	if err != nil {
		return allArtist, 404, err
	}

	locationBody, err := ioutil.ReadAll(location.Body)
	if err != nil {
		return allArtist, 500, err
	}

	jsonErrLocation := json.Unmarshal(locationBody, &allArtist[i].LocationsData)
	if jsonErrLocation != nil {
		return allArtist, 500, err
	}

	return allArtist, 0, nil
}

func (a ApiClient) ConvertOneArtist(id string) (Artists, int, error) {
	var data Artists

	idnumber, err := strconv.Atoi(id)
	if err != nil {
		return data, 500, err
	}

	if idnumber < 1 || idnumber > 52 {
		return data, 404, err
	}

	res, err := a.Client.Get(url + "/" + id)
	if err != nil {
		return data, 404, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return data, 500, err
	}

	jsonErr := json.Unmarshal(body, &data)
	if jsonErr != nil {
		return data, 500, err
	}

	relation, err := a.Client.Get(data.Relations)
	if err != nil {
		return data, 404, err
	}
	relationBody, err := ioutil.ReadAll(relation.Body)
	if err != nil {
		return data, 500, err
	}

	jsonErrRelation := json.Unmarshal(relationBody, &data.RelationsData)
	if jsonErrRelation != nil {
		return data, 500, err
	}

	location, err := a.Client.Get(data.Locations)
	if err != nil {
		return data, 404, err
	}

	locationBody, err := ioutil.ReadAll(location.Body)
	if err != nil {
		return data, 500, err
	}

	jsonErrLocation := json.Unmarshal(locationBody, &data.LocationsData)
	if jsonErrLocation != nil {
		return data, 500, err
	}

	return data, 0, nil
}
