package handlers

import (
	"net/http"
	"strconv"
	"strings"
)

func Convertstrint(str string) (int, error) {
	var res int
	if str == "" {
		return 0, nil
	}
	res, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func Filter(r *http.Request, allArtist []Artists) ([]Artists, int, error) {
	CreationDate0, err1 := Convertstrint(r.FormValue("creation-date_from"))
	CreationDate1, err2 := Convertstrint(r.FormValue("creation-date_to"))
	FirstAlbum0, err3 := Convertstrint(r.FormValue("first-album_from"))
	FirstAlbum1, err4 := Convertstrint(r.FormValue("first-album_to"))
	Members := r.PostForm["member"]
	Locations := r.FormValue("locations")
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return nil, 0, err1
	}

	var filtered []Artists

	for i := range allArtist {

		if allArtist[i].CreationDate < CreationDate0 || allArtist[i].CreationDate > CreationDate1 {
			continue
		}

		firstAlb, _ := strconv.Atoi(allArtist[i].FirstAlbum[6:])

		if firstAlb < FirstAlbum0 || firstAlb > FirstAlbum1 {
			continue
		}
		var isMemberNum bool

		if len(Members) == 0 {
			isMemberNum = true
		} else {
			for _, v := range Members {
				if vint, _ := strconv.Atoi(v); vint == len(allArtist[i].Members) {
					isMemberNum = true
					break
				}
			}
		}

		if !isMemberNum {
			continue
		}

		var isLocationFound bool

		for city := range allArtist[i].RelationsData.DatesLocations {
			location := strings.ReplaceAll(Locations, ", ", "-")
			location = strings.ReplaceAll(location, " ", "_")
			if strings.Contains(city, strings.ToLower(location)) {
				isLocationFound = true
			}
		}
		if !isLocationFound {
			continue
		}

		filtered = append(filtered, allArtist[i])

	}

	return filtered, 0, nil
}
