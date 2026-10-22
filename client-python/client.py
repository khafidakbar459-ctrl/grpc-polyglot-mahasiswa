import grpc

import mahasiswa_pb2
import mahasiswa_pb2_grpc


def main():
    # Membuka koneksi ke server Go
    with grpc.insecure_channel("localhost:50051") as channel:
        # Membuat client berdasarkan layanan di file proto
        stub = mahasiswa_pb2_grpc.MahasiswaServiceStub(channel)

        # Meminta data mahasiswa dengan ID 99
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