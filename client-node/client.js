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
    oneofs: true
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