package main

import (
	"fmt"
	"github.com/Guillaume-Boutry/grpc-backend/pkg/face_authenticator"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"os"
)

type Person struct {
	Name                string
	EnrollPhoto         string
	AuthenticationPhoto []string
	FaceCoordinates     *face_authenticator.FaceCoordinates
}

func main() {
	address := "grpc-backend.default.127.0.0.1.nip.io:80"
	//address := "127.0.0.1:9000"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()
	//img, err := ioutil.ReadFile("/home/guillaume/go/src/go-face-test/luda.jpg")
	//img, err := ioutil.ReadFile("/tmp/vin.jpeg")

	clientEnroll := face_authenticator.NewEnrollerClient(conn)
	clientAuthenticate := face_authenticator.NewAuthenticatorClient(conn)

	person := &Person{
		Name:        "Jean valjean",
		EnrollPhoto: "/home/guillaume/Downloads/105_classes_pins_dataset/pins_tom ellis/tom ellis98_4428.jpg",
		FaceCoordinates: &face_authenticator.FaceCoordinates{
			TopLeft: &face_authenticator.Point{
				X: 58,
				Y: 10,
			},
			BottomRight: &face_authenticator.Point{
				X: 323,
				Y: 398,
			},
		},
	}
	//		AuthenticationPhoto: []string{"/home/guillaume/Downloads/105_classes_pins_dataset/pins_tom ellis/tom ellis52_4393.jpg"},
	//FaceCoordinates: &face_authenticator.FaceCoordinates{
	//			TopLeft: &face_authenticator.Point{
	//				X: 94,
	//				Y: 41,
	//			},
	//			BottomRight: &face_authenticator.Point{
	//				X: 354,
	//				Y: 412,
	//			},
	//		},
	person2 := &Person{
		Name:        "Jean valjean",
		AuthenticationPhoto: []string{"/home/guillaume/Downloads/105_classes_pins_dataset/pins_tom ellis/tom ellis98_4428.jpg"},
		FaceCoordinates: &face_authenticator.FaceCoordinates{
			TopLeft: &face_authenticator.Point{
				X: 58,
				Y: 10,
			},
			BottomRight: &face_authenticator.Point{
				X: 323,
				Y: 398,
			},
		},
	}

	enroll(&clientEnroll, person)
	authenticate(&clientAuthenticate, person2)
}

func enroll(enroller *face_authenticator.EnrollerClient, person *Person) bool {
	if person == nil {
		return false
	}
	img, _ := ioutil.ReadFile(person.EnrollPhoto)

	request := &face_authenticator.EnrollRequest{
		FaceRequest: &face_authenticator.FaceRequest{
			Id:              person.Name,
			Face:            img,
			FaceCoordinates: person.FaceCoordinates,
		},
	}
	res, err := (*enroller).Enroll(context.Background(), request)
	if err != nil || res.Status == face_authenticator.EnrollStatus_ENROLL_STATUS_ERROR {
		log.Printf("%s error, %s\n", person.Name, err)
		return false
	}
	log.Printf("%s enrolled, msg: %s\n", person.Name, res.Message)
	return true
}

func authenticate(authenticator *face_authenticator.AuthenticatorClient, person *Person) []bool {
	var results []bool
	var errors []error
	for _, photo := range person.AuthenticationPhoto {
		img, _ := ioutil.ReadFile(photo)

		request := &face_authenticator.AuthenticateRequest{
			FaceRequest: &face_authenticator.FaceRequest{
				Id:              person.Name,
				Face:            img,
				FaceCoordinates: person.FaceCoordinates,
			},
		}
		res, err := (*authenticator).Authenticate(context.Background(), request)
		if err != nil || res.Status == face_authenticator.AuthenticateStatus_AUTHENTICATE_STATUS_ERROR {
			log.Printf("%s error, %s\n", person.Name, err)
			errors = append(errors, err)
			break
		}
		log.Printf("%s authenticated, score: %f, decision: %t", person.Name, res.Score, res.Decision)
		results = append(results, res.Decision)
	}
	return results
}
