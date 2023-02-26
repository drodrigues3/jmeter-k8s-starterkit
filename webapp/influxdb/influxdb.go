package influxdb 

mport (
	"fmt"
	"sync"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var token string
var once sync.Once

func createToken() error {
	client := influxdb2.NewClientWithOptions("http://localhost:8086", influxdb2.DefaultOptions().SetLogLevel(0))
	defer client.Close()

	username := "my-username"
	password := "my-password"
	org := "my-org"
	bucket := "my-bucket"

	authenticator := influxdb2.NewAuthentication(username, password)
	result, err := client.AuthorizationsAPI().CreateAuthorization(org, []influxdb2.Permission{{Action: influxdb2.ReadAction, Resource: influxdb2.Resource{Type: influxdb2.BucketsResourceType, ID: &bucket}}}, &authenticator)
	if err != nil {
		return err
	}

	token = result.Token
	return nil
}

func getToken() (string, error) {
	var err error
	once.Do(func() {
		err = createToken()
	})
	return token, err
}

func main() {
	token, err := getToken()
	if err != nil {
		fmt.Println("Error creating token:", err)
		return
	}

	fmt.Println("Token:", token)
}