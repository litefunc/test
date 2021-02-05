package benchmark

import (
	"encoding/json"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

type MediumPayload struct {
	Success bool `json:"success"`
	Data    struct {
		UpdateTime int64 `json:"update_time"`
		List       []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
			Meta  struct {
				SubLangs         []string `json:"sub_langs"`
				Summary          string   `json:"summary"`
				ShortDescription string   `json:"short_description"`
				Nation           string   `json:"nation"`
				Director         string   `json:"director"`
				ReleaseYear      string   `json:"release_year"`
				Title            string   `json:"title"`
				Type             string   `json:"type"`
				Duration         int      `json:"duration"`
				Actors           []string `json:"actors"`
				Rate             float64  `json:"rate"`
				Season           int      `json:"season"`
				Categories       []string `json:"categories"`
			} `json:"meta"`
			Episodes []struct {
				ID    string `json:"id"`
				Video struct {
					Original struct {
						UploadTime int64  `json:"uploadTime"`
						FileSize   int    `json:"fileSize"`
						FilePath   string `json:"filePath"`
					} `json:"original"`
				} `json:"video"`
				Title string `json:"title"`
				Meta  struct {
				} `json:"meta"`
			} `json:"episodes"`
			Cover struct {
				UploadTime int64  `json:"uploadTime"`
				Thumbnail  string `json:"thumbnail"`
				Original   string `json:"original"`
				Filesize   int    `json:"filesize"`
			} `json:"cover"`
		} `json:"list"`
	} `json:"data"`
}

const s = `{
    "success": true,
    "data": {
        "update_time": 1598836395509,
        "list": [
            {
                "id": "VOD-01EWJ6JPP8Q1AKAPMN862YZ724",
                "title": "Avengers: End Game test2",
                "meta": {
                    "sub_langs": [
                        "English"
                    ],
                    "summary": "Avengers: Endgame is a 2019 American superhero film based on the Marvel Comics superhero team the Avengers. Produced by Marvel Studios and distributed by Walt Disney Studios Motion Pictures, it is the direct sequel to Avengers: Infinity War (2018) and the 22nd film in the Marvel Cinematic Universe (MCU).",
                    "short_description": "Avengers: Endgame is a 2019 American superhero film based on the Marvel Comics superhero team the Avengers. Produced by Marvel Studios and distributed by Walt Disney Studios Motion Pictures, it is the direct sequel to Avengers: Infinity War (2018) and the 22nd film in the Marvel Cinematic Universe (MCU).",
                    "nation": "Earth",
                    "director": "SomeBody",
                    "release_year": "2019",
                    "title": "Avengers: End Game test2",
                    "type": "movie",
                    "duration": 120,
                    "actors": [
                        "Alice",
                        "Bob"
                    ],
                    "rate": 4.8,
                    "season": 1,
                    "categories": [
                        "Sci-FI",
                        "2019"
                    ]
                },
                "episodes": [
                    {
                        "id": "01EWJ6JRJC6MYT4TRR082R8A77",
                        "video": {
                            "original": {
                                "uploadTime": 1611224513898,
                                "fileSize": 17839845,
                                "filePath": "videos/vod/VOD-01EWJ6JPP8Q1AKAPMN862YZ724/01EWJ6JRJC6MYT4TRR082R8A77.mp4"
                            }
                        },
                        "title": "Trailor",
                        "meta": {}
                    }
                ],
                "cover": {
                    "uploadTime": 1611224947000,
                    "thumbnail": "thumb/images/vod/cover/VOD-01EWJ6JPP8Q1AKAPMN862YZ724.jpg",
                    "original": "images/vod/cover/VOD-01EWJ6JPP8Q1AKAPMN862YZ724.jpg",
                    "filesize": 227664
                }
            }
        ]
    }
}`

var mediumFixture = []byte(s)

/*
   encoding/json
*/
func BenchmarkDecodeStdStructMedium(b *testing.B) {
	b.ReportAllocs()
	var data MediumPayload
	for i := 0; i < b.N; i++ {
		json.Unmarshal(mediumFixture, &data)
	}
}

func BenchmarkEncodeStdStructMedium(b *testing.B) {
	var data MediumPayload
	json.Unmarshal(mediumFixture, &data)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		json.Marshal(data)
	}
}

func BenchmarkDecodeJsoniterStructMedium(b *testing.B) {
	b.ReportAllocs()
	var data MediumPayload
	for i := 0; i < b.N; i++ {
		jsoniter.Unmarshal(mediumFixture, &data)
	}
}

func BenchmarkEncodeJsoniterStructMedium(b *testing.B) {
	var data MediumPayload
	jsoniter.Unmarshal(mediumFixture, &data)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		jsoniter.Marshal(data)
	}
}

// func BenchmarkDecodeEasyJsonMedium(b *testing.B) {
// 	b.ReportAllocs()
// 	var data MediumPayload
// 	for i := 0; i < b.N; i++ {
// 		lexer := &jlexer.Lexer{Data: mediumFixture}
// 		data.UnmarshalEasyJSON(lexer)
// 	}
// }

// func BenchmarkEncodeEasyJsonMedium(b *testing.B) {
// 	var data MediumPayload
// 	lexer := &jlexer.Lexer{Data: mediumFixture}
// 	data.UnmarshalEasyJSON(lexer)
// 	b.ReportAllocs()
// 	buf := &bytes.Buffer{}
// 	for i := 0; i < b.N; i++ {
// 		writer := &jwriter.Writer{}
// 		data.MarshalEasyJSON(writer)
// 		buf.Reset()
// 		writer.DumpTo(buf)
// 	}
// }
