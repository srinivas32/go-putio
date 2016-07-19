package putio

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// ZipsService is the service manage zip streams.
type ZipsService struct {
	client *Client
}

// Get gives detailed information about the given zip file id.
func (z *ZipsService) Get(id int) (Zip, error) {
	if id < 0 {
		return Zip{}, errNegativeID
	}

	req, err := z.client.NewRequest("GET", "/v2/zips/"+strconv.Itoa(id), nil)
	if err != nil {
		return Zip{}, err
	}

	var r Zip
	_, err = z.client.Do(req, &r)
	if err != nil {
		return Zip{}, err
	}

	return r, nil

}

// List lists active zip files.
func (z *ZipsService) List() ([]Zip, error) {
	req, err := z.client.NewRequest("GET", "/v2/zips/list", nil)
	if err != nil {
		return nil, err
	}

	var r struct {
		Zips []Zip
	}
	_, err = z.client.Do(req, &r)
	if err != nil {
		return nil, err
	}

	return r.Zips, nil
}

// Create creates zip files for given file IDs. If the operation is successful,
// a zip ID will be returned to keep track of zip process.
func (z *ZipsService) Create(fileIDs ...int) (int, error) {
	if len(fileIDs) == 0 {
		return 0, fmt.Errorf("no file id given")
	}

	var ids []string
	for _, id := range fileIDs {
		if id < 0 {
			return 0, errNegativeID
		}
		ids = append(ids, strconv.Itoa(id))
	}

	params := url.Values{}
	params.Set("file_ids", strings.Join(ids, ","))

	req, err := z.client.NewRequest("POST", "/v2/zips/create", strings.NewReader(params.Encode()))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var r struct {
		ID int `json:"zip_id"`
	}
	_, err = z.client.Do(req, &r)
	if err != nil {
		return 0, err
	}

	return r.ID, nil
}
