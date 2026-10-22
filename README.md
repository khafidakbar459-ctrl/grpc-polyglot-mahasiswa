# Praktikum gRPC Polyglot: Server Go + Client Python & Node.js

Langkah-langkah membuat server gRPC dengan **Go** dan dua client (**Python** dan **Node.js**) yang memakai satu kontrak Protocol Buffers.

## Prasyarat

Pastikan sudah terpasang (versi saat praktikum):

- Go 1.22.2
- Python 3.12.3
- Node.js v18.19.1 dan npm 9.2.0
- protoc 3.21.12

```bash
go version && python3 --version && node -v && npm -v && protoc --version
```

## Langkah 1: Siapkan Folder Proyek

```bash
mkdir -p praktikum-rpg/{proto,server-go/mahasiswa,client-python,client-node}
cd praktikum-rpg
```

## Langkah 2: Buat Kontrak `.proto`

Buat `proto/mahasiswa.proto`:

```proto
syntax = "proto3";

package mahasiswa;

service MahasiswaService {
  rpc GetMahasiswa (MahasiswaRequest) returns (MahasiswaResponse);
}

message MahasiswaRequest {
  int32 id = 1;
}

message MahasiswaResponse {
  int32 id = 1;
  string nama = 2;
  string jurusan = 3;
}
```

## Langkah 3: Server Go

**3.1 Inisialisasi modul dan pasang dependensi**

```bash
cd server-go
go mod init mahasiswa-grpc
go get google.golang.org/grpc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
cd ..
export PATH="$PATH:$(go env GOPATH)/bin"
```

**3.2 Generate kode Go** (jalankan dari folder `praktikum-rpg`)

```bash
protoc --proto_path=proto \
  --go_out=server-go/mahasiswa \
  --go_opt=paths=source_relative \
  --go_opt=Mmahasiswa.proto=mahasiswa-grpc/mahasiswa \
  --go-grpc_out=server-go/mahasiswa \
  --go-grpc_opt=paths=source_relative \
  --go-grpc_opt=Mmahasiswa.proto=mahasiswa-grpc/mahasiswa \
  proto/mahasiswa.proto
```

Hasilnya: `mahasiswa.pb.go` dan `mahasiswa_grpc.pb.go` di `server-go/mahasiswa/`.

**3.3 Tulis `server-go/main.go`**

```go
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
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Gagal membuka port: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMahasiswaServiceServer(s, &server{})

	log.Println("Server gRPC berjalan pada port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Server gagal berjalan: %v", err)
	}
}
```

**3.4 Rapikan dependensi**

```bash
cd server-go
go mod tidy
cd ..
```

## Langkah 4: Client Python

**4.1 Buat virtual environment dan pasang pustaka**

```bash
cd client-python
python3 -m venv .venv
source .venv/bin/activate
pip install grpcio grpcio-tools
cd ..
```

**4.2 Generate kode Python** (dari folder `praktikum-rpg`, venv aktif)

```bash
python -m grpc_tools.protoc --proto_path=proto \
  --python_out=client-python \
  --grpc_python_out=client-python \
  proto/mahasiswa.proto
```

Hasilnya: `mahasiswa_pb2.py` dan `mahasiswa_pb2_grpc.py`.

**4.3 Tulis `client-python/client.py`**

```python
import grpc

import mahasiswa_pb2
import mahasiswa_pb2_grpc


def main():
    # Membuka koneksi ke server Go
    with grpc.insecure_channel("localhost:50051") as channel:
        # Membuat client berdasarkan layanan di file proto
        stub = mahasiswa_pb2_grpc.MahasiswaServiceStub(channel)

        # Meminta data mahasiswa dengan ID 1
        request = mahasiswa_pb2.MahasiswaRequest(id=1)

        try:
            response = stub.GetMahasiswa(request)

            print("=== DATA MAHASISWA ===")
            print(f"ID      : {response.id}")
            print(f"Nama    : {response.nama}")
            print(f"Jurusan : {response.jurusan}")

        except grpc.RpcError as error:
            print("Permintaan gagal!")
            print(f"Status  : {error.code()}")
            print(f"Pesan   : {error.details()}")


if __name__ == "__main__":
    main()
```

## Langkah 5: Client Node.js

**5.1 Inisialisasi dan pasang pustaka**

```bash
cd client-node
npm init -y
npm install @grpc/grpc-js @grpc/proto-loader
cd ..
```

Client Node.js **tidak perlu generate kode**, karena file `.proto` dibaca langsung saat program berjalan.

**5.2 Tulis `client-node/client.js`**

```javascript
const grpc = require("@grpc/grpc-js");
const protoLoader = require("@grpc/proto-loader");
const path = require("path");

// Lokasi file kontrak Protobuf
const PROTO_PATH = path.join(__dirname, "../proto/mahasiswa.proto");

// Membaca definisi layanan dari file .proto
const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: false,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
});

// Memuat paket mahasiswa
const mahasiswaProto = grpc.loadPackageDefinition(packageDefinition).mahasiswa;

// Membuat client yang terhubung ke server Go
const client = new mahasiswaProto.MahasiswaService(
  "localhost:50051",
  grpc.credentials.createInsecure()
);

// Meminta data mahasiswa dengan ID 1
client.getMahasiswa({ id: 1 }, (error, response) => {
  if (error) {
    console.error("Permintaan gagal!");
    console.error("Status:", error.code);
    console.error("Pesan:", error.details);
    client.close();
    return;
  }

  console.log("=== DATA MAHASISWA ===");
  console.log("ID      :", response.id);
  console.log("Nama    :", response.nama);
  console.log("Jurusan :", response.jurusan);

  client.close();
});
```

## Langkah 6: Jalankan dan Uji

Buka **tiga terminal**. Server harus berjalan lebih dulu.

**Terminal 1: server Go**

```bash
cd server-go
go run main.go
```

Keluaran: `Server gRPC berjalan pada port 50051`

**Terminal 2: client Python**

```bash
cd client-python
source .venv/bin/activate
python client.py
```

**Terminal 3: client Node.js**

```bash
cd client-node
node client.js
```

Keluaran yang diharapkan pada kedua client:

```
=== DATA MAHASISWA ===
ID      : 1
Nama    : Muhammad Khafid Akbar
Jurusan : Teknologi Rekayasa Komputer dan Jaringan
```

## Langkah 7: Uji ID Tidak Ditemukan

Ubah `id=1` menjadi `id=99` pada `client.py` (`MahasiswaRequest(id=99)`) dan `{ id: 1 }` menjadi `{ id: 99 }` pada `client.js`, lalu jalankan ulang kedua client.

- Python: `Status : StatusCode.NOT_FOUND`, `Pesan : Mahasiswa tidak ditemukan`
- Node.js: `Status: 5`, `Pesan: Mahasiswa tidak ditemukan` (5 = kode numerik `NOT_FOUND`)

## Troubleshooting

- `protoc-gen-go: program not found`: tambahkan `$(go env GOPATH)/bin` ke `PATH`.
- `ModuleNotFoundError: grpc` di Python: aktifkan venv (`source .venv/bin/activate`).
- Client gagal terhubung (`UNAVAILABLE`): pastikan server Go sudah berjalan di port 50051.

## `.gitignore`

```
.venv/
__pycache__/
node_modules/
```
