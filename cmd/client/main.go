package main

import (
	"context"
	"fmt"
	"github.com/Guillaume-Boutry/grpc-backend/pkg/face_authenticator"
	"google.golang.org/grpc"
	"io/ioutil"
	"os"
)

func main() {
	address := "127.0.0.1:9000"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	client := face_authenticator.NewEnrollerClient(conn)
	request := face_authenticator.EnrollRequest{
		FaceRequest: &face_authenticator.FaceRequest{
			Id: "test",
		},
	}

	res, err := client.Enroll(context.Background(), &request)
	fmt.Println(res)
	fmt.Println(err)

	img, err := ioutil.ReadFile("/home/guillaume/go/src/go-face-test/luda2.jpg")

	if err != nil {
		fmt.Println(err)
		return
	}

	clientAuthent := face_authenticator.NewAuthenticatorClient(conn)
	requestAuthent := face_authenticator.AuthenticateRequest{
		FaceRequest: &face_authenticator.FaceRequest{
			Id:   "Send authent",
			Face: img,
			FaceCoordinates: &face_authenticator.FaceCoordinates{
				TopLeft: &face_authenticator.Point{
					X: 362,
					Y: 404,
				},
				BottomRight: &face_authenticator.Point{
					X: 734,
					Y: 775,
				},
			},
		},
	}

	res2, err2 := clientAuthent.Authenticate(context.Background(), &requestAuthent)
	fmt.Println(res2)
	fmt.Println(err2)
}
