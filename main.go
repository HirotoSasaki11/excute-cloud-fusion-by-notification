package function

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/functions/metadata"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/urlshortener/v1"
)

// GCSEvent is the payload of a GCS event.
type GCSEvent struct {
	Kind                    string                 `json:"kind"`
	ID                      string                 `json:"id"`
	SelfLink                string                 `json:"selfLink"`
	Name                    string                 `json:"name"`
	Bucket                  string                 `json:"bucket"`
	Generation              string                 `json:"generation"`
	Metageneration          string                 `json:"metageneration"`
	ContentType             string                 `json:"contentType"`
	TimeCreated             time.Time              `json:"timeCreated"`
	Updated                 time.Time              `json:"updated"`
	TemporaryHold           bool                   `json:"temporaryHold"`
	EventBasedHold          bool                   `json:"eventBasedHold"`
	RetentionExpirationTime time.Time              `json:"retentionExpirationTime"`
	StorageClass            string                 `json:"storageClass"`
	TimeStorageClassUpdated time.Time              `json:"timeStorageClassUpdated"`
	Size                    string                 `json:"size"`
	MD5Hash                 string                 `json:"md5Hash"`
	MediaLink               string                 `json:"mediaLink"`
	ContentEncoding         string                 `json:"contentEncoding"`
	ContentDisposition      string                 `json:"contentDisposition"`
	CacheControl            string                 `json:"cacheControl"`
	Metadata                map[string]interface{} `json:"metadata"`
	CRC32C                  string                 `json:"crc32c"`
	ComponentCount          int                    `json:"componentCount"`
	Etag                    string                 `json:"etag"`
	CustomerEncryption      struct {
		EncryptionAlgorithm string `json:"encryptionAlgorithm"`
		KeySha256           string `json:"keySha256"`
	}
	KMSKeyName    string `json:"kmsKeyName"`
	ResourceState string `json:"resourceState"`
}

var config = &oauth2.Config{
	ClientID:     os.Getenv("CLIENTID"),
	ClientSecret: os.Getenv("CLIENTSECRET"),
	Endpoint:     google.Endpoint,
	Scopes:       []string{urlshortener.UrlshortenerScope},
}

//  GCS event
func Excute(ctx context.Context, e GCSEvent) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("metadata.FromContext: %v", err)
	}

	m := make(map[string]string)
	m["bucket"] = e.Bucket
	m["filename"] = meta.Resource.Name

	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	endpoint := os.Getenv("ENDPOINT")
	name := os.Getenv("PIPLINENAME")
	req, err := http.NewRequest(
		"POST",
		endpoint+"/api/v3/namespaces/default/apps/"+name+"/workflows/DataPipelineWorkflow/start",
		bytes.NewBuffer(data),
	)

	source := &oauth2.Token{
		RefreshToken: os.Getenv("RERFESH_TOKEN"),
	}

	c, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	s := config.TokenSource(c, source)
	token, err := s.Token()
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", token.AccessToken)

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return err
	} else if res.StatusCode != http.StatusOK {
		return nil
	}
	return nil
}
