package service

import (
	"context"
	authenticator "github.com/Guillaume-Boutry/face-authenticator-wrapper"
	"github.com/Guillaume-Boutry/grpc-backend/pkg/face_authenticator"
	"log"
)

type FeatureMatrix authenticator.Dlib_matrix_Sl_float_Sc_0_Sc_1_Sg_


type AuthenticateServiceGrpcImpl struct {
	face_authenticator.UnimplementedAuthenticatorServer
	jobChannel chan *work
}

//NewEnrollServiceGrpcImpl
func NewAuthenticateServiceGrpcImpl() *AuthenticateServiceGrpcImpl {
	jobChannel := make(chan *work)
	for w := 1; w <= 4; w++ {
		go worker(w, jobChannel)
	}
	return &AuthenticateServiceGrpcImpl{
		jobChannel: jobChannel,
	}
}

type work struct {
	faceRequest *face_authenticator.FaceRequest
	responseChannel     chan FeatureMatrix
}

func validRectangle(coordinates *face_authenticator.FaceCoordinates) bool {
	return coordinates.TopLeft != nil && coordinates.TopLeft.X != 0 && coordinates.TopLeft.Y != 0 && coordinates.BottomRight != nil && coordinates.BottomRight.X != 0 && coordinates.BottomRight.Y != 0
}

func worker(idThread int, jobs <-chan *work) {
	authent := authenticator.NewAuthenticator(32)
	defer authenticator.DeleteAuthenticator(authent)
	log.Printf("Thread %d: Init authenticator\n", idThread)
	authent.Init("/opt/grpc-backend/shape_predictor_5_face_landmarks.dat", "/opt/grpc-backend/dlib_face_recognition_resnet_model_v1.dat")
	log.Printf("Thread %d: Ready to authenticate\n", idThread)
	for job := range jobs {
		generateEmbeddings(&authent, job, idThread)
	}
}

func generateEmbeddings(authent *authenticator.Authenticator, work *work, idThread int) {
	facereq := work.faceRequest
	cImgData := authenticator.Load_mem_jpeg(&facereq.Face[0], len(facereq.Face))
	defer authenticator.DeleteImage(cImgData)
	var facePosition authenticator.Rectangle
	log.Printf("Thread %d: Searching for a face...\n", idThread)
	if coords := facereq.FaceCoordinates; coords == nil || !validRectangle(coords) {
		facePosition = (*authent).DetectFace(cImgData)
		defer authenticator.DeleteRectangle(facePosition)
	} else {
		facePosition = authenticator.NewRectangle()
		facePosition.SetTop(coords.TopLeft.Y)
		facePosition.SetLeft(coords.TopLeft.X)
		facePosition.SetBottom(coords.BottomRight.Y)
		facePosition.SetRight(coords.BottomRight.X)
	}
	log.Printf("Thread %d: Found face in area top_left(%d, %d), bottom_right(%d, %d)\n", idThread, facePosition.GetTop(), facePosition.GetLeft(), facePosition.GetBottom(), facePosition.GetRight(),)
	extractedFace := (*authent).ExtractFace(cImgData, facePosition)
	defer authenticator.DeleteImage(extractedFace)
	log.Printf("Thread %d: Generating embeddings\n", idThread)
	embeddings := (*authent).GenerateEmbeddings(extractedFace)
	work.responseChannel <- embeddings
}

func (serviceImpl *AuthenticateServiceGrpcImpl) Authenticate(ctx context.Context, request *face_authenticator.AuthenticateRequest) (*face_authenticator.AuthenticateResponse, error) {
	log.Println(request.FaceRequest.Id)

	if facereq := request.FaceRequest; facereq != nil {
		if face := facereq.Face; face != nil {
			responseChannel := make(chan FeatureMatrix)
			(*serviceImpl).jobChannel <- &work{
				faceRequest: facereq,
				responseChannel:     responseChannel,
			}
			embeddings := <-responseChannel
			var serialized [authenticator.EMBEDDINGS_SIZE]float32
			ptr := &serialized[0]
			authenticator.Serialize_embeddings(embeddings, ptr)
			log.Println(serialized)
			return &face_authenticator.AuthenticateResponse{
				Status:  face_authenticator.AuthenticateStatus_AUTHENTICATE_STATUS_OK,
				Message: "Worked",
			}, nil
		}
	}

	return &face_authenticator.AuthenticateResponse{
		Status:  face_authenticator.AuthenticateStatus_AUTHENTICATE_STATUS_ERROR,
		Message: "ok",
	}, nil
}
