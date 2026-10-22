package main

import (
	"context"
	"log"
	"net"

	pb "mahasiswa-grpc/mahasiswa"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedMahasiswaServiceServer
}

func (s *server) GetMahasiswa(
	ctx context.Context,
	req *pb.MahasiswaRequest,
) (*pb.MahasiswaResponse, error) {

	log.Printf("Permintaan data mahasiswa dengan ID: %d", req.GetId())

	if req.GetId() != 1 {
		return nil, status.Error(
			codes.NotFound,
			"Mahasiswa tidak ditemukan",
		)
	}

	return &pb.MahasiswaResponse{
		Id:      1,
		Nama:    "Muhammad Khafid Akbar",
		Jurusan: "Teknologi Rekayasa Komputer dan Jaringan",
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Gagal membuka port: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterMahasiswaServiceServer(
		grpcServer,
		&server{},
	)

	log.Println("Server gRPC berjalan pada port 50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Server gagal berjalan: %v", err)
	}
}