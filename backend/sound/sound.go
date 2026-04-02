package sound

import (
	"embed"
	"strings"
	"unicode"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
)

//go:embed *.mp3
var tracksFs embed.FS

type Tracks struct {
	Name     string
	FileName string
	Buffer   *beep.Buffer
	Format   beep.Format
}

func LoadTracks() (tracks []Tracks, err error) {
	entries, err := tracksFs.ReadDir(".")
	if err != nil {
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		fileName := entry.Name()
		f, err := tracksFs.Open(fileName)
		if err != nil {
			return nil, err
		}
		stream, format, decodeErr := mp3.Decode(f)
		if decodeErr != nil {
			return nil, decodeErr
		}
		buf := beep.NewBuffer(format)
		buf.Append(stream)
		stream.Close()
		tracks = append(tracks, Tracks{
			Name:     formatName(fileName),
			FileName: fileName,
			Buffer:   buf,
			Format:   format,
		})
	}
	return
}

func formatName(fileName string) string {
	name := strings.TrimSuffix(fileName, ".mp3")
	name = strings.NewReplacer("_", " ", "-", " ").Replace(name)
	words := strings.Fields(name)
	for i, word := range words {
		r := []rune(word)
		r[0] = unicode.ToUpper(r[0])
		words[i] = string(r)
	}
	return strings.Join(words, " ")
}
