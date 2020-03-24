package main

type GetInfosResponse struct {
	Current Track   `json:"current"`
	History []Track `json:"history"`
}

type Track struct {
	Comment                        string `json:"comment"`
	Genre                          string `json:"genre,omitempty"`
	Rid                            string `json:"rid"`
	OnAir                          string `json:"on_air"`
	Status                         string `json:"status"`
	InitialURI                     string `json:"initial_uri"`
	Tracknumber                    string `json:"tracknumber,omitempty"`
	Source                         string `json:"source"`
	Temporary                      string `json:"temporary"`
	Filename                       string `json:"filename"`
	Title                          string `json:"title"`
	Decoder                        string `json:"decoder"`
	Artist                         string `json:"artist"`
	Year                           string `json:"year,omitempty"`
	Kind                           string `json:"kind"`
	Date                           string `json:"date,omitempty"`
	Album                          string `json:"album,omitempty"`
	LeftTitle                      string `json:"left_title"`
	RightTitle                     string `json:"right_title"`
	Live                           int    `json:"live"`
	Mode                           string `json:"mode"`
	FullTitle                      string `json:"full_title"`
	Albumartist                    string `json:"albumartist,omitempty"`
	Discnumber                     string `json:"discnumber,omitempty"`
	Encodedby                      string `json:"encodedby,omitempty"`
	Composer                       string `json:"composer,omitempty"`
	Catalog                        string `json:"catalog #,omitempty"`
	RippingTool                    string `json:"ripping tool,omitempty"`
	URL                            string `json:"url,omitempty"`
	RipDate                        string `json:"rip date,omitempty"`
	ReleaseType                    string `json:"release type,omitempty"`
	Label                          string `json:"label,omitempty"`
	Encoding                       string `json:"encoding,omitempty"`
	Language                       string `json:"language,omitempty"`
	Supplier                       string `json:"supplier,omitempty"`
	ReplaygainAlbumPeak            string `json:"replaygain_album_peak,omitempty"`
	MusicbrainzAlbumReleaseCountry string `json:"musicbrainz album release country,omitempty"`
	Catalognumber                  string `json:"catalognumber,omitempty"`
	Artistsort                     string `json:"artistsort,omitempty"`
	Media                          string `json:"media,omitempty"`
	ReplaygainAlbumGain            string `json:"replaygain_album_gain,omitempty"`
	ReplaygainTrackGain            string `json:"replaygain_track_gain,omitempty"`
	MusicbrainzReleaseTrackID      string `json:"musicbrainz release track id,omitempty"`
	MusicbrainzAlbumid             string `json:"musicbrainz_albumid,omitempty"`
	Script                         string `json:"script,omitempty"`
	MusicbrainzReleasegroupid      string `json:"musicbrainz_releasegroupid,omitempty"`
	Originalyear                   string `json:"originalyear,omitempty"`
	Asin                           string `json:"asin,omitempty"`
	Barcode                        string `json:"barcode,omitempty"`
	Originaldate                   string `json:"originaldate,omitempty"`
	AcoustidID                     string `json:"acoustid_id,omitempty"`
	MusicbrainzArtistid            string `json:"musicbrainz_artistid,omitempty"`
	MusicbrainzAlbumStatus         string `json:"musicbrainz album status,omitempty"`
	Albumartistsort                string `json:"albumartistsort,omitempty"`
	MusicbrainzAlbumType           string `json:"musicbrainz album type,omitempty"`
	MusicbrainzAlbumartistid       string `json:"musicbrainz_albumartistid,omitempty"`
	Artists                        string `json:"artists,omitempty"`
	MusicbrainzTrackid             string `json:"musicbrainz_trackid,omitempty"`
	ReplaygainTrackPeak            string `json:"replaygain_track_peak,omitempty"`
}

func (t Track) Pretty() string {
	out := t.FullTitle
	return out
}
