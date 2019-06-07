package bumble

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Photo struct {
	ID         string `bson:"id"`
	PreviewURL string
	LargeURL   string

	// X,Y pairs for face detection.
	FaceTopLeft     [2]int
	FaceBottomRight [2]int

	// May be 0 for instagram photos.
	Width  int
	Height int
}

type Album struct {
	UID     string `bson:"uid"`
	Name    string
	Caption string
	Photos  []*Photo
}

type MusicArtist struct {
	ID   string `bson:"id"`
	Name string
}

type MusicService struct {
	ID          string `bson:"id"`
	DisplayName string
	Type        int
	TopArtists  []*MusicArtist
}

type ProfileField struct {
	ID           string `bson:"id"`
	Type         int
	Name         string
	DisplayValue string
}

type User struct {
	ID            string `bson:"id"`
	Name          string
	Age           int
	Gender        int
	Verified      bool
	DistanceLong  string
	DistanceShort string
	Albums        []*Album
	MusicServices []*MusicService
	ProfileFields []*ProfileField

	ScanDate time.Time
}

func (u *User) AllPhotos() []*Photo {
	var res []*Photo
	for _, album := range u.Albums {
		for _, ph := range album.Photos {
			res = append(res, ph)
		}
	}
	return res
}

func (u *User) Location() string {
	for _, field := range u.ProfileFields {
		if field.ID == "location" {
			return strings.Split(field.DisplayValue, "\n")[0]
		}
	}
	return "Unknown"
}

// BumbleAPI encapsulates all the required bumble calls to
// scan profiles.
type BumbleAPI struct {
	GetEncountersCall  GetEncountersCall
	DislikeCall        DislikeCall
	UpdateLocationCall UpdateLocationCall
}

// GetEncounters lists a small set of nearby users.
//
// If it returns an empty list, it likely means that no
// matches remain in the area.
func (b *BumbleAPI) GetEncounters() ([]*User, error) {
	req, err := b.GetEncountersCall.Request()
	if err != nil {
		return nil, errors.Wrap(err, "get encounters")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "get encounters")
	}
	defer resp.Body.Close()

	var responseObj struct {
		Body []struct {
			ServerErrorMessage *struct {
				ErrorMessage string `json:"error_message"`
			} `json:"server_error_message"`

			ClientEncounters struct {
				Results []struct {
					User struct {
						UserID        string `json:"user_id"`
						Name          string `json:"name"`
						Age           int    `json:"age"`
						Gender        int    `json:"gender"`
						Verified      bool   `json:"is_verified"`
						DistanceLong  string `json:"distance_long"`
						DistanceShort string `json:"distance_short"`
						Albums        []struct {
							UID     string `json:"uid"`
							Name    string `json:"name"`
							Caption string `json:"caption"`
							Photos  []struct {
								ID             string `json:"id"`
								PreviewURL     string `json:"preview_url"`
								LargeURL       string `json:"large_url"`
								LargePhotoSize struct {
									Width  int `json:"width"`
									Height int `json:"height"`
								} `json:"large_photo_size"`
								FaceTopLeft struct {
									X int `json:"x"`
									Y int `json:"y"`
								} `json:"face_top_left"`
								FaceBottomRight struct {
									X int `json:"x"`
									Y int `json:"y"`
								} `json:"face_bottom_right"`
							} `json:"photos"`
						} `json:"albums"`
						MusicServices []struct {
							Status           int `json:"status"`
							ExternalProvider struct {
								ID          string `json:"id"`
								DisplayName string `json:"display_name"`
								Type        int    `json:"type"`
							} `json:"external_provider"`
							TopArtists []struct {
								ID   string `json:"id"`
								Name string `json:"name"`
							} `json:"top_artists"`
						} `json:"music_services"`
						ProfileFields []struct {
							ID           string `json:"id"`
							Type         int    `json:"type"`
							Name         string `json:"name"`
							DisplayValue string `json:"display_value"`
						} `json:"profile_fields"`
					} `json:"user"`
				} `json:"results"`
			} `json:"client_encounters"`
		} `json:"body"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&responseObj); err != nil {
		return nil, errors.Wrap(err, "get encounters")
	}

	var users []*User
	for _, body := range responseObj.Body {
		if body.ServerErrorMessage != nil {
			return nil, errors.New(body.ServerErrorMessage.ErrorMessage)
		}
		for _, result := range body.ClientEncounters.Results {
			rawUser := result.User
			user := User{
				ID:            rawUser.UserID,
				Name:          rawUser.Name,
				Age:           rawUser.Age,
				Gender:        rawUser.Gender,
				Verified:      rawUser.Verified,
				DistanceLong:  rawUser.DistanceLong,
				DistanceShort: rawUser.DistanceShort,
				ScanDate:      time.Now(),
			}
			for _, rawAlbum := range rawUser.Albums {
				album := &Album{
					UID:     rawAlbum.UID,
					Name:    rawAlbum.Name,
					Caption: rawAlbum.Caption,
				}
				for _, rawPhoto := range rawAlbum.Photos {
					album.Photos = append(album.Photos, &Photo{
						ID:              rawPhoto.ID,
						PreviewURL:      rawPhoto.PreviewURL,
						LargeURL:        rawPhoto.LargeURL,
						FaceTopLeft:     [2]int{rawPhoto.FaceTopLeft.X, rawPhoto.FaceTopLeft.Y},
						FaceBottomRight: [2]int{rawPhoto.FaceBottomRight.X, rawPhoto.FaceBottomRight.Y},
						Width:           rawPhoto.LargePhotoSize.Width,
						Height:          rawPhoto.LargePhotoSize.Height,
					})
				}
				user.Albums = append(user.Albums, album)
			}
			for _, rawMusicService := range rawUser.MusicServices {
				musicService := &MusicService{
					ID:          rawMusicService.ExternalProvider.ID,
					DisplayName: rawMusicService.ExternalProvider.DisplayName,
					Type:        rawMusicService.ExternalProvider.Type,
				}
				for _, rawArtist := range rawMusicService.TopArtists {
					musicService.TopArtists = append(musicService.TopArtists, &MusicArtist{
						ID:   rawArtist.ID,
						Name: rawArtist.Name,
					})
				}
				user.MusicServices = append(user.MusicServices, musicService)
			}
			for _, rawProfileField := range rawUser.ProfileFields {
				user.ProfileFields = append(user.ProfileFields, &ProfileField{
					ID:           rawProfileField.ID,
					Type:         rawProfileField.Type,
					Name:         rawProfileField.Name,
					DisplayValue: rawProfileField.DisplayValue,
				})
			}
			users = append(users, &user)
		}
	}
	return users, nil
}

// Dislike dislikes a user by their ID.
func (b *BumbleAPI) Dislike(userID string) error {
	if err := b.doRequest(b.DislikeCall.Request(userID)); err != nil {
		return errors.Wrap(err, "dislike")
	}
	return nil
}

// UpdateLocation updates the user's location.
func (b *BumbleAPI) UpdateLocation(lat, lon float64) error {
	if err := b.doRequest(b.UpdateLocationCall.Request(lat, lon)); err != nil {
		return errors.Wrap(err, "update location")
	}
	return nil
}

func (b *BumbleAPI) doRequest(req *http.Request, err error) error {
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
